# gbvm

`gbvm` is a command-line tool for managing Go binaries installed in your `GOPATH/bin` directory. It lets you list installed binaries, check their versions, back up the list to JSON, and upgrade binaries to the latest releases.

## Installation

```bash
go install github.com/TBXark/gbvm@latest
```

## Commands

### list

```bash
Usage: gbvm list [options]

List all installed Go binaries

Options:
  -help      show help (default: false)
  -json      json mode (default: false)
  -verbose   show scan errors (default: false)
  -versions  show version (default: false)
```

### install

```bash
Usage: gbvm install [options] <backup file>

Install Go binaries from backup file

Options:
  -help  show help (default: false)
```

### upgrade

```bash
Usage: gbvm upgrade [options] [bin1 bin2 ...]

Upgrade Go binaries

Options:
  -help      show help (default: false)
  -skip-dev  skip dev version (default: false)
  -verbose   show scan errors (default: false)
```

## Examples

```bash
# Install binaries from backup
gbvm install backup.json

# List all installed binaries with their versions
gbvm list -versions

# List binaries in JSON format
gbvm list -json

# Upgrade specific binaries
gbvm upgrade bin1 bin2

# Upgrade all binaries except development versions
gbvm upgrade -skip-dev
```

## License

`gbvm` is released under the MIT license. See [LICENSE](LICENSE) for details.
