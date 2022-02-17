package quotes

import (
	"os"
)

// Parse quotes from given filename.
func New(filename string) (*Quotes, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer file.Close()
	quotes := &Quotes{}

	if err = quotes.Load(file); err != nil {
		return nil, err
	}

	return quotes, nil
}
