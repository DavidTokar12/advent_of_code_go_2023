package circular_buffer

type CircularTripleBuffer struct {
	data [3]*[]rune
	size int
}

func (ctb *CircularTripleBuffer) Append(value *[]rune) bool {
	var wasMooved bool = false

	if ctb.size >= 3 {
		copy(ctb.data[:ctb.size-1], ctb.data[1:ctb.size])
		ctb.size = 2
		wasMooved = true
	}

	ctb.data[ctb.size] = value
	ctb.size += 1

	return wasMooved
}

func NewCircularTripleBuffer() *CircularTripleBuffer {
	return &CircularTripleBuffer{
		data: [3]*[]rune{nil, nil, nil},
		size: 0,
	}
}

func (ctb *CircularTripleBuffer) GetAll() []*[]rune {
	var result []*[]rune

	for i := 0; i < ctb.size; i++ {
		if ctb.data[i] != nil {
			result = append(result, ctb.data[i])
		}
	}

	return result
}

func (ctb *CircularTripleBuffer) IsFull() bool {

	return ctb.size == 3
}

func (ctb *CircularTripleBuffer) GetSize() int {
	return ctb.size
}
