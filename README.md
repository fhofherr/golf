# Go Logging Facade

The Go Logging Facade (golf) is a thin abstraction over various
structured loggers heavily inspired by [GoKit](https://gokit.io/)'s
[`Logger`](https://godoc.org/github.com/go-kit/kit/log#Logger)
interface.

## Project Goals

* API compatible with GoKit `Logger`.
* Main repository depends only on Go standard library.
* Easy implementation of adapters for various loggers.
  * [stdlib](https://godoc.org/log): included in main repository
  * [logrus](https://github.com/Sirupsen/logrus): separate repository
  * [zap](https://github.com/uber-go/zap): separate repository

## License

Copyright © 2019 Ferdinand Hofherr

Distributed under the MIT License.
