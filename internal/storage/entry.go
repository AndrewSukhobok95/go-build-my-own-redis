package storage

type entryType int

const (
	stringType entryType = iota
	listType
)

type entry struct {
	typ  entryType
	data any
}

func newStringEntry(s string) *entry {
	return &entry{typ: stringType, data: s}
}

func newListEntry(l []string) *entry {
	return &entry{typ: listType, data: l}
}

func (e *entry) String() (string, error) {
	if e.typ != stringType {
		return "", ErrWrongType
	}
	return e.data.(string), nil
}

func (e *entry) List() ([]string, error) {
	if e.typ != listType {
		return []string{}, ErrWrongType
	}
	return e.data.([]string), nil
}

func (e *entry) PushLeft(values ...string) (int, error) {
	if e.typ != listType {
		return 0, ErrWrongType
	}
	e.data = append(values, e.data.([]string)...)
	return len(e.data.([]string)), nil
}

func (e *entry) PushRight(values ...string) (int, error) {
	if e.typ != listType {
		return 0, ErrWrongType
	}
	e.data = append(e.data.([]string), values...)
	return len(e.data.([]string)), nil
}

func (e *entry) PopLeft() (string, error) {
	if e.typ != listType {
		return "", ErrWrongType
	}

	l := e.data.([]string)
	if len(l) == 0 {
		return "", nil
	}

	popped := l[0]
	e.data = l[1:]
	return popped, nil
}

func (e *entry) PopRight() (string, error) {
	if e.typ != listType {
		return "", ErrWrongType
	}

	l := e.data.([]string)
	if len(l) == 0 {
		return "", nil
	}

	lstIndex := len(l) - 1
	popped := l[lstIndex]
	e.data = l[:lstIndex]
	return popped, nil
}

func (e *entry) LLen() (int, error) {
	if e.typ != listType {
		return 0, ErrWrongType
	}
	return len(e.data.([]string)), nil
}

func adjustNegIndex(idx, length int) int {
	idxAdj := idx + length
	if idxAdj < 0 {
		return 0
	}
	return idxAdj
}

func (e *entry) LRange(start, stop int) ([]string, error) {
	if e.typ != listType {
		return []string{}, ErrWrongType
	}

	l := e.data.([]string)
	if len(l) == 0 {
		return []string{}, nil
	}

	startAdj := start
	if start < 0 {
		startAdj = adjustNegIndex(start, len(l))
	}

	stopAdj := stop
	if stop < 0 {
		stopAdj = adjustNegIndex(stop, len(l))
	}

	if startAdj > stopAdj || startAdj >= len(l) {
		return []string{}, nil
	}

	if stopAdj >= len(l) {
		stopAdj = len(l) - 1
	}
	return l[startAdj : stopAdj+1], nil
}
