package lz77

import (
	"bytes"
	"testing"
)

func TestLz77Compression(t *testing.T) {
	inputs := [][]byte{
		[]byte("cadcbdbc"),
		[]byte("dbbbbacb"),
		[]byte("addadcad"),
		[]byte("abaabdad"),
		[]byte("acabdaadbbacdccb"),
		[]byte("bcacbacacbddacaa"),
		[]byte("dbbbaccdaaabbcbb"),
		[]byte("daccddcccddcabda"),
		[]byte("dcadabcbaabbcbbbbdccdcbaadcccadc"),
		[]byte("bbbcacdcabdbaabbbdcbbdbdcddbbcbc"),
		[]byte("dcabcdcbdbaabbdbbdccaacadbaabbab"),
		[]byte("abccbacdcacbacdcbbaaddcabcbdcacc"),
	}

	for _, input := range inputs {
		compressed := Compress(input)
		decompressed := Decompress(compressed)

		if !bytes.Equal(decompressed, input) {
			t.Fatalf("expected: %s found: %s", input, decompressed)
		}
	}
}
