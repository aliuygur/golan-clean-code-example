package gormdb

import (
	"app"
	"app/interfaces/errs"

	"github.com/jinzhu/gorm"
)

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db}
}

type Repo struct {
	db *gorm.DB
}

func (r *Repo) Store(data interface{}) error {
	return errs.Wrap(r.db.Create(data).Error)
}

func (r *Repo) Save(model interface{}) error {
	return errs.Wrap(r.db.Save(model).Error)
}

func (r *Repo) One(model interface{}, id interface{}) error {
	return errs.Wrap(r.db.First(model, id).Error)
}

func (r *Repo) OneBy(model interface{}, w app.DBWhere) error {
	return errs.Wrap(r.db.First(model, r.where(w)).Error)
}

func (r *Repo) FindBy(ms interface{}, w app.DBWhere, fi *app.DBFilter) error {
	if fi == nil {
		return errs.Wrap(r.db.Where(r.where(w)).Find(ms).Error)
	}
	qry := r.db.Where(r.where(w))

	qry = r.filter(qry, fi)

	return errs.Wrap(qry.Find(ms).Error)
}

func (r *Repo) FirstOrInit(m interface{}, w app.DBWhere) error {
	return errs.Wrap(r.db.FirstOrInit(m, r.where(w)).Error)
}

func (r *Repo) ExistsBy(m interface{}, w app.DBWhere) (bool, error) {
	var n uint
	err := r.db.Model(m).Where(r.where(w)).Count(&n).Error
	return n > 0, errs.Wrap(err)
}

func (r *Repo) UpdateField(m interface{}, f string, v interface{}) error {
	return errs.Wrap(r.db.Model(m).Update(f, v).Error)
}

func (r *Repo) UpdateFields(m interface{}, kv map[string]interface{}) error {
	return errs.Wrap(r.db.Model(m).Updates(kv).Error)
}

func (r *Repo) IsNotFoundErr(err error) bool {
	return errs.Cause(err) == gorm.ErrRecordNotFound
}

func (r *Repo) where(w app.DBWhere) map[string]interface{} {
	return w
}

func (r *Repo) filter(qry *gorm.DB, fi *app.DBFilter) *gorm.DB {
	if fi.Preload != nil {
		for _, p := range fi.Preload {
			qry = qry.Preload(p)
		}
	}

	if fi.Limit > 0 {
		qry = qry.Limit(fi.Limit).Offset(fi.Offset)
	}

	if fi.OrderBy != "" {
		if fi.Reverse {
			qry = qry.Order(fi.OrderBy + " desc")
		} else {
			qry = qry.Order(fi.OrderBy)
		}
	}
	return qry
}
