# yutil

[![Test Status](https://github.com/amplia-iiot/yutil/workflows/Test/badge.svg)](https://github.com/amplia-iiot/yutil/actions/workflows/test.yml)
[![Lint Status](https://github.com/amplia-iiot/yutil/workflows/Lint/badge.svg)](https://github.com/amplia-iiot/yutil/actions/workflows/lint.yml)
[![Codecov report](https://img.shields.io/codecov/c/github/amplia-iiot/yutil/main.svg)](https://codecov.io/gh/amplia-iiot/yutil)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://www.conventionalcommits.org/en/v1.0.0/)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/amplia-iiot/yutil/blob/main/LICENSE)

Common functionality for working with YAML files

## Table of contents

- [yutil](#yutil)
	- [Table of contents](#table-of-contents)
	- [Features](#features)
	- [Getting started](#getting-started)
		- [Installation](#installation)
			- [deb/rpm/apk:](#debrpmapk)
			- [Manual](#manual)
			- [Go users](#go-users)
		- [Test installation](#test-installation)
		- [Quick Start](#quick-start)
			- [Merge](#merge)
			- [External configuration](#external-configuration)
	- [Development](#development)
	- [Release Process](#release-process)
	- [CHANGELOG](#changelog)
	- [License](#license)

## Features

- [Merge](#merge) files

## Getting started

### Installation

Install `yutil` with your preferred method:

#### deb/rpm/apk:

Download the .deb, .rpm or .apk from the [latest release] and install them with the appropriate tools.

#### Manual

Download the `.tar.gz` from the [latest release] and add the binary to your path.

#### Go users

```bash
go install github.com/amplia-iiot/yutil@latest
```

### Test installation

```bash
yutil version
```

### Quick Start

```bash
yutil help
```

#### Merge

This outputs a formatted (ordered and cleaned) _YAML_ file resulting of merging the passed yaml files (or content).

The files are merged in ascending level of importance in the hierarchy. A yaml node in the last file replaces values in
any previous file. You may pass as many _YAML_ files as desired:

```bash
yutil merge base.yml changes.yml
yutil merge base.yml changes.yml important.yml
```

Use `-o` (`--output`) option if you want to output to a file instead of stdout.

```bash
yutil merge base.yml changes.yml -o merged.yml
```

By default `yutil` uses _stdin_ as the first _YAML_ content:

```bash
cat base.yml | yutil merge changes.yml > merged.yml
```

You may ignore this input if you can't control what's piped to `yutil`:

```bash
echo "this is not a yaml" | yutil --no-input merge base.yml changes.yml
```

#### External configuration

You may want to always use the same config without writting the flags, `yutil` reads a _YAML_ file to configure itself from the current folder or the user home dir in these order of precedence:
- `.yutil.yaml` in current folder
- `.yutil.yml` in current folder
- `.yutil` in current folder
- `.yutil.yaml` in user home dir
- `.yutil.yml` in user home dir
- `.yutil` in user home dir

Sample configuration file:

```yaml
# Disable stdin
no-input: true
# Merge specific config
merge:
  # Merge output file
  output: /tmp/merged.yml
```

You may pass as argument the desired config file:

```bash
# Including extension to support multiple config types
./yutil version --config config.properties
```

> Supported formats: JSON, TOML, YAML, HCL, envfile and Java properties config files

## Development

1. Use Golang version `>= 1.16`
2. Fork (https://github.com/amplia-iiot/yutil)
3. Run `make set-up` to install dev tools
4. Create a feature branch
5. Check changes (test & lint) with `make check`
6. Commit your changes following [Conventional Commits]
7. Rebase your local changes against the upstream _main_ branch
8. Create a Pull Request

You are welcome to report bugs or add feature requests and comments in [issues].

## Release Process

`make version` contains the steps to generate a new version. It uses `svu` to calculate the next version number based on the _git log_ and generates the [CHANGELOG.md] with `git-chglog`

Push the generated _tag_ and the _release_ github action will generate the release.

## CHANGELOG

See [CHANGELOG.md]

## License

[MIT © amplia-iiot](./LICENSE)

[latest release]: https://github.com/amplia-iiot/yutil/releases/latest
[Conventional Commits]: https://www.conventionalcommits.org/en/v1.0.0/
[issues]: https://github.com/amplia-iiot/yutil/issues
[CHANGELOG.md]: ./CHANGELOG.md
