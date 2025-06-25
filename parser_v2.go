package jsonflatten

import (
	"errors"
	"fmt"
	"io"

	"github.com/go-json-experiment/json/jsontext"
)

// ParserV2 implements a json value flattener using standard library tokenizer.
type ParserV2 struct {
	commonParser
}

// NewParserV2 creates a new parser using standard tokenizer. If emitter is
// nil a default printer is used.
func NewParserV2(emitter Emitter) *ParserV2 {
	return &ParserV2{
		commonParser: newCommonParser(emitter),
	}
}

// Parse json and call the provided emitter for each value.
func (p *ParserV2) Parse(r io.Reader) error {
	// dec := json.NewDecoder(r)
	dec := jsontext.NewDecoder(r)

	for {
		token, err := dec.ReadToken()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		switch token.Kind() {
		case '{':
			p.pushState(TypeObject)

		case '}':
			s := p.popState()
			if s.jsonType != TypeObject {
				return fmt.Errorf("invalid char %+v", token)
			}

			p.lastState().advance()

		case '[':
			p.pushState(TypeArray)

		case ']':
			s := p.popState()
			if s.jsonType != TypeArray {
				return fmt.Errorf("invalid char %+v", token)
			}

			p.lastState().advance()

		case '"':
			s := p.lastState()
			if s == nil {
				return fmt.Errorf("single strings not supported")
			}

			switch s.jsonType {
			case TypeObject:
				if s.key == "" {
					s.key = token.String()
				} else {
					if !p.emit(s.key, token.String()) {
						return nil
					}
					s.key = ""
				}

			case TypeArray:
				if !p.emit(s.key, token.String()) {
					return nil
				}
				s.advance()

			default:
				return fmt.Errorf("invalid type %v", s.jsonType)
			}

		case '0':
			if err := p.commonEmitter(token.Float()); err != nil {
				if errors.Is(err, errExit) {
					return nil
				}
				return err
			}

		case 't':
			if err := p.commonEmitter(true); err != nil {
				if errors.Is(err, errExit) {
					return nil
				}
				return err
			}

		case 'f':
			if err := p.commonEmitter(false); err != nil {
				if errors.Is(err, errExit) {
					return nil
				}
				return err
			}

		case 'n':
			if err := p.commonEmitter(nil); err != nil {
				if errors.Is(err, errExit) {
					return nil
				}
				return err
			}

		default:
			return fmt.Errorf("invalid type: %+v", token)
		}
	}
}
