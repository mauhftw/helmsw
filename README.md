# helmsw

A tool for switch between different version of helm written in go

# WIP

## TODO
- Manage dynamic helm/ directories --> version 0.2.0
- Implement command arguments (cobra) --> version 0.2.0
- Fix multiple versions shown
- Update variable names

## Description

helmsw allows you to download and switch helm versions easily 

## Requirements

The following packages are required by helmsw

| Package Name |         URL            | Minimum required version |
| ------------ | ---------------------- | ------------------------ |
| golang       | https://golang.org/dl/ | 1.12.x                   |


## Installation

1. Checkout repo
```bash
$ git checkout https://github.com/mauhftw/helmsw
```

2. Build helmsw
```bash
$ cd helmsw && \
  make build
```

3. Move helmsw bin to /usr/local/bin
```bash
$ sudo mv dist/helmsw /usr/local/bin
```

4. Add the .helmsw installation directory to your PATH
```bash
$ export PATH=$PATH:$HOME/.helmsw/bin
```

**NOTE:** To avoid problems with symlinks be sure not having a helm binary under your path

## Usage

```bash
$ helmsw
```
