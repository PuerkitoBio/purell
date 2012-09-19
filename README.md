# Purell

Purell is a tiny Go library to normalize URLs. It returns a pure URL. Pure-ll. Sanitizer and all. Yeah, I know...

Based on the [wikipedia paper][wiki] and the [RFC 3986 document][rfc].

## Install

`go get github.com/PuerkitoBio/purell`

## API

```go
import (
  "github.com/PuerkitoBio/purell"
)
// [...]
// Somewhere in a function
normalized, err := purell.NormalizeUrlString("hTTp://someWEBsite.com:80/Amazing%3a/url/",
  purell.FlagLowercaseScheme | purell.FlagLowercaseHost | FlagUppercaseEscapes)

// Or...
normalized := purell.MustNormalizeUrlString("hTTp://someWEBsite.com:80/Amazing%3a/url/",
  purell.FlagLowercaseScheme | purell.FlagLowercaseHost | FlagUppercaseEscapes)

// Or yet again...
u, err := url.Parse("http://someurl.com")
normalized, err := purell.NormalizeUrl(u, purell.FlagsSafe)

// And finally...
normalized := purell.MustNormalizeUrl(u, purell.FlagsSafe)

```

For convenience, the flags `FlagsSafe`, `FlagsUsuallySafe` and `FlagsUnsafe` are provided for the similarly grouped normalizations on [wikipedia's URL normalization page][wiki].

The [full godoc reference][godoc] is available on gopkgdoc.

Note that FlagDecodeUnnecessaryEscapes, FlagUppercaseEscapes and FlagRemoveEmptyQuerySeparator are always implicitly set, because internally, the URL string is parsed as an URL object, which automatically decodes unnecessary escapes and uppercases necessary ones, and removes empty query separators (an unnecessary `?` at the end of the url). So this operation cannot **not** be done. For this reason, FlagRemoveEmptyQuerySeparator has been included in the FlagsSafe convenience constant, instead of FlagsUnsafe, where Wikipedia puts it (strangely?).

The *replace IP with domain name* normalization (`http://208.77.188.166/ → http://www.example.com/`) is obviously not possible for a library without making some network requests. This is not implemented in purell.

The *remove unused query string parameters* and *remove default query parameters* are also not implemented, since this is a very case-specific normalization, and it is quite trivial to do with an URL object.

## TODOs

*    What if the source URL does not encode invalid characters? Parsing the string in a URL type automatically encodes some of them, though not all, it would seem. We'll see if it requires a normalization method.
*    Add a class/default instance to allow specifying custom directory index names?

## License

The [BSD 3-Clause license][bsd].

[bsd]: http://opensource.org/licenses/BSD-3-Clause
[wiki]: http://en.wikipedia.org/wiki/URL_normalization
[rfc]: http://tools.ietf.org/html/rfc3986#section-6
[godoc]: http://go.pkgdoc.org/github.com/puerkitobio/purell
