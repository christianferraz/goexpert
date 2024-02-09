package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	ID    string `gorm:"primaryKey"`
	Name  string
	Price float64
}

func main() {
	dsn := "root:root@tcp(localhost:3306)/goexpert?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Product{})
	// inserir os dados
	db.Create(&Product{
		Name:  "Produto 1",
		Price: 100.00,
	})
	// criar em batch
	db.Create([]Product{
		{Name: "Produto 2", Price: 200.00},
		{Name: "Produto 3", Price: 300.00},
	})
	// buscar produto - select one
	var product Product
	db.First(&product, 1)
	// ou
	db.First(&product, "name = ?", "Produto 1")
	// select all
	var products []Product
	db.Find(&products)
	for _, product := range products {
		println(product.Name)
	}
	// busca os tres primeiros produtos
	db.Limit(3).Find(&products)
	// paginação
	db.Limit(3).Offset(3).Find(&products)
	// where
	db.Where("name = ?", "Produto 1").First(&product)
	db.Where("price > ?", 200).Find(&products)

	// LIKE
	db.Where("name LIKE ?", "%Produto%").Find(&products)

	// atualizar
	db.Model(&product).Update("Price", 200)
	// ou
	db.First(&product, 1)
	product.Price = 500
	db.Save(&product)

}
