package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func indexHTTP(w http.ResponseWriter, r *http.Request) {
	index, _ := ioutil.ReadFile("index.html")
	fmt.Fprint(w, string(index))
}
