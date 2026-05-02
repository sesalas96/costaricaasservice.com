package passwords

import "testing"

func TestHashVerify(t *testing.T) {
	enc, err := Hash("hunter2-correct-horse")
	if err != nil {
		t.Fatal(err)
	}
	if err := Verify("hunter2-correct-horse", enc); err != nil {
		t.Errorf("expected verify ok, got %v", err)
	}
	if err := Verify("wrong", enc); err == nil {
		t.Error("expected mismatch for wrong password")
	}
}

func TestVerifyInvalidEncoded(t *testing.T) {
	if err := Verify("x", "not-an-argon-hash"); err == nil {
		t.Error("expected ErrInvalidEncoded")
	}
}
