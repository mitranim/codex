## Random ramblings

Imagine three ways to define a word programmatically.

    // Plain written word.
    type Word []Letter

    // Plain written word split into glyphs that correspond to phonemes, but not
    // necessarily the sequence of phonemes the speaker of the language would use.
    // Consequently, it may be a false, misleading sequence.
    type Word []Glyph

    // Spoken word represented with a sequence of phonemes. May not reflect how
    // the word would be written by a speaker of the language.
    type Word []Sound

We might define words like so:
* a spoken word is a sequence of phonemes spoken without pause;
* a written word is a sequence of letters describing a spoken word.

Commonly, the sequence of letters representing a word does not reflect its
sequence of phonemes; a speaker reading the word creates the phoneme sequence by
using implicit rules inherent in the practiced language; the rules involve
characteristics of individual letters and of their combinations.

Syllables may refer to illusory written syllables or real spoken syllables.

    type Word struct {
      Letters  string
      Glyphs   []string
      Phonemes []string
    }

When creating a word, we may ignore sounds and only use letters and fake
syllables, or fake glyphs and (still) fake syllables. Alternatively, when
creating a word, we may operate spoken phonemes and syllables, and convert them
into the written form when done.

If we could phonetise any word, we could go straight from letters to a sequence
of phonemes and operate exclusively that.

To demonstrate why it's difficult to define syllables programmatically, let's
see what happens if we try to break a word at vowel-consonant or consonant-
consonant boundaries.

    ru|di|men|tal                     |
    word                              |
    a|na|ly|sis                       | ə|nə|ləj|sis
    de|fi|nes          // inaccurate  | de|fajns
    tra|its            // inaccurate  | trəjts
    that                              | thət
    may                               | məj
    cha|rac|te|ri|se   // inaccurate  | chə|rək|te|rajs
    a                                 | ə
    and                               | ənd
    pro|vi|des         // inaccurate  | pro|vajds
    u|ti|li|ti|es      // inaccurate  | u|ti|li|tis
    for                               |
    cha|rac|te|ri|sing                | chə|rək|te|raj|sing
    e|xis|ting                        | ə|ksi|stin
    da|tes             // inaccurate  | dəjts
    so|uls                            |
    be|li|e|ve   // inaccurate twice  | be|līv
