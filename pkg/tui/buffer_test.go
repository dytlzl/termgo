package tui

import (
	"testing"

	"github.com/dytlzl/tervi/pkg/key"
)

func Test_readBuffer(t *testing.T) {
	type args struct {
		buffer []rune
	}
	tests := []struct {
		name       string
		buffer     []rune
		wantKey    rune
		wantLength int
	}{
		{
			name:       "return ascii code in buffer and 1",
			buffer:     []rune{'a', 'b', 'c'},
			wantKey:    'a',
			wantLength: 1,
		},
		{
			name:       "return Null and 0 when buffer is empty",
			buffer:     []rune{},
			wantKey:    key.Null,
			wantLength: 0,
		},
		{
			name:       "return ArrowUp in buffer and 3",
			buffer:     []rune{0x1b, 0x5b, 'A'},
			wantKey:    key.ArrowUp,
			wantLength: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotLength := readBuffer(tt.buffer)
			if gotKey != tt.wantKey {
				t.Errorf("ReadBuffer() got key = %v, want %v", gotKey, tt.wantKey)
			}
			if gotLength != tt.wantLength {
				t.Errorf("ReadBuffer() got length = %v, want %v", gotLength, tt.wantLength)
			}
		})
	}
}
