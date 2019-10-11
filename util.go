package main

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func proxiedSSHClient(proxyAddress, sshServerAddress string, sshConfig *ssh.ClientConfig) (*ssh.Client, error) {
	dialer, err := proxy.SOCKS5("tcp", proxyAddress, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}

	conn, err := dialer.Dial("tcp", sshServerAddress)
	if err != nil {
		return nil, err
	}

	c, chans, reqs, err := ssh.NewClientConn(conn, sshServerAddress, sshConfig)
	if err != nil {
		return nil, err
	}

	return ssh.NewClient(c, chans, reqs), nil
}
func ExecuteSSHCommand(serverTarget string, username string, useSocks5 bool, command string) string{
	//var hostKey ssh.PublicKey
	// A public key may be used to authenticate against the remote
	// server by using an unencrypted PEM-encoded private key file.
	//
	// If you have an encrypted private key, the crypto/x509 package
	// can be used to decrypt it.
	key, err := ioutil.ReadFile("/telus_data/id_rsa")
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}
	if err != nil {
		log.Fatal(err)
	}
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// Connect to the remote server and perform the SSH handshake.
	var client *ssh.Client
	//var err error
	if useSocks5{
		client, err = proxiedSSHClient("127.0.0.1:9090", serverTarget+":22", config)
	}else{
		client, err = ssh.Dial("tcp", serverTarget+":22", config)
	}
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer client.Close()
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()
	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
	debug.SetGCPercent(-1)
	return b.String()
}

func fetchURL(URL string){
	resp, err := http.Get(URL)

	//if there was an error, report it and exit
	if err != nil {
		//.Fatalf() prints the error and exits the process
		log.Fatalf("error fetching URL: %v\n", err)
	}

	//make sure the response body gets closed
	defer resp.Body.Close()
	//check response status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("response status code was %d\n", resp.StatusCode)
	}

	//check response content type
	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		log.Fatalf("response content type was %s not text/html\n", ctype)
	}
}