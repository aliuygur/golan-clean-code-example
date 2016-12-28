package gormdb

import (
	"app"
	"app/interfaces/errs"
)

func NewCatalog(r *Repo) *Catalog {
	return &Catalog{r}
}

type Catalog struct {
	*Repo
}

func (cr *Catalog) OneActiveProduct(id interface{}) (*app.Product, error) {
	var p app.Product
	if err := cr.db.Preload("Image").Preload("Categories").First(&p, "id=? AND is_active=?", id, true).Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return &p, nil
}

func (cr *Catalog) FindActiveProducts(f *app.DBFilter) ([]app.Product, error) {
	var ps []app.Product

	if err := cr.db.Preload("Image").Find(&ps, "is_active=?", true).Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return ps, nil
}

func (cr *Catalog) FindActiveProductsByCategory(ids []interface{}, f *app.DBFilter) ([]app.Product, error) {
	var (
		pids []int
		ps   []app.Product
	)

	if err := cr.db.Table("pivot_product_category").
		Where("category_id in (?)", ids).
		Group("product_id").
		Pluck("product_id", &pids).Error; err != nil {
		return nil, errs.Wrap(err)
	}

	if err := cr.db.Preload("Image").Where("id in (?) AND is_active=?", pids, true).Find(&ps).Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return ps, nil
}

func (cr *Catalog) FindActiveCategories(f *app.DBFilter) ([]app.Category, error) {
	var cs []app.Category

	if err := cr.db.Preload("Image").Find(&cs, "is_active=?", true).Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return cs, nil
}

func (cr *Catalog) DeleteProduct(id interface{}) error {
	// todo: send 404 error on model not found
	var p app.Product
	if err := cr.Repo.One(&p, id); err != nil {
		return err
	}

	tx := cr.db.Begin()

	// clear Associations
	if err := tx.Model(&p).Association("Categories").Clear().Error; err != nil {
		tx.Rollback()
		return errs.Wrap(err)
	}

	// todo: also delete image files
	if err := tx.Model(&p).Association("Images").Clear().Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&p).Error; err != nil {
		tx.Rollback()
		return errs.Wrap(err)
	}

	tx.Commit()
	return nil
}

func (cr *Catalog) SetProductCategories(p *app.Product, cs []app.Category) error {
	return errs.Wrap(cr.db.Model(p).Association("Categories").Replace(cs).Error)
}

func (cr *Catalog) SetProductImage(p *app.Product, img *app.Image) error {
	return errs.Wrap(cr.db.Model(p).Association("Image").Replace(img).Error)
}
