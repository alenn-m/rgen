[![Go Report Card](https://goreportcard.com/badge/github.com/alenn-m/rgen)](https://goreportcard.com/report/github.com/alenn-m/rgen)
[![codecov](https://codecov.io/gh/alenn-m/rgen/branch/master/graph/badge.svg)](https://codecov.io/gh/alenn-m/rgen)
[![Maintainability](https://api.codeclimate.com/v1/badges/85d8f959f2b9dafc56f3/maintainability)](https://codeclimate.com/github/alenn-m/rgen/maintainability)


## RGEN - Go (GoLang) REST code generator

### Installation

    go get github.com/alenn-m/rgen/v2
    go install github.com/alenn-m/rgen/v2@latest

This will install the generator in `$GOPATH/bin`. You can execute it by calling `rgen` command.

### Getting started

Here's a brief description of each command:
- `build`: This command builds the API from a YAML file.
- `generate`: This command generates API CRUD operations based on the given configuration.
- `help`: This command provides help about any other command.
- `migration`: This command manages database migrations.
- `new`: This command initializes the REST API.


### Usage

Here's a more detailed explanation of each step:
1. `cd` into your projects directory: This changes the current directory to your project's directory.
2. Run `rgen new` and answer the given questions: This initializes the REST API. The questions will ask you about your VSC, VSC domain, and package name. Your answers will determine the root package and name of your project.
3. Go into your newly created project and open `draft.yaml`: This file contains the specification of your project.
4. Run `rgen build` to create REST endpoints: This builds the API from the `draft.yaml` file.
## draft.yaml

After initialization of a project, the `draft.yaml` file will look like this:
```yaml
Models:
  User:
    Properties:
      ApiToken: string
      Email: string
      FirstName: string
      LastName: string
      Password: string
    Validation:
      Email:
      - required
      - email
      FirstName:
      - required
      LastName:
      - required
      Password:
      - required
      - min:5
    Actions: []
    Relationships: {}
    OnlyModel: false
    Public: false
```
All model definitions are under `Models` namespace. 
Each model contains:
- Properties
- Actions
- Relationships
- Validation
- OnlyModel
- Public

#### Properties
This is main portion of `draft.yaml`. In this section you setup all fields for your models.\
Fields are listed in `key:value` format. `key` represents the name of the field, while `value` represents the data type.\
You can use any valid Go data type like: int, int64, string, float64 ...

#### Actions
Actions contain the list of CRUD actions you want to have for a specific model. By default, all CRUD actions are created.\
Possible values are (case-insensitive): `index, show, create, update, delete`
#### Relationships
Relationships follow `key:value` format where `key` is the name of model the current model is related with,\
and `value` is the name of relation.\
Example: 
```
Relationships:
    User: belongsTo
    Post: hasMany
```
Possible values are: `belongsTo, manyToMany, hasMany`.\
**Warning:** IDs are not created automatically. For example if *User* has many *Posts*, you have to add *UserID*\
field to *Post* model.
#### Validation
Validation field contains list of validations for each model property. This field is optional and if you omit this field,
all fields will be optional.

Behind the scene, `rgen` is using https://github.com/go-playground/validator for validations.
Please check the package documentation to learn which fields you can use in your application.
#### OnlyModel
OnlyModel is a boolean value which indicates if you want to only create a model.\
The default value is **false**.
#### Public
Public is a boolean value which indicates if given resource is public or not. The term "public" means that users 
don't have to login to access a particular resource.
The default value is **false**.

### File structure

Here's a brief description of each file:
- `api`: This directory contains the API code.
- `config.yaml`: This file contains the configuration for the project.
- `database`: This directory contains the database code.
- `draft.yaml`: This file contains the specification of your project.
- `go.mod`: This file contains the project's dependencies.
- `main.go`: This is the main entry point for the project.
- `middleware`: This directory contains the middleware code.
- `models`: This directory contains the model code.
- `util`: This directory contains utility code.

### Next steps

Here's a more detailed explanation of each step:
- Edit `.env` file with MySQL credentials and other configurations: This sets up the database connection for the project.
- **Do not remove** config.yaml since this file contains global variables like *Package* which is used when generating new services: This ensures that the project can generate new services correctly.
- Run `go mod tidy` to install all dependencies: This installs all the dependencies that the project needs.
### *Generate* command

Here's a more detailed explanation of how to use the `rgen generate` command:
- `-n, --name string`: This specifies the name of the resource. For example, `--name='ModelName'`.
- `-f, --fields string`: This specifies the list of fields for the resource. For example, `--fields='Title:string, Description:string, UserID:int64'`.
- `-a, --actions string`: This specifies the CRUD actions for the resource. For example, `--actions='index,create,show,update,delete'`.
- `--onlyModel`: This specifies whether to only create a model. The default value is false.
- `--public`: This specifies whether the resource is public. The default value is false.
- [ ] Add support for more databases (currently only *MySQL* is supported)

**PRs are welcome**
### Contributing
We welcome contributions to this project. Here's how you can contribute:
- Submit pull requests: If you've fixed a bug or implemented a new feature, please submit a pull request.
- Report issues: If you've found a bug or have a suggestion for a new feature, please report it as an issue.
- Propose new features: If you have an idea for a new feature, please propose it as an issue.
