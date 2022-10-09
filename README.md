# maker

> GNU make snippets manager

![Status](https://img.shields.io/badge/status-active-success.svg)
[![GitHub Issues](https://img.shields.io/github/issues/wwmoraes/maker.svg)](https://github.com/wwmoraes/maker/issues)
[![GitHub Pull Requests](https://img.shields.io/github/issues-pr/wwmoraes/maker.svg)](https://github.com/wwmoraes/maker/pulls)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwwmoraes%2Fmaker.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwwmoraes%2Fmaker?ref=badge_shield)

[![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/wwmoraes/maker)](https://hub.docker.com/r/wwmoraes/maker)
[![Docker Image Version (latest semver)](https://img.shields.io/docker/v/wwmoraes/maker?label=image%20version)](https://hub.docker.com/r/wwmoraes/maker)
[![Docker Pulls](https://img.shields.io/docker/pulls/wwmoraes/maker)](https://hub.docker.com/r/wwmoraes/maker)

[![Maintainability](https://api.codeclimate.com/v1/badges/7f142f813859a82c2203/maintainability)](https://codeclimate.com/github/wwmoraes/docker-engine-plugins/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/7f142f813859a82c2203/test_coverage)](https://codeclimate.com/github/wwmoraes/docker-engine-plugins/test_coverage)

---

## 📝 Table of Contents

- [About](#-about)
- [Getting Started](#-getting-started)
- [Usage](#-usage)
- [Built Using](#-built-using)
- [TODO](./TODO.md)
- [Contributing](./CONTRIBUTING.md)
- [Authors](#-authors)
- [Acknowledgments](#-acknowledgements)

## 🧐 About

Maker provides CLI tools to manage and use GNU make-compliant snippet files,
which contain rules to deal with specific domains such as programming language
toolchains or helper tools as docker, linters, dependency injectors (google
wire, uber fx, etc) that are needed during any stage of a repository management.

Those snippet files rely on a convention-over-configuration style that uses
variables that can be either set before the include happens or overridden later
on. They also should not be modified directly on the client repositories, as
they are meant to be updated with the upstream version.

### Why?

- GNU make is true to the [Unix philosophy][unix-philosophy]: it does one thing,
and does it well (since 1976!)
- newer language tool kits (2010+ at least) often include their own tool chain,
which adds extra commands and flags to learn and repeat yourself when using them
- convention-over-configuration for common/best practices, which speed up
project creation, maintenance and usage

## 🏁 Getting Started

Clone the repository and use `make`/`make all` to build, lint, test and generate
coverage information.

To contribute please make sure you have installed the pre-commit hooks before
committing.

### Prerequisites

For building locally:

- Golang 1.15+
- GNU Make (optional)

For contributing:

- `pre-commit` tool
- `golangci-lint` binary on path
- `hadolint` binary on path

Install pre-commit using your package manager and then run `pre-commit install`
once to configure the repository hooks. You can then commit normally.

## 🔧 Running the tests

All tests can be run using `make test`. Coverage is done through `make coverage`.

## 🎈 Usage

Add notes about how to use the system.

## 🔧 Built Using

- [GNU make](https://www.gnu.org/software/make/) - 🖤
- [Golang](https://golang.org) - base language

## 🧑‍💻 Authors

- [@wwmoraes](https://github.com/wwmoraes) - Idea & Initial work

## 🎉 Acknowledgements

- Hat tip to anyone whose code was used
- Inspiration
- References

[unix-philosophy]: http://www.catb.org/esr/writings/taoup/html/ch01s06.html


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwwmoraes%2Fmaker.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwwmoraes%2Fmaker?ref=badge_large)