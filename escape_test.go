package purell

import (
	"net/url"
	"testing"
)

func TestNormalizeUsesStdlibString(t *testing.T) {
	u, err := url.Parse("HTTP://www.SRC.ca:80/to%1ato%8b%ee/OKnow%41%42%43%7e?")
	if err != nil {
		t.Fatalf("parse error: %s", err)
	}

	got := NormalizeURL(u, FlagLowercaseHost|FlagRemoveDefaultPort)
	const want = "http://www.src.ca/to%1Ato%8B%EE/OKnowABC~"
	if got != want {
		t.Fatalf("NormalizeURL() = %q, want %q", got, want)
	}
}

func TestNormalizeKeepsPurellExtraPathCharactersUnescaped(t *testing.T) {
	got, err := NormalizeURLString("http://example.com/%21%27%28%29%2A%5B%5D", FlagsSafe)
	if err != nil {
		t.Fatalf("parse error: %s", err)
	}

	const want = "http://example.com/!'()*[]"
	if got != want {
		t.Fatalf("NormalizeURLString() = %q, want %q", got, want)
	}
}

func TestNormalizeKeepsPurellExtraFragmentCharactersUnescaped(t *testing.T) {
	got, err := NormalizeURLString("http://example.com/#%27%5B%5D", FlagsSafe)
	if err != nil {
		t.Fatalf("parse error: %s", err)
	}

	const want = "http://example.com/#'[]"
	if got != want {
		t.Fatalf("NormalizeURLString() = %q, want %q", got, want)
	}
}

func TestNormalizeKeepsPurellExtraUserinfoCharactersUnescaped(t *testing.T) {
	u := &url.URL{
		Scheme: "http",
		User:   url.UserPassword("!$&'()", "*+,;="),
		Host:   "example.com",
	}

	got := NormalizeURL(u, FlagsSafe)
	const want = "http://!$&'():*+,;=@example.com"
	if got != want {
		t.Fatalf("NormalizeURL() = %q, want %q", got, want)
	}
}

func TestNormalizeKeepsPurellNonAuthoritySchemePathRendering(t *testing.T) {
	got, err := NormalizeURLString("mailto:/webmaster@golang.org", FlagsSafe)
	if err != nil {
		t.Fatalf("parse error: %s", err)
	}

	const want = "mailto:/webmaster@golang.org"
	if got != want {
		t.Fatalf("NormalizeURLString() = %q, want %q", got, want)
	}
}

func TestNormalizeMailtoQueryWithoutUser(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"mailto:?body=foo", "mailto:?body=foo"},
		{"mailto:", "mailto:"},
		{"mailto:user@example.com", "mailto:user@example.com"},
		{"TEL:+15551212", "tel:+15551212"},
	}
	for _, tt := range tests {
		got, err := NormalizeURLString(tt.in, FlagsSafe)
		if err != nil {
			t.Fatalf("NormalizeURLString(%q) error: %s", tt.in, err)
		}
		if got != tt.want {
			t.Errorf("NormalizeURLString(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestNormalizeKeepsPurellEmptyHostSchemeRendering(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"file:", "file://"},
		{"mailto:", "mailto:"},
		{"file:?q", "file://?q"},
		{"file:#frag", "file://#frag"},
	}

	for _, tt := range tests {
		got, err := NormalizeURLString(tt.in, FlagsSafe)
		if err != nil {
			t.Fatalf("NormalizeURLString(%q) error: %s", tt.in, err)
		}
		if got != tt.want {
			t.Errorf("NormalizeURLString(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestNormalizeURLKeepsUnicodeHostUnescaped(t *testing.T) {
	u := &url.URL{
		Scheme: "http",
		Host:   "é.com",
		Path:   "/p",
	}

	got := NormalizeURL(u, FlagsSafe)
	const want = "http://é.com/p"
	if got != want {
		t.Fatalf("NormalizeURL() = %q, want %q", got, want)
	}
}

func TestNormalizeURLKeepsUnicodeHostWhenUserinfoMatchesEscapedHost(t *testing.T) {
	u := &url.URL{
		Scheme: "http",
		User:   url.User("é.com"),
		Host:   "é.com",
	}

	got := NormalizeURL(u, FlagsSafe)
	const want = "http://%C3%A9.com@é.com"
	if got != want {
		t.Fatalf("NormalizeURL() = %q, want %q", got, want)
	}
}

func TestNormalizeURLKeepsRelativeColonPath(t *testing.T) {
	u := &url.URL{
		Path: "a:b",
	}

	got := NormalizeURL(u, FlagsSafe)
	const want = "a:b"
	if got != want {
		t.Fatalf("NormalizeURL() = %q, want %q", got, want)
	}
}
