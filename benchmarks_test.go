package jsonflatten

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	jsonv2 "github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/stretchr/testify/require"
)

func BenchmarkSmall(b *testing.B) {
	b.Run("parser=v1", benchmarkSmallParser)
	b.Run("parser=v2", benchmarkSmallParserV2)
	b.Run("parser=pitr", benchmarkSmallParserPitr)
	b.Run("parser=memory", benchmarkSmallMemory)
}

func BenchmarkBig(b *testing.B) {
	b.Run("parser=v1", benchmarkBigParser)
	b.Run("parser=v2", benchmarkBigParserV2)
	b.Run("parser=pitr", benchmarkBigParserPitr)
	b.Run("parser=memory", benchmarkBigMemory)
}

func BenchmarkUnmarshalSmall(b *testing.B) {
	b.Run("parser=v1", benchmarkUnmarshalSmall)
	b.Run("parser=v2", benchmarkUnmarshalSmallV2)
}

func BenchmarkUnmarshalBig(b *testing.B) {
	b.Run("parser=v1", benchmarkUnmarshalBig)
	b.Run("parser=v2", benchmarkUnmarshalBigV2)
}

func benchmarkSmallParser(b *testing.B) {
	r := strings.NewReader(testJson)

	for b.Loop() {
		_, err := r.Seek(0, io.SeekStart)
		require.NoError(b, err)

		p := new(Parser)
		p.emitter = func(k string, v any) bool {
			return true
		}

		err = p.Parse(r)
		require.NoError(b, err)
	}
}

func benchmarkSmallParserV2(b *testing.B) {
	r := strings.NewReader(testJson)

	for b.Loop() {
		_, err := r.Seek(0, io.SeekStart)
		require.NoError(b, err)

		p := new(ParserV2)
		p.emitter = func(k string, v any) bool {
			return true
		}

		err = p.Parse(r)
		require.NoError(b, err)
	}
}

func benchmarkSmallParserPitr(b *testing.B) {
	r := strings.NewReader(testJson)

	for b.Loop() {
		_, err := r.Seek(0, io.SeekStart)
		require.NoError(b, err)

		p := new(ParserPitr)
		p.emitter = func(k string, v any) bool {
			return true
		}

		err = p.Parse(r)
		require.NoError(b, err)
	}
}

func benchmarkSmallMemory(b *testing.B) {
	r := strings.NewReader(testJson)

	for b.Loop() {
		_, err := r.Seek(0, io.SeekStart)
		require.NoError(b, err)

		p := new(Memory)
		p.emitter = func(k string, v any) bool {
			return true
		}

		err = p.Parse(r)
		require.NoError(b, err)
	}
}

func benchmarkSmallMemoryV2(b *testing.B) {
	r := strings.NewReader(testJson)

	for b.Loop() {
		_, err := r.Seek(0, io.SeekStart)
		require.NoError(b, err)

		p := new(MemoryV2)
		p.emitter = func(k string, v any) bool {
			return true
		}

		err = p.Parse(r)
		require.NoError(b, err)
	}
}

func benchmarkUnmarshalSmall(b *testing.B) {
	for b.Loop() {
		var m any
		err := json.Unmarshal([]byte(testJson), &m)
		require.NoError(b, err)
	}
}

func benchmarkUnmarshalSmallV2(b *testing.B) {
	for b.Loop() {
		var m any
		err := jsonv2.Unmarshal([]byte(testJson), &m)
		require.NoError(b, err)
	}
}

func benchmarkBigParser(b *testing.B) {
	f, err := os.Open("large-file.json")
	require.NoError(b, err)
	defer f.Close()

	for b.Loop() {
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

func benchmarkBigParserV2(b *testing.B) {
	f, err := os.Open("large-file.json")
	require.NoError(b, err)
	defer f.Close()

	for b.Loop() {
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

func benchmarkBigParserPitr(b *testing.B) {
	f, err := os.Open("large-file.json")
	require.NoError(b, err)
	defer f.Close()

	for b.Loop() {
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

func benchmarkBigMemory(b *testing.B) {
	f, err := os.Open("large-file.json")
	require.NoError(b, err)
	defer f.Close()

	for b.Loop() {
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

func benchmarkUnmarshalBig(b *testing.B) {
	f, err := os.Open("large-file.json")
	require.NoError(b, err)
	defer f.Close()

	for b.Loop() {
		_, err := f.Seek(0, io.SeekStart)
		require.NoError(b, err)

		var m any
		decoder := json.NewDecoder(f)
		err = decoder.Decode(&m)
		require.NoError(b, err)
	}
}

func benchmarkUnmarshalBigV2(b *testing.B) {
	f, err := os.Open("large-file.json")
	require.NoError(b, err)
	defer f.Close()

	for b.Loop() {
		_, err := f.Seek(0, io.SeekStart)
		require.NoError(b, err)

		var m any
		decoder := jsontext.NewDecoder(f)
		err = jsonv2.UnmarshalDecode(decoder, &m)
		require.NoError(b, err)
	}
}
