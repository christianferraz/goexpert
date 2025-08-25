//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/christianferraz/goexpert/19-DI/product"
	"github.com/google/wire"
)

var setRepositoryDependency = wire.NewSet(
	product.NewProductRepository,
	// toda vez que for necessário usar o ProductRepositoryInterface,
	// o Wire irá fornecer o ProductRepository
	wire.Bind(new(product.ProductRepositoryInterface), new(*product.ProductRepository)),
)

func NewUseCase(db *sql.DB) *product.ProductUseCase {
	wire.Build(
		setRepositoryDependency,
		product.NewProductUseCase,
	)
	return &product.ProductUseCase{}
}
