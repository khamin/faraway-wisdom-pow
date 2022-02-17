package quotes

import "errors"

// Byte slice "read-as-written" JSON unmarshal helper.
// Used to avoid unnecessary conversion to a string.
type Bytes []byte

// Implements json.Unmarshaler interface.
func (m *Bytes) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("Bytes: UnmarshalJSON on nil pointer")
	}

	*m = append((*m)[0:0], data[1:len(data)-1]...)
	return nil
}
