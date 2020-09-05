package required

import "testing"

func assertError(t *testing.T, err error, expected error) {
	if err != nil {
		if IsRequiredErr(err) {
			assertRequiredError(t, err, expected)
		} else {
			assertStdError(t, err, expected)
		}
	} else {
		assertNil(t, err, expected)
	}
}

func assertNil(t *testing.T, err error, expected error) {
	if expected != err {
		t.Fatalf("expected %v, received: %v", expected, err)
	}
}

func assertStdError(t *testing.T, err, expected error) {
	if err != expected {
		t.Fatalf("expected %v, received: %v", expected, err)
	}
}

func assertRequiredError(t *testing.T, err, expected error) {
	if err.(requiredErr).err != expected {
		t.Fatalf("expected %v, received: %v", expected, err)
	}
}
