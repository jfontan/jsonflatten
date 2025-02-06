package jsonflatten

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type State int
type States []State

func (s *States) Last() State {
	if len(*s) == 0 {
		return StateUnknown
	}

	return (*s)[len(*s)-1]
}

func (s *States) Push(st State) {
	(*s) = append(*s, st)
}

func (s *States) Pop() {
	if len(*s) == 0 {
		return
	}

	*s = (*s)[:len(*s)-1]
}

const (
	StateUnknown = iota
	StateObject
	StateArray
)

type Parser struct {
	path         path
	states       *States
	lastKey      string
	arrayCounter int
}

func (p *Parser) Parse(r io.Reader) error {
	dec := json.NewDecoder(r)

	if p.states == nil {
		p.states = new(States)
	}

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
				p.states.Push(StateObject)
				if p.lastKey != "" {
					p.path = append(p.path, p.lastKey)
					p.lastKey = ""
				}
			case '}':
				if p.states.Last() != StateObject {
					return fmt.Errorf("invalid char %s", string(v))
				}
				p.states.Pop()
				if len(p.path) > 0 {
					p.path = p.path[:len(p.path)-1]
				}
			case '[':
				p.states.Push(StateArray)
				p.arrayCounter = 0
			case ']':
				if p.states.Last() != StateArray {
					return fmt.Errorf("invalid char %s", string(v))
				}
				p.states.Pop()
				p.arrayCounter = 0
			default:
				return fmt.Errorf("invalid delimiter %s", string(v))
			}

		case string:
			if p.states.Last() == StateArray {
				p.path = append(p.path, strconv.Itoa(p.arrayCounter))
				println(p.path.String(), "=", v)
				p.path = p.path[:len(p.path)-1]
				p.arrayCounter++
				break
			}

			if p.lastKey == "" {
				p.lastKey = v
			} else {
				println(p.path.StringWithKey(p.lastKey), "=", v)
				p.lastKey = ""
			}

		default:
			println("invalid type")
		}
	}
}

type path []string

func (p path) StringWithKey(k string) string {
	return strings.Join(append(p, k), ".")
}

func (p path) String() string {
	return strings.Join(p, ".")
}
