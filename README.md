# get-closer
[![GoDoc](https://godoc.org/github.com/akme/get-closer?status.svg)](https://godoc.org/github.com/akme/get-closer) 
[![Build Status](https://travis-ci.org/akme/get-closer.svg?branch=master)](https://travis-ci.org/akme/get-closer) 
[![Go Report Card](https://goreportcard.com/badge/github.com/akme/get-closer)](https://goreportcard.com/report/github.com/akme/get-closer) 
[![Maintainability](https://api.codeclimate.com/v1/badges/4c53c8de5b52f88dd2a3/maintainability)](https://codeclimate.com/github/akme/get-closer/maintainability)
[![codecov](https://codecov.io/gh/akme/get-closer/branch/master/graph/badge.svg)](https://codecov.io/gh/akme/get-closer)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

## Overview
get-closer helps you to find out closest hosts in terms of network latency or application performance.

**Warning:** project is under heavy development, feel free to open an issue
## Motivation
Network services are very sensitive to latency degradation. In era of cloud computing we can get the same services from different providers that will have different locations, so network distance will differs too.  
Choosing an optimal performing endpoint is vital part of optimization.  
We wanted an easy way to compare different endpoints in case of network latency or application performance.
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
  -c, --concurrency uint      concurrency (default 1)
      --config string         config file (default: ~/.get-closer.yaml)
      --count int             number of tests per host (default 1)
  -d, --delay int             set delay between checks (default 3)
  -r, --dns-resolver string   use custom DNS resolver
  -w, --dns-warm-up           warm up DNS cache before request (default true)
  -f, --from-file string      path to file with hosts to check
  -h, --help                  help for get-closer
  -l, --limit uint            number of hosts to return
  -b, --progress-bar          show progress bar (default true)
  -t, --timeout uint          timeout for request (default 60)
  -v, --verbose               enable verbose mode

Use "get-closer [command] --help" for more information about a command.
```

## Examples
This repo contains hosts.txt and dnsresolvers.txt as an example of hosts file that you can use for test runs.
### ICMP
On Linux, you need to enable "unprivileged" mode:
```
sudo sysctl -w net.ipv4.ping_group_range="0   2147483647"
```
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

## Publications
* [Find the best API Endpoint for your EOS dApp](https://medium.com/@akme/find-the-best-api-endpoint-for-your-eos-dapp-7b7489cb6449?source=friends_link&sk=9edf5910202b6cc2701e4fc04c08eb31) [en]
* [Выбираем ближайшие узлы в сети](https://habr.com/ru/post/451866/) [ru]

## Contributing

When contributing to this repository, please first discuss the change you wish to make via issue, email, or any other method with the owners of this repository before making a change.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/akme/get-closer/tags). 


## License
[GNU General Public License v3.0](https://github.com/akme/get-closer/blob/master/LICENSE)
