package main

import (
	"fmt"
	"regexp"
	"strings"
)

type DOI struct {
	preffix string
	suffix  string
}

var doiPattern = regexp.MustCompile(`^(10\.[0-9.]+)\/([^?#]+)(?:[?#].*)?$`)

func (doi *DOI) parseDOI(inputDoiStr string) error {
	// Trim leading/trailing whitespace from the input string.
	doiStr := strings.TrimSpace(inputDoiStr)

	// If the string is empty after trimming, it's not a valid DOI.
	if doiStr == "" {
		return fmt.Errorf("the string is empty")
	}

	var coreDOI string // This will hold the DOI string after stripping presentation prefixes.

	// Check for and strip common DOI presentation prefixes.
	// The order of checks matters for specificity (e.g., "https://" before "http://").
	if strings.HasPrefix(strings.ToLower(doiStr), "https://doi.org/") {
		coreDOI = doiStr[len("https://doi.org/"):]
	} else if strings.HasPrefix(strings.ToLower(doiStr), "http://doi.org/") {
		coreDOI = doiStr[len("http://doi.org/"):]
	} else if strings.HasPrefix(strings.ToLower(doiStr), "https://dx.doi.org/") {
		coreDOI = doiStr[len("https://dx.doi.org/"):]
	} else if strings.HasPrefix(strings.ToLower(doiStr), "http://dx.doi.org/") {
		coreDOI = doiStr[len("http://dx.doi.org/"):]
	} else if strings.HasPrefix(strings.ToLower(doiStr), "doi:") {
		// Handle "doi:", "DOI:", "Doi:", etc.
		// Slice the original string to preserve the case of the actual DOI part.
		coreDOI = doiStr[len("doi:"):]
	} else {
		// If no known prefix is found, assume it's a bare DOI string.
		coreDOI = doiStr
	}

	// After stripping presentation prefixes, the coreDOI might be empty
	// (e.g., if the input was just "https://doi.org/" or "doi:").
	// It also might have leading spaces if the original string was like "doi: 10.xxx/yyy"
	// The regex `^...$` will handle such leading spaces in coreDOI by not matching.
	if strings.TrimSpace(coreDOI) == "" { // Check again if coreDOI itself became empty
		return fmt.Errorf("invalid DOI indetifier: %v", inputDoiStr)
	}

	// The query and fragment must be ignore

	// Attempt to match the coreDOI against our defined DOI pattern.
	matches := doiPattern.FindStringSubmatch(coreDOI)

	// If there's no match or not enough capturing groups, it's not a valid DOI structure.
	// The regex expects 3 matches: the full string, the prefix, and the suffix.
	if len(matches) != 3 {
		return fmt.Errorf("unexpected number of matches for: %v", matches)
	}

	// Extract the prefix (matches[1]) and suffix (matches[2]).
	// matches[0] is the entire string matched by the regex (e.g., "10.1234/abc").
	doi.preffix = matches[1]
	doi.suffix = matches[2]

	// The regex already ensures:
	// - Prefix starts with "10." followed by digits/dots ([0-9.]+)
	// - Suffix is not empty (.+)
	return nil
}
