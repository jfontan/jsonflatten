package jsonflatten

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"

	"pitr.ca/jsontokenizer"
)

type ParserPitr struct {
	States
	emitter Emitter
}

const bufferSize = 64 * 1024

func (p *ParserPitr) Parse(r io.Reader) error {
	if p.emitter == nil {
		p.emitter = p.print
	}

	backBuffer := make([]byte, bufferSize)
	buf := bytes.NewBuffer(backBuffer)

	dec := jsontokenizer.New(r)

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

			// default:
			// 	return fmt.Errorf("invalid delimiter %s", string(v))

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
