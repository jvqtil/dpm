package main

import (
	"fmt"

	"github.com/jvqtil/view"
)

const mainHelp = `dpm - Decent Package Manager

USAGE:
  dpm <command> [args]

COMMANDS:
  i - install    Install a package
  u - update     Update packages
  r - remove     Remove a package
  l - list       List installed packages
  use            Switch to another tag
  cache          Manage cached assets

Run 'dpm <command> -h' for more information on a command.`

const installHelp = `dpm install - install a package

USAGE:
  dpm install <source> [tag] [name]

ARGUMENTS:
  source   Where to install from:
             owner/repo              github (uses configured default host)
             github.com/owner/repo   explicit github source
             ./path, /path, ~/path   local file
             example.com/binary      direct url
  tag      Optional. Choose a specific tag to install
  name     Optional. Override the installed package name (defaults to repo/file name)`

const updateHelp = `dpm update - check for and install updates

USAGE:
  dpm u [package]

Without arguments, checks all installed packages and offers to update
any that have a newer release available.
With a package name, checks and updates just the specified one.
Note: only packages from Git hostings can be updated.`

const removeHelp = `dpm remove - remove an installed package

USAGE:
  dpm remove <package>

Removes the binary from your bin directory and drops it from the registry.`

const listHelp = `dpm list - show installed packages

USAGE:
  dpm list [package]

Without arguments, lists all installed packages with their versions.
With a package name, shows details of a package.`

const fetchHelp = `dpm fetch - download a tag

USAGE:
  dpm fetch <pkg/source> <tag>

Downloads asset from the new tag to the cache directory (see 'cache_dir' in config)
and adds it to the package tags list.`

const useHelp = `dpm use - switch to a different tag

USAGE:
  dpm use <package> [tag]

Switches to a tag already present in the package's history without
fetching release info. The tag must have been installed before.`

const cacheHelp = `dpm cache - manage downloads cache

USAGE:
  dpm cache clear

Clears cached downloads from "cache_dir" in your config (default: "~/.cache/dpm").`

const helpFooter = `
Version: dpm %s
More info: https://github.com/jvqtil/dpm`

func showHelp(text string) {
	view.Show(text + "\n" + fmt.Sprintf(helpFooter, Version))
}
