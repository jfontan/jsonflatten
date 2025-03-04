package jsonflatten

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"pitr.ca/jsontokenizer"
)

// ParserPitr implements a json value flattener using Pitr tokenizer.
type ParserPitr struct {
	States
	emitter Emitter
}

const (
	readSize = 4 * 1024
)

// NewParserPitr creates a new parser using Pitr tokenizer. If emmiter is nil
// a default printer is used.
func NewParserPitr(emitter Emitter) *ParserPitr {
	return &ParserPitr{
		emitter: emitter,
	}
}

// Parse json and call the provided emitter for each value.
func (p *ParserPitr) Parse(r io.Reader) error {
	if p.emitter == nil {
		p.emitter = p.print
	}

	buf := new(strings.Builder)

	dec := jsontokenizer.NewWithSize(r, readSize)

	for {
		token, err := dec.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		switch token {
		case jsontokenizer.TokObjectOpen:
			p.pushState(TypeObject)

		case jsontokenizer.TokObjectClose:
			s := p.popState()
			if s.jsonType != TypeObject {
				return fmt.Errorf("invalid char %d", token)
			}

			p.lastState().advance()

		case jsontokenizer.TokArrayOpen:
			p.pushState(TypeArray)

		case jsontokenizer.TokArrayClose:
			s := p.popState()
			if s.jsonType != TypeArray {
				return fmt.Errorf("invalid char %d", token)
			}

			p.lastState().advance()

		case jsontokenizer.TokString:
			s := p.lastState()
			if s == nil {
				return fmt.Errorf("single strings not supported")
			}

			buf.Reset()
			_, err := dec.ReadString(buf)
			if err != nil {
				return err
			}

			v := buf.String()

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

		case jsontokenizer.TokNumber:
			s := p.lastState()
			if s == nil {
				return fmt.Errorf("single value not supported")
			}

			buf.Reset()
			_, err := dec.ReadNumber(buf)
			if err != nil {
				return err
			}

			v, err := strconv.ParseFloat(buf.String(), 64)
			if err != nil {
				return err
			}

			p.emit(s.key, v)
			s.advance()

		case jsontokenizer.TokTrue:
			if err := p.commonEmitter(true); err != nil {
				return err
			}

		case jsontokenizer.TokFalse:
			if err := p.commonEmitter(false); err != nil {
				return err
			}

		case jsontokenizer.TokNull:
			if err := p.commonEmitter(nil); err != nil {
				return err
			}

		case jsontokenizer.TokComma, jsontokenizer.TokObjectColon:

		default:
			return fmt.Errorf("invalid type: %d", token)
		}
	}
}

func (p *ParserPitr) commonEmitter(v any) error {
	s := p.lastState()
	if s == nil {
		return fmt.Errorf("single value not supported")
	}

	p.emit(s.key, v)
	s.advance()

	return nil
}

func (p *ParserPitr) emit(k string, v any) {
	var path path
	s := p.lastState()
	if s != nil {
		path = s.path
	}

	p.emitter(path.StringWithKey(k), v)
}

func (p *ParserPitr) print(k string, v any) {
	var value string
	switch nv := v.(type) {
	case string:
		value = fmt.Sprintf(`"%s"`, nv)
	default:
		value = fmt.Sprintf("%v", nv)
	}

	fmt.Println(k, "=", value)
}
