package boltdb

import (
	"app"
	"app/interfaces/errs"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
)

func NewRepo(db *storm.DB) *Repo {
	return &Repo{db}
}

type Repo struct {
	db *storm.DB
}

func (r *Repo) Store(data interface{}) error {
	return errs.Wrap(r.db.Save(data))
}

func (r *Repo) One(model interface{}, id int) error {
	return errs.Wrap(r.db.One("ID", id, model))
}

func (r *Repo) OneBy(model interface{}, field string, value interface{}) error {
	return errs.Wrap(r.db.One(field, value, model))
}

func (r *Repo) FindBy(models interface{}, f string, v interface{}, fi *app.DBFilter) error {
	if fi == nil {
		return errs.Wrap(r.db.Find(f, v, models))
	}

	qry := r.db.Select(q.Eq(f, v))
	if fi.Limit > 0 {
		qry = qry.Limit(fi.Limit).Skip(fi.Offset)
	}
	if fi.OrderBy != "" {
		if fi.Reverse {
			qry = qry.OrderBy(fi.OrderBy).Reverse()
		} else {
			qry = qry.OrderBy(fi.OrderBy)
		}
	}

	return errs.Wrap(qry.Find(models))
}

func (r *Repo) ExistsBy(b interface{}, f string, v interface{}) (bool, error) {
	n, err := r.db.Select(q.Eq(f, v)).Count(b)
	return n > 0, errs.Wrap(err)
}

func (r *Repo) UpdateField(b interface{}, f string, v interface{}) error {
	return errs.Wrap(r.db.UpdateField(b, f, v))
}

func (r *Repo) IsNotFoundErr(err error) bool {
	return errs.Cause(err) == storm.ErrNotFound
}

func (r *Repo) where(w app.DBWhere) []q.Matcher {
	var m []q.Matcher
	for f, v := range w {
		m = append(m, q.Eq(f, v))
	}
	return m
}
