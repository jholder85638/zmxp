package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"golang.org/x/net/proxy"
	"gopkg.in/ini.v1"
	"github.com/tcnksm/go-input"
	//"gopkg.in/ini.v1"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type ConnectionServerConfig struct {
	mailboxServers string
	ldapServers []string
	serverTypePreference string
	adminUsername string
	adminPassword string
	adminPort string
	adminProtocol string
	socksServerString string
	adminAuthToken string
	allServers cmap.ConcurrentMap
	skipSSLChecks bool
	useSocks5Proxy bool
}

var log = logrus.New()

func main() {
	var mode string
	ConnectionSettings :=  ConnectionServerConfig{}

	app := cli.NewApp()

	app.Action = func(ctx *cli.Context) error {
		dir, err := os.Getwd()
		if err !=nil{
			log.Fatal("Can't seem to locate the current directory. That seems odd. Exiting (sorry).")
		}
		if fileExists(dir+"/zmxp.ini") {
			log.Info("Reading zmxp.ini config...")
		} else {
			log.Warn("Configuration file is not present. Creating a default config...")
			createDefaultIniFile()
		}
		if fileExists(dir+"/input.csv") {
			log.Info("Will be using input.csv for accounts to check...")
			mode = "input"
		}else{
			log.Info("Cannot find input.csv. You will be prompted for an account to test...")
			mode = "custom"
		}
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := ini.Load("zmxp.ini")
	if err != nil {
		fmt.Printf("Fail to read configuration file: %v", err)
		os.Exit(1)
	}

	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	adminUsername := cfg.Section("ZCS Admin Settings").Key("AdminUsername").String()
	if adminUsername ==""{
		fmt.Println("The admin username in the config file is blank. ")
		query := "Please enter your Admin Username"
		name, err := ui.Ask(query, &input.Options{
			Default: "",
			Required: true,
			Loop:     true,
			HideOrder: true,
		})
		if err !=nil{
			log.Fatal(err.Error())
		}
		adminUsername = name
	}

	ConnectionSettings.adminUsername = adminUsername
	query := "Enter the password for "+adminUsername
	pw, err := ui.Ask(query, &input.Options{
		Default: "",
		Required: true,
		Loop:     true,
		HideOrder: true,
		Mask: true,
	})
	if err !=nil{
		log.Fatal(err.Error())
	}
	ConnectionSettings.adminPassword = pw

	AuthMailboxServer := cfg.Section("ZCS Admin Settings").Key("AuthMailboxServer").String()
	if AuthMailboxServer ==""{
		fmt.Println("The mailbox server used for authentication in the config file is blank.")
		query := "Please enter the IP/Hostname of the Zimbra Mailbox server which will be used for auth"
		name, err := ui.Ask(query, &input.Options{
			Default: "",
			Required: true,
			Loop:     true,
			HideOrder: true,
		})
		if err !=nil{
			log.Fatal(err.Error())
		}
		AuthMailboxServer = name
	}
	ConnectionSettings.mailboxServers = AuthMailboxServer

	AuthProtocol := cfg.Section("ZCS Admin Settings").Key("AuthProtocol").String()
	if AuthProtocol ==""{
		log.Info("The AuthProtocol used for authentication in the config file is blank. Using https")
		AuthProtocol = "https"
	}else{
		log.Info("Using AuthProtocol: "+AuthProtocol+" from the config.")
	}
	ConnectionSettings.adminProtocol = AuthProtocol

	AuthAdminPort := cfg.Section("ZCS Admin Settings").Key("AuthAdminPort").String()
	if AuthAdminPort ==""{
		log.Info("The AuthAdminPort used for authentication in the config file is blank. Using https")
		AuthAdminPort = "7071"
	}else{
		log.Info("Using AuthAdminPort: "+AuthAdminPort+" from the config.")
	}
	ConnectionSettings.adminPort = AuthAdminPort

	UseSocks5Proxy := cfg.Section("Connection").Key("UseSocksProxy").String()
	if UseSocks5Proxy ==""{
		log.Info("The UseSocks5Proxy used for connection in the config file is blank. Setting to false.")
		UseSocks5Proxy = "false"
	}else{
		log.Info("Using UseSocks5Proxy: "+UseSocks5Proxy+" from the config.")
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

	if socks5 ==true{
		Socks5URL := cfg.Section("Connection").Key("Socks5URL").String()
		if Socks5URL ==""{
			log.Info("The Socks5URL used for authentication in the config file is blank.")
			query := "Please enter the socks5 URL (ip/hostname:port)"
			socks5urlString, err := ui.Ask(query, &input.Options{
				Default: "",
				Required: true,
				Loop:     true,
				HideOrder: true,
			})
			if err !=nil{
				log.Fatal(err.Error())
			}
			Socks5URL = socks5urlString
		}else{
			log.Info("Using Socks5URL: "+Socks5URL+" from the config.")
		}
		ConnectionSettings.socksServerString = Socks5URL
	}

	SkipSSLChecks := cfg.Section("Connection").Key("SkipSSLChecks").String()
	if SkipSSLChecks ==""{
		log.Info("The SkipSSLChecks used for authentication in the config file is blank. Setting to false")
		SkipSSLChecks = "false"
	}else{
		log.Info("Using SkipSSLChecks: "+SkipSSLChecks+" from the config.")
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
	var email string
	switch mode{
	case "input":
		break
	case "custom":
		query := "Please enter an Email address to check.."
		thisUsername, err := ui.Ask(query, &input.Options{
			Default: "",
			Required: true,
			Loop:     true,
			HideOrder: true,
		})
		if err !=nil{
			log.Fatal(err.Error())
		}
		email = thisUsername
		break
	}
	ConnectionSettings.adminAuthToken = getLoginToken("admin", ConnectionSettings)
	delegateToken := delegateAuthRequest(ConnectionSettings, email)

	zimbraMailHost := strings.Split(getUserMailboxHost(ConnectionSettings, email), "home/")[0]

	ctx := context.Background()
	options := []chromedp.ExecAllocatorOption{
		chromedp.ProxyServer("socks5://127.0.0.1:9090"),
	}
	options = append(options, chromedp.DefaultExecAllocatorOptions[:]...)

	c, cc := chromedp.NewExecAllocator(ctx, options...)
	defer cc()
	// create context
	ctx, cancel := chromedp.NewContext(c)
	defer cancel()
	//var res string
	var buf []byte

	if err := chromedp.Run(ctx,fullScreenshot(zimbraMailHost+"/mail?adminPreAuth=1", 90, "ZM_AUTH_TOKEN",delegateToken,&buf)); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(email+".png", buf, 0644); err != nil {
		log.Fatal(err)
	}

}
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func fullScreenshot(urlstr string, quality int64, cookieName string, cookieValue string, res *[]byte) chromedp.Tasks {
	cookieDomain := strings.Split(urlstr, "/")[2]
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {

			success, err := network.SetCookie(cookieName, cookieValue).
				WithDomain(cookieDomain).
				WithHTTPOnly(true).
				Do(ctx)
			if err != nil {
				return err
			}
			if !success {
				return fmt.Errorf("could not set cookie %q to %q", cookieName, cookieValue)
			}

			return nil
		}),
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}

func delegateAuthRequest(config ConnectionServerConfig, accountEmail string ) string{

	body :=`<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
  <soap:Header>
    <context xmlns="urn:zimbra">
      <authToken>`+config.adminAuthToken+`</authToken>
      <nosession/>
      <userAgent name="zmmailbox" version="8.6.0_GA_1242"/>
    </context>
  </soap:Header>
  <soap:Body>
    <DelegateAuthRequest duration="86400" xmlns="urn:zimbraAdmin">
      <account by="name">`+accountEmail+`</account>
    </DelegateAuthRequest>
  </soap:Body>
</soap:Envelope>`
	return sendSoapRequest(config, "DelegateAuthRequest",body)
}

func getUserMailboxHost(config ConnectionServerConfig, accountEmail string ) string{

	body :=`<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
	  <soap:Header>
		<context xmlns="urn:zimbra">
      <authToken>`+config.adminAuthToken+`</authToken>
		  <session/>
      <account by="name">`+accountEmail+`</account>
		  <userAgent name="zclient" version="8.6.0_GA_1242"/>
		</context>
	  </soap:Header>
	  <soap:Body>
		<GetInfoRequest rights="" sections="mbox,prefs,attrs,props,idents,sigs,dsrcs,children" xmlns="urn:zimbraAccount"/>
	  </soap:Body>
	</soap:Envelope>`
	return sendSoapRequest(config, "GetInfoRequest",body)
}

func getLoginToken(sessionType string, config ConnectionServerConfig) string{
	var urn string
	switch sessionType{
	case "admin":
		urn = "zimbraAdmin"
		break
	case "user":
		urn = "zimbraMail"
		break
	}
	body :=`<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
			  <soap:Header>
				<context xmlns="urn:zimbra">
				  <nosession/>
				  <userAgent name="zm-webdriver-testing" version="0.1"/>
				</context>
			  </soap:Header>
			  <soap:Body>
				<AuthRequest xmlns="urn:`+urn+`">
				  <name>`+config.adminUsername+`</name>
				  <password>`+config.adminPassword+`</password>
				</AuthRequest>
			  </soap:Body>
			</soap:Envelope>`

	return sendSoapRequest(config, "AuthRequest",body)
}

func sendSoapRequest(config ConnectionServerConfig, requestType string, body string) string{
	var response string
	bodyRequest := strings.NewReader(body)
	httpEndPoint := config.adminProtocol+"://"+config.mailboxServers+":"+config.adminPort+"/service/admin/soap/"
	log.Info("Sending "+requestType+" request to: "+httpEndPoint)
	req, err := http.NewRequest("POST", httpEndPoint, bodyRequest)

	httpTransport := &http.Transport{}
	var httpClient *http.Client
	if config.useSocks5Proxy{

		dialer, err := proxy.SOCKS5("tcp", config.socksServerString, nil, proxy.Direct)
		if err != nil {
			fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
			os.Exit(1)
		}
		httpTransport := &http.Transport{}
		httpClient = &http.Client{Transport: httpTransport}
		httpTransport.Dial = dialer.Dial
	}else{
		httpClient = &http.Client{Transport: httpTransport}
	}

	if config.skipSSLChecks==true{
		log.Warn("Skipping SSL Verification.")
		httpClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		response = string(bodyBytes)
	}else{
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		response = string(bodyBytes)
	}
	//log.Info(response)
	response,dataType := parseResponse(requestType, response, config)
	if dataType !=0 {

	}
	return response
}

func parseResponse(requestType string, response string, config ConnectionServerConfig)  (string, int){
	checkResponseString := strings.Replace(requestType, "Request","Response",-1)
	var responseString string
	if strings.Contains(checkResponseString, "Response"){
		switch requestType {
		case "GetInfoRequest":
			//<rest>
			responseString = strings.Split(response, "<rest>")[1]
			responseString = strings.Split(responseString, "</rest>")[0]
			break
		case "DelegateAuthRequest":
			responseString = strings.Split(response, "<authToken>")[1]
			responseString = strings.Split(responseString, "</authToken>")[0]
		case "AuthRequest":
			responseString = strings.Split(response, "<authToken>")[1]
			responseString = strings.Split(responseString, "</authToken>")[0]
			break
		case "GetAllServersRequest":
			tmpSplit := strings.Split(response, "<server")

			for k,v := range tmpSplit{
				if k==0{
					continue
				}

				v = strings.Split(v, "\">")[0]
				hostname := strings.Split(v, "name=\"")[1]
				hostname = strings.Split(hostname, "\"")[0]

				serverId := strings.Split(v, "id=\"")[1]
				serverId = strings.Replace(serverId, "\">", "", -1)
				config.allServers.Set(serverId, hostname)
			}
			return "", 1

		}
	}
	return responseString,1
}

func createDefaultIniFile(){
defaultIni := `[ZCS Admin Settings]
# AdminUsername = (Username/Email) The ZCS Admin username to use to get a token on behalf of the user(s).
AdminUsername=

# AuthMailboxServer = (IP/Hostname) The target Mailbox server which will be used for authentication.
#                     This server is only used to get the auth token.
AuthMailboxServer=

# AuthProtocol = (http/https) The protocol to use for admin auth (usually https)
AuthProtocol=https

# AuthAdminPort = (int) The admin port which will be used to authenticate and get the token (usually 7071)
AuthAdminPort=7071

[Connection]
# UseSocksProxy = (true/false) Instructs zmxp to use a socks 5 proxy instead of a direct connection.
UseSocksProxy=true

# Socks5URL = (url) When using a socks 5 proxy, this is the host and port (127.0.0.1:9090)
Socks5URL=

# SkipSSLChecks = (true/false) When connecting, setting to true will ignore fatal SSL errors.
SkipSSLChecks=true

[General Settings]
# ScreenshotFilename  =
#       EmailAddress :: Email address of user beting tested.
#       Sequential :: If a list is provided, the filename will be the row number
#       Random :: The screenshot filename will be a random set of characters
ScreenshotFilename=EmailAddress

# ScreenshotQuality = (percentage) JPEG percentage quality
ScreenshotQuality=90`
	//dir, err := os.Getwd()
	f, err := os.OpenFile("zmxp.ini", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write([]byte(defaultIni))
	if err != nil {
		log.Fatal(err)
	}

	f.Close()
}