package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Prefix is a Markov chain prefix of one or more words.
type Prefix []string

// String returns the Prefix as a string (for use as a map key).
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word.
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a string of prefixLen words joined with spaces.
// A suffix is a single word. A prefix can have multiple suffixes.
type Chain struct {
	chain     map[string][]string
	prefixLen int
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string][]string), prefixLen}
}

// Build reads text from the provided Reader and
// parses it into prefixes and suffixes that are stored in Chain.
func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.prefixLen)
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}
		key := p.String()
		c.chain[key] = append(c.chain[key], s)
		p.Shift(s)
	}
}

// Generate returns a string of at most n words generated from Chain.
func (c *Chain) Generate(n int) string {
	p := make(Prefix, c.prefixLen)
	var words []string
	for i := 0; i < n; i++ {
		choices := c.chain[p.String()]
		if len(choices) == 0 {
			break
		}
		next := choices[rand.Intn(len(choices))]
		words = append(words, next)
		p.Shift(next)
	}
	return strings.Join(words, " ")
}

func main() {
	// Register command-line flags.
	numWords := flag.Int("words", 35, "maximum number of words to print")
	prefixLen := flag.Int("prefix", 2, "prefix length in words")

	flag.Parse()                     // Parse command-line flags.
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator.

	c := NewChain(*prefixLen) // Initialize a new Chain.
	scum, _ := os.Open("texts/scum.txt")
	c.Build(scum) // Build chains from standard input.

	files := []string{"fraga-ribbing-2015-08-14", "fraga-ribbing-2016-05-27", "fraga-ribbing-2015-05-22", "fraga-ribbing-2016-01-08", "fraga-ribbing-2016-03-24", "fraga-ribbing-2016-01-29", "fraga-ribbing-2015-03-06", "fraga-ribbing-2015-01-09", "fraga-ribbing-2014-12-12", "fraga-ribbing-2015-06-05", "fraga-ribbing-2015-01-16", "fraga-ribbing-2014-10-24", "fraga-ribbing-2014-06-13", "fraga-ribbing-2015-04-10", "fraga-ribbing-2015-01-30", "fraga-ribbing-2015-03-13", "fraga-ribbing-2014-08-01", "fraga-ribbing-2014-09-12", "fraga-ribbing-2015-01-23", "fraga-ribbing-2014-09-26", "fraga-ribbing-2014-10-17", "fraga-ribbing-2014-07-04"}

	for _, item := range files {
		filename := strings.Join([]string{"texts/", item}, "")
		file, _ := os.Open(string(filename))
		c.Build(file)
	}

	text := c.Generate(*numWords) // Generate text.
	strippedText := strings.Split(text, ".")
	strippedText = strippedText[:len(strippedText)-1]

	text = strings.Join(strippedText, ".")

	fmt.Print(text) // Write text to standard output.
	fmt.Print(".")
	fmt.Println("")
}
