# Contributing to GoCart

## Getting started

Clone this repo

`git clone github.com/alioygur/gocart`

And run the devserver.sh

`./devserver.sh`

that's all.

## Folder structure

This project has own GOPATH.

- $GOPATH
  - `bin`
  - `pkg`
  - `src` 
    - `app` contains the working code for the repository and represents the domain layer
      - `api` the main package.
      - `infra` represents the infrastructure layer.
      - `interface` represents the interface layer.
      - `usecases` represents the usecases layer.

Please respect the following rules and approaches while contribution.

- Uncle Bob's "The Clean Architecture" 
- Dave Cheney, "Don’t just check errors, handle them gracefully"
- SOLID design principles

Many thanks to our contributors: [contributors](https://github.com/alioygur/gocart/graphs/contributors)