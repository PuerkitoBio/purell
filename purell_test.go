package purell

import (
	"fmt"
	"net/url"
	"testing"
)

func assertResult(ex string, s string, t *testing.T) {
	if ex != s {
		t.Errorf("Expected %s, got %s.", ex, s)
	}
}

func TestLowerScheme(t *testing.T) {
	if s, e := NormalizeURLString("HTTP://www.SRC.ca", FlagLowercaseScheme); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("http://www.SRC.ca", s, t)
	}
}

func TestLowerScheme2(t *testing.T) {
	if s, e := NormalizeURLString("http://www.SRC.ca", FlagLowercaseScheme); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("http://www.SRC.ca", s, t)
	}
}

func TestLowerHost(t *testing.T) {
	if s, e := NormalizeURLString("HTTP://www.SRC.ca/", FlagLowercaseHost); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://www.src.ca/", s, t)
	}
}

func TestUpperEscapes(t *testing.T) {
	if s, e := NormalizeURLString(`http://www.whatever.com/Some%aa%20Special%8Ecases/`, FlagUppercaseEscapes); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("http://www.whatever.com/Some%AA%20Special%8Ecases/", s, t)
	}
}

func TestUnnecessaryEscapes(t *testing.T) {
	if s, e := NormalizeURLString(`http://www.toto.com/%41%42%2E%44/%32%33%52%2D/%5f%7E`, FlagDecodeUnnecessaryEscapes); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("http://www.toto.com/AB.D/23R-/_~", s, t)
	}
}

func TestRemoveDefaultPort(t *testing.T) {
	if s, e := NormalizeURLString("HTTP://www.SRC.ca:80/", FlagRemoveDefaultPort); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://www.SRC.ca/", s, t)
	}
}

func TestRemoveDefaultPort2(t *testing.T) {
	if s, e := NormalizeURLString("HTTP://www.SRC.ca:80", FlagRemoveDefaultPort); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://www.SRC.ca", s, t)
	}
}

func TestRemoveDefaultPort3(t *testing.T) {
	if s, e := NormalizeURLString("HTTP://www.SRC.ca:8080", FlagRemoveDefaultPort); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://www.SRC.ca:8080", s, t)
	}
}

func TestSafe(t *testing.T) {
	if s, e := NormalizeURLString("HTTP://www.SRC.ca:80/to%1ato%8b%ee/OKnow%41%42%43%7e", FlagsSafe); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("http://www.src.ca/to%1Ato%8B%EE/OKnowABC~", s, t)
	}
}

func TestBothLower(t *testing.T) {
	if s, e := NormalizeURLString("HTTP://www.SRC.ca:80/to%1ato%8b%ee/OKnow%41%42%43%7e", FlagLowercaseHost|FlagLowercaseScheme); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("http://www.src.ca:80/to%1Ato%8B%EE/OKnowABC~", s, t)
	}
}

func TestRemoveTrailingSlash(t *testing.T) {
	if s, e := NormalizeURLString("HTTP://www.SRC.ca:80/", FlagRemoveTrailingSlash); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://www.SRC.ca:80", s, t)
	}
}

func TestRemoveTrailingSlash2(t *testing.T) {
	if s, e := NormalizeURLString("HTTP://www.SRC.ca:80/toto/titi/", FlagRemoveTrailingSlash); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://www.SRC.ca:80/toto/titi", s, t)
	}
}

func TestRemoveTrailingSlash3(t *testing.T) {
	if s, e := NormalizeURLString("HTTP://www.SRC.ca:80/toto/titi/fin/?a=1", FlagRemoveTrailingSlash); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://www.SRC.ca:80/toto/titi/fin?a=1", s, t)
	}
}

func TestAddTrailingSlash(t *testing.T) {
	if s, e := NormalizeURLString("HTTP://www.SRC.ca:80", FlagAddTrailingSlash); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://www.SRC.ca:80/", s, t)
	}
}

func TestAddTrailingSlash2(t *testing.T) {
	if s, e := NormalizeURLString("HTTP://www.SRC.ca:80/toto/titi.html", FlagAddTrailingSlash); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://www.SRC.ca:80/toto/titi.html/", s, t)
	}
}

func TestAddTrailingSlash3(t *testing.T) {
	if s, e := NormalizeURLString("HTTP://www.SRC.ca:80/toto/titi/fin?a=1", FlagAddTrailingSlash); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://www.SRC.ca:80/toto/titi/fin/?a=1", s, t)
	}
}

func TestRemoveDotSegments(t *testing.T) {
	if s, e := NormalizeURLString("HTTP://root/a/b/./../../c/", FlagRemoveDotSegments); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://root/c/", s, t)
	}
}

func TestRemoveDotSegments2(t *testing.T) {
	if s, e := NormalizeURLString("HTTP://root/../a/b/./../c/../d", FlagRemoveDotSegments); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://root/a/d", s, t)
	}
}

func TestUsuallySafe(t *testing.T) {
	if s, e := NormalizeURLString("HTTP://www.SRC.ca:80/to%1ato%8b%ee/./c/d/../OKnow%41%42%43%7e/?a=b#test", FlagsUsuallySafeGreedy); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("http://www.src.ca/to%1Ato%8B%EE/c/OKnowABC~?a=b#test", s, t)
	}
}

func TestRemoveDirectoryIndex(t *testing.T) {
	if s, e := NormalizeURLString("HTTP://root/a/b/c/default.aspx", FlagRemoveDirectoryIndex); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://root/a/b/c/", s, t)
	}
}

func TestRemoveDirectoryIndex2(t *testing.T) {
	if s, e := NormalizeURLString("HTTP://root/a/b/c/default#a=b", FlagRemoveDirectoryIndex); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://root/a/b/c/default#a=b", s, t)
	}
}

func TestRemoveFragment(t *testing.T) {
	if s, e := NormalizeURLString("HTTP://root/a/b/c/default#toto=tata", FlagRemoveFragment); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://root/a/b/c/default", s, t)
	}
}

func TestForceHTTP(t *testing.T) {
	if s, e := NormalizeURLString("https://root/a/b/c/default#toto=tata", FlagForceHTTP); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("http://root/a/b/c/default#toto=tata", s, t)
	}
}

func TestRemoveDuplicateSlashes(t *testing.T) {
	if s, e := NormalizeURLString("https://root/a//b///c////default#toto=tata", FlagRemoveDuplicateSlashes); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("https://root/a/b/c/default#toto=tata", s, t)
	}
}

func TestRemoveDuplicateSlashes2(t *testing.T) {
	if s, e := NormalizeURLString("https://root//a//b///c////default#toto=tata", FlagRemoveDuplicateSlashes); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("https://root/a/b/c/default#toto=tata", s, t)
	}
}

func TestRemoveWWW(t *testing.T) {
	if s, e := NormalizeURLString("https://www.root/a/b/c/", FlagRemoveWWW); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("https://root/a/b/c/", s, t)
	}
}

func TestRemoveWWW2(t *testing.T) {
	if s, e := NormalizeURLString("https://WwW.Root/a/b/c/", FlagRemoveWWW); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("https://Root/a/b/c/", s, t)
	}
}

func TestAddWWW(t *testing.T) {
	if s, e := NormalizeURLString("https://Root/a/b/c/", FlagAddWWW); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("https://www.Root/a/b/c/", s, t)
	}
}

func TestSortQuery(t *testing.T) {
	if s, e := NormalizeURLString("http://root/toto/?b=4&a=1&c=3&b=2&a=5", FlagSortQuery); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("http://root/toto/?a=1&a=5&b=2&b=4&c=3", s, t)
	}
}

func TestRemoveEmptyQuerySeparator(t *testing.T) {
	if s, e := NormalizeURLString("http://root/toto/?", FlagRemoveEmptyQuerySeparator); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("http://root/toto/", s, t)
	}
}

func TestUnsafe(t *testing.T) {
	if s, e := NormalizeURLString("HTTPS://www.RooT.com/toto/t%45%1f///a/./b/../c/?z=3&w=2&a=4&w=1#invalid", FlagsUnsafeGreedy); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("http://root.com/toto/tE%1F/a/c?a=4&w=1&w=2&z=3", s, t)
	}
}

func TestSafe2(t *testing.T) {
	if s, e := NormalizeURLString("HTTPS://www.RooT.com/toto/t%45%1f///a/./b/../c/?z=3&w=2&a=4&w=1#invalid", FlagsSafe); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("https://www.root.com/toto/tE%1F///a/./b/../c/?z=3&w=2&a=4&w=1#invalid", s, t)
	}
}

func TestUsuallySafe2(t *testing.T) {
	if s, e := NormalizeURLString("HTTPS://www.RooT.com/toto/t%45%1f///a/./b/../c/?z=3&w=2&a=4&w=1#invalid", FlagsUsuallySafeGreedy); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("https://www.root.com/toto/tE%1F///a/c?z=3&w=2&a=4&w=1#invalid", s, t)
	}
}

func TestSourceModified(t *testing.T) {
	u, _ := url.Parse("HTTPS://www.RooT.com/toto/t%45%1f///a/./b/../c/?z=3&w=2&a=4&w=1#invalid")
	NormalizeURL(u, FlagsUnsafeGreedy)
	if u.Host != "http" {
		t.Logf("Expected source URL to have host http, found %s.", u.Host)
	}
	assertResult("http://root.com/toto/tE%1F/a/c?a=4&w=1&w=2&z=3", u.String(), t)
}

func TestDecodeUnnecessaryEscapesAll(t *testing.T) {
	var url = "http://host/"

	for i := 0; i < 256; i++ {
		url += fmt.Sprintf("%%%02x", i)
	}
	t.Logf("Source URL=%s", url)
	if s, e := NormalizeURLString(url, FlagDecodeUnnecessaryEscapes); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("http://host/%00%01%02%03%04%05%06%07%08%09%0A%0B%0C%0D%0E%0F%10%11%12%13%14%15%16%17%18%19%1A%1B%1C%1D%1E%1F%20%21%22%23$%25&%27%28%29%2A+,-./0123456789:;%3C=%3E%3F@ABCDEFGHIJKLMNOPQRSTUVWXYZ%5B%5C%5D%5E_%60abcdefghijklmnopqrstuvwxyz%7B%7C%7D~%7F%80%81%82%83%84%85%86%87%88%89%8A%8B%8C%8D%8E%8F%90%91%92%93%94%95%96%97%98%99%9A%9B%9C%9D%9E%9F%A0%A1%A2%A3%A4%A5%A6%A7%A8%A9%AA%AB%AC%AD%AE%AF%B0%B1%B2%B3%B4%B5%B6%B7%B8%B9%BA%BB%BC%BD%BE%BF%C0%C1%C2%C3%C4%C5%C6%C7%C8%C9%CA%CB%CC%CD%CE%CF%D0%D1%D2%D3%D4%D5%D6%D7%D8%D9%DA%DB%DC%DD%DE%DF%E0%E1%E2%E3%E4%E5%E6%E7%E8%E9%EA%EB%EC%ED%EE%EF%F0%F1%F2%F3%F4%F5%F6%F7%F8%F9%FA%FB%FC%FD%FE%FF", s, t)
	}
}
