package bitvector

type IntVector struct {
	Data           []uint64
	BitsPerElement uint64
	Count          uint64
}

func NewIntVector(count, bitsPerElement uint64) *IntVector {
	if bitsPerElement == 0 || bitsPerElement > 64 {
		panic("bitsPerElement must be between 1 and 64")
	}
	totalBits := count * bitsPerElement
	numWords := (totalBits + 63) / 64
	return &IntVector{
		Data:           make([]uint64, numWords),
		BitsPerElement: bitsPerElement,
		Count:          count,
	}
}

// Get returns the value of the i-th element.
func (iv *IntVector) Get(i uint64) uint64 {
	bitPos := i * iv.BitsPerElement
	wordIdx := bitPos / 64
	bitOffset := bitPos % 64

	var mask uint64
	if iv.BitsPerElement == 64 {
		mask = ^uint64(0)
	} else {
		mask = (1 << iv.BitsPerElement) - 1
	}

	val := (iv.Data[wordIdx] >> bitOffset) & mask

	// Handle straddling when the element crosses a word boundary.
	bitsFromFirstWord := 64 - bitOffset
	if bitsFromFirstWord < iv.BitsPerElement {
		bitsNeeded := iv.BitsPerElement - bitsFromFirstWord
		upperMask := (uint64(1) << bitsNeeded) - 1
		upperBits := iv.Data[wordIdx+1] & upperMask
		val |= upperBits << bitsFromFirstWord
	}

	return val
}

// Set writes `val` into the i-th element slot.
func (iv *IntVector) Set(i uint64, val uint64) {
	bitPos := i * iv.BitsPerElement
	wordIdx := bitPos / 64
	bitOffset := bitPos % 64

	var mask uint64
	if iv.BitsPerElement == 64 {
		mask = ^uint64(0)
	} else {
		mask = (1 << iv.BitsPerElement) - 1
	}

	val &= mask

	iv.Data[wordIdx] &= ^(mask << bitOffset)
	iv.Data[wordIdx] |= val << bitOffset

	// Handle straddling
	bitsFromFirstWord := 64 - bitOffset
	if bitsFromFirstWord < iv.BitsPerElement {
		bitsNeeded := iv.BitsPerElement - bitsFromFirstWord
		upperMask := (uint64(1) << bitsNeeded) - 1

		iv.Data[wordIdx+1] &= ^upperMask
		iv.Data[wordIdx+1] |= (val >> bitsFromFirstWord) & upperMask
	}
}