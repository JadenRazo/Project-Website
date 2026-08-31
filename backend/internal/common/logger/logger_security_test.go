package logger

import (
	"errors"
	"testing"
)

type forgedStringer struct{}

func (forgedStringer) String() string { return "trusted\nforged" }

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
		{name: "custom stringer", value: forgedStringer{}, want: `trusted\nforged`},
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

func TestSanitizeLogMessageEscapesRenderedLineBreaks(t *testing.T) {
	got := sanitizeLogMessage([]interface{}{"request ", forgedStringer{}, "\rnext"})
	want := `request trusted\nforged\rnext`
	if got != want {
		t.Fatalf("sanitizeLogMessage() = %q, want %q", got, want)
	}
}

func TestSanitizeFormattedLogMessageEscapesRenderedLineBreaks(t *testing.T) {
	got := sanitizeFormattedLogMessage("request %s", []interface{}{"first\nsecond"})
	want := `request first\nsecond`
	if got != want {
		t.Fatalf("sanitizeFormattedLogMessage() = %q, want %q", got, want)
	}
}
