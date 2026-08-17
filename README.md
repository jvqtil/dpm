# dpm — Decent Package Manager

A package manager that just works. Insanely fast

**[Features](#features)** • **[Install](#install)** • **[Build](#build)** • **[Usage](#usage)** • **[Config](#config)**

## Features
- Installing packages from
  - GitHub releases
  - Direct URLs
  - Local binaries
- Managing packages (Updating / Listing / Removing)
- Version rollback - if you want to rollback to `v1.1.0` to `v1.0.9`, you don't need to re-download binary
- Auto-matching binary from GitHub releases / unpacked archives based on OS and Architecture
- Changing install path, cache directory, etc. See [config section](#config) for details
- Automatically assumes source if not provided (`jvqtil/dpm` -> `github.com/jvqtil/dpm`). See [config section](#config) for details

> [!NOTE]
> dpm is in early stage. Please report any bugs / problems via GitHub Issues. Contributions are welcome!

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
| `fetch` | | Fetch a tag |
| `use` | | Switch to a tag |
| `update` | `u` | Update packages |
| `show` | `s` | Show package details |
| `list` | `l` | List installed packages |
| `remove` | `r` | Remove a package or a tag |
| `cache` | | Manage cache |

To see more help, use `dpm --help`

Run `dpm <command> -h` for details on any command.

## Config

`~/.config/dpm/config.toml` (optional — sensible defaults if missing):

```toml
bin_dir = "/usr/local/bin"
cache_dir = "~/.cache/dpm"
assume_source_type = "github.com"
```

- `bin_dir` — where binaries are being installed. Falls back to `sudo` automatically if the directory needs elevated permissions.
- `cache_dir` - where dpm stores downloaded assets.
- `assume_source_type` — the host used when you give a bare `owner/repo` without a domain.

## Where it stores data

- Installed packages are tracked in `~/.local/state/dpm/registry.json`.
- Downloads are cached in `cache_dir` (default: `~/.cache/dpm`), so re-installing the same version doesn't re-download. Use `dpm cache clear` to empty the directory.

## Known limitations

- Automatic release detection is only supported for GitHub. Direct URL installation works for assets hosted on GitLab, sr.ht, and Codeberg, but automatic release resolution is not yet supported for those platforms
- No Windows support. All the paths are hardcoded for Unix systems. Windows support might be added later
