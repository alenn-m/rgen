package model

const TEMPLATE = `
package models

type {{Model}}ID int64

type {{Model}} struct {
	{{Fields}}

    Base
}
`

const MIGRATIONS_TEMPLATE = `
package main

import (
    "{{Root}}/models"

    "github.com/jinzhu/gorm"
)

func RunMigrations(db *gorm.DB) {
    db.AutoMigrate({{Models}})
}
`
