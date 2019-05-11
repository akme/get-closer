# get-closer : Find out closest hosts in terms of network latency

## Overview [![GoDoc](https://godoc.org/github.com/akme/get-closer?status.svg)](https://godoc.org/github.com/akme/get-closer) [![Build Status](https://travis-ci.org/akme/get-closer.svg?branch=master)](https://travis-ci.org/akme/get-closer) [![Go Report Card](https://goreportcard.com/badge/github.com/akme/get-closer)](https://goreportcard.com/report/github.com/akme/get-closer) [![codecov](https://codecov.io/gh/akme/get-closer/branch/master/graph/badge.svg)](https://codecov.io/gh/akme/get-closer)


Network latency affects performance, so choosing closest hosts is vital part of optimization.

## Install

```
go get github.com/akme/get-closer
```

## Example
HTTP
```
get-closer http -f hosts.txt
```
DNS
```
get-closer dns -f hosts.txt
```


## Contributing

ToDo.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/akme/get-closer/tags). 


## License
[GNU General Public License v3.0](https://github.com/akme/get-closer/blob/master/LICENSE)
