<h1>Flower Loader CLI</h1>

![Language](https://img.shields.io/badge/Language-Go_1.20+-blue?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-CC_BY--NC--SA_4.0-yellowgreen?style=for-the-badge&logo=creativecommons)

**Flower Loader** is a Plugin Manager for `Creator of Another World`. With Flower Loader, you can easily manage and develop Plugins. This repository contains the command-line interface (CLI) tool to manage plugins. It is built with Go and is source-available under the Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License.

<h2>Get Involved</h2>

[![Join the Discord](https://img.shields.io/discord/1239786034561028136?color=5865F2&label=Join+The+Discord&logo=discord&style=for-the-badge)](https://discord.gg/kHSEXyawFY)

<h2>Supported Games</h2>

[![Steam](https://img.shields.io/badge/Steam-Creator_Of_Another_World-1b2838?style=for-the-badge&logo=steam)](https://store.steampowered.com/app/2761610/Creator_of_Another_World/)  
by [kuetaro (くえたろう)](https://store.steampowered.com/curator/44822906)

## Table of Contents

- [Table of Contents](#table-of-contents)
- [Installing the Tool](#installing-the-tool)
  - [From a GitHub Release](#from-a-github-release)
  - [From Source](#from-source)
- [Usage](#usage)
- [Development](#development)
  - [Prerequisites](#prerequisites)
  - [Clone the Repository](#clone-the-repository)
  - [Build the Project](#build-the-project)
  - [Run the Project](#run-the-project)
- [Plugin Ecosystem](#plugin-ecosystem)
- [Contributing](#contributing)
- [FAQ \& Troubleshooting](#faq--troubleshooting)

## Installing the Tool

### From a GitHub Release

Choose your method of installation:

<details>
<summary>Windows via PowerShell</summary>

```powershell
$INSTALL_DIR = "C:\Program Files\flower"

# Extract the archive
Expand-Archive `
  -DestinationPath $INSTALL_DIR `
  -Path flower_*.zip `
  -Force

# Add to PATH
$env:Path += ";$INSTALL_DIR"

# Check if it's installed
flower --version
```
</details>

<details>
<summary>Linux via bash (Ubuntu 20/22)</summary>

```bash
INSTALL_DIR="/usr/local/bin/flower"

# Extract the archive
sudo unzip flower_*.zip -d $INSTALL_DIR

# Add to PATH
echo "export PATH=\$PATH:$INSTALL_DIR" >> ~/.profile
source ~/.profile

# Check if it's installed
flower --version
```
</details>

### From Source

First, if you haven't, [install Go](https://golang.org/doc/install) (version 1.20 or newer.) Then, run the following:

```bash
go install github.com/flowerLoader/tool/cmd/flower@latest
```

## Usage

- **Install a plugin**: `flower add FlowerTeam.LimitBreaker`
- **Update all plugins**: `flower update all`
- **Remove a plugin**: `flower remove FlowerTeam.LimitBreaker`
- **List plugins**: `flower list`
- **Search plugins**: `flower search <partial-term>`

## Development

### Prerequisites

- [Git](https://git-scm.com/downloads)
- [Go 1.20+](https://golang.org/doc/install)

### Clone the Repository

```bash
git clone https://github.com/flowerLoader/tool flower
cd flower
```

### Build the Project

```bash
go get ./...
go build ./cmd/flower
```

### Run the Project

```bash
./flower --help
```

## Plugin Ecosystem

Plugins are hosted in GitHub repositories with the `#flower-plugin` tag. Developers can learn more about creating plugins by visiting our [plugin API documentation](https://github.com/flowerLoader/api) and [loader code](https://github.com/flowerLoader/core).

## Contributing

We welcome contributions! More information will be added soon!

## FAQ & Troubleshooting

- **How do I update all plugins at once?**
  - Use the command `flower update all`.
  
- **How do I report an issue?**
  - Please visit our [GitHub Issues page](https://github.com/flowerLoader/tool/issues).

- **How do I uninstall the tool?**
  - Delete the folder where the tool is installed. If you installed via `go install`, use `go clean -i github.com/flowerLoader/tool/cmd/flower`.

- More troubleshooting tips and frequently asked questions will be added soon.
