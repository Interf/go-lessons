package pack

import "testing"

func TestPack(t *testing.T) {
	tests := map[string]struct {
		Text string
		Want string
	}{
		"default": {
			Text: "aaaabccddddde",
			Want: "a4bc2d5e",
		},
		"empty string": {
			Text: "",
			Want: "",
		},
		"rus chars": {
			Text: "фффцц",
			Want: "ф3ц2",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := Pack(test.Text)

			if got != test.Want {
				t.Errorf("Test: %s. Got %s, want %s", name, got, test.Want)
			}
		})
	}
}
