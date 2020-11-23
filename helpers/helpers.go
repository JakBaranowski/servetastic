package helpers

import (
	"log"
	"path/filepath"
	"strings"
)

func GlobDir(dir string, pattern string) []string {
	log.Printf("Globbing files at %v", dir)
	files, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		panic("Failed globing files")
	}
	log.Printf("Files \"%v\" globbed at %v", files, dir)
	return files
}

func StripDomainName(s string) string {
	// Strip the domain name
	return s
}

func StripExtension(s string) string {
	return strings.TrimSuffix(s, filepath.Ext(s))
}

func StripPathAndExtension(s string) string {
	return strings.TrimSuffix(filepath.Base(s), filepath.Ext(s))
}
