package product

type ProductUseCase struct {
	repository ProductRepositoryInterface
}

func NewProductUseCase(repository ProductRepositoryInterface) *ProductUseCase {
	return &ProductUseCase{repository}
}

func (u *ProductUseCase) GetProduct(id string) (*Product, error) {
	return u.repository.GetProduct(id)
}
