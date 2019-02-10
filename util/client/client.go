package client

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"net"
	"os"
)

// Config is a struct which content different parameters for ssh client.
type Config struct {
	Hosts     []string
	Port      string
	Command   string
	Parallel  int
	Smooth    int
	Timeout   int
	Aggregate bool
}

type std struct {
	host   string
	stdout string
	stderr string
}

func getSSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

func exec(host string, port string, command string, result chan std) {
	sshConfig := &ssh.ClientConfig{
		User: os.Getenv("LOGNAME"),
		Auth: []ssh.AuthMethod{
			getSSHAgent(),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", host+":"+port, sshConfig)
	if err != nil {
		fmt.Printf("WARNING: %v\n", err)
		result <- std{host: host, stdout: "", stderr: ""}
		return
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("WARNING: %v\n", err)
		result <- std{host: host, stdout: "", stderr: ""}
		return
	}
	defer session.Close()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr
	if err := session.Run(command); err != nil {
		fmt.Printf("WARNING: %v\n", err)
		result <- std{host: host, stdout: "", stderr: stderr.String()}
		return
	}
	result <- std{host: host, stdout: stdout.String(), stderr: stderr.String()}
}

// New inits ssh.
func New(c *Config) *Config {
	return c
}

// Run runs a parallel command on hosts.
func (c *Config) Run() error {
	result := make(chan std, c.Parallel)
	for _, host := range c.Hosts {
		//	fmt.Println(host)
		go exec(host, c.Port, c.Command, result)
	}
	for range c.Hosts {
		fmt.Printf("%v\n", <-result)
	}
	return nil
}
