package jsonflatten

import (
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
					"GlossSee": "markup"
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
           	"embedded": [1, 2, 3]
        },
    ]
}
`

func TestTest(t *testing.T) {
	r := strings.NewReader(testJson)
	p := new(Parser)
	err := p.Parse(r)
	require.NoError(t, err)
}

func TestLarge(t *testing.T) {
	f, err := os.Open("large-file.json")
	require.NoError(t, err)

	p := new(Parser)
	err = p.Parse(f)
	require.NoError(t, err)
}
