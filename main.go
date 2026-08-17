package main

import (
	"fmt"
	"os"
)

var Version string

func main() {
	if len(os.Args) < 2 {
		printHelp(mainHelp)
		return
	}
	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		printHelp(mainHelp)
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
			printHelp(installHelp)
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

	case "fetch":
		if needsHelp(args) || len(args) < 2 {
			printHelp(fetchHelp)
			return
		}
		err = fetchTag(resolvePkgName(args[0]), args[1])

	case "use":
		if needsHelp(args) || len(args) < 1 {
			printHelp(useHelp)
			return
		}
		var tag string
		if len(args) >= 2 {
			tag = args[1]
		}
		err = useTag(resolvePkgName(args[0]), tag)

	case "u", "update":
		if needsHelp(args) {
			printHelp(updateHelp)
			return
		}
		if len(args) < 1 {
			err = updateAll()
		} else {
			err = updateTarget(resolvePkgName(args[0]))
		}

	case "s", "show":
		if needsHelp(args) || len(args) < 1 {
			printHelp(showHelp)
			return
		}

		jsonFlag, cleanArgs := wantsJSON(args)
		if len(cleanArgs) == 0 {
			printHelp(showHelp)
			return
		}

		pkg := resolvePkgName(cleanArgs[0])
		if len(cleanArgs) == 1 {
			err = showPkgInfo(pkg, jsonFlag)
		} else {
			err = showTagInfo(pkg, cleanArgs[1], jsonFlag)
		}

	case "l":
		if needsHelp(args) {
			printHelp(listHelp)
			return
		}
		err = listAll()

	case "r", "remove":
		if needsHelp(args) || len(args) < 1 {
			printHelp(removeHelp)
			return
		}
		if len(args) >= 2 {
			err = removeTag(resolvePkgName(args[0]), args[1])
		} else {
			err = removePkg(resolvePkgName(args[0]))
		}

	case "cache":
		if needsHelp(args) || len(args) < 1 {
			printHelp(cacheHelp)
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

func wantsJSON(args []string) (bool, []string) {
	var clean []string
	found := false
	for _, a := range args {
		if a == "--json" {
			found = true
		} else {
			clean = append(clean, a)
		}
	}
	return found, clean
}
