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
