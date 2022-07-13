package oktools

import "sort"

const (
	offset64 = 14695981039346656037
	prime64  = 1099511628211

	// SeparatorByte is a byte that cannot occur in valid UTF-8 sequences and is
	// used to separate label names, label values, and other strings from each other
	// when calculating their combined hash value (aka signature aka fingerprint).
	SeparatorByte byte = 255
)

// hashNew initializies a new fnv64a hash value.
func hashNew() uint64 {
	return offset64
}

// hashAdd adds a string to a fnv64a hash value, returning the updated hash.
func hashAdd(h uint64, s string) uint64 {
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= prime64
	}
	return h
}

// hashAddByte adds a byte to a fnv64a hash value, returning the updated hash.
func hashAddByte(h uint64, b byte) uint64 {
	h ^= uint64(b)
	h *= prime64
	return h
}

func GenerateHashForMapValues(labels map[string]string) uint64 {
	if len(labels) == 0 {
		emptyLabelSignature := hashNew()
		return emptyLabelSignature
	}

	labelNames := make([]string, 0, len(labels))
	for labelName := range labels {
		labelNames = append(labelNames, labelName)
	}
	sort.Strings(labelNames)

	sum := hashNew()
	for _, labelName := range labelNames {
		sum = hashAdd(sum, labelName)
		sum = hashAddByte(sum, SeparatorByte)
		sum = hashAdd(sum, labels[labelName])
		sum = hashAddByte(sum, SeparatorByte)
	}
	return sum
}
