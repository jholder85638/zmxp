package main

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func delegateAuthRequest(config ConnectionServerConfig, accountEmail string) string {

	body := `<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
  <soap:Header>
    <context xmlns="urn:zimbra">
      <authToken>` + config.adminAuthToken + `</authToken>
      <nosession/>
      <userAgent name="zmmailbox" version="8.6.0_GA_1242"/>
    </context>
  </soap:Header>
  <soap:Body>
    <DelegateAuthRequest duration="86400" xmlns="urn:zimbraAdmin">
      <account by="name">` + accountEmail + `</account>
    </DelegateAuthRequest>
  </soap:Body>
</soap:Envelope>`
	return sendSoapRequest(config, "DelegateAuthRequest", body)
}

func GetInfoRequest(config ConnectionServerConfig, accountEmail string, typeRequest string) string {

	body := `<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
	  <soap:Header>
		<context xmlns="urn:zimbra">
      <authToken>` + config.adminAuthToken + `</authToken>
		  <session/>
      <account by="name">` + accountEmail + `</account>
		  <userAgent name="zclient" version="8.6.0_GA_1242"/>
		</context>
	  </soap:Header>
	  <soap:Body>
		<GetInfoRequest rights="" sections="mbox,prefs,attrs,props,idents,sigs,dsrcs,children" xmlns="urn:zimbraAccount"/>
	  </soap:Body>
	</soap:Envelope>`
	return sendSoapRequest(config, typeRequest, body)
}

func getLoginToken(sessionType string, config ConnectionServerConfig) string {
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
				  <userAgent name="zm-webdriver-testing" version="0.1"/>
				</context>
			  </soap:Header>
			  <soap:Body>
				<AuthRequest xmlns="urn:` + urn + `">
				  <name>` + config.adminUsername + `</name>
				  <password>` + config.adminPassword + `</password>
				</AuthRequest>
			  </soap:Body>
			</soap:Envelope>`

	return sendSoapRequest(config, "AuthRequest", body)
}

func sendSoapRequest(config ConnectionServerConfig, requestType string, body string) string {
	var response string
	bodyRequest := strings.NewReader(body)
	httpEndPoint := config.adminProtocol + "://" + config.mailboxServers + ":" + config.adminPort + "/service/admin/soap/"
	log.Info("Sending " + requestType + " request to: " + httpEndPoint)
	req, err := http.NewRequest("POST", httpEndPoint, bodyRequest)

	httpTransport := &http.Transport{}
	var httpClient *http.Client
	if config.useSocks5Proxy {

		dialer, err := proxy.SOCKS5("tcp", config.socksServerString, nil, proxy.Direct)
		if err != nil {
			fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
			os.Exit(1)
		}
		httpTransport := &http.Transport{}
		httpClient = &http.Client{Transport: httpTransport}
		httpTransport.Dial = dialer.Dial
	} else {
		httpClient = &http.Client{Transport: httpTransport}
	}

	if config.skipSSLChecks == true {

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
	} else {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		response = string(bodyBytes)
	}
	//log.Info(response)
	response, dataType := parseResponse(requestType, response, config)
	if dataType != 0 {

	}
	return response
}

func parseResponse(requestType string, response string, config ConnectionServerConfig) (string, int) {
	checkResponseString := strings.Replace(requestType, "Request", "Response", -1)
	var responseString string
	if strings.Contains(checkResponseString, "Response") {
		switch requestType {
		case "host":
			//<rest>
			responseString = strings.Split(response, "<rest>")[1]
			responseString = strings.Split(responseString, "</rest>")[0]
			break
		case "folderlist":
			//<rest>
			responseString = strings.Split(response, "<rest>")[1]
			responseString = strings.Split(responseString, "</rest>")[0]
			break
		case "DelegateAuthRequest":
			tester := strings.Split(response, "<authToken>")
			if len(tester) >=2{
				responseString = strings.Split(response, "<authToken>")[1]
				responseString = strings.Split(responseString, "</authToken>")[0]
			}else{
				log.Warn("Error parsing response: ")
				log.Warn(response)
				return "", 0
			}

		case "AuthRequest":
			responseString = strings.Split(response, "<authToken>")[1]
			responseString = strings.Split(responseString, "</authToken>")[0]
			break
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
				config.allServers.Set(serverId, hostname)
			}
			return "", 1

		}
	}
	return responseString, 1
}