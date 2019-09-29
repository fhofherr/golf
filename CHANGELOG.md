# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep
a Changelog](https://keepachangelog.com/en/1.0.0/), and this project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

* `.projections.json` file. See
  [vim-projectionist](https://github.com/tpope/vim-projectionist) for
  details.
* `.editorconfig` file. See [editorconfig](https://editorconfig.org/)
  for details.

## [0.1.0] - 2019-04-22

### Added

* [README](README.md) detailing project goals and initial features.
* `log.Logger` interface. `log.Logger` has the same basic contract as
  [GoKit's `log.Logger`](https://godoc.org/github.com/go-kit/kit/log#Logger)
  interface.
* A basic logger that uses the passed writer for logging. It is intended
  for cases in which such a basic logger is all you need. Most callers
  should use one of the wrappers around other logging libraries (once
  they are ready).
* Enable contextual logging by providing the `With` function. To avoid
  excessive wrapping of contextual loggers the `Wither` interface can be
  implemented by loggers who know how to add additional context values
  to themselves.
* Adapter for the Go standard library logger.

[Unreleased]: https://github.com/fhofherr/golf/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/fhofherr/golf/releases/tag/v0.1.0
