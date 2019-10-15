package ssh

import (

	"../socks5"
	"bytes"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"runtime/debug"
)
var Log logrus.Logger
func ExecuteSSHCommand(serverTarget string, username string, useSocks5 bool, command string, sshKeyPath string) string{
	//var hostKey ssh.PublicKey
	// A public key may be used to authenticate against the remote
	// server by using an unencrypted PEM-encoded private key file.
	//
	// If you have an encrypted private key, the crypto/x509 package
	// can be used to decrypt it.
	key, err := ioutil.ReadFile(sshKeyPath)
	if err != nil {
		Log.Fatalf("unable to read private key: %v", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		Log.Fatalf("unable to parse private key: %v", err)
	}
	if err != nil {
		Log.Fatal(err)
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
		client, err = socks5.ProxiedSSHClient("127.0.0.1:9090", serverTarget+":22", config)
	}else{
		client, err = ssh.Dial("tcp", serverTarget+":22", config)
	}
	if err != nil {
		Log.Fatalf("unable to connect: %v", err)
	}
	defer client.Close()
	if err != nil {
		Log.Fatal("Failed to dial: ", err)
	}
	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		Log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()
	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		Log.Fatal("Failed to run: " + err.Error())
	}
	debug.SetGCPercent(-1)
	return b.String()
}
