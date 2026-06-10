package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Profile struct {
	User    string
	Project string
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		printHelp()
		return
	}

	switch args[0] {
	case "help":
		printHelp()
	case "profile":
		runProfile(args[1:])
	default:
		fmt.Printf("Unknown command: %s\n\n", args[0])
		printHelp()
		os.Exit(1)
	}
}

func runProfile(args []string) {
	if len(args) == 0 {
		fmt.Println("No profile subcommand specified")
		printHelp()
		os.Exit(1)
	}

	switch args[0] {
	case "create":
		flags := parseFlags(args[1:])
		create(flags["name"], flags["user"], flags["project"])
	case "get":
		flags := parseFlags(args[1:])
		get(flags["name"])
	case "list":
		list()
	case "delete":
		flags := parseFlags(args[1:])
		remove(flags["name"])
	default:
		fmt.Printf("Unknown profile subcommand: %s\n", args[0])
		os.Exit(1)
	}
}

func parseFlags(args []string) map[string]string {
	flags := make(map[string]string)
	for _, a := range args {
		if !strings.HasPrefix(a, "--") {
			continue
		}

		a = strings.TrimPrefix(a, "--")
		parts := strings.SplitN(a, "=", 2)
		if len(parts) == 2 {
			flags[parts[0]] = parts[1]
		} else {
			flags[parts[0]] = ""
		}
	}
	
	return flags
}

func fileName(name string) string {
	return name + ".yaml"
}

func create(name, user, project string) {
	if name == "" {
		fmt.Println("Flag --name is required")
		os.Exit(1)
	}

	path := fileName(name)
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("Profile %q already exists\n", name)
		os.Exit(1)
	}

	content := fmt.Sprintf("user: %s\nproject: %s\n", user, project)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Printf("Write error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Profile %q created\n", name)
}

func get(name string) {
	if name == "" {
		fmt.Println("Flag --name is required")
		os.Exit(1)
	}

	p, err := readProfile(fileName(name))
	if err != nil {
		fmt.Printf("Profile %q not found\n", name)
		os.Exit(1)
	}

	fmt.Printf("name: %s\nuser: %s\nproject: %s\n", name, p.User, p.Project)
}

func list() {
	matches, _ := filepath.Glob("*.yaml")
	if len(matches) == 0 {
		fmt.Println("No profiles found")
		return
	}

	names := make([]string, 0, len(matches))
	for _, m := range matches {
		names = append(names, strings.TrimSuffix(m, ".yaml"))
	}

	sort.Strings(names)
	fmt.Println("Available profiles:")
	for _, n := range names {
		fmt.Printf("  - %s\n", n)
	}
}

func remove(name string) {
	if name == "" {
		fmt.Println("Flag --name is required")
		os.Exit(1)
	}

	path := fileName(name)
	if _, err := os.Stat(path); err != nil {
		fmt.Printf("Profile %q not found\n", name)
		os.Exit(1)
	}

	if err := os.Remove(path); err != nil {
		fmt.Printf("Delete error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Profile %q deleted\n", name)
}

func readProfile(path string) (Profile, error) {
	var p Profile
	data, err := os.ReadFile(path)

	if err != nil {
		return p, err
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "user":
			p.User = value
		case "project":
			p.Project = value
		}
	}

	return p, nil
}

func printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  Create a new profile: profile create --name=<name> --user=<user> --project=<project>")
	fmt.Println("  Show profile contents: profile get --name=<name>")
	fmt.Println("  List all profiles: profile list")
	fmt.Println("  Delete a profile: profile delete --name=<name>")
	fmt.Println("  Show this help: help")
}
