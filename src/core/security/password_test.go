package security

import "testing"

func TestIdentityPasswordRoundTrip(t *testing.T) {
	hash, err := HashIdentityPassword("correct horse battery staple")
	if err != nil {
		t.Fatal(err)
	}
	if !VerifyIdentityPassword(hash, "correct horse battery staple") {
		t.Fatal("correct password did not verify")
	}
	if VerifyIdentityPassword(hash, "definitely wrong") {
		t.Fatal("incorrect password verified")
	}
}

func TestIdentityPasswordRejectsShortValue(t *testing.T) {
	if _, err := HashIdentityPassword("too short"); err == nil {
		t.Fatal("short password was accepted")
	}
}

func TestDummyPasswordNeverVerifies(t *testing.T) {
	if VerifyIdentityPasswordOrDummy("", "some attempted password") {
		t.Fatal("dummy password verification succeeded")
	}
}
