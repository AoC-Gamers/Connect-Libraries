package settingsruntime

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func BoolValidator() ValueValidator {
	return func(value string) error {
		if _, err := strconv.ParseBool(strings.TrimSpace(value)); err != nil {
			return fmt.Errorf("must be bool: %w", err)
		}
		return nil
	}
}

func IntValidator(min, max int) ValueValidator {
	return func(value string) error {
		parsed, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			return fmt.Errorf("must be int: %w", err)
		}
		if parsed < min {
			return fmt.Errorf("must be >= %d", min)
		}
		if max >= min && parsed > max {
			return fmt.Errorf("must be <= %d", max)
		}
		return nil
	}
}

func NonEmptyValidator() ValueValidator {
	return func(value string) error {
		if strings.TrimSpace(value) == "" {
			return errors.New("must not be empty")
		}
		return nil
	}
}
