## RGEN - Go (GoLang) REST code generator

### Installation

    go get -u github.com/alenn-m/rgen 

This will install the generator in `$GOPATH/bin`. You can execute it by calling `rgen` command.

### Getting started

Run `rgen -h` to see a list of all commands:

    Usage:
       [command]
    
    Available Commands:
      build       Builds API from YAML file
      generate    Generates API CRUD with given configuration
      help        Help about any command
      new         Initializes the REST API

### Usage

1. `cd` into the directory where you want to create new project
2. Run `rgen new` and answer the given questions (it will ask you which VSC you use, VSC domain and package name).
Your answers will determine the `gomodule` package and nam of your project.
3. `cd` into your newly create project and open `draft.yaml`. This file contains specification of your project.
4. Run `rgen build` to create REST endpoints.
