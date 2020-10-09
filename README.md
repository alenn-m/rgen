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

1. `cd` into your projects directory
2. Run `rgen new` and answer the given questions (it will ask you which VSC you use, VSC domain and package name).
Your answers will determine the root package and name of your project.
3. Go into your newly create project and open `draft.yaml`. This file contains specification of your project.
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
    SkipController: false
    Validation:
      Email:
      - Required
      - Email
      FirstName:
      - Required
      LastName:
      - Required
    Actions: []
    Relationships: {}
```
All model definitions are under `Models` namespace. 
Each model contains:
- Properties
- Validation
- Actions
- Relationships
- SkipController

#### Properties
This is main portion of `draft.yaml`. In this section you setup all fields for you models.\
Fields are listed in `key:value` format. `key` represents the name of the field, while `value` represents the data type.\
Since `rgen` uses `gorm` as ORM, you can use any valid Go data type like: int, int64, string, float64 ...
#### Validation
Validation follows `key:[]value` format, in other words for single field you can add multiple rules.\
`rgen` is using `github.com/go-ozzo/ozzo-validation` package for validation so please check the package documentation for more information.
#### Actions
Actions contains the list of CRUD actions you want to have for specific model. By default, all CRUD actions are created.\
Possible values are (not case-sensitive): `index, show, create, update, delete`
#### Relationships
Relationships follows `key:value` format where `key` is the name of model the current model is related with,\
and `value` is the name of relation.\
Example: 
```
Relationships:
    User: belongsTo
    Post: hasMany
```
Possible values are: `belongsTo, belongsToMany, hasMany`.\
**Warning:** IDs are not created automatically. For example if *User* has many *Posts*, you have to add *UserID*\
field to *Post* model.
#### SkipController
SkipController is a boolean value which indicates if you want to create a controller for specific model.\
The default value is **true**.\
If the SkipController is false then only a model will be created.

### Next steps

- Edit `.env` file with MySQL credentials and other configurations
- **Do not remove** config.yaml since this file contains global variables like *Package* which is used when generating new services.
- Run `go mod tidy` to install all dependencies, and you're good to go.

### *Generate* command

`rgen` allows you to create single service as well, this is useful when you want to update existing project.
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
```
**Example:** To create a *Comment* resource with *title, body* and *user_id* fields run the following command:
```
rgen generate -n "Comment" -f "title:string, body:string, user_id:int64" -a "index, create, delete"
```
### TODO
- [ ] Unit/integration tests
- [ ] Code refactoring
- [ ] Better documentation
- [ ] Add support for more databases (currently only *MySQL* is supported)

**PRs are welcome**
