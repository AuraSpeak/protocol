package protocol

type Flag uint8

// Has returns true if the given flag bit is set in f.
func (f Flag) Has(flag Flag) bool {
	return f&flag != 0
}

func (f *Flag) Set(flag Flag) {
	*f |= flag
}

func (f *Flag) Clear(flag Flag) {
	*f &^= flag
}

func (f *Flag) Toggle(flag Flag) {
	*f ^= flag
}
