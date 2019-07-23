package ssh

import (
	"fmt"
	"github.com/lin1heart/spider/go/src/util"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"time"
)

var SshClient *ssh.Client

func init() {
	fmt.Println("ssh init ")
	const (
		username = "root"
		password = ""
		ip       = "39.104.226.149"
		port     = 22
		key      = "/Users/hikaruamano/.ssh/id_rsa"
	)
	ciphers := []string{}
	var err error
	SshClient, err = connect(username, password, ip, key, port, ciphers)
	util.CheckError(err)
}

func connect(user, password, host, key string, port int, cipherList []string) (*ssh.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		config       ssh.Config
		//session      *ssh.Session
		err error
	)

	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	if key == "" {
		auth = append(auth, ssh.Password(password))
	} else {
		pemBytes, err := ioutil.ReadFile(key)
		if err != nil {
			return nil, err
		}

		var signer ssh.Signer
		if password == "" {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(password))
		}
		if err != nil {
			return nil, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}

	if len(cipherList) == 0 {
		config = ssh.Config{
			Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
		}
	} else {
		config = ssh.Config{
			Ciphers: cipherList,
		}
	}

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	client, err = ssh.Dial("tcp", addr, clientConfig)
	return client, err
}
