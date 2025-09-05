package codec

import (
	"testing"
)

// TestStringCodec tests the StringCodec implementation
func TestStringCodec(t *testing.T) {
	codec := StringCodec{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"simple string", "hello", "hello"},
		{"unicode string", "世界", "世界"},
		{"special characters", "!@#$%^&*()", "!@#$%^&*()"},
		{"numbers as string", "12345", "12345"},
		{"whitespace", "  \t\n  ", "  \t\n  "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Encode
			encoded, err := codec.Encode(tt.input)
			if err != nil {
				t.Errorf("StringCodec.Encode() error = %v, want nil", err)
			}
			if encoded != tt.expected {
				t.Errorf("StringCodec.Encode() = %v, want %v", encoded, tt.expected)
			}

			// Test Decode
			decoded, err := codec.Decode(tt.input)
			if err != nil {
				t.Errorf("StringCodec.Decode() error = %v, want nil", err)
			}
			if decoded != tt.expected {
				t.Errorf("StringCodec.Decode() = %v, want %v", decoded, tt.expected)
			}
		})
	}
}

// TestStringCodecRoundTrip tests that encoding and then decoding returns the original value
func TestStringCodecRoundTrip(t *testing.T) {
	codec := StringCodec{}

	testValues := []string{"", "hello", "世界", "!@#$%^&*()", "12345", "  \t\n  "}

	for _, input := range testValues {
		t.Run("roundtrip_"+input, func(t *testing.T) {
			// Encode
			encoded, err := codec.Encode(input)
			if err != nil {
				t.Errorf("StringCodec.Encode() error = %v, want nil", err)
			}

			// Decode
			decoded, err := codec.Decode(encoded)
			if err != nil {
				t.Errorf("StringCodec.Decode() error = %v, want nil", err)
			}

			// Verify round trip
			if decoded != input {
				t.Errorf("StringCodec round trip failed: input = %q, encoded = %q, decoded = %q", input, encoded, decoded)
			}
		})
	}
}

// TestStringCodecInterface tests that StringCodec implements the Codec interface
func TestStringCodecInterface(t *testing.T) {
	var _ Codec[string] = StringCodec{}
}

// TestStringCodecGenerics tests that the generic interface works correctly with StringCodec
func TestStringCodecGenerics(t *testing.T) {
	stringCodec := StringCodec{}

	// Test string codec through generic interface
	var stringCodecGeneric Codec[string] = stringCodec
	encoded, err := stringCodecGeneric.Encode("test")
	if err != nil {
		t.Errorf("Generic string codec encode error = %v, want nil", err)
	}
	if encoded != "test" {
		t.Errorf("Generic string codec encode = %v, want 'test'", encoded)
	}

	decoded, err := stringCodecGeneric.Decode("test")
	if err != nil {
		t.Errorf("Generic string codec decode error = %v, want nil", err)
	}
	if decoded != "test" {
		t.Errorf("Generic string codec decode = %v, want 'test'", decoded)
	}
}
