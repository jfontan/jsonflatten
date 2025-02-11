package jsonflatten

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

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

func TestPrint(t *testing.T) {
	t.Skip()
	r := strings.NewReader(testJson)
	p := new(Parser)
	err := p.Parse(r)
	require.NoError(t, err)
}

func TestMap(t *testing.T) {
	r := strings.NewReader(testJson)
	m := make(map[string]any)

	p := new(Parser)
	p.emitter = func(k string, v any) {
		m[k] = v
	}

	err := p.Parse(r)
	require.NoError(t, err)

	expected := map[string]any{
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

func BenchmarkSmall(b *testing.B) {
	for range b.N {
		r := strings.NewReader(testJson)
		p := new(Parser)
		p.emitter = func(k string, v any) {
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

func BenchmarkBig(b *testing.B) {
	f, err := os.Open("large-file.json")
	require.NoError(b, err)
	defer f.Close()

	for range b.N {
		_, err := f.Seek(0, io.SeekStart)
		require.NoError(b, err)

		p := new(Parser)
		p.emitter = func(k string, v any) {
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
