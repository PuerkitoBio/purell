package purell

import (
	"net/url"
	"regexp"
	"strings"
)

type NormalizationFlags int

const (
	FlagLowercaseScheme NormalizationFlags = 1 << iota
	FlagLowercaseHost
	FlagUppercaseEscapes
	FlagDecodeUnnecessaryEscapes
	FlagRemoveDefaultPort

	FlagsSafe NormalizationFlags = FlagLowercaseHost | FlagLowercaseScheme | FlagUppercaseEscapes | FlagDecodeUnnecessaryEscapes | FlagRemoveDefaultPort
)

//var rxEscape = regexp.MustCompile(`(%[0-9a-fA-F]{2})`)
var rxPort = regexp.MustCompile(`(:\d+)/?$`)

func MustNormalizeUrlString(u string, f NormalizationFlags) string {
	if parsed, e := url.Parse(u); e != nil {
		panic(e)
	} else {
		return MustNormalizeUrl(parsed, f)
	}
	panic("Unreachable code.")
}

func MustNormalizeUrl(u *url.URL, f NormalizationFlags) string {
	if res, e := NormalizeUrl(u, f); e != nil {
		panic(e)
	} else {
		return res
	}
	panic("Unreachable code.")
}

func NormalizeUrlString(u string, f NormalizationFlags) (string, error) {
	if parsed, e := url.Parse(u); e != nil {
		return "", e
	} else {
		return NormalizeUrl(parsed, f)
	}
	panic("Unreachable code.")
}

func NormalizeUrl(u *url.URL, f NormalizationFlags) (string, error) {
	var normalized *url.URL = u
	var e error

	// FlagDecodeUnnecessaryEscapes has no action, since it is done automatically
	// by parsing the string as an URL. Same for FlagUppercaseEscapes.
	flags := map[NormalizationFlags]func(*url.URL) (*url.URL, error){
		FlagLowercaseScheme:   lowercaseScheme,
		FlagLowercaseHost:     lowercaseHost,
		FlagRemoveDefaultPort: removeDefaultPort,
	}

	for k, v := range flags {
		if f&k == k {
			if normalized, e = v(normalized); e != nil {
				return "", e
			}
		}
	}
	return normalized.String(), e
}

func lowercaseScheme(u *url.URL) (*url.URL, error) {
	if len(u.Scheme) > 0 {
		u.Scheme = strings.ToLower(u.Scheme)
	}
	return u, nil
}

func lowercaseHost(u *url.URL) (*url.URL, error) {
	if len(u.Host) > 0 {
		u.Host = strings.ToLower(u.Host)
	}
	return u, nil
}

func removeDefaultPort(u *url.URL) (*url.URL, error) {
	u.Host = rxPort.ReplaceAllStringFunc(u.Host, func(val string) string {
		if val == ":80" {
			return ""
		}
		return val
	})
	return u, nil
}
