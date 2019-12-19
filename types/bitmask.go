package types

// Our set of bytes we are working with
type Bitmask uint32

// Checks if a specific Flag is set in the bytes
func (f Bitmask) HasFlag(flag Bitmask) bool {
	return f&flag != 0
}

// Sets a Flag to 1
func (f *Bitmask) AddFlag(flag Bitmask) {
	*f = (*f) | flag
}

// Sets a Flag to 0
func (f *Bitmask) RemoveFlag(flag Bitmask) {
	*f = (*f) &^ flag
}

// Sets flag when it is not set, unset when it is set
func (f *Bitmask) ToggleFlag(flag Bitmask) {
	*f ^= flag
}
