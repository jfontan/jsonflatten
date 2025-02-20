package jsonflatten

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Emitter func(string, any)

type Parser struct {
	States
	emitter Emitter
}

func (p *Parser) Parse(r io.Reader) error {
	if p.emitter == nil {
		p.emitter = p.print
	}

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
			if err := p.commonEmitter(v); err != nil {
				return err
			}

		case bool:
			if err := p.commonEmitter(v); err != nil {
				return err
			}

		case nil:
			if err := p.commonEmitter(v); err != nil {
				return err
			}

		default:
			return fmt.Errorf("invalid type: %+v", v)
		}
	}
}

func (p *Parser) commonEmitter(v any) error {
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

	p.emitter(path.StringWithKey(k), v)
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
