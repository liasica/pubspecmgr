# pubspecmgr

## install

### Linux / macOS

```bash
curl -fsSL https://raw.githubusercontent.com/liasica/pubspecmgr/master/install.sh | sh
```

### Windows (PowerShell)

```powershell
iwr -useb https://raw.githubusercontent.com/liasica/pubspecmgr/master/install.ps1 | iex
```

### Options

Install a specific version:

```bash
# Linux / macOS
curl -fsSL https://raw.githubusercontent.com/liasica/pubspecmgr/master/install.sh | VERSION=1.2.3 sh
```

```powershell
# Windows
& ([scriptblock]::Create((iwr -useb https://raw.githubusercontent.com/liasica/pubspecmgr/master/install.ps1))) -Version 1.2.3
```

Override the install directory with `INSTALL_DIR` (sh) or `-InstallDir` (ps1). The PowerShell script also reads `$env:PUBSPECMGR_VERSION` and `$env:PUBSPECMGR_INSTALL_DIR`.

Pre-built binaries are also available on the [Releases page](https://github.com/liasica/pubspecmgr/releases).

## usage

```bash
pubspecmgr <command> [arguments]

pubspecmgr -h | --help
```

### commands

#### config

- `config print` - print current configuration
- `config create` - create a new configuration file in the current directory

#### upgrade

- `upgrade` - upgrade dependencies in pubspec.yaml according to the configuration file
