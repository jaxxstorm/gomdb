The Golang Omdb API [![Build Status](https://travis-ci.org/agnivade/gomdb.svg?branch=master)](https://travis-ci.org/agnivade/gomdb) [![GoDoc](https://godoc.org/github.com/agnivade/gomdb?status.svg)](https://godoc.org/github.com/agnivade/gomdb)
===================

This API uses the wonderful [omdbapi.com](http://omdbapi.com/) API by Brian Fritz. It is an implementation of that API in golang.

Usage
-----

```go
package main

import (
	"fmt"
	"github.com/agnivade/gomdb"
)

func main() {
	query := &gomdb.QueryData{Title: "Macbeth", SearchType: gomdb.MovieSearch}
	res, err := gomdb.Search(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.Search)

	query = &gomdb.QueryData{Title: "Macbeth", Year: "2015"}
	res2, err := gomdb.LookupByTitle(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res2)

	res3, err := gomdb.LookupByImdbID("tt2884018")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res3)
}
```

Kindly look into the godocs for documentation on the response objects.

