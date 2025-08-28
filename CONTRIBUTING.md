# Contributing to bank-data

First off, thank you for considering contributing to bank-data! It's people like you that make open source such a great community.

## Where do I go from here?

If you've noticed a bug or have a feature request, [make one](https://github.com/SamyRai/bank-data/issues/new)! It's generally best if you get confirmation of your bug or approval for your feature request this way before starting to code.

### Fork & create a branch

If this is something you think you can fix, then [fork the repository](https://github.com/SamyRai/bank-data/fork) and create a branch with a descriptive name.

A good branch name would be (where issue #123 is the ticket you're working on):

```sh
git checkout -b 123-add-a-bigger-boat
```

### Get the code

```sh
git clone https://github.com/your-username/bank-data.git
cd bank-data
```

### Run the tests

```sh
make test
```

This will run all the tests, so you can see if your changes have introduced any regressions.

### Make your changes

Make your changes to the code. Please follow the coding style of the project.

### Commit your changes

Commit your changes with a descriptive commit message.

```sh
git commit -m "feat: add a bigger boat"
```

### Push your changes

Push your changes to your fork.

```sh
git push origin 123-add-a-bigger-boat
```

### Open a pull request

Open a pull request to the `master` branch of the `SamyRai/bank-data` repository.

## Development

### Prerequisites

- Go 1.24 or later

### Getting Started

1.  Fork the repository.
2.  Clone your fork: `git clone https://github.com/your-username/bank-data.git`
3.  Create a new branch: `git checkout -b my-feature-branch`
4.  Make your changes.
5.  Run the tests: `make test`
6.  Commit your changes: `git commit -am 'Add some feature'`
7.  Push to the branch: `git push origin my-feature-branch`
8.  Submit a pull request.

### Coding Style

Please follow the existing coding style. We use `gofmt` to format our code. You can run `make fmt` to format your code before committing.

We also use `golangci-lint` for linting. You can run `make lint` to check your code for any linting issues.

### Testing

Please add tests for any new features or bug fixes. We use the built-in Go testing framework.

### Issues

If you find a bug or have a feature request, please [open an issue](https://github.com/SamyRai/bank-data/issues/new). Please provide as much information as possible, including the version of Go you are using, the version of the library you are using, and a code sample that reproduces the issue.
