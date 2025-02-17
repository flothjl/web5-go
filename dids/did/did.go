package did

import (
	"fmt"
	"regexp"
	"strings"
)

// DID provides a way to parse and handle Decentralized Identifier (DID) URIs
// according to the W3C DID Core specification (https://www.w3.org/TR/did-core/).
type DID struct {
	// URI represents the complete Decentralized Identifier (DID) URI.
	// Spec: https://www.w3.org/TR/did-core/#did-syntax
	URI string

	// URL represents the DID URI + A network location identifier for a specific resource
	// Spec: https://www.w3.org/TR/did-core/#did-url-syntax
	URL string

	// Method specifies the DID method in the URI, which indicates the underlying
	// method-specific identifier scheme (e.g., jwk, dht, key, etc.).
	// Spec: https://www.w3.org/TR/did-core/#method-schemes
	Method string

	// ID is the method-specific identifier in the DID URI.
	// Spec: https://www.w3.org/TR/did-core/#method-specific-id
	ID string

	// Params is a map containing optional parameters present in the DID URI.
	// These parameters are method-specific.
	// Spec: https://www.w3.org/TR/did-core/#did-parameters
	Params map[string]string

	// Path is an optional path component in the DID URI.
	// Spec: https://www.w3.org/TR/did-core/#path
	Path string

	// Query is an optional query component in the DID URI, used to express a request
	// for a specific representation or resource related to the DID.
	// Spec: https://www.w3.org/TR/did-core/#query
	Query string

	// Fragment is an optional fragment component in the DID URI, used to reference
	// a specific part of a DID document.
	// Spec: https://www.w3.org/TR/did-core/#fragment
	Fragment string
}

func (d DID) String() string {
	return d.URL
}

func (d DID) MarshalText() (text []byte, err error) {
	return []byte(d.String()), nil
}

func (d *DID) UnmarshalText(text []byte) error {
	did, err := Parse(string(text))
	if err != nil {
		return err
	}
	*d = did
	return nil
}

// relevant ABNF rules: https://www.w3.org/TR/did-core/#did-syntax
var (
	pctEncodedPattern = `(?:%[0-9a-fA-F]{2})`
	idCharPattern     = `(?:[a-zA-Z0-9._-]|` + pctEncodedPattern + `)`
	methodPattern     = `([a-z0-9]+)`
	methodIDPattern   = `((?:` + idCharPattern + `*:)*(` + idCharPattern + `+))`
	paramCharPattern  = `[a-zA-Z0-9_.:%-]`
	paramPattern      = `;` + paramCharPattern + `+=` + paramCharPattern + `*`
	paramsPattern     = `((` + paramPattern + `)*)`
	pathPattern       = `(/[^#?]*)?`
	queryPattern      = `(\?[^\#]*)?`
	fragmentPattern   = `(\#.*)?`
	didURIPattern     = regexp.MustCompile(`^did:` + methodPattern + `:` + methodIDPattern + paramsPattern + pathPattern + queryPattern + fragmentPattern + `$`)
)

// Parse parses a DID URI in accordance to the ABNF rules specified in the
// specification here: https://www.w3.org/TR/did-core/#did-syntax. Returns
// a DIDURI instance if parsing is successful. Otherwise, returns an error.
func Parse(input string) (DID, error) {
	match := didURIPattern.FindStringSubmatch(input)

	if match == nil {
		return DID{}, fmt.Errorf("invalid DID URI")
	}

	did := DID{
		URI:    "did:" + match[1] + ":" + match[2],
		URL:    input,
		Method: match[1],
		ID:     match[2],
	}

	if len(match[4]) > 0 {
		params := strings.Split(match[4][1:], ";")
		parsedParams := make(map[string]string)
		for _, p := range params {
			kv := strings.Split(p, "=")
			parsedParams[kv[0]] = kv[1]
		}
		did.Params = parsedParams
	}

	if match[6] != "" {
		did.Path = match[6]
	}
	if match[7] != "" {
		did.Query = match[7][1:]
	}
	if match[8] != "" {
		did.Fragment = match[8][1:]
	}

	return did, nil
}
