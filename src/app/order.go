package app

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Order struct {
	Model
	UserID          int        `json:"-"`
	StatusID        int        `json:"-"`
	PaymentMethodID int        `json:"-"`
	PatmentDetails  string     `json:"-"`
	Total           float32    `json:"total"`
	CustomerNote    string     `json:"customerNote"`
	DeliveryTime    *time.Time `json:"deliveryTime"`

	Status        *OrderStatus   `json:"status"`
	Address       *OrderAddress  `json:"address"`
	PaymentMethod *PaymentMethod `json:"paymentMethod"`
	Products      []OrderProduct `json:"items"`
}

func (o *Order) SetTotal() {
	for _, op := range o.Products {
		o.Total += op.Total
	}
}

type OrderProduct struct {
	Model
	OrderID   int     `json:"-"`
	ProductID int     `json:"-"`
	Qty       int     `json:"qty"`
	Price     float32 `json:"price"`
	Total     float32 `json:"total"`
	TaxRate   float32 `json:"taxRate"`
	Options   string  `json:"options"`

	Product *Product `json:"product"`
}

func (op *OrderProduct) GetTotal() float32 {
	return float32(op.Qty) * op.Price
}

func (op *OrderProduct) SetTotal() {
	op.Total = float32(op.Qty) * op.Price
}

type OrderAddress struct {
	Model
	OrderID int `json:"-"`
	AddressBody
}

type OrderHistory struct {
	Model
	OrderID  int
	UserID   int
	StatusID int16
	Note     string
}

type OrderStatus struct {
	Model
	Name        string `json:"name"`
	Description string `json:"description"`
	SortNumber  int16  `json:"-"`
	Status      bool   `json:"-"`
}

type PaymentMethod struct {
	Model
	Name        string `json:"name"`
	Description string `json:"description"`
	SortNumber  int16  `json:"-"`
	Status      bool   `json:"-"`
}

type AddressBody struct {
	Name        string `json:"name"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Tel         string `json:"tel"`
	Tel2        string `json:"tel2"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	City        string `json:"city"`
	District    string `json:"district"`
	Description string `json:"description"`
}

type Address struct {
	Model
	UserID  int  `json:"-"`
	Default bool `json:"default"`
	AddressBody
}

func (a *Address) BeforeCreate(tx *gorm.DB) error {
	var count int

	err := tx.
		Table("addresses").
		Where("user_id=?", a.UserID).
		Where("`default`=?", true).
		Count(&count).
		Error

	if count == 0 {
		a.Default = true
	}

	return err
}

type Product struct {
	Model
	Title       string  `json:"title" fako:"product_name"`
	Description string  `json:"description" gorm:"size:1024" fako:"paragraph"`
	Price       float32 `json:"price"`
	IsActive    bool    `json:"isActive"`

	Categories []Category `gorm:"many2many:pivot_product_category" json:"categories,omitempty"`
	Image      *Image     `json:"defaultImage,omitempty"`
	ImageID    int        `json:"-"`
}

func (p *Product) AddCategory(c Category) {
	for _, v := range p.Categories {
		if c.ID == v.ID {
			return
		}
	}
	p.Categories = append(p.Categories, c)
}

type Category struct {
	Model
	Title       string `json:"title" fako:"title"`
	Description string `json:"description" gorm:"size:1024" fako:"paragraph"`
	IsActive    bool   `json:"isActive"`

	Image    *Image    `json:"image,omitempty"`
	ImageID  int       `json:"-"`
	Products []Product `gorm:"many2many:pivot_product_category" json:"products,omitempty"`
}
