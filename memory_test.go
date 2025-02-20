package jsonflatten

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemory(t *testing.T) {
	t.Skip()
	r := strings.NewReader(testJson)
	m := new(Memory)
	err := m.Parse(r)
	require.NoError(t, err)
}
