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
