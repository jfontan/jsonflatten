package jsonflatten

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// Emitter is a function that is called for each value.
type Emitter func(string, any)

// Parser implements a json value flattener using standard library tokenizer.
type Parser struct {
	commonParser
}

// NewParser creates a new parser using standard tokenizer. If emitter is
// nil a default printer is used.
func NewParser(emitter Emitter) *Parser {
	return &Parser{
		commonParser: newCommonParser(emitter),
	}
}

// Parse json and call the provided emitter for each value.
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
