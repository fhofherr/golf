# Go Logging Facade

The Go Logging Facade (golf) is a thin abstraction over various
structured loggers.

[![Go Reference](https://pkg.go.dev/badge/github.com/fhofherr/golf.svg)](https://pkg.go.dev/github.com/fhofherr/golf)
[![Go Report Card](https://goreportcard.com/badge/github.com/fhofherr/golf)](https://goreportcard.com/report/github.com/fhofherr/golf)
[![Coverage Status](https://coveralls.io/repos/github/fhofherr/golf/badge.svg?branch=main)](https://coveralls.io/github/fhofherr/golf?branch=main)
![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/fhofherr/golf.svg)

`golf` adheres to semantic versioning. Any breaking changes before version
1.0 or between major versions are listed under [Breaking
Changes](#breaking-changes).

## Philosophy

This package strives to provide a facade over as many loggers as
possible. This allows to keep your code independent from those loggers
and to exchange them as needed. By using the `golf.Logger` interface
library code can accept an optional logger which allows it to use
whatever logger the calling application uses.

Additionally using different loggers during test execution can be handy,
for example if the fact that a message was logged, or its contents need
to be asserted.

Versions before 0.3.0 used to define a Logger interfaces that was
identical to go-kit's
[Logger](https://pkg.go.dev/github.com/go-kit/log#Logger) interface.
However, I found that calling code almost never checks the error
returned by Log. In most cases it is just ignored. Usually this goes
along with configuring [`errcheck`](https://github.com/kisielk/errcheck)
to ignore the `Logger` interface.

In most cases code calling the `Log` method should not be concerned by
an error in the logging library. Ideally the operation should continue
to be executed. If it is really necessary to handle an error that
occurred during logging, golf provides an
[`Error`](https://pkg.go.dev/github.com/fhofherr/golf#Error) method.
This method allows to access the last error that occurred during
logging, if the underlying logger adapter supports it.

## Known Adapters

* [stdlib](https://pkg.go.dev/log): included in main repository
* [zap](https://github.com/uber-go/zap): https://github.com/fhofherr/golf-zap

## Breaking Changes

The following lists all breaking changes between versions.

### [0.3.0]

* Moved `log.Logger` to `golf.Logger`
* Moved `log.Log` to `golf.Log`
* Moved `log.NewStdlib` to `golf.NewStdlib`
* Moved `log.Formatter` to `golf.Formatter`
* Moved `log.PlainTextFormatter` to `golf.PlainTextFormatter`
* Moved `log.JSONFormatter` to `golf.JSONFormatter`
* Removed remaining `log` package

## License

Copyright © 2021 Ferdinand Hofherr

Distributed under the MIT License.

[0.3.0]: https://github.com/fhofherr/golf/compare/v0.2.0...v0.3.0
