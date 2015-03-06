[![GoDoc](https://godoc.org/github.com/Mitranim/codex?status.svg)](https://godoc.org/github.com/Mitranim/codex)

## Description

Generator of random synthetic words or names. Takes a sample provided by the
user, analyses it, and produces a set of similar derived words.

Example program using `codex`:

```golang
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
/*
  {
    "inari", "tikarik", "karinat", "ariko", "minatik", "ikasmin",
    "kasmine", "katiko", "rikasmi", "mikatin", "natie", "natika",
  }
  total: 180
  <nil>
*/
```

## Contents

* [Description](#description)
* [Contents](#contents)
* [Installation](#installation)
* [API Reference](#api-reference)
  * [Words()](#wordsstring-set-error)
  * [WordsN()](#wordsnstring-int-set-error)
  * [type Traits](#type-traits)
    * [NewTraits()](#newtraitsstring-traits-error)
    * [Traits.Words()](#traitswords-set-error)
    * [Traits.Examine()](#traitsexaminestring-error)
  * [type State](#type-state)
    * [NewState()](#newstatestring-state-error)
    * [State.Words()](#statewords-set)
    * [State.WordsN()](#statewordsnint-set)
  * [type Set](#type-set)
    * [Set.New()](#setnew_-string-set)
    * [Set.Has()](#sethasstring-bool)
    * [Set.Add()](#setaddstring)
    * [Set.Del()](#setdelstring)
* [ToDo / WIP](#todo--wip)

## Installation

In a shell:

```shell
go get github.com/Mitranim/codex
```

In your Go files:

```golang
import (
  "fmt"
  "github.com/Mitranim/codex"
)

func main() {
  words, err := codex.Words([]string{"sample", "pair"})
  fmt.Println(words, err)
}
```

To test the package, go into the package directory and run:

```shell
# Just tests
go test
# With benchmarks
go test -bench .
```

## API Reference

Most public functions exposed by the package take existing words as input. Words
must consist of known glyphs, as defined by the sound sets in
[`sounds.go`](sounds.go) or by custom sets passed into a Traits object (see the
[reference](#type-traits)). If an invalid word is encountered, an error is
returned.

### `Words([]string) (Set, error)`

Returns the entire set of synthetic words that may be derived from the given
sample. Beware: passing more than just a handful of words leads to a
combinatorial explosion and takes forever to calculate. This must only be used
with miniscule datasets.

This function is pure, meaning that repeated calls with the same dataset will
yield the same (unordered) result.

```golang
words, err := Words([]string{"goblin", "smoke"})
fmt.Println(words)
// {"smobli", "smoblin", "smoke", "goke", "gobli", "moke", "mobli", "goblin", "obli", "oblin", "oke", "moblin"}
```

See the [`Set`](#type-set) reference for how to handle the results.

### `WordsN([]string, int) (Set, error)`

Returns a random sample from the set of synthetic words that may be derived from
the given words, limited to the given count. The sequence is guaranteed to be
duplicate-free.

Unlike `Words()`, this remains very fast even for large source datasets, and is
suitable for use on a web server or another application where responses must be
quick.

```golang
words, err := WordsN([]string{"goblin", "smoke"}, 4)
fmt.Println(words)
// {"mobli", "smobli", "obli", "smoblin"}
```

See the [`Set`](#type-set) reference for how to handle the results.

### `type Traits`

```golang
type Traits struct {
  // Minimum and maximum number of sounds.
  MinNSounds int
  MaxNSounds int
  // Minimum and maximum number of vowels.
  MinNVowels int
  MaxNVowels int
  // Maximum number of consequtive vowels.
  MaxConseqVow int
  // Maximum number of consequtive consonants.
  MaxConseqCons int
  // Set of sounds that occur in the words.
  SoundSet Set
  // Set of pairs of sounds that occur in the words.
  PairSet PairSet

  // Replacement sound set to use instead of the default `knownSounds`.
  KnownSounds Set
  // Replacement sound set to use instead of the default `knownVowels`.
  KnownVowels Set
}
```

`Traits` represent rudimental characteristics of a word or group of words, and
are central to the package's functionality. Word generation always begins by
examining the source words and extracting their shared traits.

A traits object unequivocally defines a set of synthetic words that may be
derived from them. This set may be retrieved with `Traits.Words()`.

A traits object is stateless, and the `Traits.Words()` method is pure, meaning
that is has no side effects and is guaranteed to produce the same (unordered)
set on repeated calls for the same combination of traits. A transient traits
object is used internally by the static `Words()` function.

The fields `Traits.KnownSounds` and `Traits.KnownVowels` let you specify custom
sets of sounds and vowels to recognise in words. This lets you make the package
compatible with any character set, including non-Latin alphabets. See
`Traits.Examine()`.

#### `NewTraits([]string) (*Traits, error)`

Analyses the given group of sample words and returns a `Traits` object with
their shared characteristics. After getting hold of a traits object, you can
apply custom restrictions to its derived words by editing its fields. This is
also true for a traits object embedded in a `State`.

```golang
traits, err := NewTraits([]string{"mountain", "waterfall", "grotto"})
```

#### `Traits.Words() (Set, error)`

Generates and returns the entire set of synthetic words defined by the traits.
See the notes to the [`Words()`](#wordsstring-set-error) function. This method
is pure.

```golang
traits, err := NewTraits([]string{"goblin", "smoke"})
words := traits.Words()
fmt.Println(words)
// {"smobli", "smoblin", "smoke", "goke", "gobli", "moke", "mobli", "goblin", "obli", "oblin", "oke", "moblin"}
```

#### `Traits.Examine([]string) error`

Analyses the given group of sample words and merges their traits into the given
traits object. This is useful when you want to create a traits object with
custom `Traits.KnownSounds` and `Traits.KnownVowels` before analysing the source
words.

```golang
traits := &Traits{KnownSounds: Set.New(nil, "ε", "λ", "η", "ν", "ι", "κ", "ά")}
err := traits.Examine([]string{"ελ", "νικά"})
fmt.Println(err)
// <nil>
```

### `type State`

```golang
type State struct {
  Traits *Traits
  // unexported fields
}
```

A `State` object is a superset of `Traits` that maintains an internal state.
It's used for generating small samples from the set of synthetic words defined
by its traits through its `State.WordsN()` method. Statefulness allows it to
guarantee that no word is ever repeated. A state's generator methods share the
same virtual pool of words, and may be used interchangeably until the entire set
has been exhausted.

A state must always be obtained through a `NewState()` call, or given a valid
`Traits` object if created manually. Its behaviour without traits is undefined.
Its internal `Traits` object may be edited to apply custom restrictions on the
derived words.

A transient state object is used internally by the static `WordsN()` function.

Example of creating a new state with custom traits:

```golang
traits := &Traits{
  KnownSounds: /* <custom sound glyphs> */,
  KnownVowels: /* <custom vowel glyphs> */,
}
err := traits.Examine([]string{/* <sample words> */})
state := &State{Traits: traits}
```

#### `NewState([]string) (*State, error)`

Takes a group of sample words, generates their shared characteristics via
`NewTraits()`, and makes a `State` object that encapsulates those traits.

```golang
state, err := NewState([]string{"lava", "ridge", "rock"})
```

#### `State.Words() Set`

Generates and returns the remainder of the set of synthetic words defined by the
state's traits. Any words previously returned by the state's `State.WordsN()`
method are withheld. If called immediately after creating the state, the result
is guaranteed to be equivalent to `Words()` or `Traits.Words()` for the same
source data, with roughly 2/3d the performance, give or take. It's also equally
problematic for large datasets.

This method exhausts the remainder of the state's word set, and subsequent calls
to `State.Words()` or `State.WordsN()` return empty results.

```golang
state, err := NewState([]string{"goblin", "smoke"})
fmt.Println(state.Words())
// {"mobli", "smoke", "gobli", "smoblin", "goblin", "moblin", "moke", "obli", "oblin", "oke", "smobli", "goke"}
```

#### `State.WordsN(int) Set`

Generates and returns a random sample from the set of synthetic words defined by
the state's traits. Any words returned by this method are guaranteed to never
repeat in subsequent calls to the state's `State.Words()` and `State.WordsN()`
methods. If called enough times, this eventually exhausts the entire set of
words defined by the state's traits, and subsequent calls return empty results.

This method remains fast even for large source datasets, and is suitable for use
on web servers and in other applications where responses must be quick.

```golang
state, err := NewState([]string{"goblin", "smoke"})
fmt.Println(state.WordsN(7))
fmt.Println(state.WordsN(7))
fmt.Println(state.WordsN(7))

// {"smoblin", "mobli", "smobli", "smoke", "moke", "goke", "oblin"}
// {"moblin", "goblin", "obli", "oke", "gobli"}
// {}
```

### `type Set`

```golang
type Set map[string]struct{}
```

Represents a set of strings. Generated words are always returned as a `Set`.

Because it's a map, iterating over a `Set` is dead simple:

```golang
for word := range Set{} {
  // do stuff
}
```

A `Set` is unordered. In fact, Go actively randomises map iteration order, but
this is not always random enough. If you want to iterate over a set of words
randomly, make a `State` object and use its `State.WordsN()` method, which is
guaranteed to return random samples with no repeats.

#### `Set.New(_, ...string) Set`

Creates a Set with the given strings.

```golang
set := Set.New(nil, "one", "other" /*, ... */)
```

#### `Set.Has(string) bool`

Checks if the set has the given string.

```golang
set := Set.New(nil, "icecream")
set.Has("icecream") // true
```

#### `Set.Add(string)`

Adds the given string to the set.

```golang
set := Set{}
set.Add("sledges")
set.Has("sledges") // true
```

#### `Set.Del(string)`

Deletes the given string from the set.

```golang
set := Set.New(nil, "polaris")
set.Del("polaris")
set.Has("polaris") // false
```

## ToDo / WIP

### Investigation

Consider providing an option to enable reverse pairs for the `WordsN()` static
function. Enabling it for `Words()` or `Traits` or `State` objects (where
`Words()` could be called) is too hazardous, the combinatorial explosion goes
beyond any reasonable measure.

### Algorithms

Perhaps Traits.validPart() should also forbid repeated triples.

### Tests

Random distribution test for `State.WordsN()` should verify that preceding calls
may return words that contain (starting at index 0) words returned from later
calls.

### Readme

Include examples of:
  * using custom sets of known sounds and vowels, particularly non-Latin;
  * modifying Traits fields to restrict word characteristics.

Document what kind of input data is allowed.
