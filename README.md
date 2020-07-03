# migrate-pullreqs

**migrate-pullreqs** is a utility that changes the base branch of a
repository's pull requests in bulk.

## Background

Like many other projects, *collectd* used the name `master` for the default
development branch. We made the decision to use `main` instead going forward.
Creating a new branch and changing the default is easy, but there were over 100
open pull requests referencing `master` as their base branch, i.e. the branch
they wanted to get merged into.

This tool was written to change the base branch of all these pull requests
efficiently.

## Usage

*   Get a Github OAuth access token
    * Open the settings menu (your icon in the upper right-hand corner → Settings)
    * Developer Settings → Personal access tokens → Generate new token
    * Add a note that lets you identify the use of this token, e.g. "Migrate pull requests from master to main".
    * Under "Select scopes", check "public\_repo – Access public repositories"
    * Copy the token
*   *Optional:* Implement command line flags for this tool and send a PR ;-)
*   *Alternative:* Update the constants at the beginning of the `main.go` file.
    * Repository owner and repository name
    * Source and destination branch name
    * OAuth access token
*   Compile, run, profit!

## License

ISC License

## Author

Florian Forster &lt;ff at octo.it&gt;
