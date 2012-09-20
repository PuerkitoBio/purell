package purell

import (
	"log"
	"testing"
)

// This tests normalization to a unicode representation
// precent escapes for unreserved values are unescaped to their unicode value
// tests normalization to idna domains
// test ip word handling, ipv6 address handling, and trailing domain periods
// in general, this matches google chromes unescaping for things in the address bar.
// spaces are converted to '+' (perhaphs controversial)
// http://code.google.com/p/google-url/ probably is another good reference for this approach
func TestUrlnorm(t *testing.T) {
	// from https://github.com/jehiah/urlnorm/blob/master/test_urlnorm.py
	testcases := map[string]string{
		"http://1113982867/":                       "http://66.102.7.147/",                    //ip dword encoding
		"http://www.thedraymin.co.uk:/main/?p=308": "http://www.thedraymin.co.uk/main/?p=308", //empty port
		"http://www.foo.com:80/foo":                "http://www.foo.com/foo",
		"http://www.foo.com:8000/foo":              "http://www.foo.com:8000/foo",
		"http://www.foo.com./foo/bar.html":         "http://www.foo.com/foo/bar.html",
		"http://www.foo.com.:81/foo":               "http://www.foo.com:81/foo",
		"http://www.foo.com/%7ebar":                "http://www.foo.com/~bar",
		"http://www.foo.com/%7Ebar":                "http://www.foo.com/~bar",
		"ftp://user:pass@ftp.foo.net/foo/bar":      "ftp://user:pass@ftp.foo.net/foo/bar",
		"http://USER:pass@www.Example.COM/foo/bar": "http://USER:pass@www.example.com/foo/bar",
		"http://www.example.com./":                 "http://www.example.com/",
		"http://test.example/?a=%26&b=1":           "http://test.example/?a=%26&b=1",         //should not un-encode the & that is part of a parameter value
		"http://test.example/?a=%e3%82%82%26":      "http://test.example/?a=\xe3\x82\x82%26", //should return a unicode character

		"http://s.xn--q-bga.de/": "http://s.q\xc3\xa9.de/", //should be in idna format
		"http://test.example/?":  "http://test.example/",   //no trailing ?
		"http://test.example?":   "http://test.example/",   //with trailing /
		"http://a.COM/path/?b&a": "http://a.com/path/?b&a",
		//test utf8 and unicode
		"http://XBLA\u306eXbox.com": "http://xbla\xe3\x81\xaexbox.com/",
		//test idna + utf8 domain
		"http://xn--q-bga.XBLA\u306eXbox.com":                                                                                     "http://q\xc3\xa9.xbla\xe3\x81\xaexbox.com",
		"http://ja.wikipedia.org/wiki/%E3%82%AD%E3%83%A3%E3%82%BF%E3%83%94%E3%83%A9%E3%83%BC%E3%82%B8%E3%83%A3%E3%83%91%E3%83%B3": "http://ja.wikipedia.org/wiki/\xe3\x82\xad\xe3\x83\xa3\xe3\x82\xbf\xe3\x83\x94\xe3\x83\xa9\xe3\x83\xbc\xe3\x82\xb8\xe3\x83\xa3\xe3\x83\x91\xe3\x83\xb3",
		"http://test.example/\xe3\x82\xad":                                                                                        "http://test.example/\xe3\x82\xad",

		//check that %23 (#) is not escaped where it shouldn"t be
		"http://test.example/?p=%23val#test-%23-val%25": "http://test.example/?p=%23val#test-%23-val%25",
		//check that %20 or %25 is not unescaped to " " or %
		"http://test.example/%25/?p=%20val%20%25":                                       "http://test.example/%25/?p=%20val%20%25",
		"http://test.domain/I%C3%B1t%C3%ABrn%C3%A2ti%C3%B4n%EF%BF%BDliz%C3%A6ti%C3%B8n": "http://test.domain/I\xc3\xb1t\xc3\xabrn\xc3\xa2ti\xc3\xb4n\xef\xbf\xbdliz\xc3\xa6ti\xc3\xb8n",
		//check that spaces are collated to "+"
		"http://test.example/path/with a%20space+/":     "http://test.example/path/with%20a%20space+/",
		"http://[2001:db8:1f70::999:de8:7648:6e8]/test": "http://[2001:db8:1f70::999:de8:7648:6e8]/test", // ipv6 address
		"http://[::ffff:192.168.1.1]/test":              "http://[::ffff:192.168.1.1]/test",              //ipv4 address in ipv6 notation
		"http://[::ffff:192.168.1.1]:80/test":           "http://[::ffff:192.168.1.1]/test",              //ipv4 address in ipv6 notation
		"htTps://[::fFff:192.168.1.1]:443/test":         "https://[::ffff:192.168.1.1]/test",             //ipv4 address in ipv6 notation
	}
	
	for bad, good := range testcases {
		s, e := NormalizeUrlString(bad, FlagsSafe | FlagRemoveDotSegments);
		if e != nil {
			log.Printf("%s normalizing %v to %v", e.Error(), bad, good)
			t.Fail()
		} else {
			if s != good {
				log.Printf("expected: %v got: %v", good, s)
				t.Fail()
			}
		}
	}
}

func TestSlashes(t *testing.T) {
	// some taken from RFC1808, rfc3986
	testcases := map[string]string{
		"http://test.example/foo/bar/.":               "http://test.example/foo/bar/",
		"http://test.example/foo/bar/./":              "http://test.example/foo/bar/",
		"http://test.example/foo/bar/..":              "http://test.example/foo/",
		"http://test.example/foo/bar/../":             "http://test.example/foo/",
		"http://test.example/foo/bar/../baz":          "http://test.example/foo/baz",
		"http://test.example/foo/bar/../..":           "http://test.example/",
		"http://test.example/foo/bar/../../":          "http://test.example/",
		"http://test.example/foo/bar/../../baz":       "http://test.example/baz",
		"http://test.example/foo/bar/../../../baz":    "http://test.example/baz",
		"http://test.example/foo/bar/../../../../baz": "http://test.example/baz",
		"http://test.example/./foo":                   "http://test.example/foo",
		"http://test.example/../foo":                  "http://test.example/foo",
		"http://test.example/foo.":                    "http://test.example/foo.",
		"http://test.example/.foo":                    "http://test.example/.foo",
		"http://test.example/foo..":                   "http://test.example/foo..",
		"http://test.example/..foo":                   "http://test.example/..foo",
		"http://test.example/./../foo":                "http://test.example/foo",
		"http://test.example/./foo/.":                 "http://test.example/foo/",
		"http://test.example/foo/./bar":               "http://test.example/foo/bar",
		"http://test.example/foo/../bar":              "http://test.example/bar",
		"http://test.example/foo//":                   "http://test.example/foo/",
		"http://test.example/foo///bar//":             "http://test.example/foo/bar/",
	}

	for bad, good := range testcases {
		s, e := NormalizeUrlString(bad, FlagsSafe | FlagRemoveDotSegments);
		if e != nil {
			log.Printf("%s normalizing %v to %v", e.Error(), bad, good)
			t.Fail()
		} else {
			if s != good {
				log.Printf("expected: %v got: %v", good, s)
				t.Fail()
			}
		}
	}
}
