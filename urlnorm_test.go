package purell

import (
	"testing"
)

// Test cases merged from PR #1
// Originally from https://github.com/jehiah/urlnorm/blob/master/test_urlnorm.py

func assertMap(t *testing.T, cases map[string]string, f NormalizationFlags) {
	for bad, good := range cases {
		s, e := NormalizeURLString(bad, f)
		if e != nil {
			t.Errorf("%s normalizing %v to %v", e.Error(), bad, good)
		} else {
			if s != good {
				t.Errorf("source: %v expected: %v got: %v", bad, good, s)
			}
		}
	}
}

func TestIPv6(t *testing.T) {
	testcases := map[string]string{
		"http://[2001:db8:1f70::999:de8:7648:6e8]/test": "http://[2001:db8:1f70::999:de8:7648:6e8]/test", // ipv6 address
		"http://[::ffff:192.168.1.1]/test":              "http://[::ffff:192.168.1.1]/test",              //ipv4 address in ipv6 notation
		"http://[::ffff:192.168.1.1]:80/test":           "http://[::ffff:192.168.1.1]/test",              //ipv4 address in ipv6 notation
		"htTps://[::fFff:192.168.1.1]:443/test":         "https://[::ffff:192.168.1.1]/test",             //ipv4 address in ipv6 notation
	}

	assertMap(t, testcases, FlagsSafe|FlagRemoveDotSegments)
}

func TestFtp(t *testing.T) {
	testcases := map[string]string{
		"ftp://user:pass@ftp.foo.net/foo/bar": "ftp://user:pass@ftp.foo.net/foo/bar",
	}

	assertMap(t, testcases, FlagsSafe|FlagRemoveDotSegments)
}

func TestStandardCases(t *testing.T) {
	testcases := map[string]string{
		"http://www.foo.com:80/foo":                "http://www.foo.com/foo",
		"http://www.foo.com:8000/foo":              "http://www.foo.com:8000/foo",
		"http://www.foo.com/%7ebar":                "http://www.foo.com/~bar",
		"http://www.foo.com/%7Ebar":                "http://www.foo.com/~bar",
		"http://USER:pass@www.Example.COM/foo/bar": "http://USER:pass@www.example.com/foo/bar",
		"http://test.example/?a=%26&b=1":           "http://test.example/?a=%26&b=1", //should not un-encode the & that is part of a parameter value
		//check that %20 or %25 is not unescaped to " " or %
		"http://test.example/%25/?p=%20val%20%25": "http://test.example/%25/?p=%20val%20%25",
		//check that spaces are collated to "+"
		"http://test.example/path/with a%20space+/": "http://test.example/path/with%20a%20space+/",
		"http://test.example/?":                     "http://test.example/", //no trailing ?
		"http://a.COM/path/?b&a":                    "http://a.com/path/?b&a",
	}

	assertMap(t, testcases, FlagsSafe|FlagRemoveDotSegments)
}

func TestStandardCasesAddTrailingSlash(t *testing.T) {
	testcases := map[string]string{
		"http://test.example?": "http://test.example/", //with trailing /
	}

	assertMap(t, testcases, FlagsSafe|FlagAddTrailingSlash)
}

func TestOctalIP(t *testing.T) {
	testcases := map[string]string{
		"http://0123.011.0.4/":                  "http://0123.011.0.4/",             //NOT octal encoding
		"http://0102.0146.07.0223/":             "http://66.102.7.147/",             //ip octal encoding
		"http://0102.0146.07.0223.:23/":         "http://66.102.7.147.:23/",         //ip octal encoding
		"http://USER:pass@0102.0146.07.0223../": "http://USER:pass@66.102.7.147../", //ip octal encoding
	}

	assertMap(t, testcases, FlagsSafe|FlagDecodeOctalHost)
}

func TestDWORDIP(t *testing.T) {
	testcases := map[string]string{
		"http://123.1113982867/":         "http://123.1113982867/",           //NOT ip dword encoding
		"http://1113982867/":             "http://66.102.7.147/",             //ip dword encoding
		"http://1113982867.:23/":         "http://66.102.7.147.:23/",         //ip dword encoding
		"http://USER:pass@1113982867../": "http://USER:pass@66.102.7.147../", //ip dword encoding
	}

	assertMap(t, testcases, FlagsSafe|FlagDecodeDWORDHost)
}

func TestHexIP(t *testing.T) {
	testcases := map[string]string{
		"http://0x123.1113982867/":       "http://0x123.1113982867/",         //NOT ip hex encoding
		"http://0x42660793/":             "http://66.102.7.147/",             //ip hex encoding
		"http://0x42660793.:23/":         "http://66.102.7.147.:23/",         //ip hex encoding
		"http://USER:pass@0x42660793../": "http://USER:pass@66.102.7.147../", //ip hex encoding
	}

	assertMap(t, testcases, FlagsSafe|FlagDecodeHexHost)
}

func TestUnnecessaryHostDots(t *testing.T) {
	testcases := map[string]string{
		"http://.www.foo.com../foo/bar.html": "http://www.foo.com/foo/bar.html",
		"http://www.foo.com./foo/bar.html":   "http://www.foo.com/foo/bar.html",
		"http://www.foo.com.:81/foo":         "http://www.foo.com:81/foo",
		"http://www.example.com./":           "http://www.example.com/",
	}

	assertMap(t, testcases, FlagsSafe|FlagRemoveUnnecessaryHostDots)
}

func TestEmptyPort(t *testing.T) {
	testcases := map[string]string{
		"http://www.thedraymin.co.uk:/main/?p=308": "http://www.thedraymin.co.uk/main/?p=308", //empty port
		"http://www.src.ca:":                       "http://www.src.ca",                       //empty port
	}

	assertMap(t, testcases, FlagsSafe|FlagRemoveEmptyPortSeparator)
}

// This tests normalization to a unicode representation
// precent escapes for unreserved values are unescaped to their unicode value
// tests normalization to idna domains
// test ip word handling, ipv6 address handling, and trailing domain periods
// in general, this matches google chromes unescaping for things in the address bar.
// spaces are converted to '+' (perhaphs controversial)
// http://code.google.com/p/google-url/ probably is another good reference for this approach
func xTestUrlnorm(t *testing.T) {
	testcases := map[string]string{
		"http://test.example/?a=%e3%82%82%26": "http://test.example/?a=\xe3\x82\x82%26", //should return a unicode character
		"http://s.xn--q-bga.de/":              "http://s.q\xc3\xa9.de/",                 //should be in idna format
		"http://XBLA\u306eXbox.com":           "http://xbla\xe3\x81\xaexbox.com/",       //test utf8 and unicode
		"http://xn--q-bga.XBLA\u306eXbox.com": "http://q\xc3\xa9.//test idna + utf8 domainxbla\xe3\x81\xaexbox.com",

		"http://ja.wikipedia.org/wiki/%E3%82%AD%E3%83%A3%E3%82%BF%E3%83%94%E3%83%A9%E3%83%BC%E3%82%B8%E3%83%A3%E3%83%91%E3%83%B3": "http://ja.wikipedia.org/wiki/\xe3\x82\xad\xe3\x83\xa3\xe3\x82\xbf\xe3\x83\x94\xe3\x83\xa9\xe3\x83\xbc\xe3\x82\xb8\xe3\x83\xa3\xe3\x83\x91\xe3\x83\xb3",

		"http://test.example/\xe3\x82\xad":              "http://test.example/\xe3\x82\xad",
		"http://test.example/?p=%23val#test-%23-val%25": "http://test.example/?p=%23val#test-%23-val%25", //check that %23 (#) is not escaped where it shouldn"t be

		"http://test.domain/I%C3%B1t%C3%ABrn%C3%A2ti%C3%B4n%EF%BF%BDliz%C3%A6ti%C3%B8n": "http://test.domain/I\xc3\xb1t\xc3\xabrn\xc3\xa2ti\xc3\xb4n\xef\xbf\xbdliz\xc3\xa6ti\xc3\xb8n",
	}

	assertMap(t, testcases, FlagsSafe|FlagRemoveDotSegments)
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

	assertMap(t, testcases, FlagsSafe|FlagRemoveDotSegments|FlagRemoveDuplicateSlashes)
}
