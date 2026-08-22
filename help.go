package main

import (
	"fmt"

	"github.com/jvqtil/view"
)

const mainHelp = `dpm - Decent Package Manager

USAGE:
  dpm <command> [args]

COMMANDS:
  install    Install a package
  fetch      Fetch a specific tag
  use        Switch to another tag
  update     Update packages
  show       Show package or tag info
  list       List installed packages
  remove     Remove a package or a tag
  cache      Manage cached assets
  migrate    Migrate data from older versions

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

const fetchHelp = `dpm fetch - download a tag

  USAGE:
    dpm fetch <pkg> <tag>

  Downloads asset from the new tag to the cache directory (see 'cache_dir' in config)
  and adds it to the package tags list.`

const useHelp = `dpm use - switch to a different tag

  USAGE:
    dpm use <package> [tag]

  Switches to a tag already present in the package's history without
  fetching release info. The tag must have been installed before.`

const updateHelp = `dpm update - check for and install updates

USAGE:
  dpm update [package]

Without arguments, checks all installed packages and offers to update
any that have a newer release available.
With a package name, checks and updates just the specified one.
Note: only packages from Git hostings can be updated.`

const showHelp = `dpm show - show package info

USAGE:
  dpm show <package> [tag] [--json]

Shows detailed information about an installed package or specific tag.
With --json, outputs the data in JSON format.`

const listHelp = `dpm list - list installed packages

USAGE:
  dpm list

Lists all installed packages with their versions.`

const removeHelp = `dpm remove - remove an installed package

USAGE:
  dpm remove <package> [tag]

With a tag removes only it, else removes the binary
from your bin directory and drops it from the registry.`

const cacheHelp = `dpm cache - manage downloads cache

USAGE:
  dpm cache clear                  Clear all cache
  dpm cache clear <package>        Clear cache for a package (all tags)
  dpm cache clear <package> <tag>  Clear cache for a specific tag

Cache directory is set by "cache_dir" in config (default: "~/.cache/dpm").`

const migrationHelp = `dpm migration - migrate data from older versions

USAGE:
dpm migration <version>

MIGRATIONS:
v0.0.4 Migrate dpm data after the v0.0.4 path changes

The v0.0.4 migration moves existing dpm configuration, registry,
and cache data from the old paths to their current locations.

The migration is safe to run when the new locations already exist;
existing files and directories are not overwritten.

NOTE:
The migration suite is temporary and will be removed in v1.0.0.`

const helpFooter = `
Version: dpm %s
More info: https://github.com/jvqtil/dpm`

func printHelp(text string) {
	view.Show(text + "\n" + fmt.Sprintf(helpFooter, Version))
}
