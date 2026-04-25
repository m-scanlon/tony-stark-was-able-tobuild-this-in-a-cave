package being

import "testing"

func TestParseImpulseRejectsMultipleDividers(t *testing.T) {
	_, err := ParseImpulse("skyra prefrontal hello | first | second")
	if err == nil {
		t.Fatalf("ParseImpulse() error = nil, want error")
	}
	if err.Error() != "being: invalid impulse: expected exactly one | divider" {
		t.Fatalf("ParseImpulse() error = %q, want exact divider error", err.Error())
	}
}
