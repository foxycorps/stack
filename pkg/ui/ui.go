package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func SelectBranch(branches []string) (string, error) {

	var branchOptions []huh.Option[string]
	var selectedBranch string

	for _, branch := range branches {
		branchOptions = append(branchOptions, huh.NewOption(branch, branch))
	}

	theme := huh.ThemeCatppuccin()
	form := huh.NewSelect[string]().
		Title("Select a branch").
		Options(
			branchOptions...,
		).
		Value(&selectedBranch).
		WithTheme(theme)

	err := form.Run()

	return selectedBranch, err
}

func SelectStack(branches []string) (string, error) {
	var stackOptions []huh.Option[string]
	var selectedBranch string

	for _, branch := range branches {
		stackOptions = append(stackOptions, huh.NewOption(branch, branch))
	}

	theme := huh.ThemeCatppuccin()
	form := huh.NewSelect[string]().
		Title("Select a Stack from remote").
		Options(
			stackOptions...,
		).
		Value(&selectedBranch).
		WithTheme(theme)

	err := form.WithTheme(theme).Run()

	return selectedBranch, err
}

type PRDetails struct {
	Title     string
	Body      string
	Labels    []string
	Reviewers []string
	Draft     bool
	Create    bool
}

func CreatePR(branchName, commitTitle string) *PRDetails {

	var title, body string
	var tmpLabels, tmpReviewers string
	var draft, create bool

	var labels, reviewers []string

	theme := huh.ThemeBase16()
	theme.Focused.FocusedButton.Background(lipgloss.Color("117"))
	theme.Focused.TextInput.Cursor.Foreground(lipgloss.Color("117"))
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().Title(fmt.Sprintf("Open a pull request for '%s%s%s'?", White, branchName, Reset)).Value(&create),
		),
		huh.NewGroup(
			huh.NewInput().Inline(true).
				Title("Enter a title:").
				Placeholder(commitTitle).
				Value(&title),
			huh.NewText().
				Title("Enter a body:").
				Value(&body),
			huh.NewInput().Inline(true).
				Title("Enter reviewers (comma-separated):").
				Value(&tmpReviewers),
			huh.NewInput().Inline(true).
				Title("Enter labels (comma-separated):").
				Value(&tmpLabels),
		).
			WithHideFunc(func() bool {
				return !create
			}),
	).WithTheme(theme)

	form.WithTheme(theme).Run()

	labels = strings.Split(tmpLabels, ",")
	reviewers = strings.Split(tmpReviewers, ",")

	if !create {
		return &PRDetails{
			Create: create,
		}
	}

	if title == "" {
		title = commitTitle
	}

	return &PRDetails{
		Title:     title,
		Body:      body,
		Labels:    labels,
		Reviewers: reviewers,
		Draft:     draft,
		Create:    true,
	}

}
