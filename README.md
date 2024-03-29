[![Go Report Card](https://goreportcard.com/badge/github.com/alenn-m/rgen)](https://goreportcard.com/report/github.com/alenn-m/rgen)
[![codecov](https://codecov.io/gh/alenn-m/rgen/branch/master/graph/badge.svg)](https://codecov.io/gh/alenn-m/rgen)
[![Maintainability](https://api.codeclimate.com/v1/badges/85d8f959f2b9dafc56f3/maintainability)](https://codeclimate.com/github/alenn-m/rgen/maintainability)


## RGEN - Go (GoLang) REST code generator

### Installation

    go get github.com/alenn-m/rgen/v2
    go install github.com/alenn-m/rgen/v2@latest

This will install the generator in `$GOPATH/bin`. You can execute it by calling `rgen` command.

### Getting started

Run `rgen -h` to see a list of all commands:

    Usage:
       [command]
    
    Available Commands:
      build       Builds API from YAML file
      generate    Generates API CRUD with given configuration
      help        Help about any command
      migration   Manages database migrations
      new         Initializes the REST API


### Usage

1. `cd` into your projects directory
2. Run `rgen new` and answer the given questions (it will ask you which VSC you use, VSC domain and package name).
Your answers will determine the root package and name of your project.
3. Go into your newly created project and open `draft.yaml`. This file contains specification of your project.
4. Run `rgen build` to create REST endpoints.

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

```
├── api
│   └── auth
│       ├── controller.go
│       ├── repositories
│       │   └── mysql
│       │       └── auth.go
│       ├── repository.go
│       └── transport.go
├── config.yaml
├── database
│   └── seeds
│       ├── DatabaseSeeder.go
│       └── UserSeeder.go
├── draft.yaml
├── go.mod
├── main.go
├── middleware
│   ├── AuthMiddleware.go
│   └── ExampleMiddleware.go
├── models
│   └── Base.go
└── util
    ├── auth
    │   ├── auth.go
    │   └── interface.go
    ├── cache
    │   ├── memory
    │   │   └── memory.go
    │   └── service.go
    ├── paginate
    │   └── paginate.go
    ├── req
    │   └── req.go
    ├── resp
    │   └── response.go
    └── validators
        ├── Base.go
        ├── Equals.go
        ├── RecordExists.go
        ├── RecordsExists.go
        └── Unique.go

```

### Next steps

- Edit `.env` file with MySQL credentials and other configurations
- **Do not remove** config.yaml since this file contains global variables like *Package* which is used when generating new services.
- Run `go mod tidy` to install all dependencies, and you're good to go.

### *Generate* command

`rgen` allows you to create single service as well, this is useful when you want to update existing project.<br/><br/>
**WARNING:** `rgen` relies on various markers and file paths to add new services, if you want to use `rgen generate` command,
then **do not remove** markers *[services], [protected routes] and [public routes]* inside `main.go`.
```
rgen generate -h
-------------------------------------------
Generates API CRUD with given configuration

Usage:
   generate [flags]

Flags:
  -a, --actions string   CRUD actions --actions='index,create,show,update,delete'
  -f, --fields string    List of fields (required) --fields='Title:string, Description:string, UserID:int64'
  -h, --help             help for generate
  -n, --name string      Resource name (required) --name='ModelName'
      --onlyModel        Create only model (default = false)
      --public           Public resource (default = false)
```
**Example:** To create a *Comment* resource with *title, body* and *user_id* fields run the following command:
```
rgen generate -n "Comment" -f "title:string#required|min:5, body:string#required, user_id:int64#required" -a "index, create, delete"
```
### TODO
- [ ] Add support for more databases (currently only *MySQL* is supported)

**PRs are welcome**
