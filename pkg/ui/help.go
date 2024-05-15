package ui

import (
	"fmt"
	"strings"
)

func ColorHeadings(templateString string) string {
	// Headings to colorize
	headings := []string{
		"Usage:",
		"Examples:",
		"Available Commands:",
		"Flags:",
		"Aliases:",
		"Additional Commands:",
	}

	// Replace each heading with its colorized version
	for _, heading := range headings {
		templateString = strings.ReplaceAll(templateString, heading, fmt.Sprintf("%s%s%s%s", White, Bold, heading, Reset))
	}
	return templateString
}
