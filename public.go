/*

Generator of random synthetic words or names. Takes a sample provided by the
user, analyses it, and produces a set of similar derived words.

See the readme and a better organised API reference on github:
https://github.com/Mitranim/codex.

Example program:

    package main

    import (
      "fmt"
      "github.com/Mitranim/codex"
    )

    func main() {
      source := []string{"jasmine", "katie", "nariko"}

      sample, err := codex.WordsN(source, 12)
      total, err := codex.Words(source)

      fmt.Println(sample)
      fmt.Println("total:", len(total))
      fmt.Println(err)
    }

    // Printed results:
    //   {
    //     "inari", "tikarik", "karinat", "ariko", "minatik", "ikasmin",
    //     "kasmine", "katiko", "rikasmi", "mikatin", "natie", "natika",
    //   }
    //   total: 180
    //   <nil>

*/
package codex

// Exposes type-independent public functions and hosts an introductory overview.

/********************************* Functions *********************************/

// Takes a sample group of words, analyses their traits, and builds a set of all
// synthetic words that may be derived from those traits. This should only be
// used for very small samples. More than just a handful of sample words causes
// a combinatorial explosion, takes a lot of time to calculate, and produces too
// many results to be useful. The number of results can easily reach hundreds of
// thousands for just a dozen of sample words.
func Words(words []string) (Set, error) {
	traits, err := NewTraits(words)
	if err != nil {
		return nil, err
	}
	return traits.Words(), nil
}

// Takes a sample group of words and a count limiter. Analyses the words and
// builds a random sample of synthetic words that may be derived from those
// traits, limited to the given count.
func WordsN(words []string, num int) (Set, error) {
	state, err := NewState(words)
	if err != nil {
		return nil, err
	}
	return state.WordsN(num), nil
}
