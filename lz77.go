package lz77

import (
	"bytes"
	"encoding/binary"
	"math"
)

const SearchBufferSize = 32767  // 15 bits
const LookaheadBufferSize = 255 // 8 bits

type slidingWindow struct {
	inputStreamLength    int
	cursor               int
	searchBufferIndex    uint
	lookaheadBufferIndex uint
	nextSymbol           byte
}

func newSlidingWindow(inputStreamLength int) *slidingWindow {
	return &slidingWindow{
		inputStreamLength:    inputStreamLength,
		cursor:               0,
		searchBufferIndex:    0,
		lookaheadBufferIndex: uint(math.Min(float64(LookaheadBufferSize), float64(inputStreamLength))),
		nextSymbol:           0,
	}
}

// Shifts the sliding window. Moves the cursor length + 1 and adjusts the search and lookahead buffers.
func (sw *slidingWindow) shift(length int) {
	sw.cursor += length + 1

	sw.searchBufferIndex = uint(math.Max(float64(sw.cursor-SearchBufferSize), 0))
	sw.lookaheadBufferIndex = uint(math.Min(float64(sw.cursor+LookaheadBufferSize), float64(sw.inputStreamLength)))
}

// Finds the longest repeated occurence of input that begins in the lookahead buffer.
func (sw *slidingWindow) match(inputStream []byte) (uint16, int, byte) {
	lookaheadCursor := sw.cursor
	searchBuffer := inputStream[sw.searchBufferIndex:sw.cursor]
	searchBufferLength := len(searchBuffer)

	maxOffset := 0
	maxLength := 0
	nextSymbol := inputStream[sw.cursor]

	for lookaheadCursor < int(sw.lookaheadBufferIndex) {
		candidate := inputStream[sw.cursor : lookaheadCursor+1]

		offset, length := Search(candidate, searchBuffer)

		if offset != -1 {
			lookaheadCursor++

			if length > maxLength {
				maxOffset = searchBufferLength - offset
				maxLength = length

				if sw.cursor+length >= sw.inputStreamLength {
					nextSymbol = 0
				} else {
					nextSymbol = inputStream[sw.cursor+length]
				}
			}
		} else {
			break
		}
	}

	return uint16(maxOffset), maxLength, byte(nextSymbol)
}

type token struct {
	offset     uint16
	length     int
	nextSymbol byte
}

func newToken(offset uint16, length int, nextSymbol byte) *token {
	return &token{
		offset,
		length,
		nextSymbol,
	}
}

// Encodes the token as a 4 byte chunk that will be emitted into the compressed output.
func (token *token) encode() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.BigEndian, token.offset)
	binary.Write(buf, binary.BigEndian, uint8(token.length))
	binary.Write(buf, binary.BigEndian, token.nextSymbol)

	return buf.Bytes()
}

// Compresses an input stream.
func Compress(inputStream []byte) []byte {
	inputStreamLength := len(inputStream)
	slidingWindow := newSlidingWindow(inputStreamLength)
	outputBuffer := make([]byte, 0)

	for slidingWindow.cursor < inputStreamLength {
		offset, length, nextSymbol := slidingWindow.match(inputStream)

		token := newToken(offset, length, nextSymbol)
		outputBuffer = append(outputBuffer, token.encode()...)

		slidingWindow.shift(length)
	}

	return outputBuffer
}

// Decompresses a compressed input stream.
func Decompress(inputStream []byte) []byte {
	inputStreamLength := len(inputStream)
	outputBuffer := make([]byte, 0)

	for chunkIndex := 0; chunkIndex < inputStreamLength; chunkIndex += 4 {
		chunk := inputStream[chunkIndex:(chunkIndex + 4)]

		offset := binary.BigEndian.Uint16(chunk[0:2])
		length := chunk[2]
		nextSymbol := chunk[3]

		if length == 0 {
			outputBuffer = append(outputBuffer, nextSymbol)
		} else {
			lower := int(math.Max(0, float64(len(outputBuffer)-int(offset))))
			upper := lower + int(length)
			outputBuffer = append(outputBuffer, outputBuffer[lower:upper]...)
			if nextSymbol != 0 {
				outputBuffer = append(outputBuffer, nextSymbol)
			}
		}
	}

	return outputBuffer
}
