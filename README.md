# ztock

[![GitHub release](https://img.shields.io/github/release/ztock/ztock.svg)](https://github.com/ztock/ztock/releases)
[![Github Build Status](https://github.com/ztock/ztock/workflows/Go/badge.svg?branch=main)](https://github.com/ztock/ztock/actions?query=workflow%3AGo+branch%3Amain)
[![GoDoc](https://godoc.org/github.com/ztock/ztock?status.svg)](https://godoc.org/github.com/ztock/ztock)

Command-line for displaying real-time stock data.

## Installation

```shell
$ go get github.com/ztock/ztock
```

## Usage

```shell
$ ztock 600000
NUMBER  CURRENT PRICE  PERCENTAGE CHANGE  OPENING PRICE  PREVIOUS CLOSING PRICE  HIGH PRICE  LOW PRICE  DATE
600660  47.200         0.75%              47.570         46.850                  47.900      46.520     16007-03-16 15:00:02
```

### CLI

Here are the options you can give ztock:

```shell
Usage:
  ztock [flags]

Flags:
      --config string       config file (default is $HOME/.ztock/config.yaml)
  -h, --help                help for ztock
  -i, --index index         set the stock market index (default sh)
      --log-format string   set the format that is used for logging (default "text")
      --log-level string    set the level that is used for logging (default "warn")
  -p, --platform platform   set the source platform for stock data (default sina)
  -v, --version             version for ztock
```

### Configuration

The command will look for the configuration file `config.yaml` in `$HOME/.ztock`, unless overridden by the `--config` option.
The following settings can be configured:

```yaml
# platform for stock data
index: sh
# stock market index
platform: sina
# log level
log_level: debug
# log format
log_format: text
```

## Issues

- [Open an issue in GitHub](https://github.com/ztock/ztock/issues)

## License

[MIT](LICENSE)
