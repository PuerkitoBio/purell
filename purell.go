/*
Package purell offers URL normalization as described on the wikipedia page:
http://en.wikipedia.org/wiki/URL_normalization
*/
package purell

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"sort"
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
	FlagRemoveEmptyQuerySeparator

	// Usually safe normalizations
	FlagRemoveTrailingSlash // Should choose one or the other (in add-remove slash)
	FlagAddTrailingSlash
	FlagRemoveDotSegments

	// Unsafe normalizations
	FlagRemoveDirectoryIndex
	FlagRemoveFragment
	FlagForceHttp
	FlagRemoveDuplicateSlashes
	FlagRemoveWww // Should choose one or the other (in add-remove www)
	FlagAddWww
	FlagSortQuery

	FlagsSafe NormalizationFlags = FlagLowercaseHost | FlagLowercaseScheme | FlagUppercaseEscapes | FlagDecodeUnnecessaryEscapes | FlagRemoveDefaultPort | FlagRemoveEmptyQuerySeparator

	FlagsUsuallySafe NormalizationFlags = FlagsSafe | FlagRemoveTrailingSlash | FlagRemoveDotSegments

	FlagsUnsafe NormalizationFlags = FlagsUsuallySafe | FlagRemoveDirectoryIndex | FlagRemoveFragment | FlagForceHttp | FlagRemoveDuplicateSlashes | FlagRemoveWww | FlagSortQuery
)

var rxPort = regexp.MustCompile(`(:\d+)/?$`)
var rxDirIndex = regexp.MustCompile(`(^|/)((?:default|index)\.\w{1,4})$`)
var rxDupSlashes = regexp.MustCompile(`/{2,}`)

// MustNormalizeUrlString() returns the normalized string, and panics if an error occurs.
// It takes an URL string as input, as well as the normalization flags.
func MustNormalizeUrlString(u string, f NormalizationFlags) string {
	if parsed, e := url.Parse(u); e != nil {
		panic(e)
	} else {
		return MustNormalizeUrl(parsed, f)
	}
	panic("Unreachable code.")
}

// MustNormalizeUrl() returns the normalized string, and panics if an error occurs.
// It takes a parsed URL object as input, as well as the normalization flags.
func MustNormalizeUrl(u *url.URL, f NormalizationFlags) string {
	if res, e := NormalizeUrl(u, f); e != nil {
		panic(e)
	} else {
		return res
	}
	panic("Unreachable code.")
}

// NormalizeUrlString() returns the normalized string, or an error.
// It takes an URL string as input, as well as the normalization flags.
func NormalizeUrlString(u string, f NormalizationFlags) (string, error) {
	if parsed, e := url.Parse(u); e != nil {
		return "", e
	} else {
		return NormalizeUrl(parsed, f)
	}
	panic("Unreachable code.")
}

// NormalizeUrl() returns the normalized string, or an error.
// It takes a parsed URL object as input, as well as the normalization flags.
func NormalizeUrl(u *url.URL, f NormalizationFlags) (string, error) {
	var normalized *url.URL = u
	var e error

	// FlagDecodeUnnecessaryEscapes has no action, since it is done automatically
	// by parsing the string as an URL. Same for FlagUppercaseEscapes.
	flags := map[NormalizationFlags]func(*url.URL) (*url.URL, error){
		FlagLowercaseScheme:        lowercaseScheme,
		FlagLowercaseHost:          lowercaseHost,
		FlagRemoveDefaultPort:      removeDefaultPort,
		FlagRemoveTrailingSlash:    removeTrailingSlash,
		FlagRemoveDirectoryIndex:   removeDirectoryIndex, // Must be before add trailing slash
		FlagAddTrailingSlash:       addTrailingSlash,
		FlagRemoveDotSegments:      removeDotSegments,
		FlagRemoveFragment:         removeFragment,
		FlagForceHttp:              forceHttp,
		FlagRemoveDuplicateSlashes: removeDuplicateSlashes,
		FlagRemoveWww:              removeWww,
		FlagAddWww:                 addWww,
		FlagSortQuery:              sortQuery,
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
	if len(u.Host) > 0 {
		u.Host = rxPort.ReplaceAllStringFunc(u.Host, func(val string) string {
			if val == ":80" {
				return ""
			}
			return val
		})
	}
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

func removeDotSegments(u *url.URL) (*url.URL, error) {
	var dotFree []string

	if len(u.Path) > 0 {
		sections := strings.Split(u.Path, "/")
		for _, s := range sections {
			if s == ".." {
				if len(dotFree) > 0 {
					dotFree = dotFree[:len(dotFree)-1]
				}
			} else if s != "." {
				dotFree = append(dotFree, s)
			}
		}
		// Special case if host does not end with / and new path does not begin with /
		u.Path = strings.Join(dotFree, "/")
		if !strings.HasSuffix(u.Host, "/") && !strings.HasPrefix(u.Path, "/") {
			u.Path = "/" + u.Path
		}
	}
	return u, nil
}

func removeDirectoryIndex(u *url.URL) (*url.URL, error) {
	if len(u.Path) > 0 {
		u.Path = rxDirIndex.ReplaceAllString(u.Path, "$1")
	}
	return u, nil
}

func removeFragment(u *url.URL) (*url.URL, error) {
	u.Fragment = ""
	return u, nil
}

func forceHttp(u *url.URL) (*url.URL, error) {
	if strings.ToLower(u.Scheme) == "https" {
		u.Scheme = "http"
	}
	return u, nil
}

func removeDuplicateSlashes(u *url.URL) (*url.URL, error) {
	if len(u.Path) > 0 {
		u.Path = rxDupSlashes.ReplaceAllString(u.Path, "/")
	}
	return u, nil
}

func removeWww(u *url.URL) (*url.URL, error) {
	if len(u.Host) > 0 && strings.HasPrefix(strings.ToLower(u.Host), "www.") {
		u.Host = u.Host[4:]
	}
	return u, nil
}

func addWww(u *url.URL) (*url.URL, error) {
	if len(u.Host) > 0 && !strings.HasPrefix(strings.ToLower(u.Host), "www.") {
		u.Host = "www." + u.Host
	}
	return u, nil
}

func sortQuery(u *url.URL) (*url.URL, error) {
	q := u.Query()

	if len(q) > 0 {
		arKeys := make([]string, len(q))
		i := 0
		for k, _ := range q {
			arKeys[i] = k
			i++
		}
		sort.Strings(arKeys)
		buf := new(bytes.Buffer)
		for _, k := range arKeys {
			sort.Strings(q[k])
			for _, v := range q[k] {
				if buf.Len() > 0 {
					buf.WriteRune('&')
				}
				buf.WriteString(fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
			}
		}

		// Rebuild the raw query string
		u.RawQuery = buf.String()
	}
	return u, nil
}
