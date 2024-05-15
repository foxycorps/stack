package github

import (
	"context"
	"fmt"
	"os"

	"github.com/foxycorps/stack/pkg/git"
	"github.com/foxycorps/stack/pkg/keyring"
	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

type Client struct {
	client   *github.Client
	repo     *git.Repository
	username string
	email    string
}

func NewClient(repo *git.Repository) *Client {
	ctx := context.Background()
	token, err := keyring.Get()
	if err != nil {
		fmt.Printf("failed to get token: %v\n", err)
		os.Exit(1)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	username, _ := repo.GetUsername()
	email, _ := repo.GetEmail()

	return &Client{client: client, repo: repo, username: username, email: email}
}

func (gh *Client) ValidateToken() bool {
	ctx := context.Background()
	_, _, err := gh.client.Users.Get(ctx, "")
	return err == nil
}
