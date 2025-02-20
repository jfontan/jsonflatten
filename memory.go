package jsonflatten

import (
	"encoding/json"
	"fmt"
	"io"
)

type Memory struct {
	States
	emitter Emitter
}

func (m *Memory) Parse(r io.Reader) error {
	dec := json.NewDecoder(r)

	if m.emitter == nil {
		m.emitter = m.print
	}

	var d any
	err := dec.Decode(&d)
	if err != nil {
		return err
	}

	switch v := d.(type) {
	case map[string]any:
		m.pushState(TypeObject)
		return m.parseMap(v)
	case []any:
		m.pushState(TypeArray)
		return m.parseArray(v)
	default:
		return fmt.Errorf("unknown type %+v", v)
	}
}

func (m *Memory) parseAny(a any) error {
	switch v := a.(type) {
	case map[string]any:
		return m.parseMap(v)
	case []any:
		return m.parseArray(v)
	case string, float64, bool, nil:
		return m.parseValue(v)

	default:
		return fmt.Errorf("invalid type: %+v", v)
	}
}

func (m *Memory) parseMap(v map[string]any) error {
	m.pushState(TypeObject)

	for k, v := range v {
		s := m.lastState()
		s.key = k
		err := m.parseAny(v)
		if err != nil {
			return err
		}
	}

	m.popState()

	return nil
}

func (m *Memory) parseArray(a []any) error {
	m.pushState(TypeArray)

	s := m.lastState()

	for _, v := range a {
		err := m.parseAny(v)
		if err != nil {
			return err
		}
		s.advance()
	}

	m.popState()

	return nil
}

func (m *Memory) parseValue(v any) error {
	s := m.lastState()
	if s == nil {
		return fmt.Errorf("single strings not supported")
	}

	m.emit(s.key, v)

	return nil
}

func (p *Memory) emit(k string, v any) {
	var path path
	s := p.lastState()
	if s != nil {
		path = s.path
	}

	p.emitter(path.StringWithKey(k), v)
}

func (p *Memory) print(k string, v any) {
	var value string
	switch nv := v.(type) {
	case string:
		value = fmt.Sprintf(`"%s"`, nv)
	default:
		value = fmt.Sprintf("%v", nv)
	}

	fmt.Println(k, "=", value)
}
