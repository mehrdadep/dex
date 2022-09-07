package dex

import (
	"log"
	"testing"
)

var (
	cache = "/tmp/tld.cache"
	ex    *Tld
	err   error
)

func init() {
	ex, err = New(cache, true)
	if err != nil {
		log.Fatal(err)
	}
}

func assert(url string, expected *Result, returned *Result, t *testing.T) {
	if expected.IsPrivate == returned.IsPrivate && expected.IsIcann == returned.IsIcann && expected.Root == returned.Root && expected.Subdomain == returned.Subdomain && expected.Tld == returned.Tld {
		return
	}
	t.Errorf("%s: expected:%+v => returned:%+v", url, expected, returned)
}

func TestAll(t *testing.T) {
	cases := map[string]*Result{
		"http://www.google.com": {IsPrivate: false, IsIcann: true, IsIpV6: false, IsIpV4: false, Subdomain: "www", Root: "google", Tld: "com"},
		"https://www.google.co.uk/search?q=mehrdad+esmaeilpour&oq=mehrdad+esmaeilpour&aqs=chrome.0.35i39j0i512j0i22i30l2j69i61l2j69i60j69i65.4902j0j4&sourceid=chrome&ie=UTF-8": {IsPrivate: false, IsIcann: true, IsIpV6: false, IsIpV4: false, Subdomain: "www", Root: "google", Tld: "co.uk"},
		"http://mehrdadep.blogspot.ca":              {IsPrivate: true, IsIcann: false, IsIpV6: false, IsIpV4: false, Subdomain: "mehrdadep", Root: "blogspot", Tld: "ca"},
		"ftp://mehrdadep:password@1992.ftp.com:21/": {IsPrivate: false, IsIcann: true, IsIpV6: false, IsIpV4: false, Subdomain: "1992", Root: "ftp", Tld: "com"},
		"git+ssh://www.github.com:8443/":            {IsPrivate: false, IsIcann: true, IsIpV6: false, IsIpV4: false, Subdomain: "www", Root: "github", Tld: "com"},
		"http://www.!github.com:8443/":              {IsPrivate: false, IsIcann: false, IsIpV6: false, IsIpV4: false, Subdomain: "", Root: "", Tld: ""},
		"http://www.theregister.co.uk":              {IsPrivate: false, IsIcann: true, IsIpV6: false, IsIpV4: false, Subdomain: "www", Root: "theregister", Tld: "co.uk"},
		"http://www.mehrdad.eu.org":                 {IsPrivate: true, IsIcann: false, IsIpV6: false, IsIpV4: false, Subdomain: "www.mehrdad", Root: "eu", Tld: "org"},
		"192.168.0.103":                             {IsPrivate: false, IsIcann: false, IsIpV6: false, IsIpV4: true, Subdomain: "", Root: "192.168.0.103", Tld: ""},
		"http://192.168.0.103":                      {IsPrivate: false, IsIcann: false, IsIpV6: false, IsIpV4: true, Subdomain: "", Root: "192.168.0.103", Tld: ""},
		"http://34.22.nice.coop/":                   {IsPrivate: false, IsIcann: true, IsIpV6: false, IsIpV4: false, Subdomain: "34.22", Root: "nice", Tld: "coop"},
		"http://Gmail.org/":                         {IsPrivate: false, IsIcann: true, IsIpV6: false, IsIpV4: false, Subdomain: "", Root: "gmail", Tld: "org"},
		"http://wiki.info/":                         {IsPrivate: false, IsIcann: true, IsIpV6: false, IsIpV4: false, Subdomain: "", Root: "wiki", Tld: "info"},
		"http://wiki.information/":                  {IsPrivate: false, IsIcann: false, IsIpV6: false, IsIpV4: false, Subdomain: "", Root: "", Tld: ""},
		"http://wiki/":                              {IsPrivate: false, IsIcann: false, IsIpV6: false, IsIpV4: false, Subdomain: "", Root: "", Tld: ""},
		"http://258.15.32.876":                      {IsPrivate: false, IsIcann: false, IsIpV6: false, IsIpV4: false, Subdomain: "", Root: "", Tld: ""},
		"http://www.ai.act.edu.au/":                 {IsPrivate: false, IsIcann: true, IsIpV6: false, IsIpV4: false, Subdomain: "www", Root: "ai", Tld: "act.edu.au"},
		"http://net.cn":                             {IsPrivate: false, IsIcann: false, IsIpV6: false, IsIpV4: false, Subdomain: "", Root: "", Tld: ""},
		"http://google.com?q=marvel":                {IsPrivate: false, IsIcann: true, IsIpV6: false, IsIpV4: false, Subdomain: "", Root: "google", Tld: "com"},
		"ir":                                        {IsPrivate: false, IsIcann: false, IsIpV6: false, IsIpV4: false, Subdomain: "", Root: "", Tld: ""},
		"c.ir":                                      {IsPrivate: false, IsIcann: true, IsIpV6: false, IsIpV4: false, Subdomain: "", Root: "c", Tld: "ir"},
		"b.c.ir":                                    {IsPrivate: false, IsIcann: true, IsIpV6: false, IsIpV4: false, Subdomain: "b", Root: "c", Tld: "ir"},
		"a.b.c.ir":                                  {IsPrivate: false, IsIcann: true, IsIpV6: false, IsIpV4: false, Subdomain: "a.b", Root: "c", Tld: "ir"},
		"c.b.ide.kyoto.jp":                          {IsPrivate: false, IsIcann: true, IsIpV6: false, IsIpV4: false, Subdomain: "c", Root: "b", Tld: "ide.kyoto.jp"},
	}
	for url, expected := range cases {
		returned := ex.Parse(url)
		assert(url, expected, returned, t)
	}
}
