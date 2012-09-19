# Purell

Purell is a tiny Go library to normalize URLs. It returns a pure URL. Pure-ell. Sanitizer and all. Yeah, I know...

Based on the [wikipedia paper][wiki] and the [RFC 3986 document][rfc].

## Install

`go get github.com/PuerkitoBio/purell`

## Examples

From `example_test.go` (note that in your code, you would import "github.com/PuerkitoBio/purell", and would prefix references to its methods and constants with "purell."):

```go
package purell

import (
  "fmt"
  "net/url"
)

func ExampleNormalizeUrlString() {
  if normalized, err := NormalizeUrlString("hTTp://someWEBsite.com:80/Amazing%3f/url/",
    FlagLowercaseScheme|FlagLowercaseHost|FlagUppercaseEscapes); err != nil {
    panic(err)
  } else {
    fmt.Print(normalized)
  }
  // Output: http://somewebsite.com:80/Amazing%3F/url/
}

func ExampleMustNormalizeUrlString() {
  normalized := MustNormalizeUrlString("hTTpS://someWEBsite.com:80/Amazing%fa/url/",
    FlagsUnsafe)
  fmt.Print(normalized)

  // Output: http://somewebsite.com/Amazing%FA/url
}

func ExampleNormalizeUrl() {
  if u, err := url.Parse("Http://SomeUrl.com:8080/a/b/.././c///g?c=3&a=1&b=9&c=0#target"); err != nil {
    panic(err)
  } else {
    if normalized, err := NormalizeUrl(u, FlagsUsuallySafe|FlagRemoveDuplicateSlashes|FlagRemoveFragment); err != nil {
      panic(err)
    } else {
      fmt.Print(normalized)
    }
  }

  // Output: http://someurl.com:8080/a/c/g?c=3&a=1&b=9&c=0
}

func ExampleMustNormalizeUrl() {
  if u, err := url.Parse("Http://SomeUrl.com:8080/a/b/.././c///g?c=3&a=1&b=9&c=0#target"); err != nil {
    panic(err)
  } else {
    normalized := MustNormalizeUrl(u, FlagsUnsafe&^FlagRemoveDotSegments)
    fmt.Print(normalized)
  }

  // Output: http://someurl.com:8080/a/b/.././c/g?a=1&b=9&c=0&c=3
}
```

## API

For convenience, the flags `FlagsSafe`, `FlagsUsuallySafe` and `FlagsUnsafe` are provided for the similarly grouped normalizations on [wikipedia's URL normalization page][wiki]. You can add (using the bitwise OR `|` operator) or remove (using the bitwise AND NOT `&^` operator) individual flags from the sets if required.

The [full godoc reference][godoc] is available on gopkgdoc.

Note that FlagDecodeUnnecessaryEscapes, FlagUppercaseEscapes and FlagRemoveEmptyQuerySeparator are always implicitly set, because internally, the URL string is parsed as an URL object, which automatically decodes unnecessary escapes, uppercases necessary ones, and removes empty query separators (an unnecessary `?` at the end of the url). So this operation cannot **not** be done. For this reason, FlagRemoveEmptyQuerySeparator has been included in the FlagsSafe convenience constant, instead of FlagsUnsafe, where Wikipedia puts it (strangely?).

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
