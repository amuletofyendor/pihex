package pihex

import (
	"bytes"
	"encoding/hex"
	"testing"
)

const (
	Digits0To49OfPi        = "243f6a8885a308d313198a2e03707344a4093822299f31d008"
	Digits49950To49999OfPi = "08ce5db76425c7b4bc661001cbc30e805c6db26c2a35ab5135"
	Digits99950To99999OfPi = "443388751069558b3e62e612bc302ec487aa9a6ea22673c1a5"
)

func testDigitRange(rangeStart, span int64, expectedHexStr string, t *testing.T) {
	var digitBuf bytes.Buffer
	DigitRange(rangeStart, span, &digitBuf)
	resultString := hex.EncodeToString(digitBuf.Bytes())

	if resultString != expectedHexStr {
		t.Errorf("Expected %s, got %s", expectedHexStr, resultString)
	}
}

func TestDigits0To49(t *testing.T) {
	testDigitRange(0, 50, Digits0To49OfPi, t)
}

func TestDigits49950To49999(t *testing.T) {
	testDigitRange(49950, 50, Digits49950To49999OfPi, t)
}

func TestDigits99950To99999(t *testing.T) {
	testDigitRange(99950, 50, Digits99950To99999OfPi, t)
}

func TestDigitRangeErrors(t *testing.T) {
	var digitBuf bytes.Buffer
	err := DigitRange(-1, 2, &digitBuf)

	if err == nil {
		t.Errorf("Expected an error when index < 0")
	}

	err = DigitRange(0, 1, &digitBuf)

	if err == nil {
		t.Errorf("Expected an error when span < 2")
	}

	err = DigitRange(0, 3, &digitBuf)

	if err == nil {
		t.Errorf("Expected an error when span odd")
	}
}
