package main

import (
	"fmt"

	"github.com/jvqtil/view"
)

const mainHelp = `dpm - Decent Package Manager

USAGE:
  dpm <command> [args]

COMMANDS:
  i, install    Install a package
  u, update     Update packages
  r, remove     Remove a package
  l, list       List installed packages

Run 'dpm <command> -h' for more information on a command.`

const installHelp = `dpm install - install a package

USAGE:
  dpm i <source> [name]
  dpm install <source> [name]

ARGUMENTS:
  source   Where to install from:
             owner/repo              github (uses configured default host)
             github.com/owner/repo   explicit github source
             ./path, /path, ~/path   local file
           Append @tag to pin a specific version (github) or label a local file:
             owner/repo@v1.2.3
             ./mybinary@v1.0.0
  name     Optional. Override the installed package name (defaults to repo/file name)

EXAMPLES:
  dpm i cli/cli
  dpm i cli/cli@v2.40.0
  dpm i github.com/cli/cli gh-cli
  dpm i ./mybinary
  dpm i ~/Downloads/tool@v2`

const updateHelp = `dpm update - check for and install updates

USAGE:
  dpm u
  dpm u <package>

Without arguments, checks all installed packages and offers to update
any that have a newer release available. With a package name, checks
and updates just the specified one.
Note: local packages can't be updated. Reinstall instead`

const removeHelp = `dpm remove - remove an installed package

USAGE:
  dpm r <package>
  dpm remove <package>

Removes the binary from your bin directory and drops it from the registry.`

const listHelp = `dpm list - show installed packages

USAGE:
  dpm l
  dpm l <package>
  dpm list <package>

Without arguments, lists all installed packages with their versions.
With a package name, shows full details: source, binary path, install date, etc.`

const helpFooter = `
Version: dpm %s
More info: https://github.com/jvqtil/dpm`

func showHelp(text string) {
	view.Show(text + "\n" + fmt.Sprintf(helpFooter, Version))
}
