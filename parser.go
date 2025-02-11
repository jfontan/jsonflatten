package jsonflatten

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

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

type Parser struct {
	states States
}

func (p *Parser) Parse(r io.Reader) error {
	dec := json.NewDecoder(r)

	for {
		token, err := dec.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		switch v := token.(type) {
		case json.Delim:
			switch v {
			case '{':
				p.pushState(TypeObject)

			case '}':
				s := p.popState()
				if s.jsonType != TypeObject {
					return fmt.Errorf("invalid char %s", string(v))
				}

				p.lastState().advance()

			case '[':
				p.pushState(TypeArray)

			case ']':
				s := p.popState()
				if s.jsonType != TypeArray {
					return fmt.Errorf("invalid char %s", string(v))
				}

				p.lastState().advance()

			default:
				return fmt.Errorf("invalid delimiter %s", string(v))
			}

		case string:
			s := p.lastState()
			if s == nil {
				return fmt.Errorf("single strings not supported")
			}

			switch s.jsonType {
			case TypeObject:
				if s.key == "" {
					s.key = v
				} else {
					p.emit(s.key, v)
					s.key = ""
				}

			case TypeArray:
				p.emit(s.key, v)
				s.advance()

			default:
				return fmt.Errorf("invalid type %v", s.jsonType)
			}

		case float64:
			if err := p.commonEmiter(v); err != nil {
				return err
			}

		case bool:
			if err := p.commonEmiter(v); err != nil {
				return err
			}

		case nil:
			if err := p.commonEmiter(v); err != nil {
				return err
			}

		default:
			return fmt.Errorf("invalid type: %+v", v)
		}
	}
}

func (p *Parser) pushState(t Type) {
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

	p.states = append(p.states, NewState(t, path))
}

func (p *Parser) popState() State {
	if len(p.states) == 0 {
		return State{}
	}

	l := len(p.states) - 1
	s := p.states[l]
	p.states = p.states[:l]

	return s
}

func (p *Parser) lastState() *State {
	if len(p.states) == 0 {
		return &State{}
	}

	l := len(p.states) - 1
	return &p.states[l]
}

func (p *Parser) commonEmiter(v any) error {
	s := p.lastState()
	if s == nil {
		return fmt.Errorf("single value not supported")
	}

	p.emit(s.key, v)
	s.advance()

	return nil
}

func (p *Parser) emit(k string, v any) {
	var path path
	s := p.lastState()
	if s != nil {
		path = s.path
	}
	p.print(path.StringWithKey(k), v)
}

func (p *Parser) print(k string, v any) {
	var value string
	switch nv := v.(type) {
	case string:
		value = fmt.Sprintf(`"%s"`, nv)
	default:
		value = fmt.Sprintf("%v", nv)
	}

	fmt.Println(k, "=", value)
}

type path []string

func (p path) StringWithKey(k string) string {
	return strings.Join(append(p, k), ".")
}

func (p path) String() string {
	return strings.Join(p, ".")
}
