[![GoDoc](https://godoc.org/github.com/Mitranim/codex?status.svg)](https://godoc.org/github.com/Mitranim/codex)

## Description

Generator of random synthetic words or names. Takes sample words, analyses them,
and lazily produces a set of similar derived words. Works for
[any language](#traitsexaminestring-error).

Example program using `codex`:

```golang
package main

import (
  "fmt"
  "github.com/Mitranim/codex"
)

func main() {
  source := []string{"jasmine", "katie", "nariko", "karen"}

  traits, err := codex.NewTraits(source)
  if err != nil {
    panic(err)
  }
  gen := traits.Generator()

  // Print twelve random words.
  for i := 0; i < 12; i++ {
    fmt.Println(gen())
  }

  // Printed (your result will be different):
  //   jarik smiko ikatik arinat nasmin katie
  //   rikatin smikas minena ikatin jasmika rinaren

  // Find out how many words can be generated from this sample.
  gen = traits.Generator()
  i := 0
  for gen() != "" {i++}
  fmt.Println("total:", i)

  // Printed:
  //   total: 392
}
```

## Contents

* [Description](#description)
* [Contents](#contents)
* [Installation](#installation)
* [API Reference](#api-reference)
  * [type Traits](#type-traits)
    * [NewTraits()](#newtraitsstring-traits-error)
    * [Traits.Examine()](#traitsexaminestring-error)
    * [Traits.Generator()](#traitsgenerator-func-string)
* [ToDo / WIP](#todo--wip)

## Installation

In a shell:

```sh
go get github.com/Mitranim/codex
```

In your Go files:

```golang
import (
  "fmt"
  "github.com/Mitranim/codex"
)

func main() {
  traits, err := codex.NewTraits([]string{"sample", "pair"})
  if err != nil {
    panic(err)
  }

  gen := traits.Generator()
  for word := gen(); word != ""; word = gen() {
    fmt.Println(word)
  }
}
```

To test the package, `cd` into the package directory and run:

```sh
# Just tests
go test
```

To run benchmarks:

```sh
# With benchmarks
go test -bench .
```

## API Reference

The entry point for everything is a `Traits` object. It takes existing words as
input. Words must consist of known glyphs, as defined by the sound sets in
[`sounds.go`](sounds.go) or by custom sets assigned to a traits struct (see the
[reference](#type-traits)). If an invalid word is encountered, an error is
returned.

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

  // Optional custom set of known sounds.
  KnownSounds Set
  // Optional custom set of known vowels.
  KnownVowels Set
}
```

`Traits` represent rudimental characteristics of a word or group of words. A
traits object unequivocally defines a set of synthetic words that may be derived
from them. They're produced by a generator function made with
[`Traits.Generator()`](#traitsgenerator-func-string).

The optional fields `KnownSounds` and `KnownVowels` specify custom sets of
sounds and vowels. This lets you use `codex` for any character set, including
non-Latin alphabets. See
[`Traits.Examine()`](#traitsexaminestring-error).

#### `NewTraits([]string) (*Traits, error)`

Shortcut for creating a `Traits` object and calling its `Examine()` method.
These are equivalent:

```golang
traits, err := NewTraits([]string{"mountain", "waterfall", "grotto"})

traits := &Traits{}
err := traits.Examine([]string{"mountain", "waterfall", "grotto"})
```

Ignore this if you're using custom sound sets (e.g. non-Latin).

#### `Traits.Examine([]string) error`

Analyses the given words and merges their attributes into self.

```golang
traits := &Traits{}
err := traits.Examine([]string{"mountain", "waterfall", "grotto"})
```

By default, this uses sets of known sounds and vowels defined in
[`sounds.go`](sounds.go). This includes the 26 letters of the standard US
English alphabet and some common digraphs like `th`, which are treated as single
phonemes.

However, `codex` is language-independent. Assign custom `KnownSounds` and
`KnownVowels` to teach it a sound system of your choosing. It can be Greek or
Cyrillic or Elvish or Clingon — doesn't matter as long as the given sounds and
vowels cover the words in your input. Refer to [`sounds.go`](sounds.go) as an
example.

Here's how to teach it Greek:

```golang
traits := &codex.Traits{
  KnownSounds: codex.Set.New(nil,
    "α", "β", "γ", "δ", "ε", "ζ", "η", "θ", "ι", "κ", "λ", "μ",
    "ν", "ξ", "ο", "π", "ρ", "σ", "ς", "τ", "υ", "φ", "χ", "ψ", "ω"),
  KnownVowels: codex.Set.New(nil, "α", "ε", "η", "ι", "ο", "υ", "ω"),
}

traits.Examine([]string{"ελ", "διδασκω", "ελληνικο", "αλφαβητο"})

gen := traits.Generator()
for word := gen(); word != ""; word = gen() {
  fmt.Println(word)
}

// "ιδαλφ"
// "κο"
// "ηνικο"
// ...
```

#### `Traits.Generator() func() string`

Creates a generator function that yields a new random synthetic word on each
call. The words are guaranteed to never repeat, and to be randomly distributed
across the total set of possible words for these traits.

After a generator is exhausted, subsequent calls return `""`.

A traits object is stateless, and `Generator()` produces a completely new
generator on each call. Generators don't affect each other.

This remains fast even for large source datasets, and is suitable for use on web
servers and in other applications where responses must be quick.

```golang
traits, err := codex.NewTraits([]string{"goblin", "smoke"})
gen := traits.Generator()

for word := gen(); word != ""; word = gen() {
  fmt.Print(word, " ")
}

// moblin oblin mobli goblin smobli gobli smoke
// this generator is exhausted
```

## ToDo / WIP

### Investigation

Consider providing an option to enable reverse pairs in `Traits.Examine()`.
Check the performance impact, particularly with large datasets.

### Algorithms

Perhaps Traits.validPart() should also forbid repeated triples.

### Tests

Random distribution test for the generators should verify that preceding calls
may return words that contain (starting at index 0) words returned from later
calls.

### Readme

* Include examples of modifying Traits fields to restrict word characteristics.
* Document what kind of input data is allowed.
