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
		var name string
		var tag string
		if len(args) > 1 {
			tag = args[1]

			if len(args) > 2 {
				name = args[2]
			}
		}
		err = installPkg(args[0], tag, name)
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
		err = removePkg(resolvePkgName(args[0]))
	case "l":
		if needsHelp(args) {
			showHelp(listHelp)
			return
		}
		if len(args) < 1 {
			err = listAll()
		} else if len(args) >= 2 && args[1] == "--json" {
			err = showPkgInfoJSON(resolvePkgName(args[0]))
		} else {
			err = showPkgInfo(resolvePkgName(args[0]))
		}
	case "fetch":
		if needsHelp(args) || len(args) < 2 {
			showHelp(fetchHelp)
			return
		}
		err = fetchTag(resolvePkgName(args[0]), args[1])
	case "use":
		if needsHelp(args) || len(args) < 1 {
			showHelp(useHelp)
			return
		}
		var tag string
		if len(args) >= 2 {
			tag = args[1]
		}
		err = useTag(resolvePkgName(args[0]), tag)
	case "cache":
		if needsHelp(args) || len(args) < 1 {
			showHelp(cacheHelp)
			return
		}
		switch args[0] {
		case "clear", "c":
			err = clearCache()
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
