package jsonflatten

import (
	"errors"
	"fmt"
	"strings"
)

var errExit = errors.New("exit")

type commonParser struct {
	States

	emitter Emitter
}

func newCommonParser(emitter Emitter) commonParser {
	c := commonParser{
		emitter: emitter,
	}

	if emitter == nil {
		c.emitter = c.print
	}

	return c
}

func (p *commonParser) commonEmitter(v any) error {
	s := p.lastState()
	if s == nil {
		return fmt.Errorf("single value not supported")
	}

	ok := p.emit(s.key, v)
	if !ok {
		return errExit
	}

	s.advance()

	return nil
}

func (p *commonParser) emit(k string, v any) bool {
	var path path
	s := p.lastState()
	if s != nil {
		path = s.path
	}

	// spew.Print(path)
	return p.emitter(path.StringWithKey(k), v)
}

func (p *commonParser) print(k string, v any) bool {
	var value string
	switch nv := v.(type) {
	case string:
		value = fmt.Sprintf(`"%s"`, nv)
	default:
		value = fmt.Sprintf("%v", nv)
	}

	fmt.Println(k, "=", value)

	return true
}

type path []string

func (p path) StringWithKey(k string) string {
	return strings.Join(append(p, k), ".")
}

func (p path) String() string {
	return strings.Join(p, ".")
}
