# party-aas

Takes and image and turns it into a party (gif)! By default reads a jpg or png
from stdin and writes it to stdout. Can also be run as a web service by passing
the `-addr` option.

![PARTY](/out.gif)

# Install

```
$ go get github.com/disintegration/imaging
$ go get github.com/akenn/graphics-go/graphics
```

# Usage

```
$ go run index.go main.go --addr=localhost:8000
```

