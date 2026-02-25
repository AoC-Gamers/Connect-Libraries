package migrate

import (
	"fmt"
	"regexp"
)

var identifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

func quoteIdentifier(identifier string) (string, error) {
	if !identifierPattern.MatchString(identifier) {
		return "", fmt.Errorf("invalid SQL identifier: %q", identifier)
	}

	return `"` + identifier + `"`, nil
}
