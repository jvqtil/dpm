package main

import (
	"fmt"
	"os"
)

var Version string

func main() {
	if len(os.Args) < 2 {
		showHelp(mainHelp)
		return
	}
	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		showHelp(mainHelp)
		return
	}
	if os.Args[1] == "-v" || os.Args[1] == "--version" {
		fmt.Printf("dpm %s\n", Version)
		return
	}
	if err := initConfig(); err != nil {
		fmt.Println("failed to load config:", err)
		return
	}

	args := os.Args[2:]
	var err error
	switch os.Args[1] {
	case "i", "install":
		if needsHelp(args) || len(args) < 1 {
			showHelp(installHelp)
			return
		}
		var pkgName string
		if len(args) > 1 {
			pkgName = args[1]
		}
		err = install(args[0], pkgName)
	case "u":
		if needsHelp(args) {
			showHelp(updateHelp)
			return
		}
		if len(args) < 1 {
			err = updateAll()
		} else {
			err = updateTarget(resolvePkgName(args[0]))
		}
	case "r", "remove":
		if needsHelp(args) || len(args) < 1 {
			showHelp(removeHelp)
			return
		}
		err = remove(resolvePkgName(args[0]))
	case "l":
		if needsHelp(args) {
			showHelp(listHelp)
			return
		}
		if len(args) < 1 {
			err = list()
		} else {
			err = listTarget(resolvePkgName(args[0]))
		}
	case "cache":
		if needsHelp(args) || len(args) < 1 {
			showHelp(cacheHelp)
			return
		}
		switch args[0] {
		case "clear", "clean", "wipe", "c":
			clearCache()
		}
	default:
		fmt.Println("Unknown command:", os.Args[1])
		return
	}

	if err != nil {
		fmt.Println(err)
	}
}

func needsHelp(args []string) bool {
	for _, a := range args {
		if a == "-h" || a == "--help" {
			return true
		}
	}
	return false
}
