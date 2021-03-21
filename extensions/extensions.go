package extensions

import (
	"fmt"
	"strings"
)

func EnableExtensionsInterface() string {
	ifce := []string{
		"Enable extension",
		fmt.Sprintf("%s", "ad"),
	}
	return strings.Join(ifce, "\n")
}
