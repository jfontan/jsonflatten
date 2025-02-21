package jsonflatten

import "strconv"

type Type int
type Types []Type

func (s *Types) Last() Type {
	if len(*s) == 0 {
		return TypeUnknown
	}

	return (*s)[len(*s)-1]
}

func (s *Types) Push(st Type) {
	(*s) = append(*s, st)
}

func (s *Types) Pop() {
	if len(*s) == 0 {
		return
	}

	*s = (*s)[:len(*s)-1]
}

const (
	TypeUnknown = iota
	TypeObject
	TypeArray
)

type State struct {
	path         path
	jsonType     Type
	key          string
	arrayCounter int
}

func NewState(t Type, p path) State {
	key := ""
	if t == TypeArray {
		key = "0"
	}

	if len(p) == 0 {
		p = make(path, 0, 64)
	}
	return State{
		jsonType: t,
		path:     p,
		key:      key,
	}
}

func (s *State) advance() {
	if s == nil {
		return
	}

	switch s.jsonType {
	case TypeObject:
		s.key = ""
	case TypeArray:
		s.arrayCounter++
		s.key = strconv.Itoa(s.arrayCounter)
	}
}

type States []State

func (p *States) pushState(t Type) {
	var path path
	var key string

	s := p.lastState()
	if s != nil {
		path = s.path
		key = s.key
	}

	if key != "" {
		path = append(path, key)
	}

	*p = append(*p, NewState(t, path))
}

func (p *States) popState() State {
	if len(*p) == 0 {
		return State{}
	}

	l := len(*p) - 1
	s := (*p)[l]
	*p = (*p)[:l]

	return s
}

func (s *States) lastState() *State {
	if len(*s) == 0 {
		return &State{}
	}

	l := len(*s) - 1
	return &(*s)[l]
}
