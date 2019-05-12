# get-closer : find out closest hosts in terms of network latency or application performance

## Overview [![GoDoc](https://godoc.org/github.com/akme/get-closer?status.svg)](https://godoc.org/github.com/akme/get-closer) [![Build Status](https://travis-ci.org/akme/get-closer.svg?branch=master)](https://travis-ci.org/akme/get-closer) [![Go Report Card](https://goreportcard.com/badge/github.com/akme/get-closer)](https://goreportcard.com/report/github.com/akme/get-closer) [![codecov](https://codecov.io/gh/akme/get-closer/branch/master/graph/badge.svg)](https://codecov.io/gh/akme/get-closer)


Network latency affects performance, so choosing closest hosts is vital part of monitoring and optimization.

## Install

```
go get -u github.com/akme/get-closer
```

## How to use
```
Find out closest hosts in terms of network latency.

Usage:
  get-closer [command]

Available Commands:
  dns         Measuring domain resolve time via DNS resolver
  help        Help about any command
  http        Measuring time for HTTP request.
  icmp        Measuring RTT for hosts from list.
  tcp         Measuring time for connecting to open TCP port.

Flags:
  -c, --concurrency uint    Concurrency (default 1)
      --config string       config file (default: ~/.get-closer.yaml)
      --count int           number of tests per host (default 1)
      --dns-server string   use custom DNS resolver
  -w, --dns-warm-up         warm up DNS cache before request (default true)
  -f, --from-file string    Path to file with hosts to check
  -h, --help                help for get-closer
  -l, --limit uint          number of hosts to return
  -b, --progress-bar        show progress bar (default true)
  -t, --timeout uint        Timeout for request (default 60)
  -v, --verbose             enable verbose mode

Use "get-closer [command] --help" for more information about a command.
```

## Examples
### ICMP
```
get closer ping -f hosts.txt
```
### TCP
```
get closer tcp -f hosts.txt
```
### HTTP
```
get-closer http -f hosts.txt
```
### DNS
```
get-closer dns -f hosts.txt
```


## Contributing

When contributing to this repository, please first discuss the change you wish to make via issue, email, or any other method with the owners of this repository before making a change.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/akme/get-closer/tags). 


## License
[GNU General Public License v3.0](https://github.com/akme/get-closer/blob/master/LICENSE)
