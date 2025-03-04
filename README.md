# JSON Flatten

> **NOTE:** this is still an experiment to test several ways extract values with a single level of keys. It lacks tests to check how it behaves with different objects and with errors.

This project converts a JSON object into a flat one that consists on an object with a single level of depths and the keys contain the path separated by `.`.

For example, this JSON:

```json
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
```

becomes this with the standard print emitter:

```python
glossary.title = "example glossary"
glossary.GlossDiv.title = "S"
glossary.GlossDiv.GlossList.GlossEntry.ID = "SGML"
glossary.GlossDiv.GlossList.GlossEntry.SortAs = "SGML"
glossary.GlossDiv.GlossList.GlossEntry.GlossTerm = "Standard Generalized Markup Language"
glossary.GlossDiv.GlossList.GlossEntry.Acronym = "SGML"
glossary.GlossDiv.GlossList.GlossEntry.Abbrev = "ISO 8879:1986"
glossary.GlossDiv.GlossList.GlossEntry.GlossDef.para = "A meta-markup language, used to create markup languages such as DocBook."
glossary.GlossDiv.GlossList.GlossEntry.GlossDef.GlossSeeAlso.0 = "GML"
glossary.GlossDiv.GlossList.GlossEntry.GlossDef.GlossSeeAlso.1 = "XML"
glossary.GlossDiv.GlossList.GlossEntry.GlossSee = "markup"
glossary.GlossDiv.GlossList.GlossEntry.float64 = 42
glossary.GlossDiv.GlossList.GlossEntry.bool = true
glossary.GlossDiv.GlossList.GlossEntry.null = <nil>
array.0.one = 1
array.0.two = 2
array.1.three = 1
array.1.four = 2
array.1.embedded.0 = 1
array.1.embedded.1 = 2
array.1.embedded.2 = 3
array.1.embedded.3 = true
array.1.embedded.4 = <nil>
array.1.embedded.5 = "string"
```

## Versions

- `Parser`: this version uses the standard json package tokenizer. Emits all the values with the key that represents the path to them. It is done in an stream fashion so the values are emitted as they are found.
- `ParserPitr`: does the same as `Parser` but uses another tokenizer: https://pkg.go.dev/pitr.ca/jsontokenizer
- `Memory`: this one unmarshals the whole JSON object in memory using standard json package and iterates over all the values in it. It is used to test the difference with the other parsers.

## Benchmark

There are two sizes of objects tested:

- `Small`: 49 lines of JSON. Also has arrays to check that it properly handles them.
- `Big`: the file `large-file.json` that is 25 Mb of JSON.

Besides the parsers described before it also benchmarks just unmarshalling the object to memory.

```
$ go version
go version go1.24.0 linux/amd64

$ go test -bench=. -benchmem ./...
goos: linux
goarch: amd64
pkg: github.com/jfontan/jsonflatten
cpu: AMD Ryzen 7 5800X 8-Core Processor
BenchmarkSmallParser-16            56736             20951 ns/op           11640 B/op        416 allocs/op
BenchmarkSmallParserPitr-16       205270              5576 ns/op            7832 B/op         97 allocs/op
BenchmarkSmallMemory-16            95386             12413 ns/op           11064 B/op        130 allocs/op
BenchmarkUnmarshalSmall-16        139206              8475 ns/op            5536 B/op         93 allocs/op
BenchmarkBigParser-16                  2         542336878 ns/op        192280732 B/op  10093540 allocs/op
BenchmarkBigParserPitr-16              9         121950704 ns/op        54910155 B/op    2100604 allocs/op
BenchmarkBigMemory-16                  4         279605370 ns/op        183342498 B/op   2328719 allocs/op
BenchmarkUnmarshalBig-16               5         226514305 ns/op        164635795 B/op   1775777 allocs/op
PASS
ok      github.com/jfontan/jsonflatten  12.690s
```
