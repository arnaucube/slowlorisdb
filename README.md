# slowlorisdb [![Go Report Card](https://goreportcard.com/badge/github.com/arnaucube/slowlorisdb)](https://goreportcard.com/report/github.com/arnaucube/slowlorisdb) [![Build Status](https://travis-ci.org/arnaucube/slowlorisdb.svg?branch=master)](https://travis-ci.org/arnaucube/slowlorisdb) [![GoDoc](https://godoc.org/github.com/arnaucube/slowlorisdb?status.svg)](https://godoc.org/github.com/arnaucube/slowlorisdb)

Slow, decentralized and cryptographically consistent database

Basically this repo is a blockchain written from scratch, that allows to launch multiple simultaneous blockchains.

**Warning**: this project was started in the free time of a long travel, not having much more free time to continue developing it.

Watch the blockchain in action: http://www.youtubemultiplier.com/5ca9c1a540b31-slowlorisdb-visual-representation.php

![slowloris](https://04019a5a-a-62cb3a1a-s-sites.googlegroups.com/site/jchristensensdigitalportfolio/slow-loris/IO-moth-eating-frozen-apple-sauce.jpg "slowloris")


## Run
The repo is under construction

- create node
```sh
# node0
go run main.go --config config0.yaml create

# node1
go run main.go --config config1.yaml create
```

- run node
```sh
# node0
go run main.go --config config0.yaml start

# node1
go run main.go --config config1.yaml start
```
