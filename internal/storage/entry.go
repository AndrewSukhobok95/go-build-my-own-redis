package storage

type entryType int

const (
	stringType entryType = iota
	listType
	setType
	hashType
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

func newSetEntry(members ...string) *entry {
	s := make(map[string]struct{})
	for _, m := range members {
		s[m] = struct{}{}
	}
	return &entry{typ: setType, data: s}
}

func newHashEntry() *entry {
	m := make(map[string]string)
	return &entry{typ: hashType, data: m}
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

func (e *entry) SAdd(members ...string) (int, error) {
	if e.typ != setType {
		return 0, ErrWrongType
	}
	set := e.data.(map[string]struct{})
	cnt := 0
	for _, m := range members {
		if _, exists := set[m]; !exists {
			set[m] = struct{}{}
			cnt++
		}
	}
	return cnt, nil
}

func (e *entry) SMembers() ([]string, error) {
	if e.typ != setType {
		return []string{}, ErrWrongType
	}

	set := e.data.(map[string]struct{})
	members := make([]string, 0, len(set))
	for member := range set {
		members = append(members, member)
	}
	return members, nil
}

func (e *entry) SIsMember(member string) (bool, error) {
	if e.typ != setType {
		return false, ErrWrongType
	}

	set := e.data.(map[string]struct{})
	_, exists := set[member]
	if !exists {
		return false, nil
	}
	return true, nil
}

func (e *entry) SRem(members ...string) (int, error) {
	if e.typ != setType {
		return 0, ErrWrongType
	}
	cnt := 0
	set := e.data.(map[string]struct{})
	for _, m := range members {
		if _, exists := set[m]; exists {
			delete(set, m)
			cnt++
		}
	}
	return cnt, nil
}

func (e *entry) SLen() (int, error) {
	if e.typ != setType {
		return 0, ErrWrongType
	}
	return len(e.data.(map[string]struct{})), nil
}

func (e *entry) HSet(field, value string) (bool, error) {
	if e.typ != hashType {
		return false, ErrWrongType
	}
	m := e.data.(map[string]string)
	_, exists := m[field]
	m[field] = value
	return !exists, nil
}

func (e *entry) HGet(field string) (string, bool, error) {
	if e.typ != hashType {
		return "", false, ErrWrongType
	}
	m := e.data.(map[string]string)
	value, exists := m[field]
	return value, exists, nil
}

func (e *entry) HGetAll() (map[string]string, error) {
	if e.typ != hashType {
		return map[string]string{}, ErrWrongType
	}
	return e.data.(map[string]string), nil
}
