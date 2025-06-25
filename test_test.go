package jsonflatten

import (
	"encoding/json"
	"io"
	"os"
	"slices"
	"strings"
	"testing"

	jsonv2 "github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/stretchr/testify/require"
)

const testJson = `
{
    "glossary": {
        "title": "example glossary",
		"GlossDiv": {
            "title": "S",
			"GlossList": {
                "GlossEntry": {
                    "ID": "SGML",
					"SortAs": "SGML",
					"GlossTerm": "Standard Generalized Markup Language",
					"Acronym": "SGML",
					"Abbrev": "ISO 8879:1986",
					"GlossDef": {
                        "para": "A meta-markup language, used to create markup languages such as DocBook.",
						"GlossSeeAlso": ["GML", "XML"]
                    },
					"GlossSee": "markup",
					"float64": 42,
					"bool": true,
					"null": null
                }
            }
        }
    },
    "array": [
    	{
     		"one": 1,
       		"two": 2
        },
       	{
      		"three": 1,
      		"four": 2,
           	"embedded": [1, 2, 3, true, null, "string"]
        }
    ]
}
`

var expected = map[string]any{
	"glossary.title":                                                 "example glossary",
	"glossary.GlossDiv.title":                                        "S",
	"glossary.GlossDiv.GlossList.GlossEntry.ID":                      "SGML",
	"glossary.GlossDiv.GlossList.GlossEntry.SortAs":                  "SGML",
	"glossary.GlossDiv.GlossList.GlossEntry.GlossTerm":               "Standard Generalized Markup Language",
	"glossary.GlossDiv.GlossList.GlossEntry.Acronym":                 "SGML",
	"glossary.GlossDiv.GlossList.GlossEntry.Abbrev":                  "ISO 8879:1986",
	"glossary.GlossDiv.GlossList.GlossEntry.GlossDef.para":           "A meta-markup language, used to create markup languages such as DocBook.",
	"glossary.GlossDiv.GlossList.GlossEntry.GlossDef.GlossSeeAlso.0": "GML",
	"glossary.GlossDiv.GlossList.GlossEntry.GlossDef.GlossSeeAlso.1": "XML",
	"glossary.GlossDiv.GlossList.GlossEntry.GlossSee":                "markup",
	"glossary.GlossDiv.GlossList.GlossEntry.float64":                 float64(42),
	"glossary.GlossDiv.GlossList.GlossEntry.bool":                    true,
	"glossary.GlossDiv.GlossList.GlossEntry.null":                    nil,
	"array.0.one":        float64(1),
	"array.0.two":        float64(2),
	"array.1.three":      float64(1),
	"array.1.four":       float64(2),
	"array.1.embedded.0": float64(1),
	"array.1.embedded.1": float64(2),
	"array.1.embedded.2": float64(3),
	"array.1.embedded.3": true,
	"array.1.embedded.4": nil,
	"array.1.embedded.5": "string",
}

func TestPrint(t *testing.T) {
	t.Skip()
	r := strings.NewReader(testJson)
	p := new(Parser)
	err := p.Parse(r)
	require.NoError(t, err)
}

func TestPrintV2(t *testing.T) {
	// t.Skip()
	r := strings.NewReader(testJson)
	p := new(ParserV2)
	p.emitter = p.print
	err := p.Parse(r)
	require.NoError(t, err)
}

func TestPrintPitr(t *testing.T) {
	t.Skip()
	r := strings.NewReader(testJson)
	p := new(ParserPitr)
	err := p.Parse(r)
	require.NoError(t, err)
}

func TestMap(t *testing.T) {
	r := strings.NewReader(testJson)
	m := make(map[string]any)

	p := new(Parser)
	p.emitter = func(k string, v any) bool {
		m[k] = v
		return true
	}

	err := p.Parse(r)
	require.NoError(t, err)

	require.Equal(t, expected, m)
}

func TestMapV2(t *testing.T) {
	r := strings.NewReader(testJson)
	m := make(map[string]any)

	p := new(ParserV2)
	p.emitter = func(k string, v any) bool {
		m[k] = v
		return true
	}

	err := p.Parse(r)
	require.NoError(t, err)

	require.Equal(t, expected, m)
}

func TestMapPitr(t *testing.T) {
	r := strings.NewReader(testJson)
	m := make(map[string]any)

	p := new(Parser)
	p.emitter = func(k string, v any) bool {
		m[k] = v
		return true
	}

	err := p.Parse(r)
	require.NoError(t, err)

	require.Equal(t, expected, m)
}

func TestMapPitrStop(t *testing.T) {
	r := strings.NewReader(testJson)
	m := make(map[string]any)

	keys := []string{
		"array.1.embedded.0",
		"glossary.GlossDiv.GlossList.GlossEntry.GlossDef.GlossSeeAlso.0",
	}

	p := NewParserPitr(func(k string, v any) bool {
		if !slices.Contains(keys, k) {
			return true
		}

		m[k] = v
		return len(m) < len(keys)
	})

	err := p.Parse(r)
	require.NoError(t, err)

	expected := map[string]any{
		"glossary.GlossDiv.GlossList.GlossEntry.GlossDef.GlossSeeAlso.0": "GML",
		"array.1.embedded.0": float64(1),
	}

	require.Equal(t, expected, m)
}

func TestMapMemory(t *testing.T) {
	r := strings.NewReader(testJson)
	m := make(map[string]any)

	p := new(Memory)
	p.emitter = func(k string, v any) bool {
		m[k] = v
		return true
	}

	err := p.Parse(r)
	require.NoError(t, err)

	require.Equal(t, expected, m)
}

func TestMapMemoryV2(t *testing.T) {
	r := strings.NewReader(testJson)
	m := make(map[string]any)

	p := new(MemoryV2)
	p.emitter = func(k string, v any) bool {
		m[k] = v
		return true
	}

	err := p.Parse(r)
	require.NoError(t, err)

	require.Equal(t, expected, m)
}

func TestLarge(t *testing.T) {
	t.Skip()
	f, err := os.Open("large-file.json")
	require.NoError(t, err)

	p := new(Parser)
	err = p.Parse(f)
	require.NoError(t, err)
}

func TestLargeV2(t *testing.T) {
	// t.Skip()
	f, err := os.Open("large-file.json")
	require.NoError(t, err)

	p := new(ParserV2)
	p.emitter = func(s string, a any) bool { return true }
	err = p.Parse(f)
	require.NoError(t, err)
}

func TestLargePitr(t *testing.T) {
	f, err := os.Open("large-file.json")
	require.NoError(t, err)

	p := new(ParserPitr)
	p.emitter = func(s string, a any) bool { return true }
	err = p.Parse(f)
	require.NoError(t, err)
}

func BenchmarkSmallParser(b *testing.B) {
	for range b.N {
		r := strings.NewReader(testJson)
		p := new(Parser)
		p.emitter = func(k string, v any) bool {
			return true
		}

		err := p.Parse(r)
		require.NoError(b, err)
	}
}

func BenchmarkSmallParserV2(b *testing.B) {
	for range b.N {
		r := strings.NewReader(testJson)
		p := new(ParserV2)
		p.emitter = func(k string, v any) bool {
			return true
		}

		err := p.Parse(r)
		require.NoError(b, err)
	}
}

func BenchmarkSmallParserPitr(b *testing.B) {
	for range b.N {
		r := strings.NewReader(testJson)
		p := new(ParserPitr)
		p.emitter = func(k string, v any) bool {
			return true
		}

		err := p.Parse(r)
		require.NoError(b, err)
	}
}

func BenchmarkSmallMemory(b *testing.B) {
	for range b.N {
		r := strings.NewReader(testJson)
		p := new(Memory)
		p.emitter = func(k string, v any) bool {
			return true
		}

		err := p.Parse(r)
		require.NoError(b, err)
	}
}

func BenchmarkSmallMemoryV2(b *testing.B) {
	for range b.N {
		r := strings.NewReader(testJson)
		p := new(MemoryV2)
		p.emitter = func(k string, v any) bool {
			return true
		}

		err := p.Parse(r)
		require.NoError(b, err)
	}
}

func BenchmarkUnmarshalSmall(b *testing.B) {
	for range b.N {
		var m any
		err := json.Unmarshal([]byte(testJson), &m)
		require.NoError(b, err)
	}
}

func BenchmarkUnmarshalSmallV2(b *testing.B) {
	for range b.N {
		var m any
		err := jsonv2.Unmarshal([]byte(testJson), &m)
		require.NoError(b, err)
	}
}

func BenchmarkBigParser(b *testing.B) {
	f, err := os.Open("large-file.json")
	require.NoError(b, err)
	defer f.Close()

	for range b.N {
		_, err := f.Seek(0, io.SeekStart)
		require.NoError(b, err)

		p := new(Parser)
		p.emitter = func(k string, v any) bool {
			return true
		}

		err = p.Parse(f)
		require.NoError(b, err)
	}
}

func BenchmarkBigParserV2(b *testing.B) {
	f, err := os.Open("large-file.json")
	require.NoError(b, err)
	defer f.Close()

	for range b.N {
		_, err := f.Seek(0, io.SeekStart)
		require.NoError(b, err)

		p := new(ParserV2)
		p.emitter = func(k string, v any) bool {
			return true
		}

		err = p.Parse(f)
		require.NoError(b, err)
	}
}

func BenchmarkBigParserPitr(b *testing.B) {
	f, err := os.Open("large-file.json")
	require.NoError(b, err)
	defer f.Close()

	for range b.N {
		_, err := f.Seek(0, io.SeekStart)
		require.NoError(b, err)

		p := new(ParserPitr)
		p.emitter = func(k string, v any) bool {
			return true
		}

		err = p.Parse(f)
		require.NoError(b, err)
	}
}

func BenchmarkBigMemory(b *testing.B) {
	f, err := os.Open("large-file.json")
	require.NoError(b, err)
	defer f.Close()

	for range b.N {
		_, err := f.Seek(0, io.SeekStart)
		require.NoError(b, err)

		p := new(Memory)
		p.emitter = func(k string, v any) bool {
			return true
		}

		err = p.Parse(f)
		require.NoError(b, err)
	}
}

func BenchmarkUnmarshalBig(b *testing.B) {
	f, err := os.Open("large-file.json")
	require.NoError(b, err)
	defer f.Close()

	for range b.N {
		_, err := f.Seek(0, io.SeekStart)
		require.NoError(b, err)

		var m any
		decoder := json.NewDecoder(f)
		err = decoder.Decode(&m)
		require.NoError(b, err)
	}
}

func BenchmarkUnmarshalBigV2(b *testing.B) {
	f, err := os.Open("large-file.json")
	require.NoError(b, err)
	defer f.Close()

	for range b.N {
		_, err := f.Seek(0, io.SeekStart)
		require.NoError(b, err)

		var m any
		decoder := jsontext.NewDecoder(f)
		err = jsonv2.UnmarshalDecode(decoder, &m)
		require.NoError(b, err)
	}
}
