package main

import (
	"bufio"
	"fmt"
	"github.com/integrii/flaggy"
	concur "github.com/korovkin/limiter"
	"github.com/olekukonko/tablewriter"
	cmap "github.com/orcaman/concurrent-map"
	"runtime/debug"

	"github.com/gookit/color"
	//"github.com/schollz/progressbar"
	"github.com/sirupsen/logrus"
	"github.com/tcnksm/go-input"
	"gopkg.in/ini.v1"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	//"github.com/asticode/go-astilectron"
)

const version = "0.2a"

type ConnectionServerConfig struct {
	mailboxServers       string
	ldapServers          []string
	serverTypePreference string
	adminUsername        string
	adminPassword        string
	adminPort            string
	adminProtocol        string
	socksServerString    string
	adminAuthToken       string
	allServers           cmap.ConcurrentMap
	skipSSLChecks        bool
	useSocks5Proxy       bool
	printedSSLSkipNotice bool
}

var inputFile = ""
var serverStats = false
var password = ""
var module = ""
var task = ""
var auditType = ""
var auditValue = ""
var log = logrus.New()
var logTimestamp string
var problemAccounts = cmap.New()
var threads = ""
var auditHeaderWritten = false
var averageResponseTimeMaker = cmap.New()
var userServerMapping = cmap.New()
var noSuchItem = cmap.New()
var messageScanLimit string
var scanLimit int
var alreadySeen []string
var serverTracker = cmap.New()

func init() {

	flaggy.SetName("Zimbra Experience (zmxp)")
	flaggy.SetDescription("A utility to gather information about the quality of the end user Zimbra Experience.")
	flaggy.DefaultParser.ShowHelpOnUnexpected = true
	flaggy.DefaultParser.AdditionalHelpPrepend = "http://github.com/zimbra/zmxp"
	flaggy.String(&threads, "th", "threads", "The Max number of threads to spawn for checking via soap and screenshots.")
	flaggy.String(&inputFile, "f", "file", "Provide an input file with a list of accounts to test (--file=myfile.txt). The file should consist of user email addresses, one per line.")
	flaggy.Bool(&serverStats, "ss", "server-stats", "Gather server stats for success vs failure.")
	flaggy.String(&password, "p", "password", "Manually define a password.")
	flaggy.String(&module, "m", "module", "The module to use.")
	flaggy.String(&task, "t", "task", "The task to to perform.")
	flaggy.String(&auditType, "at", "audit-type", "When 'mail-audit' is specified as the task, this is the audit type.")
	flaggy.String(&auditValue, "av", "audit-value", "When 'mail-audit' is specified as the task, this is the audit value.")
	flaggy.String(&messageScanLimit, "sl", "scan-limit", "When 'mail-audit' is specified as the task, this is the maximum amount of messages which will be scanned per account.")

	flaggy.SetVersion("zmxp " + version)
	flaggy.Parse()
}

func main() {

	logTimestamp = strconv.FormatInt(time.Now().Unix(), 10)

	concurrency := 5
	var wg sync.WaitGroup
	wg.Add(concurrency)
	log.Out = os.Stdout
	log.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02-01-2006 15:04:05",
	}
	knownTasks := []string{"mail-audit", "screenshot"}
	knownAuditTypes := []string{"SearchByDate"}
	log.Info("ZMXP " + version + ": starting up.")
	if runtime.GOOS == "windows" {
		fmt.Println("Hello from Windows")
	}
	taskSpecified := false
	auditSpecified := false
	taskStringBuilder := ""
	auditStringBuilder := ""
	//var scanLimit int
	var maxThreads int
	if threads == "" {
		threads = "15"
		maxThreads = 15
	} else {
		if n, err := strconv.Atoi(threads); err == nil {
			maxThreads = n
			if maxThreads > 20 {
				log.Fatal("You are insane. No.")
			}
		} else {
			log.Fatal(threads + " is not a valid thread count.")
		}
	}
	log.Info("zmX will be using " + threads + " threads.")
	if task == "" {
		log.Fatal("No Task was given.")
	} else {
		for _, v := range knownTasks {
			if strings.ToLower(task) == v {
				taskSpecified = true
			}
			taskStringBuilder += "--task=" + v + ", "
		}
	}

	if !taskSpecified {
		taskStringBuilder = strings.TrimRight(taskStringBuilder, ", ")
		log.Fatal("An unknown task was specified. The available tasks are: " + taskStringBuilder)
	} else {
		//log.Fatal(auditType)
		if task == "mail-audit" {
			for _, v := range knownAuditTypes {
				if strings.ToLower(auditType) == strings.ToLower(v) {
					auditSpecified = true
				}
				auditStringBuilder += "--audit-type=" + v + ", "
			}
			if !auditSpecified {
				auditStringBuilder = strings.TrimRight(auditStringBuilder, ", ")
				log.Fatal("An unknown audit type was specified. The available types are: " + auditStringBuilder)
			} else {
				if auditValue == "" {
					log.Fatal("When using a mail-audit, you must also specify an '--audit-value=[something]'. Example: '--audit-value=before:10/10/2019'")
				}
			}
			if messageScanLimit == "" {
				scanLimit = 1000
			}
		}

		var mode = "file"
		var saveToConfig bool

		ui := &input.UI{
			Writer: os.Stdout,
			Reader: os.Stdin,
		}

		ConnectionSettings := ConnectionServerConfig{}
		ConnectionSettings.adminAuthToken = "NONE"
		ConnectionSettings.printedSSLSkipNotice = false
		if inputFile != "" {
			if fileExists(inputFile) {
				log.Info("Using the input file: " + inputFile)
			} else {
				log.Fatal("Cannot find the file: " + inputFile + ". If the file path is not specified, then it is relevant to the current working directory.")
			}
		} else {
			log.Warn("An input file was not specified with --file=(filename). You will be prompted for an account to test.")
			mode = "input"
		}
		dir, err := os.Getwd()

		if fileExists(dir + "/zmxp.ini") {

		} else {
			log.Warn("Configuration file is not present. Creating a default config...")
			createDefaultIniFile()
			log.Warn("You will be prompted for details.")
			query := "Would you like to save your responses to the configuration file? (Passwords are not saved)."
			saveToConfigAnswer, err := ui.Ask(query, &input.Options{
				Default:   "No",
				Required:  true,
				Loop:      true,
				HideOrder: true,
			})
			if err != nil {
				log.Fatal(err.Error())
			}

			saveToConfigAnswer = strings.ToLower(saveToConfigAnswer)
			switch saveToConfigAnswer {
			case "y":
				saveToConfig = true
				log.Info("Will save your answers to zmxp.ini")
				break
			case "yes":
				saveToConfig = true
				log.Info("Will save your answers to zmxp.ini")
				break
			case "n":
				saveToConfig = false
				log.Info("Will not save your answers to zmxp.ini")
				break
			case "no":
				saveToConfig = false
				log.Info("Will nosave your answers to zmxp.ini")
				break
			}
		}

		cfg, err := ini.Load("zmxp.ini")
		if err != nil {
			fmt.Printf("Fail to read configuration file: %v", err)
			os.Exit(1)
		}

		adminUsername := cfg.Section("ZCS Admin Settings").Key("AdminUsername").String()
		if adminUsername == "" {
			fmt.Println("The admin username in the config file is blank. ")
			query := "Please enter your Admin Username"
			name, err := ui.Ask(query, &input.Options{
				Default:   "",
				Required:  true,
				Loop:      true,
				HideOrder: true,
			})
			if err != nil {
				log.Fatal(err.Error())
			}
			adminUsername = name
			if saveToConfig {
				writeAnswerToFile(cfg, "zmxp.ini", "ZCS Admin Settings", "AdminUsername", adminUsername)
			}
		}

		ConnectionSettings.adminUsername = adminUsername
		if password == "" {
			query := "Enter the password for " + adminUsername
			pw, err := ui.Ask(query, &input.Options{
				Default:   "",
				Required:  true,
				Loop:      true,
				HideOrder: true,
				Mask:      true,
			})
			if err != nil {
				log.Fatal(err.Error())
			}
			ConnectionSettings.adminPassword = pw
		} else {
			log.Warn("Password has been set via the command line... (This is NOT recommended.)")
			ConnectionSettings.adminPassword = password
		}

		AuthMailboxServer := cfg.Section("ZCS Admin Settings").Key("AuthMailboxServer").String()
		if AuthMailboxServer == "" {
			if !saveToConfig {
				fmt.Println("The mailbox server used for authentication in the config file is blank.")
			}
			query := "Please enter the IP/Hostname of the Zimbra Mailbox server which will be used for auth"
			name, err := ui.Ask(query, &input.Options{
				Default:   "",
				Required:  true,
				Loop:      true,
				HideOrder: true,
			})
			if err != nil {
				log.Fatal(err.Error())
			}
			AuthMailboxServer = name
			if saveToConfig {
				writeAnswerToFile(cfg, "zmxp.ini", "ZCS Admin Settings", "AuthMailboxServer", AuthMailboxServer)
			}

		}
		ConnectionSettings.mailboxServers = AuthMailboxServer

		AuthProtocol := cfg.Section("ZCS Admin Settings").Key("AuthProtocol").String()
		if AuthProtocol == "" {
			log.Info("The AuthProtocol used for authentication in the config file is blank. Using https")
			AuthProtocol = "https"
			if saveToConfig {
				writeAnswerToFile(cfg, "zmxp.ini", "ZCS Admin Settings", "AuthProtocol", AuthProtocol)
			}
		} else {
			log.Info("Using AuthProtocol: " + AuthProtocol + " from the config.")
		}
		ConnectionSettings.adminProtocol = AuthProtocol

		AuthAdminPort := cfg.Section("ZCS Admin Settings").Key("AuthAdminPort").String()
		if AuthAdminPort == "" {
			log.Info("The AuthAdminPort used for authentication in the config file is blank. Using https")
			AuthAdminPort = "7071"
			if saveToConfig {
				writeAnswerToFile(cfg, "zmxp.ini", "ZCS Admin Settings", "AuthAdminPort", AuthAdminPort)
			}
		} else {
			log.Info("Using AuthAdminPort: " + AuthAdminPort + " from the config.")
		}
		ConnectionSettings.adminPort = AuthAdminPort

		UseSocks5Proxy := cfg.Section("Connection").Key("UseSocksProxy").String()
		if UseSocks5Proxy == "" {
			log.Info("The UseSocks5Proxy used for connection in the config file is blank. Setting to false.")
			UseSocks5Proxy = "false"
			if saveToConfig {
				writeAnswerToFile(cfg, "zmxp.ini", "Connection", "UseSocks5Proxy", UseSocks5Proxy)
			}
		} else {
			log.Info("Using UseSocks5Proxy: " + UseSocks5Proxy + " from the config.")
		}
		var socks5 bool
		switch UseSocks5Proxy {
		case "true":
			socks5 = true
			break
		case "false":
			socks5 = false
			break
		}

		ConnectionSettings.useSocks5Proxy = socks5

		if socks5 == true {
			Socks5URL := cfg.Section("Connection").Key("Socks5URL").String()
			if Socks5URL == "" {
				log.Info("The Socks5URL used for authentication in the config file is blank.")
				query := "Please enter the socks5 URL (ip/hostname:port)"
				socks5urlString, err := ui.Ask(query, &input.Options{
					Default:   "",
					Required:  true,
					Loop:      true,
					HideOrder: true,
				})
				if err != nil {
					log.Fatal(err.Error())
				}
				Socks5URL = socks5urlString
				if saveToConfig {
					writeAnswerToFile(cfg, "zmxp.ini", "Connection", "Socks5URL", Socks5URL)
				}
			} else {
				log.Info("Using Socks5URL: " + Socks5URL + " from the config.")
			}
			ConnectionSettings.socksServerString = Socks5URL
		}

		SkipSSLChecks := cfg.Section("Connection").Key("SkipSSLChecks").String()
		if SkipSSLChecks == "" {
			log.Info("The SkipSSLChecks used for authentication in the config file is blank. Setting to false")
			SkipSSLChecks = "false"
			if saveToConfig {
				writeAnswerToFile(cfg, "zmxp.ini", "Connection", "SkipSSLChecks", SkipSSLChecks)
			}
		} else {
			log.Info("Using SkipSSLChecks: " + SkipSSLChecks + " from the config.")
		}

		var SkipSSL bool
		switch SkipSSLChecks {
		case "true":
			SkipSSL = true
			break
		case "false":
			SkipSSL = false
			break
		}
		ConnectionSettings.skipSSLChecks = SkipSSL
		if ConnectionSettings.skipSSLChecks == true {
			if ConnectionSettings.printedSSLSkipNotice == false {
				log.Warn("Skipping SSL Verification. (this will only be printed once).")
				ConnectionSettings.printedSSLSkipNotice = true
			}
		}

		var email string
		var accountsToTest []string
		var accountsToSkip []string
		var domains []string
		switch mode {
		case "file":
			serverTable := tablewriter.NewWriter(os.Stdout)
			table := tablewriter.NewWriter(os.Stdout)
			log.Info("Parsing input file....")
			re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

			file, err := os.Open(inputFile)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				if re.MatchString(scanner.Text()) {
					text := strings.ToLower(strings.Replace(scanner.Text(), " ", "", -1))
					accountsToTest = append(accountsToTest, text)
					domain := strings.Split(text, "@")[1]
					if !contains(domains, domain) {
						domains = append(domains, domain)
					}

				} else {
					if scanner.Text() == "" {
						continue
					}
					accountsToSkip = append(accountsToSkip, scanner.Text())
				}
			}
			if len(accountsToTest) == 0 {
				log.Fatal("The file did not contain any valid accounts to test. Cannot continue.")
			} else {
				valid := strconv.Itoa(len(accountsToTest))
				invalid := strconv.Itoa(len(accountsToSkip))
				domains := strconv.Itoa(len(domains))

				skipPlurality := "address"
				linePlurality := "line"
				if len(accountsToSkip) > 1 {
					skipPlurality = "addresses"
					linePlurality = "lines"
				}

				domainPlurality := "domain"
				if len(domains) > 1 {
					domainPlurality = "domains"
				}

				accountsPlurality := "account"
				if len(valid) > 1 {
					accountsPlurality = "accounts"
				}

				summary := "There are " + valid + " " + accountsPlurality + " to test over " + domains + " " + domainPlurality + ". "
				summary += "Skipping " + invalid + " " + linePlurality + " (invalid email " + skipPlurality + "). Use --debug=true to see details."
				//bar := uiprogress.AddBar(len(accountsToTest))
				log.Info(summary)
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
			if task == "screenshot" {
				for _, v := range accountsToTest {
					takeScreenshot(ConnectionSettings, v)
				}
			} else if task == "mail-audit" {
				limit := concur.NewConcurrencyLimiter(maxThreads)
				for i := 0; i < len(accountsToTest); i++ {

					thisAccountToTest := accountsToTest[i]
					limit.Execute(func() {
						auditAccount(ConnectionSettings, thisAccountToTest, auditType, auditValue)
					})
				}
				limit.Wait()

				table.SetHeader([]string{"Email", "Server", "Error Message"})
				renewToken := true

				log.Error("The following accounts were identified as having issues. Please investigate...")
				keys := problemAccounts.Keys()

				for _, v := range keys {
					var errorMessage string

					if tmp, ok := problemAccounts.Get(v); ok {
						value := tmp.(string)
						if strings.Contains(value, "zmerror") {
							errorMessage = strings.Split(value, "||")[0]
						}
					}

					v = strings.Split(v, "|||")[0]
					var server interface{}
					if userServerMapping.Has(v) {
						server, _ = userServerMapping.Get(v)
					} else {
						if renewToken == true {
							ConnectionSettings.adminAuthToken = getLoginToken("admin", ConnectionSettings)
							renewToken = false
						}
						infoRequest := GetInfoRequest(ConnectionSettings, v, "host")
						if strings.Contains(infoRequest, "ZMERROR") {
							errorMessage = strings.Split(infoRequest, "||ZMERROR||")[1]
							errorMessage = strings.Split(errorMessage, ": ")[0]
							server = errorMessage
						} else {
							server = infoRequest
						}

					}
					var svr string
					if !strings.Contains(server.(string), "http") {
						svr = "----"
					} else {
						svr = server.(string)
						svr = strings.Split(svr, ":")[1]
						svr = strings.Replace(svr, "//", "", -1)
						if tmp, ok := serverTracker.Get(svr); ok {
							existingCount := tmp.(int)
							serverTracker.Set(svr, existingCount+1)
						} else {
							serverTracker.Set(svr, 1)
						}

					}

					if strings.Contains(errorMessage, "invalid metadata: ") {
						errorMessage = "invalid metadata"
					}
					table.Append([]string{v, svr, errorMessage})
				}

				serverTable.SetHeader([]string{"Server", "Issue", "ZCS Version", "Installed", "Load Avg", "Uptime", "RAM Usage"})
				keys = serverTracker.Keys()
				allCommands := "w; echo \"ZMDELIM\";rpm -qa --last;echo \"ZMDELIM\";free -m; echo \"ZMDELIM\";dmesg"
				limit = concur.NewConcurrencyLimiter(maxThreads)
				for i := 0; i < len(keys); i++ {
					thisKey := serverTracker.Keys()[i]
					log.Info("Connecting to: " + thisKey + " (" + strconv.Itoa(i+1) + " of " + strconv.Itoa(len(keys)) + ")")
					value, err := serverTracker.Get(thisKey)
					if err != false {
						limit.Execute(func() {
							GetServerIntel(thisKey, "root", true, allCommands, serverTable, value.(int), i+1)
						})
					}

				}
				limit.Wait()
			}
			serverTable.Render()
			table.Render()

			break
		case "input":
			query := "Please enter an Email address to check.."
			thisUsername, err := ui.Ask(query, &input.Options{
				Default:   "",
				Required:  true,
				Loop:      true,
				HideOrder: true,
			})
			if err != nil {
				log.Fatal(err.Error())
			}
			email = thisUsername
			if task == "screenshot" {
				takeScreenshot(ConnectionSettings, email)
			} else if task == "mail-audit" {
				//runWorker(count, totalTasks/concurrency,ConnectionSettings, email)
				//
				//for i := 0; i < concurrency; i++ {
				//	go func(count int) {
				//		wg.Done()
				//	}(i)
				//}
				//	//go auditAccount(ConnectionSettings, email, auditType, auditValue)
				//	//time.Sleep(1 * time.Second)

			}
			break
		}
	}
}

//GetServerIntel(server, "root", true, allCommands,serverTable)
func GetServerIntel(server string, username string, useSocks5 bool, commands string, serverTable *tablewriter.Table, value int, slot int) {
	debug.SetGCPercent(-1)
	serverData := ExecuteSSHCommand(server, username, useSocks5, commands)
	serverDataArray := strings.Split(serverData, "ZMDELIM")

	ut := serverDataArray[0]
	var version string
	var installedDate string
	rpms := strings.Split(serverDataArray[1], "\n")
	for _, v := range rpms {
		if strings.Contains(v, "zimbra-core") {
			version = strings.Split(v, "-")[2]
			installedDate = strings.Split(v, "x86_64 ")[1]
			installedDate = strings.Replace(installedDate, " MDT", "", -1)
			installedDate = strings.TrimRight(installedDate, " ")
		}
	}

	loadAverage := strings.Split(ut, "load average: ")[1]
	loadAverage = strings.Split(loadAverage, "\n")[0]
	uptime := strings.Split(ut, "load average: ")[0]
	uptime = strings.Split(uptime, ",")[0]
	uptime = strings.Split(uptime, " up ")[1]

	memoryArray := strings.Split(serverDataArray[2], "\n")
	RAMRow := memoryArray[2]
	space := regexp.MustCompile(`\s+`)
	RAMRow = space.ReplaceAllString(RAMRow, " ")
	RAMArray := strings.Split(RAMRow, " ")
	RAMInstalled, err := strconv.Atoi(RAMArray[1])
	printRAM := true
	if err != nil {
		printRAM = false
	}
	RAMUsed, err := strconv.Atoi(RAMArray[2])
	if err != nil {
		printRAM = false
	}
	if printRAM {
		var fgColor uint8
		var bgColor uint8
		RAMPercentUsed := RAMUsed * 100
		RAMPercentUsed = RAMPercentUsed / RAMInstalled

		if RAMPercentUsed > 50 {
			bgColor = 0
			fgColor = 220
			if RAMPercentUsed > 80 {
				bgColor = 0
				fgColor = 9
				if RAMPercentUsed > 90 {
					bgColor = 9
					fgColor = 15
				}
			}
			s := color.S256(fgColor, bgColor)
			RAMString := strconv.Itoa(RAMPercentUsed) + "%"
			serverTable.Append([]string{server, strconv.Itoa(value), version, installedDate, loadAverage, uptime, s.Sprint(RAMString)})
		}else{
			RAMString := strconv.Itoa(RAMPercentUsed) + "%"
			serverTable.Append([]string{server, strconv.Itoa(value), version, installedDate, loadAverage, uptime, RAMString})

		}
	} else {
		serverTable.Append([]string{server, strconv.Itoa(value), version, installedDate, loadAverage, uptime, "N/A"})
	}
	log.Info("Server: " + server + " (Thread: " + strconv.Itoa(slot) + ") has completed.")
}

func auditAccount(ConnectionSettings ConnectionServerConfig, email string, auditType string, auditValue string) {
	for _, v := range alreadySeen {
		if v == email {
			log.Warn("Skipping already seen user: " + email)
			return
		}
	}
	alreadySeen = append(alreadySeen, email)
	f, err := os.OpenFile("zmx-mail-audit-"+logTimestamp+".log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if !auditHeaderWritten {
		if _, err := f.WriteString("Email,Host,Earliest Message Seen,Total Messages Seen,Messages Received after outage,Messages Before Outage\n"); err != nil {
			log.Println(err)
		}
		auditHeaderWritten = true
	}
	log.Info("[" + email + "-audit] Getting Messages.")
	initialOffset := 0
	shouldContinue := false
	if ConnectionSettings.adminAuthToken == "NONE" {
		ConnectionSettings.adminAuthToken = getLoginToken("admin", ConnectionSettings)
	}

	response := SearchDateRequest(ConnectionSettings, email, auditValue, strconv.Itoa(initialOffset))
	infoRequest := GetInfoRequest(ConnectionSettings, email, "host")
	host := strings.Split(infoRequest, ":")[1]
	host = strings.Replace(host, "//", "", -1)
	if !strings.Contains(response, "|ZM|") {
		stringCheck := strings.ToLower(response)

		if problemAccounts.Has(email) {

		} else {
			problemAccounts.Set(email, stringCheck)
		}

		if strings.Contains(stringCheck, "no_such") {
			go LogError(email, host, "no-such-"+logTimestamp+".log")
		} else {
			response = strings.Split(response, ":")[0]
			if _, err := f.WriteString(email + "," + host + ",ERROR " + response + "\n"); err != nil {
				log.Println(err)
			}

		}

		return
	}

	data := strings.Split(response, "|ZM|")
	messageCounter := len(data) - 1
	beforeOutage := 0
	afterOutage := 0
	outageDate, _ := time.Parse(time.RFC822, "12 Aug 19 00:00 UTC")
	var earliestSeen time.Time
	earliestSeenSet := false
	dateMap := make(map[string]int)
	log.Info("[" + email + "] Parsing message dates.")

	for _, v := range data {

		dateTime, err := msToTime(v)

		if err != nil {

		} else {
			if !earliestSeenSet {
				earliestSeen = dateTime
				earliestSeenSet = true
			}
			if dateTime.Before(earliestSeen) {
				earliestSeen = dateTime
			}
			before := dateTime.Before(outageDate)
			if before {
				beforeOutage += 1
			}
			After := dateTime.After(outageDate)
			if After {
				afterOutage += 1
			}
			date := strings.Split(dateTime.String(), " ")[0]
			if _, ok := dateMap[date]; ok {
				dateMap[date] += 1
			} else {
				dateMap[date] = 1
			}

		}

	}
	if data[0] == "MORE" {
		shouldContinue = true
	}
	if scanLimit <= initialOffset {
		for {
			if shouldContinue == false {
				break
			} else {
				initialOffset += 1000
				currenStatusString := 1000 + initialOffset
				log.Info("[" + email + "] Getting more mail (" + strconv.Itoa(currenStatusString) + " messages so far...)")
				response := SearchDateRequest(ConnectionSettings, email, auditValue, strconv.Itoa(initialOffset))

				data := strings.Split(response, "|ZM|")

				messageCounter += len(data) - 1
				if data[0] == "MORE" {
					shouldContinue = true
				} else {
					shouldContinue = false
				}
				for _, v := range data {
					dateTime, err := msToTime(v)
					if err != nil {

					} else {
						if !earliestSeenSet {
							earliestSeen = dateTime
							earliestSeenSet = true
						}
						if dateTime.Before(earliestSeen) {
							earliestSeen = dateTime
						}
						before := dateTime.Before(outageDate)
						if before {
							beforeOutage += 1
						}
						After := dateTime.After(outageDate)
						if After {
							afterOutage += 1
						}
						date := strings.Split(dateTime.String(), " ")[0]
						if _, ok := dateMap[date]; ok {
							dateMap[date] += 1
						} else {
							dateMap[date] = 1
						}

					}

				}
			}
		}
	}

	totalCounter := 0
	for _, v := range dateMap {
		totalCounter += v
	}

	if _, err := f.WriteString(email + "," + host + "," + earliestSeen.String() + "," + strconv.Itoa(totalCounter) + "," + strconv.Itoa(afterOutage) + "," + strconv.Itoa(beforeOutage) + "\n"); err != nil {
		log.Println(err)
	}
	log.Info("[" + email + "] Taking screenshot of inbox...")
	go takeScreenshot(ConnectionSettings, email)
}

func LogError(s string, s2 string, s3 string) {
	f, err := os.OpenFile(s3,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	log.Info("[" + s + "-error-logging] Recording error...")
	if _, err := f.WriteString(s + ",ERROR: " + s2 + "\n"); err != nil {
		log.Println(err)
	}
}

func msToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(0, msInt*int64(time.Millisecond)), nil
}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}
