package codec

import (
	"fmt"
	"testing"
)

// TestBoolCodec tests the BoolCodec implementation
func TestBoolCodec(t *testing.T) {
	codec := BoolCodec{}

	tests := []struct {
		name        string
		input       bool
		expectedStr string
	}{
		{"true value", true, "true"},
		{"false value", false, "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Encode
			encoded, err := codec.Encode(tt.input)
			if err != nil {
				t.Errorf("BoolCodec.Encode() error = %v, want nil", err)
			}
			if encoded != tt.expectedStr {
				t.Errorf("BoolCodec.Encode() = %v, want %v", encoded, tt.expectedStr)
			}
		})
	}
}

// TestBoolCodecDecode tests the BoolCodec.Decode method with various string inputs
func TestBoolCodecDecode(t *testing.T) {
	codec := BoolCodec{}

	tests := []struct {
		name        string
		input       string
		expected    bool
		expectError bool
	}{
		{"true string", "true", true, false},
		{"false string", "false", false, false},
		{"TRUE string", "TRUE", true, false},
		{"FALSE string", "FALSE", false, false},
		{"True string", "True", true, false},
		{"False string", "False", false, false},
		{"1 string", "1", true, false},
		{"0 string", "0", false, false},
		{"t string", "t", true, false},
		{"f string", "f", false, false},
		{"T string", "T", true, false},
		{"F string", "F", false, false},
		{"invalid string", "invalid", false, true},
		{"empty string", "", false, true},
		{"number string", "123", false, true},
		{"special chars", "!@#", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoded, err := codec.Decode(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("BoolCodec.Decode() expected error for input %q, but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("BoolCodec.Decode() unexpected error for input %q: %v", tt.input, err)
				}
				if decoded != tt.expected {
					t.Errorf("BoolCodec.Decode() = %v, want %v for input %q", decoded, tt.expected, tt.input)
				}
			}
		})
	}
}

// TestBoolCodecRoundTrip tests that encoding and then decoding returns the original value
func TestBoolCodecRoundTrip(t *testing.T) {
	codec := BoolCodec{}

	testValues := []bool{true, false}

	for _, input := range testValues {
		t.Run("roundtrip_"+fmt.Sprintf("%v", input), func(t *testing.T) {
			// Encode
			encoded, err := codec.Encode(input)
			if err != nil {
				t.Errorf("BoolCodec.Encode() error = %v, want nil", err)
			}

			// Decode
			decoded, err := codec.Decode(encoded)
			if err != nil {
				t.Errorf("BoolCodec.Decode() error = %v, want nil", err)
			}

			// Verify round trip
			if decoded != input {
				t.Errorf("BoolCodec round trip failed: input = %v, encoded = %q, decoded = %v", input, encoded, decoded)
			}
		})
	}
}

// TestBoolCodecInterface tests that BoolCodec implements the Codec interface
func TestBoolCodecInterface(t *testing.T) {
	var _ Codec[bool] = BoolCodec{}
}

// TestBoolCodecGenerics tests that the generic interface works correctly with BoolCodec
func TestBoolCodecGenerics(t *testing.T) {
	boolCodec := BoolCodec{}

	// Test bool codec through generic interface
	var boolCodecGeneric Codec[bool] = boolCodec
	encodedBool, err := boolCodecGeneric.Encode(true)
	if err != nil {
		t.Errorf("Generic bool codec encode error = %v, want nil", err)
	}
	if encodedBool != "true" {
		t.Errorf("Generic bool codec encode = %v, want 'true'", encodedBool)
	}

	decodedBool, err := boolCodecGeneric.Decode("true")
	if err != nil {
		t.Errorf("Generic bool codec decode error = %v, want nil", err)
	}
	if decodedBool != true {
		t.Errorf("Generic bool codec decode = %v, want true", decodedBool)
	}
}
