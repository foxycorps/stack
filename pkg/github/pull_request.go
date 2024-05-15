package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v60/github"
)

func (c *Client) CreatePullRequest(head, base, title, description string) (*github.PullRequest, error) {
	owner, err := c.repo.Owner()
	if err != nil {
		return nil, err
	}

	repo, err := c.repo.Name()
	if err != nil {
		return nil, err
	}

	newPR := &github.NewPullRequest{
		Title:               github.String(title),
		Body:                github.String(description),
		Head:                github.String(head),
		Base:                github.String(base),
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := c.client.PullRequests.Create(context.Background(), owner, repo, newPR)
	return pr, err
}

func (c *Client) FindPrForBranch(branchName string) (*github.PullRequest, error) {
	owner, err := c.repo.Owner()
	if err != nil {
		return nil, err
	}

	repo, err := c.repo.Name()
	if err != nil {
		return nil, err
	}

	prs, _, err := c.client.PullRequests.List(context.Background(), owner, repo, nil)
	if err != nil {
		return nil, err
	}

	for _, pr := range prs {
		if pr.GetHead().GetRef() == branchName {
			return pr, nil
		}
	}

	return nil, nil
}

func (c *Client) WriteReviewComment(pr *github.PullRequest, comment string, event, sha string) error {
	owner, err := c.repo.Owner()
	if err != nil {
		return err
	}

	repo, err := c.repo.Name()
	if err != nil {
		return err
	}

	review := &github.PullRequestReviewRequest{
		CommitID: github.String(sha),
		Body:     github.String(comment),
		Event:    &event,
	}

	_, _, err = c.client.PullRequests.CreateReview(context.Background(), owner, repo, pr.GetNumber(), review)
	return err
}

func (c *Client) IsPRAwaitingReview(pr *github.PullRequest) bool {
	return pr.GetState() == "open" && !pr.GetDraft() && pr.GetReviewComments() == 0
}

func (c *Client) IsPRApproved(pr *github.PullRequest) (bool, string, error) {
	owner, err := c.repo.Owner()
	if err != nil {
		return false, "", err
	}

	repo, err := c.repo.Name()
	if err != nil {
		return false, "", err
	}

	reviews, _, err := c.client.PullRequests.ListReviews(context.Background(), owner, repo, pr.GetNumber(), nil)
	if err != nil {
		return false, "", err
	}

	for _, review := range reviews {
		fmt.Println("Review state: ", review.GetState())
		if review.GetState() == "APPROVED" {
			return true, review.GetUser().GetLogin(), nil
		}
	}

	return false, "", nil
}
