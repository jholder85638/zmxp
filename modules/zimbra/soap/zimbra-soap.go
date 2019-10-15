package soap

import (
	"../../common"
	"crypto/tls"
	b64 "encoding/base64"
	_ "github.com/go-xmlfmt/xmlfmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)
var Log logrus.Logger

func SearchDateRequest(config common.ConnectionServerConfig, accountEmail string, dateBefore string, offset string)string{
	//10/10/19
	//<session id="158423761"/>
	//<sessionId id="158423761"/>
	body :=`<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
  <soap:Header>
    <context xmlns="urn:zimbra">
    <authToken>` + config.AdminAuthToken + `</authToken>
    <account by="name">` + accountEmail + `</account>
	   <userAgent name="zmxp"  version="`+common.Version+`"/>
    </context>
  </soap:Header>
  <soap:Body>
    <SearchRequest types="conversation" limit="1000" offset="`+offset+`" sortBy="dateDesc" xmlns="urn:zimbraMail">
      <query>`+dateBefore+`</query>
    </SearchRequest>
  </soap:Body>
</soap:Envelope>`
	return SendSoapRequest(config, "SearchRequest", body)
}

func DelegateAuthRequest(config common.ConnectionServerConfig, accountEmail string) string {

	body := `<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
  <soap:Header>
    <context xmlns="urn:zimbra">
      <authToken>` + config.AdminAuthToken + `</authToken>
      <nosession/>
				  <userAgent name="zmxp"  version="`+common.Version+`"/>
    </context>
  </soap:Header>
  <soap:Body>
    <DelegateAuthRequest duration="86400" xmlns="urn:zimbraAdmin">
      <account by="name">` + accountEmail + `</account>
    </DelegateAuthRequest>
  </soap:Body>
</soap:Envelope>`
	return SendSoapRequest(config, "DelegateAuthRequest", body)
}

func GetInfoRequest(config common.ConnectionServerConfig, accountEmail string, typeRequest string) string {

	body := `<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
	  <soap:Header>
		<context xmlns="urn:zimbra">
      <authToken>` + config.AdminAuthToken + `</authToken>
		  <session/>
      <account by="name">` + accountEmail + `</account>
				  <userAgent name="zmxp"  version="`+common.Version+`"/>
		</context>
	  </soap:Header>
	  <soap:Body>
		<GetInfoRequest rights="" sections="mbox,prefs,attrs,props,idents,sigs,dsrcs,children" xmlns="urn:zimbraAccount"/>
	  </soap:Body>
	</soap:Envelope>`
	return SendSoapRequest(config, typeRequest, body)
}

func ModifyZimbraMailhost(config common.ConnectionServerConfig, accountEmail string, key string, value string) string {
	body := `<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
	  <soap:Header>
		<context xmlns="urn:zimbra">
      <authToken>` + config.AdminAuthToken + `</authToken>
		  <session/>
      <account by="name">` + accountEmail + `</account>
				  <userAgent name="zmxp"  version="`+common.Version+`"/>
		</context>
	  </soap:Header>
	  <soap:Body>
    <ModifyAccountRequest xmlns="urn:zimbraAdmin">
      <a n="zimbraMailHost">mtlp000035.email.telus.net</a>
    </ModifyAccountRequest>
	  </soap:Body>
	</soap:Envelope>`

	Log.Fatal("TODO: this function is incomplete...")

	return SendSoapRequest(config, "AuthRequest", body)
}

func GetLoginToken(sessionType string, config common.ConnectionServerConfig) string {
	var urn string
	switch sessionType {
	case "admin":
		urn = "zimbraAdmin"
		break
	case "user":
		urn = "zimbraMail"
		break
	}

	body := `<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
			  <soap:Header>
				<context xmlns="urn:zimbra">
				  <nosession/>
				  <userAgent name="zmxp"  version="`+common.Version+`"/>
				</context>
			  </soap:Header>
			  <soap:Body>
				<AuthRequest xmlns="urn:` + urn + `">
				  <name>` + config.AdminUsername + `</name>
				  <password>` + config.AdminPassword + `</password>
				</AuthRequest>
			  </soap:Body>
			</soap:Envelope>`

	return SendSoapRequest(config, "AuthRequest", body)
}

func SendSoapRequest(config common.ConnectionServerConfig, requestType string, body string) string {
	var response string
	bodyRequest := strings.NewReader(body)
	httpEndPoint := config.AdminProtocol + "://" + config.MailboxServers + ":" + config.AdminPort + "/service/admin/soap/"
	req, err := http.NewRequest("POST", httpEndPoint, bodyRequest)

	httpTransport := &http.Transport{}
	var httpClient *http.Client
	if config.UseSocks5Proxy {
		connected := false
		dialer, err := proxy.SOCKS5("tcp", config.SocksServerString, nil, proxy.Direct)
		if err != nil {
			for {
				if connected==true{
					break
				}
				dialer, err = proxy.SOCKS5("tcp", config.SocksServerString, nil, proxy.Direct)
				if err !=nil{
					Log.Error("Connection issue, throttling to ensure stability...("+err.Error()+")")
					time.Sleep(2*time.Second)
				}else{
					connected = true
				}
			}
		}
		httpTransport := &http.Transport{}
		httpClient = &http.Client{Transport: httpTransport}
		httpTransport.Dial = dialer.Dial
	} else {
		httpClient = &http.Client{Transport: httpTransport}
	}

	if config.SkipSSLChecks == true {

		httpClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		for {
			resp, err = httpClient.Do(req)
			if err !=nil{
				httpTransport.CloseIdleConnections()
				httpClient.CloseIdleConnections()
				Log.Error("[SOAP-connection] Connection Issue, throttling...")
				time.Sleep(2*time.Second)
			}else{
				break
			}
		}
	}
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			Log.Fatal(err)
		}
		response = string(bodyBytes)
	} else {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			Log.Fatal(err)
		}
		response = string(bodyBytes)
	}

	response, dataType := ParseResponse(requestType, response, config)
	if dataType != 0 {
		//if dataTyppe
	}
	_ = resp.Body.Close()
	httpTransport.CloseIdleConnections()
	httpClient.CloseIdleConnections()
	return response
}

func ParseResponse(requestType string, response string, config common.ConnectionServerConfig) (string, int) {
	if strings.Contains(response, "soap:Fault") {
		//x := xmlfmt.FormatXML(response, "\t", "  ")
		//Log.Error("There was a problem with the " + requestType + " request. Use --debug=true to see there request and response.")
		//x = "Debug Response:\n" + x
		//Log.Warn(x)
		errorText := strings.Split(response, "<soap:Text>")[1]
		errorText = strings.Split(errorText, "</soap:Text>")[0]
		errorCode := strings.Split(response, "<Code>")[1]
		errorCode = strings.Split(errorCode, "</Code>")[0]
		errorTrace := strings.Split(response, "<Trace>")[1]
		errorTrace = strings.Split(errorTrace, "</Trace>")[0]
		Log.Error("Error: " + errorText+"; Code: " + errorCode+"; Trace Details: " + errorTrace)
		//if strings.Contains(errorTrace, "qtp") {
		//	Log.Error("This is an Queued Thread Pool (QTP) failure, which means it's a failure using HTTP(s).")
		//} else if strings.Contains(errorTrace, "Imap") {
		//	Log.Error("This is an IMAP failure, which means it's a failure using IMAP.")
		//} else if strings.Contains(errorTrace, "Pop") {
		//	Log.Error("This is an IMAP failure, which means it's a failure using IMAP.")
		//}
		if errorCode =="account.AUTH_FAILED"{
			Log.Fatal("Cannot proceed without valid auth.")
		}
		errorString := errorText+"||ZMERROR||"+errorCode+"||ZMERROR||base64_encoded_response||ZMERROR||"+b64.URLEncoding.EncodeToString([]byte(strings.Replace(response, "\n","",-1)))
		return errorString, 1
	}
	var responseString string
	switch requestType {
	case "host":
		responseString = strings.Split(response, "<rest>")[1]
		responseString = strings.Split(responseString, "</rest>")[0]
		break
	case "folderlist":
		responseString = strings.Split(response, "<rest>")[1]
		responseString = strings.Split(responseString, "</rest>")[0]
		break
	case "DelegateAuthRequest":
		tester := strings.Split(response, "<authToken>")
		if len(tester) >= 2 {
			responseString = strings.Split(response, "<authToken>")[1]
			responseString = strings.Split(responseString, "</authToken>")[0]
		} else {
			if strings.Contains(response, "account.NO_SUCH_ACCOUNT") {
				errorResponse := strings.Split(response, "<soap:Reason><soap:Text>")[1]
				errorResponse = strings.Split(errorResponse, "</soap:Text>")[0]
				Log.Error(errorResponse)
			}
			Log.Warn(response)
			return "", 0
		}

	case "AuthRequest":
		responseString = strings.Split(response, "<authToken>")[1]
		responseString = strings.Split(responseString, "</authToken>")[0]
		break
	case "SearchRequest":
		//<SearchResponse
		var returnResponse string
		var more string
		if strings.Contains(response,"<SearchResponse "){
			items := strings.Split(response, "more=\"")[1]
			item := strings.Split(items, "\"")[0]
			if item=="1"{
				more = "MORE"
			}else{
				more = "NOMORE"
			}
			messageArray := strings.Split(response, "<c sf")
			for _,v := range messageArray{
				tmp := strings.Split(v, " ")[0]
				tmp = strings.Replace(tmp,"\"","",-1)
				tmp = strings.Replace(tmp,"=","",-1)
				returnResponse +=tmp+"|ZM|"
			}
		}
		returnResponse = strings.Replace(returnResponse,"<soap:Envelope","",-1)
		return more+"|ZM|"+returnResponse, 0
	case "GetAllServersRequest":
		tmpSplit := strings.Split(response, "<server")

		for k, v := range tmpSplit {
			if k == 0 {
				continue
			}

			v = strings.Split(v, "\">")[0]
			hostname := strings.Split(v, "name=\"")[1]
			hostname = strings.Split(hostname, "\"")[0]

			serverId := strings.Split(v, "id=\"")[1]
			serverId = strings.Replace(serverId, "\">", "", -1)
			config.AllServers.Set(serverId, hostname)
		}
		return "", 1

	}
	return responseString, 1
}
