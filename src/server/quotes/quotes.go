package quotes

import (
	"encoding/json"
	"io"
	"math/rand"
)

type Quotes []Quote

// Decode JSON from reader and load quotes.
func (quotes *Quotes) Load(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(quotes)
}

// Pick random quote.
func (quotes *Quotes) Pick() Quote {
	i := rand.Intn(len(*quotes))
	return (*quotes)[i]
}
