# party-aas

Takes and image and turns it into a party (gif)! By default reads a jpg or png
from stdin and writes it to stdout. Can also be run as a web service by passing
the `-addr` option.

![PARTY](/out.gif)

# Install

```
$ go get
```

# Exampe Usages

party-aas can be run as either a webserver or a normal command-line utility that reads from stdin.

Here are example usages of each:

## Command line tool

```
$ cat fire.png | go run index.go main.go -counterclockwise=true > fire-party.gif
```

## Webserver

```
$ go run index.go main.go -addr=localhost:8000
```

