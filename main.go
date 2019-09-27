package main

import (
	"bufio"
	"fmt"
	"github.com/integrii/flaggy"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/sirupsen/logrus"
	"github.com/tcnksm/go-input"
	"gopkg.in/ini.v1"
	"os"
	"regexp"
	"strconv"
	"strings"
)

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

var log = logrus.New()
var version = "0.1-Alpha"
var inputFile = ""
var serverStats = false

func init() {
	flaggy.SetName("Zimbra Experience (zmxp)")
	flaggy.SetDescription("A utility to gather information about the quality of the end user Zimbra Experience.")
	flaggy.DefaultParser.ShowHelpOnUnexpected = true
	flaggy.DefaultParser.AdditionalHelpPrepend = "http://github.com/zimbra/zmxp"
	flaggy.String(&inputFile, "f", "file", "Provide an input file with a list of accounts to test (--file=myfile.txt). The file should consist of user email addresses, one per line.")

	flaggy.Bool(&serverStats, "ss", "server-stats", "Gather server stats for success vs failure.")

	flaggy.SetVersion(version)
	flaggy.Parse()
}

func main() {
	var mode = "file"
	var saveToConfig bool
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}
	ConnectionSettings := ConnectionServerConfig{}
	ConnectionSettings.printedSSLSkipNotice = false
	if inputFile !=""{
		if fileExists(inputFile) {
			log.Info("Using the input file: "+inputFile)
		} else {
			log.Fatal("Cannot find the file: "+inputFile+". If the file path is not specified, then it is relevant to the current working directory.")
		}
	}else{
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
		if saveToConfig{
			writeAnswerToFile(cfg, "zmxp.ini", "ZCS Admin Settings", "AdminUsername", adminUsername)
		}
	}

	ConnectionSettings.adminUsername = adminUsername
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

	AuthMailboxServer := cfg.Section("ZCS Admin Settings").Key("AuthMailboxServer").String()
	if AuthMailboxServer == "" {
		if !saveToConfig{
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
		if saveToConfig{
			writeAnswerToFile(cfg, "zmxp.ini", "ZCS Admin Settings", "AuthMailboxServer", AuthMailboxServer)
		}


	}
	ConnectionSettings.mailboxServers = AuthMailboxServer

	AuthProtocol := cfg.Section("ZCS Admin Settings").Key("AuthProtocol").String()
	if AuthProtocol == "" {
		log.Info("The AuthProtocol used for authentication in the config file is blank. Using https")
		AuthProtocol = "https"
		if saveToConfig{
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
		if saveToConfig{
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
		if saveToConfig{
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
			if saveToConfig{
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
		if saveToConfig{
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
		log.Info("Parsing input file....")
		re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

		file, err := os.Open(inputFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if re.MatchString(scanner.Text()){
				accountsToTest = append(accountsToTest,scanner.Text())
				domain := strings.Split(scanner.Text(), "@")[1]
				if !contains(domains, domain){
					domains = append(domains, domain)
				}

			}else{
				if scanner.Text()==""{
					continue
				}
				accountsToSkip = append(accountsToSkip, scanner.Text())
			}
		}
		if len(accountsToTest) ==0{
			log.Fatal("The file did not contain any valid accounts to test. Cannot continue.")
		}else{
			valid := strconv.Itoa(len(accountsToTest))
			invalid := strconv.Itoa(len(accountsToSkip))
			domains := strconv.Itoa(len(domains))

			skipPlurality := "address"
			linePlurality := "line"
			if len(accountsToSkip) >1{
				skipPlurality = "addresses"
				linePlurality = "lines"
			}

			domainPlurality := "domain"
			if len(domains)>1{
				domainPlurality = "domains"
			}

			accountsPlurality := "account"
			if len(valid)>1{
				accountsPlurality = "accounts"
			}
			summary := "There are "+valid+" "+accountsPlurality+" to test over "+domains+" "+domainPlurality+". "
			summary += "Skipping "+invalid+" "+linePlurality+" (invalid email "+skipPlurality+"). Use --debug=true to see details."

			log.Info(summary)
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		for _,v := range accountsToTest{
			takeScreenshot(ConnectionSettings, v)
		}
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
		takeScreenshot(ConnectionSettings, email)
		break
	}
}