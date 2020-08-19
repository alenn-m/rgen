package controller

const TEMPLATE = `
package {{Package}}

import (
    "context"

    "{{Root}}/models"
)

func (u *{{Controller}}) Index(c context.Context) ([]models.{{Model}}, error) {
    items, err := u.db.ListAll()
    if err != nil {
        return nil, err
    }

    return items, nil
}

func (u *{{Controller}}) Store(c context.Context, r *StoreReq) (int64, error) {
    id, err := u.db.Insert(models.{{Model}}{
        {{Fields}}
    })

    return id, err
}

func (u *{{Controller}}) Show(c context.Context, id int64) (*models.{{Model}}, error) {
    return u.db.FindByID(id)
}

func (u *{{Controller}}) Update(c context.Context, r *UpdateReq, id int64) error {
    return u.db.Update(models.{{Model}}{
        ID:       id,
        {{Fields}}
    })
}

func (u *{{Controller}}) Delete(c context.Context, id int64) error {
    return u.db.Delete(id)
}
`
