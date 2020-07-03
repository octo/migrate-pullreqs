package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/octo/retry"
	"golang.org/x/oauth2"
)

// accessToken holds an OAuth bearer token.
const accessToken = "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"

const (
	sourceBranch = "master"
	destBranch   = "main"
	owner        = "collectd"
	repository   = "collectd"
)

func newClient(ctx context.Context) *github.Client {
	t := &retry.Transport{
		RoundTripper: &oauth2.Transport{
			Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken}),
			Base:   http.DefaultTransport,
		},
	}

	return github.NewClient(&http.Client{
		Transport: t,
	})
}

func migrateOne(ctx context.Context, c *github.Client, pr *github.PullRequest) error {
	baseBranch := pr.GetBase().GetRef()
	if baseBranch != sourceBranch {
		return nil
	}

	fmt.Printf("Updating %s ... ", pr.GetHTMLURL())

	repo := pr.GetBase().GetRepo()
	_, _, err := c.PullRequests.Edit(ctx, repo.GetOwner().GetLogin(), repo.GetName(), pr.GetNumber(), &github.PullRequest{
		Base: &github.PullRequestBranch{
			Ref: github.String(destBranch),
		},
	})
	if err != nil {
		return fmt.Errorf("Edit(%q, %q, %d, base.ref=%q): %w", repo.GetOwner(), repo.GetName(), pr.GetNumber(), destBranch, err)
	}

	fmt.Println("done")
	return nil
}

func migrateAll(ctx context.Context, c *github.Client) error {
	opts := github.PullRequestListOptions{
		Base: sourceBranch,
	}
	for {
		prs, res, err := c.PullRequests.List(ctx, owner, repository, &opts)
		if err != nil {
			return fmt.Errorf("List(%q, %q, %v): %w", owner, repository, &opts, err)
		}

		for _, pr := range prs {
			if err := migrateOne(ctx, c, pr); err != nil {
				return err
			}
		}

		if res.NextPage == 0 {
			break
		}
		opts.Page = res.NextPage
	}

	return nil
}

func main() {
	ctx := context.Background()

	c := newClient(ctx)

	if err := migrateAll(ctx, c); err != nil {
		log.Fatal(err)
	}
}
