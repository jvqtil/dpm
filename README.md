# dpm — Decent Package Manager

A package manager that just works. Insanely fast

**[Features](#features)** • **[Install](#install)** • **[Build](#build)** • **[Usage](#usage)** • **[Config](#config)**

## Features
- Installing packages from
  - Github releases (binaries only)
  - Local binaries
- Managing packages (Updating / Listing / Removing)
- Chaning install path, cache directory, etc. See [Config](#config) for details

> [!NOTE]
> dpm is in early stage. Please report any bugs / problems via Github Issues. Contributions are welcome!

## Install

```sh
go install github.com/jvqtil/dpm@latest
```

Or download a binary from [Releases](https://github.com/jvqtil/dpm/releases).

## Build

```sh
git clone https://github.com/jvqtil/dpm
cd dpm
go build
```

## Usage

```sh
dpm <command> [args]
```

| Command | Alias | Description |
|---|---|---|
| `install` | `i` | Install a package |
| `update` | `u` | Update packages |
| `remove` | `r` | Remove a package |
| `list` | `l` | List installed packages |

Run `dpm <command> -h` for details on any command.

### Install

```sh
dpm i cli/cli                       # github.com/cli/cli, latest release
dpm i cli/cli@v2.40.0               # pin a specific tag
dpm i github.com/cli/cli gh-cli     # explicit source + custom name
dpm i ./mybinary                    # take a local file under management
dpm i ~/Downloads/tool@v2           # local file with a version label
```

`dpm` looks at the release assets, tries to auto-match one for your OS/arch, and falls back to asking you to pick if it can't decide.

### Update

```sh
dpm u               # check every installed package for updates
dpm u gh-cli        # check and update just one
```

Local packages are skipped — reinstall them manually with `dpm i` when you have a new file.

### Remove

```sh
dpm r gh-cli
```

### List

```sh
dpm l            # everything installed
dpm l gh-cli     # source, binary path, install date, etc.
```

## Config

`~/.config/dpm/config.toml` (optional — sensible defaults if missing):

```toml
bin_dir = "/usr/local/bin"
assume_source_type = "github.com"
```

- `bin_dir` — where binaries get installed. Falls back to `sudo` automatically if the directory needs elevated permissions.
- `cache_dir` - where dpm stores downloaded assets.
- `assume_source_type` — the host used when you give a bare `owner/repo` without a domain.

## Where it stores data

- Installed packages are tracked in `~/.local/state/dpm/registry.json`.
- Downloads are cached in `cache_dir` (default: `~/.cache/dpm`), so re-installing the same version doesn't re-download. Use `dpm cache clear` to empty the directory.

## Known limitations

- No archive support yet - assets need to be a bare binary for now. This is the #1 priority
- Only GitHub releases and local files are supported for now. Very soon we'll have support for Codeberg, GitLab and Direct URL.
- No Windows support. All the paths are hardcoded for Unix systems. Windows support might be added soon
