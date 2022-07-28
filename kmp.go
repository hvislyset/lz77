package lz77

// An implementation of the Knuth-Morris-Pratt string-searching algorithm.
func Search(needle []byte, haystack []byte) (int, int) {
	offset := -1
	length := -1

	if len(needle) == 0 || len(haystack) == 0 {
		return offset, length
	}

	haystackPointer := 0
	needlePointer := 0

	failureFunction := computeFailureFunction(needle)

	for haystackPointer < len(haystack) {
		if needle[needlePointer] == haystack[haystackPointer] {
			haystackPointer++
			needlePointer++
			if needlePointer == len(needle) {
				offset = haystackPointer - needlePointer
				length = len(needle)

				return offset, length
			}
		} else {
			needlePointer = failureFunction[needlePointer]

			if needlePointer < 0 {
				haystackPointer++
				needlePointer++
			}
		}
	}

	return offset, length
}

// Computes the "partial match" table.
func computeFailureFunction(needle []byte) []int {
	needleLength := len(needle)
	failureFunction := make([]int, needleLength+1)

	position := 1
	candidate := 0

	failureFunction[0] = -1

	for position < needleLength {
		if needle[position] == needle[candidate] {
			failureFunction[position] = failureFunction[candidate]
		} else {
			failureFunction[position] = candidate

			for candidate >= 0 && needle[position] != needle[candidate] {
				candidate = failureFunction[candidate]
			}
		}

		position++
		candidate++
	}

	failureFunction[position] = candidate

	return failureFunction
}
