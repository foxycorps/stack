package ui

import "fmt"

const (
	Reset         = "\033[0m"
	Bold          = "\033[1m"
	Underline     = "\033[4m"
	StrikeThrough = "\033[9m"
	White         = "\033[97m"
	Dim           = "\033[2m"

	// Pastel colors
	PastelRed    = "\033[38;5;203m"
	PastelGreen  = "\033[38;5;156m"
	PastelBlue   = "\033[38;5;117m"
	PastelYellow = "\033[38;5;221m"
	PastelPurple = "\033[38;5;183m"
	PastelCyan   = "\033[38;5;159m"
	PastelPink   = "\033[38;5;218m"
	PastelOrange = "\033[38;5;208m"
)

func CreateHyperlink(url, text string) string {
	return fmt.Sprintf("\033]8;;%s\a%s\033]8;;\a", url, text)
}
