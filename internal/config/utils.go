package config

import (
	"path/filepath"
	"strings"
)

func trimExtension(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}
