package utils_test

import (
	"testing"

	"github.com/jponc/domain-crawler/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		out  []string
	}{
		{
			name: "removes duplicates from slice",
			in:   []string{"a", "b", "a", "c", "b"},
			out:  []string{"a", "b", "c"},
		},
		{
			name: "handle empty slice",
			in:   []string{},
			out:  []string{},
		},
		{
			name: "handle slice with one element",
			in:   []string{"a"},
			out:  []string{"a"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.RemoveDuplicates(tt.in)
			require.Equal(t, tt.out, result)
		})
	}
}
