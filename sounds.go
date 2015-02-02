package codex

// Defines known sounds.

// Glyphs used in English to directly represent some phonemes. This does not
// reflect all phonemes used in English (not even near...).
var knownSounds = Set.New(nil,
	// Digraphs
	"ae", "ch", "sh", "th", "zh",
	// English alphabet monographs
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
	"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	// Some other Latin monographs
	"æ",
)

// Vowel glyphs used in English.
var knownVowels = Set.New(nil,
	// Digraphs
	"ae",
	// English alphabet monographs
	"a", "e", "i", "o", "u", "y",
	// Some other Latin monographs
	"æ",
)
