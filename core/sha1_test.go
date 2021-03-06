package core

import (
	"testing"
)

var (
	_fixtureSha Sha1 = Sha1{
		0xbd, 0x9d, 0xbf, 0x5a, 0xae, 0x1a, 0x38, 0x62, 0xdd, 0x15,
		0x26, 0x72, 0x32, 0x46, 0xb2, 0x02, 0x06, 0xe5, 0xfc, 0x37,
	}

	_fixtureShaEmpty Sha1 = Sha1{}
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

func TestSha1_Compare_Equal(t *testing.T) {
	if _fixtureSha.Compare(_fixtureSha) != 0 {
		t.Errorf("%#v does not equal itself", _fixtureSha)
	}
}

func TestSha1_Compare_GreaterThan(t *testing.T) {
	zeroSha1 := Sha1{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	if _fixtureSha.Compare(zeroSha1) <= 0 {
		t.Errorf("(%#v).Compare(%#v) = %d, expected > 0", _fixtureSha, zeroSha1)
	}
}

func TestSha1_Compare_LesserThan(t *testing.T) {
	fullSha1 := Sha1{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	}

	if _fixtureSha.Compare(fullSha1) >= 0 {
		t.Errorf("(%#v).Compare(%#v) = %d, expected < 0", _fixtureSha, fullSha1)
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

func TestSha1FromByteSlice_Exact(t *testing.T) {
	v := []byte{
		0xbd, 0x9d, 0xbf, 0x5a, 0xae, 0x1a, 0x38, 0x62, 0xdd, 0x15,
		0x26, 0x72, 0x32, 0x46, 0xb2, 0x02, 0x06, 0xe5, 0xfc, 0x37,
	}
	var actual Sha1 = Sha1FromByteSlice(v)
	var expected Sha1 = Sha1{
		0xbd, 0x9d, 0xbf, 0x5a, 0xae, 0x1a, 0x38, 0x62, 0xdd, 0x15,
		0x26, 0x72, 0x32, 0x46, 0xb2, 0x02, 0x06, 0xe5, 0xfc, 0x37,
	}

	if actual != expected {
		t.Errorf("Sha1FromByteSlice(%#v) = %#v != %#v", v, actual, expected)
	}
}

func TestSha1FromByteSlice_Under(t *testing.T) {
	v := []byte{
		0xbd, 0x9d, 0xbf, 0x5a, 0xae, 0x1a, 0x38, 0x62, 0xdd, 0x15,
	}
	var actual Sha1 = Sha1FromByteSlice(v)
	var expected Sha1 = Sha1{
		0xbd, 0x9d, 0xbf, 0x5a, 0xae, 0x1a, 0x38, 0x62, 0xdd, 0x15,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	if actual != expected {
		t.Errorf("Sha1FromByteSlice(%#v) = %#v != %#v", v, actual, expected)
	}
}

func TestSha1FromByteSlice_Over(t *testing.T) {
	v := []byte{
		0xbd, 0x9d, 0xbf, 0x5a, 0xae, 0x1a, 0x38, 0x62, 0xdd, 0x15,
		0x26, 0x72, 0x32, 0x46, 0xb2, 0x02, 0x06, 0xe5, 0xfc, 0x37,
		0xde, 0xad, 0xbe, 0xef, 0x00, 0xca, 0xfe, 0xba, 0xbe, 0x00,
	}
	var actual Sha1 = Sha1FromByteSlice(v)
	var expected Sha1 = Sha1{
		0xbd, 0x9d, 0xbf, 0x5a, 0xae, 0x1a, 0x38, 0x62, 0xdd, 0x15,
		0x26, 0x72, 0x32, 0x46, 0xb2, 0x02, 0x06, 0xe5, 0xfc, 0x37,
	}

	if actual != expected {
		t.Errorf("Sha1FromByteSlice(%#v) = %#v != %#v", v, actual, expected)
	}
}
