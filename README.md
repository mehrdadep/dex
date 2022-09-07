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
	cache := "/tmp/list.cache"
	extract, _ := dex.New(cache)

	for _, u := range (urls) {
		result:=extract.Parse(u)
		fmt.Printf("%s => %+v\n",u,result)
	}
}

```

Output will look like:

```plain
www.github.com;git+ssh://www.github.com:8443/
git+ssh://www.github.com:8443/ => &{IsIcann:true IsIpV4:false IsIpV6:false IsPrivate:false Subdomain:www Root:github Tld:com}
media.wiki.marvel.co.uk;http://media.wiki.marvel.co.uk
http://media.wiki.marvel.co.uk => &{IsIcann:true IsIpV4:false IsIpV6:false IsPrivate:false Subdomain:media.wiki Root:marvel Tld:co.uk}
210.15.45.32;http://210.15.45.32
http://210.15.45.32 => &{IsIcann:false IsIpV4:true IsIpV6:false IsPrivate:false Subdomain: Root:210.15.45.32 Tld:}
bing.com;http://bing.com?q=dc
http://bing.com?q=dc => &{IsIcann:true IsIpV4:false IsIpV6:false IsPrivate:false Subdomain: Root:bing Tld:com}
```