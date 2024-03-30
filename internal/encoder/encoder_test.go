package encoder

import (
	"testing"
)

func checkResult(t *testing.T, got string, expected string) {
	if got != expected {
		t.Errorf("expected %s but got %s instead", expected, got)
	}
}

// func TestUrlEncoder(t *testing.T) {
// 	prefix := "https://shortly.io/"

// 	testUrl1 := "https://www.google.com"
// 	testUrl2 := "http://www.testurl.com"
// 	testUrl3 := "http://wassup.com"

// 	expected1 := prefix + strconv.Itoa(0)
// 	expected2 := prefix + strconv.Itoa(1)

// 	alias := "abc"
// 	expected3 := prefix + alias

// 	enc := MakeSimple("https://shortly.io/", 0)

// 	res1 := enc.Encode(testUrl1, "")
// 	checkResult(t, res1, expected1)

// 	res2 := enc.Encode(testUrl2, "")
// 	checkResult(t, res2, expected2)

// 	res3 := enc.Encode(testUrl3, "abc")
// 	checkResult(t, res3, expected3)
// }
