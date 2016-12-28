package interfaces

import (
	"app"
	"fmt"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/wawandco/fako"
)

// InitDB creates tables
func InitDB(db *gorm.DB) error {
	return db.Set("gorm:table_options", "CHARSET=utf8").AutoMigrate(
		&app.User{},
		&app.Address{},
		&app.Product{},
		&app.Category{},
		&app.Image{},
		&app.Order{},
		&app.OrderProduct{},
		&app.OrderHistory{},
		&app.OrderAddress{},
		&app.OrderStatus{},
		&app.PaymentMethod{},
	).Error
}

func SeedDB(db *gorm.DB) (err error) {

	fns := []func(*gorm.DB) error{
		truncateTables,
		seedMaster,
		seedUsers,
		seedImages,
		seedCategories,
		seedProducts,
	}

	for _, fn := range fns {
		err = fn(db)

		if err != nil {
			return
		}
	}
	return
}

func truncateTables(db *gorm.DB) (err error) {
	tables := []string{
		"addresses",
		"categories",
		"images",
		"order_addresses",
		"order_histories",
		"order_products",
		"order_statuses",
		"orders",
		"payment_methods",
		"pivot_product_category",
		"pivot_product_image",
		"products",
		"users",
	}

	for _, t := range tables {
		err = db.Exec(fmt.Sprintf("TRUNCATE %s", t)).Error
		if err != nil {
			return
		}
	}
	return
}

func dropTables(db *gorm.DB) (err error) {
	tables := []string{
		"addresses",
		"categories",
		"images",
		"order_addresses",
		"order_histories",
		"order_products",
		"order_statuses",
		"orders",
		"payment_methods",
		"pivot_product_category",
		"pivot_product_image",
		"products",
		"users",
	}

	for _, t := range tables {
		err = db.Exec(fmt.Sprintf("DROP TABLE %s", t)).Error
		if err != nil {
			return
		}
	}
	return
}

func seedProducts(db *gorm.DB) (err error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		var p app.Product
		fako.Fill(&p)
		p.Price = r.Float32()
		p.IsActive = true
		db.Create(&p)
	}

	var categories []app.Category
	db.Find(&categories)

	for _, c := range categories {
		var products []app.Product
		db.Order("RAND()").Limit(8).Find(&products)
		db.Model(&c).Association("Products").Append(products)
	}

	var products []app.Product

	db.Find(&products)

	for _, p := range products {
		var images []app.Image
		db.Order("RAND()").Limit(3).Find(&images)
		db.Model(&p).Association("Image").Replace(images[0])

	}

	return
}

func seedCategories(db *gorm.DB) (err error) {

	names := []string{
		"Sağlıklı Beslenme",
		"Zayıflama",
		"Sporcu",
		"Çorba & Atıştırma",
		"Tatlılar",
		"İçecekler",
	}
	for _, name := range names {
		var i app.Image
		db.Order("RAND()").First(&i)

		c := new(app.Category)
		fako.Fill(c)
		c.ImageID = i.ID
		c.Title = name
		db.Create(c)
	}
	return
}

func seedImages(db *gorm.DB) (err error) {
	images := []string{"special_meals_tl9dvc", "Chilli-Pork-Ribeye-Meal-Deal_dhavoe", "chicken-meal_z7eefg", "G15022_KFC_71-big-box-meal-Enviro_1212_isqval"}
	for _, i := range images {
		err = db.Create(&app.Image{PublicID: i, ResourceType: "image"}).Error
		if err != nil {
			return
		}
	}
	return
}

func seedUsers(db *gorm.DB) error {
	u := app.User{
		FirstName:   "Ali",
		LastName:    "OYGUR",
		Email:       "alioygur@gmail.com",
		IsActivated: true,
		IsAdmin:     true,
	}
	u.SetPassword("password")

	return db.Create(&u).Error
}

func seedMaster(db *gorm.DB) (err error) {
	oss := []app.OrderStatus{{Name: "Isleme Alindi"}, {Name: "Tamamlandi"}}
	pms := []app.PaymentMethod{
		{Name: "Nakit Ödeme", SortNumber: 1, Status: true},
		{Name: "Kredi Kartı / Banka Kartı", SortNumber: 2, Status: true},
		{Name: "Ticket Restaurant Kartı ile Ödeme", SortNumber: 3, Status: true},
		{Name: "Ticket Restaurant Çeki ile Ödeme", SortNumber: 4, Status: true},
		{Name: "Sodexo Yemek Kartı ile Ödeme", SortNumber: 5, Status: true},
		{Name: "Sodexo Yemek Çeki ile Ödeme", SortNumber: 6, Status: true},
		{Name: "Multinet ile Ödeme", SortNumber: 7, Status: true},
		{Name: "SetCard ile Ödeme", SortNumber: 8, Status: true},
	}

	for _, os := range oss {
		db.Create(&os)
	}

	for _, pm := range pms {
		db.Create(&pm)
	}

	return
}
