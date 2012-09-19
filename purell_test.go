package purell

import (
	"testing"
)

func assertResult(ex string, s string, t *testing.T) {
	if ex != s {
		t.Errorf("Expected %s, got %s.", ex, s)
	}
}

func TestLowerScheme(t *testing.T) {
	if s, e := NormalizeUrlString("HTTP://www.SRC.ca", FlagLowercaseScheme); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("http://www.SRC.ca", s, t)
	}
}

func TestLowerScheme2(t *testing.T) {
	if s, e := NormalizeUrlString("http://www.SRC.ca", FlagLowercaseScheme); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("http://www.SRC.ca", s, t)
	}
}

func TestLowerHost(t *testing.T) {
	if s, e := NormalizeUrlString("HTTP://www.SRC.ca/", FlagLowercaseHost); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://www.src.ca/", s, t)
	}
}

func TestUpperEscapes(t *testing.T) {
	if s, e := NormalizeUrlString(`http://www.whatever.com/Some%aa%20Special%8Ecases/`, FlagUppercaseEscapes); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("http://www.whatever.com/Some%AA%20Special%8Ecases/", s, t)
	}
}

func TestUnnecessaryEscapes(t *testing.T) {
	if s, e := NormalizeUrlString(`http://www.toto.com/%41%42%2E%44/%32%33%52%2D/%5f%7E`, FlagDecodeUnnecessaryEscapes); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("http://www.toto.com/AB.D/23R-/_~", s, t)
	}
}

func TestRemoveDefaultPort(t *testing.T) {
	if s, e := NormalizeUrlString("HTTP://www.SRC.ca:80/", FlagRemoveDefaultPort); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://www.SRC.ca/", s, t)
	}
}

func TestRemoveDefaultPort2(t *testing.T) {
	if s, e := NormalizeUrlString("HTTP://www.SRC.ca:80", FlagRemoveDefaultPort); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://www.SRC.ca", s, t)
	}
}

func TestRemoveDefaultPort3(t *testing.T) {
	if s, e := NormalizeUrlString("HTTP://www.SRC.ca:8080", FlagRemoveDefaultPort); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://www.SRC.ca:8080", s, t)
	}
}

func TestSafe(t *testing.T) {
	if s, e := NormalizeUrlString("HTTP://www.SRC.ca:80/to%1ato%8b%ee/OKnow%41%42%43%7e", FlagsSafe); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("http://www.src.ca/to%1Ato%8B%EE/OKnowABC~", s, t)
	}
}

func TestBothLower(t *testing.T) {
	if s, e := NormalizeUrlString("HTTP://www.SRC.ca:80/to%1ato%8b%ee/OKnow%41%42%43%7e", FlagLowercaseHost|FlagLowercaseScheme); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("http://www.src.ca:80/to%1Ato%8B%EE/OKnowABC~", s, t)
	}
}

func TestRemoveTrailingSlash(t *testing.T) {
	if s, e := NormalizeUrlString("HTTP://www.SRC.ca:80/", FlagRemoveTrailingSlash); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://www.SRC.ca:80", s, t)
	}
}

func TestRemoveTrailingSlash2(t *testing.T) {
	if s, e := NormalizeUrlString("HTTP://www.SRC.ca:80/toto/titi/", FlagRemoveTrailingSlash); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://www.SRC.ca:80/toto/titi", s, t)
	}
}

func TestRemoveTrailingSlash3(t *testing.T) {
	if s, e := NormalizeUrlString("HTTP://www.SRC.ca:80/toto/titi/fin/?a=1", FlagRemoveTrailingSlash); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://www.SRC.ca:80/toto/titi/fin?a=1", s, t)
	}
}

func TestAddTrailingSlash(t *testing.T) {
	if s, e := NormalizeUrlString("HTTP://www.SRC.ca:80", FlagAddTrailingSlash); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://www.SRC.ca:80/", s, t)
	}
}

func TestAddTrailingSlash2(t *testing.T) {
	if s, e := NormalizeUrlString("HTTP://www.SRC.ca:80/toto/titi.html", FlagAddTrailingSlash); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://www.SRC.ca:80/toto/titi.html/", s, t)
	}
}

func TestAddTrailingSlash3(t *testing.T) {
	if s, e := NormalizeUrlString("HTTP://www.SRC.ca:80/toto/titi/fin?a=1", FlagAddTrailingSlash); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://www.SRC.ca:80/toto/titi/fin/?a=1", s, t)
	}
}

func TestRemoveDotSegments(t *testing.T) {
	if s, e := NormalizeUrlString("HTTP://root/a/b/./../../c/", FlagRemoveDotSegments); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://root/c/", s, t)
	}
}

func TestRemoveDotSegments2(t *testing.T) {
	if s, e := NormalizeUrlString("HTTP://root/../a/b/./../c/../d", FlagRemoveDotSegments); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://root/a/d", s, t)
	}
}

func TestUsuallySafe(t *testing.T) {
	if s, e := NormalizeUrlString("HTTP://www.SRC.ca:80/to%1ato%8b%ee/./c/d/../OKnow%41%42%43%7e/?a=b#test", FlagsUsuallySafe); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("http://www.src.ca/to%1Ato%8B%EE/c/OKnowABC~?a=b#test", s, t)
	}
}

func TestRemoveDirectoryIndex(t *testing.T) {
	if s, e := NormalizeUrlString("HTTP://root/a/b/c/default.aspx", FlagRemoveDirectoryIndex); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://root/a/b/c/", s, t)
	}
}

func TestRemoveDirectoryIndex2(t *testing.T) {
	if s, e := NormalizeUrlString("HTTP://root/a/b/c/default#a=b", FlagRemoveDirectoryIndex); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://root/a/b/c/default#a=b", s, t)
	}
}

func TestRemoveFragment(t *testing.T) {
	if s, e := NormalizeUrlString("HTTP://root/a/b/c/default#toto=tata", FlagRemoveFragment); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("HTTP://root/a/b/c/default", s, t)
	}
}

func TestForceHttp(t *testing.T) {
	if s, e := NormalizeUrlString("https://root/a/b/c/default#toto=tata", FlagForceHttp); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("http://root/a/b/c/default#toto=tata", s, t)
	}
}

func TestRemoveDuplicateSlashes(t *testing.T) {
	if s, e := NormalizeUrlString("https://root/a//b///c////default#toto=tata", FlagRemoveDuplicateSlashes); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("https://root/a/b/c/default#toto=tata", s, t)
	}
}

func TestRemoveDuplicateSlashes2(t *testing.T) {
	if s, e := NormalizeUrlString("https://root//a//b///c////default#toto=tata", FlagRemoveDuplicateSlashes); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("https://root/a/b/c/default#toto=tata", s, t)
	}
}

func TestRemoveWww(t *testing.T) {
	if s, e := NormalizeUrlString("https://www.root/a/b/c/", FlagRemoveWww); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("https://root/a/b/c/", s, t)
	}
}

func TestRemoveWww2(t *testing.T) {
	if s, e := NormalizeUrlString("https://WwW.Root/a/b/c/", FlagRemoveWww); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("https://Root/a/b/c/", s, t)
	}
}

func TestAddWww(t *testing.T) {
	if s, e := NormalizeUrlString("https://Root/a/b/c/", FlagAddWww); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("https://www.Root/a/b/c/", s, t)
	}
}

func TestSortQuery(t *testing.T) {
	if s, e := NormalizeUrlString("http://root/toto/?b=4&a=1&c=3&b=2&a=5", FlagSortQuery); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("http://root/toto/?a=1&a=5&b=2&b=4&c=3", s, t)
	}
}

func TestRemoveEmptyQuerySeparator(t *testing.T) {
	if s, e := NormalizeUrlString("http://root/toto/?", FlagRemoveEmptyQuerySeparator); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("http://root/toto/", s, t)
	}
}

func TestUnsafe(t *testing.T) {
	if s, e := NormalizeUrlString("HTTPS://RooT/toto/t%45%1f///a/./b/../c/?z=3&w=2&a=4&w=1#invalid", FlagsUnsafe); e != nil {
		t.Errorf("Got error %s", e.Error())
	} else {
		assertResult("http://root/toto/tE%1F/a/c?a=4&w=1&w=2&z=3", s, t)
	}
}
