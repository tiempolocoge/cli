## Contributing

[legal]: https://help.github.com/articles/github-terms-of-service/#6-contributions-under-repository-license
[license]: ../LICENSE
[code-of-conduct]: CODE-OF-CONDUCT.md
[bug issues]: https://github.com/cli/cli/issues?q=is%3Aopen+is%3Aissue+label%3Abug
[feature request issues]: https://github.com/cli/cli/issues?q=is%3Aopen+is%3Aissue+label%3Aenhancement

Hi! Thanks for your interest in contributing to the GitHub CLI!

We accept pull requests for bug fixes and features where we've discussed the approach in an issue and given the go-ahead for a community member to work on it. We'd also love to hear about ideas for new features as issues.

Please do:

* Check existing issues to verify that the [bug][bug issues] or [feature request][feature request issues] has not already been submitted.
* Open an issue if things aren't working as expected.
* Open an issue to propose a significant change.
* Open a pull request to fix a bug.
* Open a pull request to fix documentation about a command.
* Open a pull request for an issue with the `help-wanted` label and leave a comment claiming it.

Please avoid:

* Opening pull requests for issues marked `needs-design`, `needs-investigation`, `needs-user-input`, or `blocked`.
* Adding installation instructions specifically for your OS/package manager.
* Opening pull requests for any issue marked `core`. These issues require additional context from
  the core CLI team at GitHub and any external pull requests will not be accepted.

## Building the project

Prerequisites:
- Go 1.13+ for building the binary
- Go 1.15+ for running the test suite

Build with: `make` or `go build -o bin/gh ./cmd/gh`

Run the new binary as: `./bin/gh`

Run tests with: `make test` or `go test ./...`

## Submitting a pull request

1. Create a new branch: `git checkout -b my-branch-name`
1. Make your change, add tests, and ensure tests pass
1. Submit a pull request: `gh pr create --web`

Contributions to this project are [released][legal] to the public under the [project's open source license][license].

Please note that this project adheres to a [Contributor Code of Conduct][code-of-conduct]. By participating in this project you agree to abide by its terms.

We generate manual pages from source on every release. You do not need to submit pull requests for documentation specifically; manual pages for commands will automatically get updated after your pull requests gets accepted.

## Resources

- [How to Contribute to Open Source](https://opensource.guide/how-to-contribute/)
- [Using Pull Requests](https://help.github.com/articles/about-pull-requests/)
- [GitHub Help](https://help.github.com)
