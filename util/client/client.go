package client

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"net"
	"os"
	"time"
)

// Config is a struct which content different parameters for ssh client.
type Config struct {
	Hosts     []string
	Port      string
	Command   string
	Parallel  int
	Smooth    int
	Timeout   time.Duration
	Aggregate bool
}

type std struct {
	host   string
	stdout string
	stderr string
	err    string
}

func getSSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

func exec(host string, port string, timeout time.Duration, command string) std {
	sshConfig := &ssh.ClientConfig{
		User: os.Getenv("LOGNAME"),
		Auth: []ssh.AuthMethod{
			getSSHAgent(),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
	}
	client, err := ssh.Dial("tcp", host+":"+port, sshConfig)
	if err != nil {
		fmt.Printf("WARNING: %v\n", err)
		return std{host: host, stdout: "", stderr: "", err: "ERROR Dial()"}
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("WARNING: %v\n", err)
		return std{host: host, stdout: "", stderr: "", err: "ERROR NewSession()"}
	}
	defer session.Close()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr
	if err := session.Run(command); err != nil {
		fmt.Printf("WARNING: %v\n", err)
		return std{host: host, stdout: "", stderr: stderr.String(), err: "ERROR Run()"}
	}
	return std{host: host, stdout: stdout.String(), stderr: stderr.String(), err: "OK"}
}

func (c *Config) worker(id int, jobs <-chan int, results chan<- std) {
	for j := range jobs {
		results <- exec(c.Hosts[j], c.Port, c.Timeout, c.Command)
	}
}

// New inits ssh.
func New(c *Config) *Config {
	return c
}

// Run runs a parallel command on hosts.
func (c *Config) Run() {
	jobs := make(chan int, len(c.Hosts))
	results := make(chan std, len(c.Hosts))
	for w := 1; w <= c.Parallel; w++ {
		go c.worker(w, jobs, results)
	}
	for j := 0; j < len(c.Hosts); j++ {
		jobs <- j
	}
	close(jobs)
	uniqErr := make(map[string]int, 0)
	for i := 0; i < len(c.Hosts); i++ {
		r := <-results
		uniqErr[r.err]++
		fmt.Printf("%v\n", r)
	}
	fmt.Printf("%v\n", uniqErr)
}
