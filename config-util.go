package main

import (
	"gopkg.in/ini.v1"
	"os"
)

func writeAnswerToFile(file *ini.File, filename string, s string, s2 string, s3 string) {
	file.Section(s).Key(s2).SetValue(s3)
	err := file.SaveTo(filename)
	if err !=nil{
		log.Fatal(err.Error())
	}else{
		log.Info("Saved "+s+" => "+s2+" => "+s3+" to config file: "+filename)
	}
}


func createDefaultIniFile() {
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
ScreenshotQuality=90

# ServerStats = (true/false) Gather server stats for success/failure
ServerStats = false
`
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
