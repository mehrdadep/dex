# Dex (Domain Extractor)


Extract root domain, subdomain name, tld from an url, using [the Public Suffix List](http://www.publicsuffix.org).

## Installation

Install dex:
```sh
go get github.com/mehrdadep/dex

```

To run unit tests, run this command  in dex's source directory($GOPATH/src/github.com/mehrdadep/dex):

```sh
go test
```

## Example

```go
package main

import (
"fmt"
"github.com/mehrdadep/dex"
)


func main() {
	urls := []string{"git+ssh://www.github.com:8443/", "http://media.wiki.marvel.co.uk", "http://210.15.45.32", "http://bing.com?q=dc"}
	cache := "/path/to/public/suffix/list.cache"
	extract, _ := dex.New(cache,false)

	for _, u := range (urls) {
		result:=extract.Parse(u)
		fmt.Printf("%s => %+v\n",u,result)
	}
}

```
Output will look like:
```plain
