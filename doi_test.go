package main

import (
	//"errors" // Used for example error checking
	"testing"
)

func TestParseDOI(t *testing.T) {
	testCases := []struct {
		name            string
		input           string
		expectedPreffix string // Corresponds to DOI.preffix
		expectedSuffix  string
		expectError     bool
	}{
		// --- Happy Path: Valid URL forms ---
		{
			name:            "Valid HTTPS URL (doi.org)",
			input:           "https://doi.org/10.1000/xyz123",
			expectedPreffix: "10.1000",
			expectedSuffix:  "xyz123",
			expectError:     false,
		},
		{
			name:            "Valid HTTP URL (doi.org)",
			input:           "http://doi.org/10.1000/xyz123",
			expectedPreffix: "10.1000",
			expectedSuffix:  "xyz123",
			expectError:     false,
		},
		{
			name:            "Valid HTTPS URL (dx.doi.org)",
			input:           "https://dx.doi.org/10.1234/AB-CD",
			expectedPreffix: "10.1234",
			expectedSuffix:  "AB-CD",
			expectError:     false,
		},
		{
			name:            "Valid HTTP URL (dx.doi.org) with suffix containing slashes",
			input:           "http://dx.doi.org/10.5678/SuffixWith/Slashes",
			expectedPreffix: "10.5678",
			expectedSuffix:  "SuffixWith/Slashes",
			expectError:     false,
		},
		{
			name:            "URL with mixed case scheme (HtTpS)",
			input:           "HtTpS://doi.org/10.1000/MixedCaseScheme",
			expectedPreffix: "10.1000",
			expectedSuffix:  "MixedCaseScheme",
			expectError:     false,
		},
		{
			name:            "URL with mixed case domain (DoI.oRg)",
			input:           "https://DoI.oRg/10.1000/MixedCaseDomain",
			expectedPreffix: "10.1000",
			expectedSuffix:  "MixedCaseDomain",
			expectError:     false,
		},
		/*{
			name:            "URL with port number",
			input:           "https://doi.org:8080/10.1101/2020.03.13.991234",
			expectedPreffix: "10.1101",
			expectedSuffix:  "2020.03.13.991234",
			expectError:     false,
		},*/
		{
			name:            "URL with query parameters (should be ignored)",
			input:           "https://doi.org/10.1000/xyz123?foo=bar&baz=qux",
			expectedPreffix: "10.1000",
			expectedSuffix:  "xyz123",
			expectError:     false,
		},
		{
			name:            "URL with fragment (should be ignored)",
			input:           "https://doi.org/10.1000/xyz123#section1",
			expectedPreffix: "10.1000",
			expectedSuffix:  "xyz123",
			expectError:     false,
		},
		{
			name:            "URL with query and fragment (dx.doi.org)",
			input:           "http://dx.doi.org/10.1234/foobar?key=value#title",
			expectedPreffix: "10.1234",
			expectedSuffix:  "foobar",
			expectError:     false,
		},

		// --- Happy Path: Valid DOI URI scheme forms ---
		{
			name:            "Valid doi: scheme (lowercase)",
			input:           "doi:10.5555/123456789",
			expectedPreffix: "10.5555",
			expectedSuffix:  "123456789",
			expectError:     false,
		},
		{
			name:            "Valid DOI: scheme (uppercase)",
			input:           "DOI:10.5555/abcdef",
			expectedPreffix: "10.5555",
			expectedSuffix:  "abcdef",
			expectError:     false,
		},
		{
			name:            "doi: scheme with complex suffix",
			input:           "doi:10.123/Suffix.With-Special_Chars:(Parentheses)/And/Slashes",
			expectedPreffix: "10.123",
			expectedSuffix:  "Suffix.With-Special_Chars:(Parentheses)/And/Slashes",
			expectError:     false,
		},

		// --- Happy Path: Valid Bare identifier forms ---
		{
			name:            "Valid bare identifier",
			input:           "10.1016/j.physletb.2003.10.071",
			expectedPreffix: "10.1016",
			expectedSuffix:  "j.physletb.2003.10.071",
			expectError:     false,
		},
		{
			name:            "Bare identifier with short prefix registrant code and short suffix",
			input:           "10.1/a",
			expectedPreffix: "10.1",
			expectedSuffix:  "a",
			expectError:     false,
		},
		{
			name:            "Bare identifier with complex suffix containing various allowed characters",
			input:           "10.1002/(SICI)1097-4628(19970321)63:12<1457::AID-APP2>3.0.CO;2-I",
			expectedPreffix: "10.1002",
			expectedSuffix:  "(SICI)1097-4628(19970321)63:12<1457::AID-APP2>3.0.CO;2-I",
			expectError:     false,
		},
		{
			name:            "Bare identifier, case preserved in prefix registrant and suffix",
			input:           "10.1234/EfGh-123.XyZ",
			expectedPreffix: "10.1234",
			expectedSuffix:  "EfGh-123.XyZ",
			expectError:     false,
		},
		{
			name:            "Bare DOI prefix with dots in registrant code",
			input:           "10.100.123/suffixpart.v2",
			expectedPreffix: "10.100.123",
			expectedSuffix:  "suffixpart.v2",
			expectError:     false,
		},

		// --- Edge Cases & Whitespace (Assuming trim) ---
		{
			name:            "Leading whitespace bare identifier",
			input:           "  10.2000/b",
			expectedPreffix: "10.2000",
			expectedSuffix:  "b",
			expectError:     false,
		},
		{
			name:            "Trailing whitespace bare identifier",
			input:           "10.2000/c  ",
			expectedPreffix: "10.2000",
			expectedSuffix:  "c",
			expectError:     false,
		},
		{
			name:            "Leading and Trailing whitespace bare identifier",
			input:           "  10.2000/d_ ", // Also testing suffix with trailing space before overall trim
			expectedPreffix: "10.2000",
			expectedSuffix:  "d_", // Assuming suffix itself might have spaces if they are not trailing the whole string
			expectError:     false,
		},
		{
			name:            "Leading whitespace URL",
			input:           "  https://doi.org/10.3000/e",
			expectedPreffix: "10.3000",
			expectedSuffix:  "e",
			expectError:     false,
		},
		{
			name:            "Trailing whitespace URL",
			input:           "https://doi.org/10.3000/f  ",
			expectedPreffix: "10.3000",
			expectedSuffix:  "f",
			expectError:     false,
		},
		{
			name:            "Leading whitespace DOI scheme",
			input:           "  doi:10.4000/g",
			expectedPreffix: "10.4000",
			expectedSuffix:  "g",
			expectError:     false,
		},
		{
			name:            "Trailing whitespace DOI scheme",
			input:           "doi:10.4000/h  ",
			expectedPreffix: "10.4000",
			expectedSuffix:  "h",
			expectError:     false,
		},
		{
			name:            "URL with percent-encoded slash %2F in suffix (treated literally)",
			input:           "https://doi.org/10.1000/suffix%2Fpart",
			expectedPreffix: "10.1000",
			expectedSuffix:  "suffix%2Fpart",
			expectError:     false,
		},

		// --- Invalid Inputs & Error Conditions ---
		{name: "Invalid: Empty string", input: "", expectError: true},
		{name: "Invalid: Only whitespace", input: "   \t\n ", expectError: true},
		{name: "Invalid: Bare identifier missing prefix", input: "/suffixonly", expectError: true},
		{name: "Invalid: Bare identifier missing suffix", input: "10.1000/", expectError: true},
		{name: "Invalid: Bare identifier prefix not starting with 10.", input: "11.1000/suffix", expectError: true},
		{name: "Invalid: Bare identifier prefix is just '10.' (empty registrant code)", input: "10./suffix", expectError: true},
		{name: "Invalid: Bare identifier prefix is just '10' (no dot, no registrant code)", input: "10/suffix", expectError: true},
		{name: "Invalid: Bare identifier no slash separator", input: "10.1000suffix", expectError: true},
		{name: "Invalid: URL (doi.org) with no path (no prefix/suffix)", input: "https://doi.org/", expectError: true},
		{name: "Invalid: URL (dx.doi.org) with no path", input: "http://dx.doi.org", expectError: true},
		{name: "Invalid: URL with only prefix, no slash separating suffix", input: "https://doi.org/10.1000", expectError: true},
		{name: "Invalid: URL with prefix and slash, but empty suffix", input: "https://doi.org/10.1000/", expectError: true},
		{name: "Invalid: URL with unrecognized domain", input: "https://example.com/10.1000/suffix", expectError: true},
		{name: "Invalid: DOI scheme (doi:) with no value", input: "doi:", expectError: true},
		{name: "Invalid: DOI scheme (DOI:) with no value", input: "DOI:", expectError: true},
		{name: "Invalid: DOI scheme with only prefix, no slash separating suffix", input: "doi:10.1000", expectError: true},
		{name: "Invalid: DOI scheme with prefix and slash, but empty suffix", input: "doi:10.1000/", expectError: true},
		{name: "Invalid: Malformed DOI scheme (doi//)", input: "doi://10.1000/suffix", expectError: true},
		{name: "Invalid: Malformed DOI scheme (doi:/)", input: "doi:/10.1000/suffix", expectError: true},
		{name: "Invalid: URL prefix not starting with 10.", input: "https://doi.org/prefix/suffix", expectError: true},
		{name: "Invalid: DOI scheme prefix not starting with 10.", input: "doi:prefix/suffix", expectError: true},
		{name: "Invalid: URL prefix is just '10.' (empty registrant code)", input: "https://doi.org/10./suffix", expectError: true},
		{name: "Invalid: DOI scheme prefix is just '10.' (empty registrant code)", input: "doi:10./suffix", expectError: true},
		{name: "Invalid: Random string not matching any pattern", input: "ThisIsNotADoi", expectError: true},
		{name: "Invalid: Bare identifier with internal spaces around slash", input: "10.123 / suffix", expectError: true}, // Assumes not allowed
		{name: "Invalid: URL with multiple slashes after domain before prefix", input: "https://doi.org//10.1234/suffix", expectError: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var d DOI // Create a new DOI struct for each test case

			err := d.parseDOI(tc.input)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error for input '%s', but got nil. Parsed to preffix='%s', suffix='%s'", tc.input, d.preffix, d.suffix)
				}
				// Optional: Add assertions here about the state of d.preffix and d.suffix on error,
				// e.g., if they are expected to be cleared:
				if d.preffix != "" || d.suffix != "" {
					t.Errorf("On error for input '%s', expected fields to be cleared, but got preffix: '%s', suffix: '%s'", tc.input, d.preffix, d.suffix)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for input '%s', but got: %v", tc.input, err)
				}
				if d.preffix != tc.expectedPreffix {
					t.Errorf("For input '%s', expected preffix '%s', but got '%s'", tc.input, tc.expectedPreffix, d.preffix)
				}
				if d.suffix != tc.expectedSuffix {
					t.Errorf("For input '%s', expected suffix '%s', but got '%s'", tc.input, tc.expectedSuffix, d.suffix)
				}
			}
		})
	}
}
