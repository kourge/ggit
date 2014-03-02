package core

import (
	"testing"
)

var (
	_fixtureSha Sha1 = Sha1([20]byte{
		0xbd, 0x9d, 0xbf, 0x5a, 0xae, 0x1a, 0x38, 0x62, 0xdd, 0x15,
		0x26, 0x72, 0x32, 0x46, 0xb2, 0x02, 0x06, 0xe5, 0xfc, 0x37,
	})

	_fixtureShaEmpty Sha1 = Sha1([20]byte{})
)

func TestSha1_String(t *testing.T) {
	var actual string = _fixtureSha.String()
	var expected string = "bd9dbf5aae1a3862dd1526723246b20206e5fc37"

	if actual != expected {
		t.Errorf("sha.String() = %v, want %v", actual, expected)
	}
}

func TestSha1_Split(t *testing.T) {
	var actual1, actual2 string = _fixtureSha.Split(2)
	var expected1, expected2 string = "bd", "9dbf5aae1a3862dd1526723246b20206e5fc37"

	if actual1 != expected1 {
		t.Errorf("sha.Split(2)[0] = %v, want %v", actual1, expected1)
	}

	if actual2 != expected2 {
		t.Errorf("sha.Split(2)[1] = %v, want %v", actual2, expected2)
	}
}

func TestSha1_IsEmpty(t *testing.T) {
	var actual1, actual2 bool = _fixtureSha.IsEmpty(), _fixtureShaEmpty.IsEmpty()
	var expected1, expected2 bool = false, true

	if actual1 != expected1 {
		t.Errorf("sha.IsEmpty() = %v, want %v", actual1, expected1)
	}

	if actual2 != expected2 {
		t.Errorf("emptySha.IsEmpty() = %v, want %v", actual2, expected2)
	}
}

func TestSha1FromString(t *testing.T) {
	s := "bd9dbf5aae1a3862dd1526723246b20206e5fc37"
	actual, err := Sha1FromString(s)
	expected := _fixtureSha

	if actual != expected {
		t.Errorf("Sha1FromString(\"%s\") = (%v, %v), want %v", s, actual, err, expected, nil)
	}
}

func TestSha1FromString_Invalid(t *testing.T) {
	s := "foobar"
	_, err := Sha1FromString(s)

	if err == nil {
		t.Errorf("Sha1FromString(\"%s\")[1] != nil", s)
	}
}
