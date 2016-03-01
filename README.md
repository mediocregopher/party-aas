# party-aas

Takes and image and turns it into a party (gif)! By default reads a jpg or png
from stdin and writes it to stdout. Can also be run as a web service by passing
the `-addr` option.

![PARTY](/out.gif)

# Building

```
$ cd party-aas
$ go get
$ go install github.com/mediocregopher/varembed
$ go generate
$ go build
```

# Exampe Usages

party-aas can be run as either a webserver or a normal command-line utility that reads from stdin.

Here are example usages of each:

## Command line tool

```
$ cat fire.png | ./party-aas -counterclockwise=true > fire-party.gif
```

## Webserver

```
$ ./party-aas -addr=localhost:8000
```

