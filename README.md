# e-shop-gopack

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/WildEgor/e-shop-gopack)](https://goreportcard.com/report/github.com/WildEgor/e-shop-gopack)

Contains shared code for eShop demo

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)

## Installation
```shell
go get github.com/WildEgor/e-shop-gopack
```

## Usage

- .github
    - workflows
        - release.yml - run semantic release
        - testing.yml - run checks and tests
    - pkg
        - core - core domain shared entities/models etc.
            - dtos
        - libs
            - logger - slog wrapper
            - notifier - client for e-shop-notifier service

## Contributing

Please, use git cz for commit messages!

```shell
git clone https://github.com/WildEgor/e-shop-gopack
cd e-shop-gopack
git checkout -b feature-or-fix-branch
git add .
git cz
git push --set-upstream-to origin/feature-or-fix-branch
```
