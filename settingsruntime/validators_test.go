package settingsruntime

import "testing"

func TestBoolValidator(t *testing.T) {
	validator := BoolValidator()

	if err := validator("true"); err != nil {
		t.Fatalf("expected true to be valid, got error: %v", err)
	}

	if err := validator("not-bool"); err == nil {
		t.Fatal("expected invalid bool to fail")
	}
}

func TestIntValidator(t *testing.T) {
	validator := IntValidator(1, 5)

	if err := validator("3"); err != nil {
		t.Fatalf("expected value in range, got error: %v", err)
	}

	if err := validator("0"); err == nil {
		t.Fatal("expected value below min to fail")
	}

	if err := validator("6"); err == nil {
		t.Fatal("expected value above max to fail")
	}
}

func TestNonEmptyValidator(t *testing.T) {
	validator := NonEmptyValidator()

	if err := validator("value"); err != nil {
		t.Fatalf("expected non-empty value to pass, got error: %v", err)
	}

	if err := validator("   "); err == nil {
		t.Fatal("expected blank value to fail")
	}
}
