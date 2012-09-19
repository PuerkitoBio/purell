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

func MustNormalizeUrlString(u string, f NormalizationFlags) string {
	if parsed, e := url.Parse(u); e != nil {
		panic(e.Error())
	} else {
		return MustNormalizeUrl(parsed, f)
	}
	panic("Unreachable code.")
}

func MustNormalizeUrl(u *url.URL, f NormalizationFlags) string {
	if res, e := NormalizeUrl(u, f); e != nil {
		panic(e.Error())
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

	flags := map[NormalizationFlags]func(*url.URL) (*url.URL, error){
		FlagLowercaseScheme:  lowercaseScheme,
		FlagLowercaseHost:    lowercaseHost,
		FlagUppercaseEscapes: uppercaseEscapes,
	}

	for k, v := range flags {
		if f|k == k {
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

func uppercaseEscapes(u *url.URL) (*url.URL, error) {
	rx := regexp.MustCompile(`%[0-9a-fA-F]{2}`)
	s := u.String()
	s = rx.ReplaceAllStringFunc(s, func(val string) string {
		return strings.ToUpper(val)
	})
	return url.Parse(s)
}
