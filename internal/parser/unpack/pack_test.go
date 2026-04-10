package unpack

import (
	"testing"
)

func TestUnpack(t *testing.T) {

	tests := map[string]struct {
		Text string
		Want string
	}{
		"default": {
			Text: "a4bc2d5e",
			Want: "aaaabccddddde",
		},
		"without digit": {
			Text: "abcd",
			Want: "abcd",
		},
		"first Digit": {
			Text: "3abc",
			Want: "",
		},
		"first number": {
			Text: "45",
			Want: "",
		},
		"empty string": {
			Text: "",
			Want: "",
		},
		"escape symbol": {
			Text: "d\\n5abc",
			Want: "d\\n\\n\\n\\n\\nabc",
		},
		"number": {
			Text: "f35",
			Want: "",
		},
		"zero digit": {
			Text: "qwe0",
			Want: "qw",
		},
		"raw string": {
			Text: "`qwe\\4\\5`",
			Want: "qwe45",
		},
		"rus chars": {
			Text: "а\\3",
			Want: "а\\\\\\",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, _ := Unpack(test.Text)

			if got != test.Want {
				t.Errorf("Test: %s. got %s, want %s", name, got, test.Want)
			}

		})
	}

}

func TestUnpackRaw(t *testing.T) {
	tests := map[string]struct {
		Text string
		Want string
	}{
		"default": {
			Text: "`qwe\\4\\5`",
			Want: "qwe45",
		},
		"escape number": {
			Text: "`qwe\\45`",
			Want: "qwe44444",
		},
		"escape slash": {
			Text: "`qwe\\\\5`",
			Want: "qwe\\\\\\\\\\",
		},
		"escape without digit": {
			Text: "`qw\\\ne`",
			Want: "",
		},
		"empty string": {
			Text: "",
			Want: "",
		},
		"first Digit": {
			Text: "`3abc`",
			Want: "",
		},
		"number": {
			Text: "`f35`",
			Want: "",
		},
		"zero digit": {
			Text: "`qwe0`",
			Want: "qw",
		},
		"rus chars": {
			Text: "`а\\3`",
			Want: "а3",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, _ := unpackRaw(test.Text)

			if got != test.Want {
				t.Errorf("Test: %s. got %s, want %s", name, got, test.Want)
			}
		})
	}
}
