package validate

import (
	"strings"
	"testing"
)

// SampleObject exercises: plain "required", the custom "nohtml" rule, the
// custom "safeurl" rule (chained after the built-in "url" rule), and a field
// with no json tag at all (to check the field-name fallback behaviour).
type sampleObject struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"omitempty,nohtml"`
	RedirectUrl string `json:"redirectUrl" validate:"omitempty,url,safeurl"`
	Code        string `validate:"required"`
}

func validObject() sampleObject {
	return sampleObject{
		Name:        "Pastor Matey Aryeh",
		Description: "Thank you God.",
		RedirectUrl: "https://example.com/callback",
		Code:        "ABC123",
	}
}

func TestValidateObject_ValidObjectPasses(t *testing.T) {
	message, err := ValidateObject(validObject())
	if err != nil {
		t.Fatalf("expected no error, got %v (message: %q)", err, message)
	}
	if message != "" {
		t.Fatalf("expected empty message, got %q", message)
	}
}

func TestValidateObject_MissingRequiredField(t *testing.T) {
	object := validObject()
	object.Name = ""

	message, err := ValidateObject(object)
	if err == nil {
		t.Fatal("expected an error for missing required field, got nil")
	}
	// the message must use the JSON tag name ("name"), not the Go field name ("Name")
	if !strings.Contains(message, "[name]") {
		t.Errorf("expected message to reference json tag %q, got %q", "[name]", message)
	}
}

func TestValidateObject_NoHtmlRejectsAngleBrackets(t *testing.T) {
	cases := []string{
		`<a href="https://evil.example">scamzzzzzz!!!</a>`,
		`<script>alert(1)</script>`,
		`plain text <img src=x onerror=alert(1)>`,
	}
	for _, description := range cases {
		object := validObject()
		object.Description = description

		message, err := ValidateObject(object)
		if err == nil {
			t.Errorf("expected nohtml to reject %q, got no error", description)
			continue
		}
		if !strings.Contains(message, "[description]") {
			t.Errorf("expected message to reference [description] for input %q, got %q", description, message)
		}
	}
}

func TestValidateObject_NoHtmlAllowsOrdinaryPunctuation(t *testing.T) {
	// guard against the rule being overly aggressive and flagging text that
	// merely contains other special characters, not angle brackets
	object := validObject()
	object.Description = `Waakye 2 cedis, "great deal"!`

	message, err := ValidateObject(object)
	if err != nil {
		t.Fatalf("expected ordinary punctuation to pass, got error %v (message: %q)", err, message)
	}
}

func TestValidateObject_SafeUrlRejectsDangerousSchemes(t *testing.T) {
	cases := []string{
		"javascript:alert(document.cookie)",
		"data:text/html,<script>alert(1)</script>",
		"file:///etc/passwd",
	}
	for _, redirectUrl := range cases {
		object := validObject()
		object.RedirectUrl = redirectUrl

		message, err := ValidateObject(object)
		if err == nil {
			t.Errorf("expected safeurl to reject %q, got no error", redirectUrl)
			continue
		}
		if !strings.Contains(message, "[redirectUrl]") {
			t.Errorf("expected message to reference [redirectUrl] for input %q, got %q", redirectUrl, message)
		}
	}
}

func TestValidateObject_SafeUrlAllowsHttpAndHttps(t *testing.T) {
	for _, redirectUrl := range []string{"http://example.com", "https://example.com/path?query=1"} {
		object := validObject()
		object.RedirectUrl = redirectUrl

		_, err := ValidateObject(object)
		if err != nil {
			t.Errorf("expected %q to be accepted, got error", redirectUrl)
		}
	}
}

func TestValidateObject_SafeUrlSkippedWhenEmptyAndOmitempty(t *testing.T) {
	object := validObject()
	object.RedirectUrl = ""

	_, err := ValidateObject(object)
	if err != nil {
		t.Fatalf("expected empty redirectUrl to be skipped via omitempty, got error %v", err)
	}
}

func TestValidateObject_MultipleFailuresAreAllReported(t *testing.T) {
	object := validObject()
	object.Name = ""
	object.Description = "<script>alert(1)</script>"

	message, err := ValidateObject(object)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(message, "[name]") {
		t.Errorf("expected message to contain [name], got %q", message)
	}
	if !strings.Contains(message, "[description]") {
		t.Errorf("expected message to contain [description], got %q", message)
	}
}

func TestValidateObject_FieldNameFallsBackToGoNameWithoutJsonTag(t *testing.T) {
	object := validObject()
	object.Code = ""

	message, err := ValidateObject(object)
	if err == nil {
		t.Fatal("expected an error for missing Code, got nil")
	}
	if !strings.Contains(message, "[Code]") {
		t.Errorf("expected message to fall back to Go field name [Code], got %q", message)
	}
}
