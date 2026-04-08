<h1 align="center">
<strong>Koma コマ</strong>
</h1>

<p align="center">
    <img alt="Linux" src="https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black">
    <img alt="macOS" src="https://img.shields.io/badge/mac%20os-000000?style=for-the-badge&logo=macos&logoColor=F0F0F0">
    <img alt="Windows" src="https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white">
    <img alt="Termux" src="https://img.shields.io/badge/Termux-000000?style=for-the-badge&logo=GNOME%20Terminal&logoColor=white">
</p>

<h3 align="center">
    The most advanced CLI manga downloader in the entire universe!
</h3>

<p align="center">
    <em>Named after コマ (koma) — the panel, the smallest storytelling unit of manga.</em>
</p>

<p align="center">
    <img alt="Koma TUI" src="assets/tui.gif">
</p>

> [!NOTE]
> Koma is an actively maintained fork of [metafates/mangal](https://github.com/metafates/mangal), which was archived in April 2025. This fork includes security fixes, dependency upgrades, bug fixes, and new features. See [Changes from upstream](#changes-from-upstream) below.

## Table of contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Custom scrapers](#custom-scrapers)
- [Anilist](#anilist)
- [Honorable mentions](#honorable-mentions)

## Features

- __Lua Scrapers!!!__ You can add any source you want by creating your own _(or using someone's else)_ scraper with
  __Lua 5.1__. See [mangal-scrapers repository](https://github.com/metafates/mangal-scrapers)
- __2 Built-in sources__ - [Mangadex](https://mangadex.org) & [Mangapill](https://mangapill.com). Additional sources via Lua scrapers
- __Download & Read Manga__ - I mean, it would be strange if you couldn't, right?
- __Caching__ - Koma will cache as much data as possible, so you don't have to wait for it to download the same data over and over again.
- __4 Different export formats__ - PDF, CBZ, ZIP and plain images
- __TUI ✨__ - You already know how to use it! (ﾉ>ω<)ﾉ :｡･::･ﾟ’★,｡･:･ﾟ’☆
- __Scriptable__ - You can use Koma in your scripts, it's just a CLI app after all. [Examples](https://github.com/metafates/mangal/wiki/Inline-mode)
- __History__ - Resume your reading from where you left off!
- __Fast?__ - YES.
- __Monolith__ - ZERO runtime dependencies. Even Lua is built in. Easy to install and use.
- __Cross-Platform__ - Linux, macOS, Windows, Termux, even your toaster. (¬‿¬ )
- __Anilist integration__ - Koma will collect additional data from Anilist and use it to improve your reading experience. It can also sync your progress!

## Installation

### Homebrew (macOS / Linux)

```shell
brew tap JaneWoo0068/homebrew-koma
brew install koma
```

### Pre-compiled binaries

Download the latest release for your platform from the [Releases page](https://github.com/JaneWoo0068/koma/releases/latest).

**macOS / Linux:**
```shell
# Example for macOS Apple Silicon (arm64)
tar xzf koma_*_Darwin_arm64.tar.gz
sudo mv koma /usr/local/bin/
koma
```

**Windows:**
1. Download the `.zip` for your architecture
2. Extract it
3. Open Command Prompt or PowerShell in that folder
4. Run `koma.exe`

**Debian / Ubuntu:**
```shell
sudo dpkg -i koma_*_amd64.deb
koma
```

**Fedora / RHEL:**
```shell
sudo rpm -i koma-*-1.x86_64.rpm
koma
```

### Docker

```shell
docker pull ghcr.io/janewoo0068/koma:latest
docker run --rm -ti -v $(pwd)/downloads:/downloads ghcr.io/janewoo0068/koma
```

> **Tip for macOS users:** [OrbStack](https://orbstack.dev/) is a lightweight and fast alternative to Docker Desktop. Install it, then use the same `docker` commands above.

### From source

Visit this link to install [Go](https://go.dev/doc/install).

Clone the repo
```shell
git clone --depth 1 https://github.com/JaneWoo0068/mangal.git koma
cd koma
```

GNU Make **(Recommended)**
```shell
make install # if you want to compile and install koma to path
make build # if you want to just build the binary
```

<details>
<summary>If you don't have GNU Make use this</summary>

```shell
# To build
go build -ldflags "-X 'github.com/metafates/mangal/constant.BuiltAt=$(date -u)' -X 'github.com/metafates/mangal/constant.BuiltBy=$(whoami)' -X 'github.com/metafates/mangal/constant.Revision=$(git rev-parse --short HEAD)' -s -w"

# To install
go install -ldflags "-X 'github.com/metafates/mangal/constant.BuiltAt=$(date -u)' -X 'github.com/metafates/mangal/constant.BuiltBy=$(whoami)' -X 'github.com/metafates/mangal/constant.Revision=$(git rev-parse --short HEAD)' -s -w"
```

</details>

If you want to build for other architectures, set env variables `GOOS` and `GOARCH`

```shell
GOOS=linux GOARCH=arm64 make build
```

[Available GOOS and GOARCH combinations](https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63)

<details>
<summary><strong>Legacy installation methods (original mangal, not maintained for Koma)</strong></summary>

The following installation methods were available for the original [metafates/mangal](https://github.com/metafates/mangal) project. They install the original mangal, **not Koma**, and may no longer be functional since the upstream project was archived.

- **Script (Linux, MacOS, Termux):** `curl -sSL mangal.metafates.one/install | sh`
- **Arch Linux:** [AUR package](https://aur.archlinux.org/packages/mangal-bin)
- **MacOS (Homebrew):** `brew tap metafates/mangal && brew install mangal`
- **Windows (Scoop):** `scoop bucket add extras && scoop install mangal`
- **Termux:** `pkg install mangal`
- **Gentoo:** via [raiagent](https://github.com/leycec/raiagent) overlay
- **Nix:** `nix-env -iA nixpkgs.mangal`
- **Docker:** `docker pull metafates/mangal`

</details>

## Usage

### TUI

Just run `koma` and you're ready to go.

<details>
<summary>Keybinds</summary>

| Bind                                                        | Description                          |
|-------------------------------------------------------------|--------------------------------------|
| <kbd>?</kbd>                                                | Show help                            |
| <kbd>↑/j</kbd> <kbd>↓/k</kbd> <kbd>→/l</kbd> <kbd>←/h</kbd> | Navigate                             |
| <kbd>g</kbd>                                                | Go to first                          |
| <kbd>G</kbd>                                                | Go to last                           |
| <kbd>/</kbd>                                                | Filter                               |
| <kbd>esc</kbd>                                              | Back                                 |
| <kbd>space</kbd>                                            | Select one                           |
| <kbd>tab</kbd>                                              | Select all                           |
| <kbd>v</kbd>                                                | Select volume                        |
| <kbd>backspace</kbd>                                        | Unselect all                         |
| <kbd>enter</kbd>                                            | Confirm                              |
| <kbd>o</kbd>                                                | Open URL                             |
| <kbd>r</kbd>                                                | Read                                 |
| <kbd>q</kbd>                                                | Quit                                 |
| <kbd>ctrl+c</kbd>                                           | Force quit                           |
| <kbd>a</kbd>                                                | Select Anilist manga (chapters list) |
| <kbd>d</kbd>                                                | Delete single history entry          |

</details>

![TUI](https://user-images.githubusercontent.com/62389790/198830334-fd85c74f-cf3b-4e56-9262-5d62f7f829f4.png)

> If you wonder what those icons mean - `D` stands for "downloaded", `*` shows that chapter is marked to be downloaded.
> You can choose different icons, e.g. nerd font ones - just run mangal with `--icons nerd`.
> Available options are `nerd`, `emoji`, `kaomoji` and `squares`

### Mini

Mini mode tries to mimic [ani-cli](https://github.com/pystardust/ani-cli)

To run: `koma mini`

![mini](https://user-images.githubusercontent.com/62389790/198830544-f2005ec4-c206-4fe0-bd08-862ffd08320e.png)

### Inline

Inline mode is intended for use with other scripts.

Type `koma help inline` for more information.

See [Wiki](https://github.com/metafates/mangal/wiki/Inline-mode) for more examples.

<p align="center">
    <img alt="Koma Inline" src="assets/inline.gif">
</p>

### Other

See `koma help` for more information

## Configuration

Koma uses [TOML](https://toml.io) format for configuration under the `koma.toml` filename.
Config path depends on the OS.
To find yours, use `koma where --config`.
For example, on __Linux__ it would be `~/.config/koma/koma.toml`.

Use env variable `KOMA_CONFIG_PATH` to set custom config path.
> See `koma env` to show all available env variables.

| Command              | Description                                      |
|----------------------|--------------------------------------------------|
| `koma config get`    | Get config value for specific key                |
| `koma config set`    | Set config value for specific key                |
| `koma config reset`  | Reset config value for specific key              |
| `koma config info`   | List all config fields with description for each |
| `koma config write`  | Write current config to a file                   |

## Custom scrapers

TLDR; To browse and install a custom scraper
from [mangal-scrapers repository](https://github.com/metafates/mangal-scrapers) run

    koma sources install

Koma has a Lua5.1 VM built-in + some useful libraries, such as headless chrome, http client, html parser and so on...

Check the [defined modules](https://github.com/metafates/mangal-lua-libs) for more information.

For scrapers examples, check the [mangal-scrapers repository](https://github.com/metafates/mangal-scrapers)

### Creating a custom scraper

This command will create `example.lua` file in the `koma where --sources` directory.

    koma sources gen --name example --url https://example.com

Open the file and edit it as you wish.
Take a look at the comments for more information.
See [mangal-scrapers repository](https://github.com/metafates/mangal-scrapers) for examples.

You can test it by running `koma run <filepath>`

It should automatically appear in the list of available scrapers.

> New to Lua? [Quick start guide](https://learnxinyminutes.com/docs/lua/)

## Anilist

Koma also supports integration with Anilist.

Besides fetching metadata for each manga when downloading,
Koma can also mark chapters as read on your Anilist profile when you read them inside Koma.

For more information see [wiki](https://github.com/metafates/mangal/wiki/Anilist-Integration)

## Changes from upstream

This fork includes the following changes compared to the original [metafates/mangal](https://github.com/metafates/mangal):

### Security
- Upgraded Go 1.18 to 1.23.0
- Fixed 9 CVEs in `golang.org/x/net`, `golang.org/x/image`, `google.golang.org/protobuf`

### Bug fixes
- Fixed MangaDex chapter pagination skipping pages on language filter (#172)
- Fixed MangaDex chapter index out of range crash (#150, #196)
- Fixed MangaDex search fatal exit on API error
- Fixed Mangapill scraper (CSS selector, URL encoding, Host header)
- Fixed config `set` crash on empty string value (#147)
- Fixed nil pointer panic during chapter download (#135)
- Fixed empty volume folders when downloading CBZ (#183)
- Fixed long vertical pages clipped in PDF export (#192)
- Fixed headless browser memory leak in Lua scrapers (#171)
- Fixed TUI spinner not animating during loading
- Fixed "Anilsit" typo in Anilist integration (#209)

### New features
- MangaDex chapter deduplication across scanlation groups, preferring official translations (#162)
- MangaDex API rate limiting (configurable, default 200ms/req) (#152)
- Chapter publish date from source used in ComicInfo metadata (#164)
- Download on enter when `tui.read_on_enter` is false (#156)
- `downloader.escape_whitespace` option to control filename whitespace handling (#159)

### Removed
- Manganelo and Manganato built-in scrapers (domains defunct, Cloudflare-blocked). Use Lua scrapers for alternative sources.

## Roadmap

Planned for future releases:

- [ ] Inline mode output directory override (#186)
- [ ] Add new manga sources (WeebCentral, etc.)
- [ ] Headless browser support for built-in scrapers (enables Cloudflare-protected sites)
- [ ] GitHub Releases with GoReleaser
- [ ] Homebrew tap
- [ ] Docker image
- [ ] Go module path migration from `metafates/mangal` to `JaneWoo0068/koma`
- [ ] Modernize TUI interface
- [ ] Updated Lua scraper documentation and examples

## Honorable mentions

### Projects using mangal/Koma

- [kaizoku](https://github.com/oae/kaizoku) - Self-hosted manga downloader with mangal as its core

### Similar Projects

- [mangadesk](https://github.com/darylhjd/mangadesk) - Terminal client for MangaDex
- [ani-cli](https://github.com/pystardust/ani-cli) - A cli tool to browse and play anime
- [manga-py](https://github.com/manga-py/manga-py) - Universal manga downloader
- [animdl](https://github.com/justfoolingaround/animdl) - A highly efficient, fast, powerful and light-weight anime
  downloader and streamer
- [tachiyomi](https://github.com/tachiyomiorg/tachiyomi) - Free and open source manga reader for Android

### Libraries

- [bubbletea](https://github.com/charmbracelet/bubbletea), [bubbles](https://github.com/charmbracelet/bubbles)
  & [lipgloss](https://github.com/charmbracelet/lipgloss) - TUI framework
- [gopher-lua](https://github.com/yuin/gopher-lua) - Lua VM for custom scrapers
- [cobra](https://github.com/spf13/cobra) and [viper](https://github.com/spf13/viper) - CLI & config
- [pdfcpu](https://github.com/pdfcpu/pdfcpu) - PDF processor in pure Go
- _And many others!_

### Contributors

Thanks to all [original contributors](https://github.com/metafates/mangal/graphs/contributors) of mangal!
