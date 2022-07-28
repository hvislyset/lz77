package lz77

import "testing"

func TestSubstringSearch(t *testing.T) {
	cases := []struct {
		desc     string
		needle   []byte
		haystack []byte
		offset   int
		length   int
	}{
		{"Empty needle", []byte(""), []byte("cadcbdbc"), -1, -1},
		{"Empty haystack", []byte("cadcbdbc"), []byte(""), -1, -1},
		{"Substring exists within target", []byte("dc"), []byte("cadcbdbc"), 2, 2},
		{"Substring does not exist within target", []byte("xyz"), []byte("cadcbdbc"), -1, -1},
	}

	for _, testCase := range cases {
		offset, length := Search(testCase.needle, testCase.haystack)

		if offset != testCase.offset || length != testCase.length {
			t.Fatalf("%s: expected: (%d, %d) found: (%d, %d)", testCase.desc, testCase.offset, testCase.length, offset, length)
		}
	}
}
