package controller

const CONTROLLER_HEADER = `
package {{Package}}

import (
    "context"

    "{{Root}}/models"
)`
const CONTROLLER_INDEX = `func (u *{{Controller}}) Index(c context.Context) ([]models.{{Model}}, error) {
    items, err := u.db.ListAll()
    if err != nil {
        return nil, err
    }

    return items, nil
}`

const CONTROLLER_CREATE = `func (u *{{Controller}}) Store(c context.Context, r *StoreReq) (int64, error) {
    id, err := u.db.Insert(models.{{Model}}{
        {{Fields}}
    })

    return id, err
}`

const CONTROLLER_SHOW = `func (u *{{Controller}}) Show(c context.Context, id int64) (*models.{{Model}}, error) {
    return u.db.FindByID(id)
}`

const CONTROLLER_UPDATE = `func (u *{{Controller}}) Update(c context.Context, r *UpdateReq, id int64) error {
    return u.db.Update(models.{{Model}}{
        ID:       id,
        {{Fields}}
    })
}`

const CONTROLLER_DELETE = `func (u *{{Controller}}) Delete(c context.Context, id int64) error {
    return u.db.Delete(id)
}`
