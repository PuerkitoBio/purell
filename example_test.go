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
