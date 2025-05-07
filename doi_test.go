package main

import (
	//"errors" // Used for example error checking
	"testing"
)

func TestParseDOISuccess(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantPreffix string
		wantSuffix  string
	}{
		{
			name:        "Bare DOI simple",
			input:       "10.1000/xyz123",
			wantPreffix: "10.1000",
			wantSuffix:  "xyz123",
		},
		{
			name:        "Bare DOI with complex suffix (no slash)",
			input:       "10.12345/foo.bar-Baz.qux_123",
			wantPreffix: "10.12345",
			wantSuffix:  "foo.bar-Baz.qux_123",
		},
		{
			name:        "Bare DOI with longer prefix number",
			input:       "10.100000/short",
			wantPreffix: "10.100000",
			wantSuffix:  "short",
		},
		{
			name:        "HTTPS doi.org URL",
			input:       "https://doi.org/10.2000/suffixA",
			wantPreffix: "10.2000",
			wantSuffix:  "suffixA",
		},
		{
			name:        "HTTP doi.org URL",
			input:       "http://doi.org/10.3000/SUFFIX_B",
			wantPreffix: "10.3000",
			wantSuffix:  "SUFFIX_B",
		},
		{
			name:        "HTTPS dx.doi.org URL",
			input:       "https://dx.doi.org/10.4000/sometype-12345",
			wantPreffix: "10.4000",
			wantSuffix:  "sometype-12345",
		},
		{
			name:        "HTTP dx.doi.org URL",
			input:       "http://dx.doi.org/10.5000/another.suffix",
			wantPreffix: "10.5000",
			wantSuffix:  "another.suffix",
		},
		{
			name:        "HTTPS DOI.ORG URL (case insensitive domain/scheme)",
			input:       "HTTPS://DOI.ORG/10.6000/MixedCaseSuffix123",
			wantPreffix: "10.6000",
			wantSuffix:  "MixedCaseSuffix123", // Suffix case preserved
		},
		{
			name:        "hTTp dx.DOI.org URL with mixed case scheme/domain",
			input:       "hTTp://dX.dOi.OrG/10.7000/CASE",
			wantPreffix: "10.7000",
			wantSuffix:  "CASE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d DOI
			err := d.parseDOI(tt.input)

			if err != nil {
				t.Errorf("parseDOI(%q) expected no error, but got: %v", tt.input, err)
				return // No point checking fields if error occurred unexpectedly
			}
			if d.preffix != tt.wantPreffix {
				t.Errorf("parseDOI(%q) got preffix %q, want %q", tt.input, d.preffix, tt.wantPreffix)
			}
			if d.suffix != tt.wantSuffix {
				t.Errorf("parseDOI(%q) got suffix %q, want %q", tt.input, d.suffix, tt.wantSuffix)
			}
		})
	}
}

func TestParseDOIFailure(t *testing.T) {
	tests := []struct {
		name  string
		input string
		// For more specific error checking, you could add:
		// expectedErrorMsgContains string
	}{
		// General Invalid Inputs
		{name: "Empty string", input: ""},
		{name: "Completely random string", input: "this is not a doi at all"},
		{name: "Invalid URL (not a DOI resolver)", input: "http://example.com/path"},
		{name: "URL path resembling DOI but wrong domain", input: "https://example.org/10.1234/abcdef"},
		{name: "Only DOI host (https), no DOI string", input: "https://doi.org/"},
		{name: "Only DOI host (http), no DOI string", input: "http://dx.doi.org/"},
		{name: "Invalid URL scheme", input: "ftp://doi.org/10.1000/xyz123"},
		{name: "Relative URL looking like DOI", input: "/10.1234/asdf"},

		// Prefix Rule Violations (must start with "10.")
		{name: "Bare DOI prefix not starting with 10.", input: "1.1234/abcdef"},
		{name: "Bare DOI prefix missing dot after 10", input: "101234/abcdef"},
		{name: "Bare DOI prefix is just '10.' (empty actual prefix part)", input: "10./abcdef"},
		{name: "Bare DOI prefix is just '10' (no dot)", input: "10/abcdef"},
		{name: "URL DOI prefix not starting with 10.", input: "https://doi.org/12.1234/abcdef"},
		{name: "URL DOI prefix missing dot after 10", input: "http://doi.org/101234/abcdef"},
		{name: "URL DOI prefix is just '10.'", input: "https://dx.doi.org/10./abcdef"},

		// Slash Rule Violations (only one slash in DOI part)
		{name: "Bare DOI with multiple slashes", input: "10.1234/abc/def"},
		{name: "Bare DOI with multiple slashes at start", input: "10.1234//def"},
		{name: "URL DOI with multiple slashes in DOI part", input: "https://doi.org/10.1234/abc/def"},
		{name: "Bare DOI no slash separator", input: "10.1234abcdef"},
		{name: "URL DOI no slash separator in DOI part", input: "http://doi.org/10.1234abcdef"},
		{name: "URL DOI with slash only at the end of URL", input: "https://doi.org/10.1234/"}, // This also means empty suffix

		// Missing Parts
		{name: "Bare DOI missing suffix (ends with slash)", input: "10.1234/"},
		{name: "Bare DOI missing prefix (starts with slash)", input: "/abcdef"}, // Also violates "starts with 10."
		{name: "URL DOI missing suffix (ends with slash)", input: "https://doi.org/10.5678/"},
		{name: "URL DOI missing actual DOI string (just spaces after host)", input: "https://doi.org/   "},
		{name: "DOI string is just a slash", input: "/"},
		{name: "DOI string is just '10.' and a slash", input: "10./"}, // Empty actual prefix and empty suffix

		// Whitespace Violations (assuming no automatic trim of the DOI part)
		//{name: "Bare DOI with leading whitespace", input: " 10.1000/xyz123"},
		//{name: "Bare DOI with trailing whitespace", input: "10.1000/xyz123 "},
		{name: "Bare DOI with internal whitespace in prefix", input: "10.10 00/xyz123"},
		{name: "Bare DOI with internal whitespace in suffix", input: "10.1000/xyz 123"},
		//{name: "URL with leading whitespace before scheme", input: " https://doi.org/10.1000/xyz123"},
		//{name: "URL with trailing whitespace after DOI", input: "https://doi.org/10.1000/xyz123 "},
		{name: "URL with whitespace between host and DOI string", input: "https://doi.org/ 10.1000/xyz123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d DOI
			initialPreffix := d.preffix // Capture initial state if needed for stricter checks
			initialSuffix := d.suffix

			err := d.parseDOI(tt.input)

			if err == nil {
				t.Errorf("parseDOI(%q) expected an error, but got nil. Parsed to preffix=%q, suffix=%q", tt.input, d.preffix, d.suffix)
			} else {
				// Optionally, check if fields are unmodified or reset on error
				if d.preffix != initialPreffix || d.preffix != "" { // Check against initial or known zero value
					t.Logf("parseDOI(%q) on error, d.preffix became %q, expected it to be empty or unchanged (%q)", tt.input, d.preffix, initialPreffix)
				}
				if d.suffix != initialSuffix || d.suffix != "" { // Check against initial or known zero value
					t.Logf("parseDOI(%q) on error, d.suffix became %q, expected it to be empty or unchanged (%q)", tt.input, d.suffix, initialSuffix)
				}
				// If you add `expectedErrorMsgContains` to the struct:
				// if tt.expectedErrorMsgContains != "" && !strings.Contains(err.Error(), tt.expectedErrorMsgContains) {
				//  t.Errorf("parseDOI(%q) got error %q, want error containing %q", tt.input, err.Error(), tt.expectedErrorMsgContains)
				// }
			}
		})
	}
}
