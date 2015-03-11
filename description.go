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

*/
package codex
