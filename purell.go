package purell

import (
	"net/url"
	"regexp"
	"strings"
)

type NormalizationFlags int

const (
	// Safe normalizations
	FlagLowercaseScheme NormalizationFlags = 1 << iota
	FlagLowercaseHost
	FlagUppercaseEscapes
	FlagDecodeUnnecessaryEscapes
	FlagRemoveDefaultPort

	// Usually safe normalizations

	// Should choose one or the other (add-remove slash)
	FlagRemoveTrailingSlash
	FlagAddTrailingSlash

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
		FlagLowercaseScheme:     lowercaseScheme,
		FlagLowercaseHost:       lowercaseHost,
		FlagRemoveDefaultPort:   removeDefaultPort,
		FlagRemoveTrailingSlash: removeTrailingSlash,
		FlagAddTrailingSlash:    addTrailingSlash,
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

func removeTrailingSlash(u *url.URL) (*url.URL, error) {
	if l := len(u.Path); l > 0 && strings.HasSuffix(u.Path, "/") {
		u.Path = u.Path[:l-1]
	} else if l = len(u.Host); l > 0 && strings.HasSuffix(u.Host, "/") {
		u.Host = u.Host[:l-1]
	}
	return u, nil
}

func addTrailingSlash(u *url.URL) (*url.URL, error) {
	if l := len(u.Path); l > 0 && !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	} else if l = len(u.Host); l > 0 && !strings.HasSuffix(u.Host, "/") {
		u.Host += "/"
	}
	return u, nil
}
