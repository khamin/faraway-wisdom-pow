package quotes

import (
	"bytes"
	"fmt"
)

type Quote struct {
	Author Bytes `json:"author"`
	Text   Bytes `json:"text"`
}

// Format quote.
func (quote *Quote) Format() []byte {
	buf := bytes.Buffer{}

	fmt.Fprintf(&buf, "%s - %s", quote.Text, quote.Author)

	return buf.Bytes()
}
