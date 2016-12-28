package usecases

import "app"

type cRepo interface {
	app.IDatabase
	OneActiveProduct(interface{}) (*app.Product, error)
	FindActiveProducts(*app.DBFilter) ([]app.Product, error)
	FindActiveProductsByCategory([]interface{}, *app.DBFilter) ([]app.Product, error)
	FindActiveCategories(*app.DBFilter) ([]app.Category, error)
	DeleteProduct(id interface{}) error
	SetProductCategories(*app.Product, []app.Category) error
	SetProductImage(*app.Product, *app.Image) error
}

func NewCatalog(r cRepo) *Catalog {
	return &Catalog{r}
}

type Catalog struct {
	cRepo
}

// todo: validate input
func (cs *Catalog) CreateProduct(f *ProductForm) (*app.Product, error) {
	var p app.Product
	p.Title = f.Title
	p.Description = f.Description
	p.Price = *f.Price
	p.IsActive = *f.IsActive

	if f.Image != "" {
		var img app.Image
		if err := cs.FirstOrInit(&img, app.DBWhere{"public_id": f.Image}); err != nil {
			return nil, err
		}
		p.Image = &img
	}

	if len(f.Categories) > 0 {
		for _, id := range f.Categories {
			var cat app.Category
			if err := cs.One(&cat, id); err != nil {
				return nil, err
			}
			p.AddCategory(cat)
		}
	}

	return &p, cs.Store(&p)
}

// todo: validate input
func (cs *Catalog) UpdateProduct(f *ProductForm) (*app.Product, error) {
	var p app.Product
	p.ID = f.ID

	kv := make(map[string]interface{})

	if f.Title != "" {
		kv["Title"] = f.Title
	}
	if f.Description != "" {
		kv["Description"] = f.Description
	}
	if f.Price != nil {
		kv["Price"] = *f.Price
	}
	if f.IsActive != nil {
		kv["IsActive"] = *f.IsActive
	}

	if f.Image != "" {
		var img app.Image
		if err := cs.FirstOrInit(&img, app.DBWhere{"public_id": f.Image}); err != nil {
			return nil, err
		}

		if p.ImageID == 0 || img.ID != p.ImageID {
			if err := cs.SetProductImage(&p, &img); err != nil {
				return nil, err
			}
		}
	}

	if len(f.Categories) > 0 {
		for _, id := range f.Categories {
			var cat app.Category
			if err := cs.One(&cat, id); err != nil {
				return nil, err
			}
			p.AddCategory(cat)
		}
		if len(p.Categories) > 0 {
			if err := cs.SetProductCategories(&p, p.Categories); err != nil {
				return nil, err
			}
		}
	}

	return &p, cs.UpdateFields(&p, kv)
}

type ProductForm struct {
	ID          int      `json:"-"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       *float32 `json:"price"`
	IsActive    *bool    `json:"isActive"`
	Image       string   `json:"image"`
	Categories  []int    `json:"categories"`
}
