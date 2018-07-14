package main

import (
	"os/exec"
	"os"
	"log"
)

var cmdChain = []*exec.Cmd {
	exec.Command("synonyms"),
	exec.Command("sprinkle"),
	exec.Command("coolify"),
	exec.Command("domainify"),
	exec.Command("available"),
}

func main() {
	cmdChain[0].Stdin = os.Stdin
	cmdChain[len(cmdChain)-1].Stdout = os.Stdout
	for i := 0; i < len(cmdChain)-1; i++ { // chain the commands
		thisCmd := cmdChain[i]
		nextCmd := cmdChain[i+1]
		stdout, err := thisCmd.StdoutPipe()
		if err != nil {
			log.Fatalln(err)
		}
		nextCmd.Stdin = stdout
	}
	for _, cmd := range cmdChain {
		if err := cmd.Start(); err != nil {
			log.Fatalln(err)
		}
		defer cmd.Process.Kill()
	}
	for _, cmd := range cmdChain {
		if err := cmd.Wait(); err != nil {
			log.Fatalln(err)
		}
	}
}