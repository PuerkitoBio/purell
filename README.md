# Purell

Purell is a tiny Go library to normalize URLs. It returns a pure URL. Pure-ll. Sanitizer and all. Yeah, I know...

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

## License

The [BSD 3-Clause license][bsd].

[bsd]: http://opensource.org/licenses/BSD-3-Clause
