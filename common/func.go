package common

import (
	"strings"
	"errors"
)

// pascalcase 2 camelcase
func TranslatePascal(name string) (string, error) {
	if len(name) <= 0 {
		return "", errors.New("Param is empty.")
	}
	name = strings.Title(name)
	for strings.Index(name, "_") != -1 {
		idx := strings.Index(name, "_")
		name = name[:idx] + strings.Title(name[idx+1:])
	}

	return name, nil
}


