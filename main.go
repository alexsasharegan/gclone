package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Program struct {
	rawPath  string
	resource string
	owner    string
	name     string
}

func init() {
	flag.Parse()
}

func main() {
	program := Program{}
	args := flag.Args()
	if len(args) == 0 {
		log.Fatalln("Missing url to clone")
	}

	program.rawPath = args[0]
	if !strings.HasSuffix(program.rawPath, ".git") {
		log.Fatalln("Not a valid git url.")
	}
	program.rawPath = program.rawPath[:len(program.rawPath)-len(".git")]

	if strings.HasPrefix(program.rawPath, "git@") {
		// git@github.com:pictalk/pictalk.git
		s := program.rawPath[len("git@"):]
		parts := strings.Split(s, ":")
		if len(parts) != 2 {
			log.Fatalln("Not a valid ssh git url.")
		}

		program.resource = parts[0]
		i := strings.IndexByte(parts[1], '/')
		if i == -1 {
			log.Fatalln("Not a valid ssh git url.")
		}

		program.owner = parts[1][:i]
		program.name = parts[1][i+1:]
	} else {
		// https://github.com/pictalk/pictalk.git
		u, err := url.Parse(program.rawPath)
		if err != nil {
			log.Fatalf("Failed to parse git url: %v\n", err)
		}

		program.resource = u.Hostname()
		parts := strings.Split(u.Path, "/")
		if parts[0] == "" {
			parts = parts[1:]
		}

		program.owner = parts[0]
		program.name = parts[1]
	}

	repo := filepath.Join(os.Getenv("HOME"), program.resource, program.owner, program.name)

	cmd := exec.Command("git", "clone", program.rawPath, repo)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed clone: %v", err)
	}

	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr)

	fmt.Println(repo)
}
