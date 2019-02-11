package main

import (
	"./util/client"
	"./util/resolver"
	"flag"
	"fmt"
	"os"
)

func runAction() {
	var (
		f    = flag.String("f", "", "f is a script file path for xecute remotely.")
		c    = flag.String("c", "", "c is a command for execute remotely.")
		p    = flag.Int("p", 1, "p is a number of parallel ssh session.")
		port = flag.String("P", "22", "P is a server ssh port.")
		s    = flag.Int("s", 0, "s is a duration of run command or script - smooth ssh session on time.")
		t    = flag.Duration("t", 20000000000, "t is a time (a number) of connection timeout.")
		l    = flag.Bool("l", false, "l is a line mode without aggregate results.")
	)
	flag.CommandLine.Parse(os.Args[2:])
	command := ""
	if *f != "" {
		command = *f
	}
	if *c != "" {
		command = *c
	}
	if command == "" || flag.NArg() == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	hosts, err := resolver.GetHosts(flag.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	config := client.Config{
		Hosts:     hosts,
		Port:      *port,
		Command:   command,
		Parallel:  *p,
		Smooth:    *s,
		Timeout:   *t,
		Aggregate: !*l,
	}
	ssh := client.New(&config)
	fmt.Println("Run")
	ssh.Run()
}
