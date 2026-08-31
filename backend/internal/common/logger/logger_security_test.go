package logger

import (
	"errors"
	"testing"
)

func TestSanitizeLogValueEscapesLineBreaks(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  interface{}
	}{
		{name: "string", value: "first\nsecond\rthird", want: `first\nsecond\rthird`},
		{name: "error", value: errors.New("failed\nforged"), want: `failed\nforged`},
		{name: "bytes", value: []byte("payload\r\nnext"), want: `payload\r\nnext`},
		{name: "number", value: 42, want: 42},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := sanitizeLogValue(test.value)
			if got != test.want {
				t.Fatalf("sanitizeLogValue() = %#v, want %#v", got, test.want)
			}
		})
	}
}
