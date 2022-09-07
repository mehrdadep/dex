package dex

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
)

type Result struct {
	IsIcann   bool
	IsIpV4    bool
	IsIpV6    bool
	IsPrivate bool
	Subdomain string
	Root      string
	Tld       string
}

type Tld struct {
	CacheFile  string
	rootNode   *Trie
	debug      bool
	noValidate bool // do not validate URL schema
	noStrip    bool // do not strip .html suffix from URL
}

type Trie struct {
	ExceptRule bool
	ValidTld   bool
	IsIcann    bool
	IsPrivate  bool
	matches    map[string]*Trie
}

var (
	schemaRegex = regexp.MustCompile(`^([[:lower:]\d\+\-\.]+:)?//`)
	domainRegex = regexp.MustCompile(`^[a-z0-9-\p{L}]{1,63}$`)
	ip4Regex    = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])`)
)

// New creates a new *Tld, it may be shared between goroutines, we usually need a single instance in an application.
func New(cacheFile string, debug bool) (*Tld, error) {
	data, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		data, err = readFromUrl()
		if err != nil {
			return &Tld{}, err
		}
		if err = ioutil.WriteFile(cacheFile, data, 0644); err != nil {
			return &Tld{}, err
		}
	}
	ts := strings.Split(string(data), "\n")
	newMap := make(map[string]*Trie)
	rootNode := &Trie{ExceptRule: false, ValidTld: false, IsIcann: false, IsPrivate: false, matches: newMap}
	isIcann := false
	isPrivate := false
	for _, t := range ts {
		if t != "" && !strings.HasPrefix(t, "// ===BEGIN ICANN DOMAINS===") {
			isIcann = true
			isPrivate = false
		}
		if t != "" && !strings.HasPrefix(t, "// ===BEGIN PRIVATE DOMAINS===") {
			isIcann = false
			isPrivate = true
		}
		if t != "" && !strings.HasPrefix(t, "//") {
			t = strings.TrimSpace(t)
			exceptionRule := t[0] == '!'
			if exceptionRule {
				t = t[1:]
			}
			addTldRule(rootNode, strings.Split(t, "."), exceptionRule, isIcann, isPrivate)
		}
	}

	return &Tld{CacheFile: cacheFile, rootNode: rootNode, debug: debug}, nil
}

// SetNoValidate disables schema check in order to increase performance.
func (tEx *Tld) SetNoValidate() {
	tEx.noValidate = true
}

// SetNoStrip disables URL stripping in order to increase performance.
func (tEx *Tld) SetNoStrip() {
	tEx.noStrip = true
}

func addTldRule(rootNode *Trie, labels []string, ex, icann, private bool) {
	n := len(labels)
	t := rootNode
	for i := n - 1; i >= 0; i-- {
		lab := labels[i]
		m, found := t.matches[lab]
		if !found {
			except := ex
			valid := !ex && i == 0
			newMap := make(map[string]*Trie)
			t.matches[lab] = &Trie{ExceptRule: except, ValidTld: valid, matches: newMap, IsIcann: icann, IsPrivate: private}
			m = t.matches[lab]
		}
		t = m
	}
}

func (tEx *Tld) Parse(u string) *Result {
	input := u
	u = strings.ToLower(u)
	if !tEx.noValidate {
		u = schemaRegex.ReplaceAllString(u, "")
		i := strings.Index(u, "@")
		if i != -1 {
			u = u[i+1:]
		}

		index := strings.IndexFunc(u, func(r rune) bool {
			switch r {
			case '&', '/', '?', ':', '#':
				return true
			}
			return false
		})
		if index != -1 {
			u = u[0:index]
		}
	}
	if !tEx.noStrip {
		if strings.HasSuffix(u, ".html") {
			u = u[0 : len(u)-len(".html")]
		}
	}
	if tEx.debug {
		fmt.Printf("%s;%s\n", u, input)
	}
	return tEx.extract(u)
}

func (tEx *Tld) extract(url string) *Result {
	domain, tld, private, icann := tEx.extractTld(url)
	if tld == "" {
		ip := net.ParseIP(url)
		if ip != nil {
			if ip4Regex.MatchString(url) {
				return &Result{IsIpV4: true, Root: url, IsIcann: false, IsPrivate: false}
			}
			return &Result{IsIpV6: true, Root: url, IsIcann: false, IsPrivate: false}
		}
		return &Result{IsIpV6: false, IsIpV4: false, IsIcann: false, IsPrivate: false}
	}
	sub, root := subdomain(domain)
	if domainRegex.MatchString(root) {
		return &Result{Root: root, Subdomain: sub, Tld: tld, IsIcann: icann, IsPrivate: private, IsIpV4: false, IsIpV6: false}
	}
	return &Result{IsIpV6: false, IsIpV4: false, IsIcann: false, IsPrivate: false}
}

func (tEx *Tld) extractTld(url string) (domain, tld string, private, icann bool) {
	spl := strings.Split(url, ".")
	tldIndex, validTld, private, icann := tEx.getIndex(spl)
	if validTld {
		domain = strings.Join(spl[:tldIndex], ".")
		tld = strings.Join(spl[tldIndex:], ".")
	} else {
		domain = url
	}
	return
}

func (tEx *Tld) getIndex(labels []string) (int, bool, bool, bool) {
	t := tEx.rootNode
	parentValid := false
	private := false
	icann := false
	for i := len(labels) - 1; i >= 0; i-- {
		lab := labels[i]
		n, found := t.matches[lab]
		_, star := t.matches["*"]

		switch {
		case found && !n.ExceptRule:
			parentValid = n.ValidTld
			private = n.IsPrivate
			icann = n.IsIcann
			t = n
		// Found an exception rule
		case found:
			fallthrough
		case parentValid:
			return i + 1, true, private, icann
		case star:
			parentValid = true
		default:
			return -1, false, private, icann
		}
	}
	return -1, false, private, icann
}

//return sub domain,root domain
func subdomain(d string) (string, string) {
	ps := strings.Split(d, ".")
	l := len(ps)
	if l == 1 {
		return "", d
	}
	return strings.Join(ps[0:l-1], "."), ps[l-1]
}

func readFromUrl() ([]byte, error) {
	u := "https://publicsuffix.org/list/public_suffix_list.dat"
	resp, err := http.Get(u)
	if err != nil {
		return []byte(""), err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	lines := strings.Split(string(body), "\n")
	var buffer bytes.Buffer

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "//") {
			buffer.WriteString(line)
			buffer.WriteString("\n")
		}
	}

	return buffer.Bytes(), nil
}
