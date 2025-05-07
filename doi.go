package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"unicode"
)

type DOI struct {
	preffix string
	suffix  string
}

func (doi *DOI) parseDOI(doiID string) error {

	doiID = strings.TrimSpace(doiID)

	if strings.ContainsFunc(doiID, unicode.IsSpace) {
		return errors.New("DOI strings (or URL) cannot contains blank space")
	}

	if len(doiID) == 0 {
		return errors.New("expected a DOI code, but received an empty string")
	}

	link, err := url.Parse(doiID)

	doiParts := []string{}

	if err == nil && link.Hostname() == "" {

		doiParts = strings.Split(link.Path, "/")

	} else if err == nil && (strings.ToLower(link.Hostname()) == "doi.org" || strings.ToLower(link.Hostname()) == "dx.doi.org") {

		if strings.ToLower(link.Scheme) != "http" && strings.ToLower(link.Scheme) != "https" {
			return fmt.Errorf("the scheme \"%v\" in not permited in a DOI identifier", link.Scheme)
		}

		doiParts = strings.Split(link.Path[1:], "/")

	} else if err == nil && link.Hostname() != "" {

		return fmt.Errorf("the hostame \"%v\" in not permited in a DOI identifier", link.Hostname())

	} else {
		doiParts = strings.Split(doiID, "/")
	}

	if len(doiParts) != 2 {
		return fmt.Errorf("the preffix \"%v\" is malformed", doiParts)
	}

	//fmt.Printf("Good DOI identifier: preffix: %v; suffix: %v\n", doiParts[0], doiParts[1])

	if len(doiParts[0]) < 7 {
		return fmt.Errorf("malformed DOI identifier ( missing or too short preffix): preffix: %v; suffix: %v", doiParts[0], doiParts[1])
	}

	if doiParts[0][0:3] != "10." {
		return fmt.Errorf("malformed DOI identifier (preffix must start with \"10.0\"): preffix: %v; suffix: %v", doiParts[0], doiParts[1])
	}

	if len(doiParts[1]) == 0 {
		return fmt.Errorf("malformed DOI identifier ( missing suffix): preffix: %v; suffix: %v", doiParts[0], doiParts[1])
	}

	doi.preffix = doiParts[0]
	doi.suffix = doiParts[1]
	return nil
}
