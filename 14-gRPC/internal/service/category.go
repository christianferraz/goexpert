package service

import (
	"context"
	"io"

	"github.com/christianferraz/goexpert/14-gRPC/internal/database"
	"github.com/christianferraz/goexpert/14-gRPC/internal/pb"
)

type CategoryService struct {
	CategoryDB database.Category
	pb.UnimplementedCategoryServiceServer
}

func NewCategoryService(db database.Category) *CategoryService {
	return &CategoryService{CategoryDB: db}
}

func (c *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.Category, error) {
	category, err := c.CategoryDB.CreateCategory(in.Name, in.Description)
	if err != nil {
		return nil, err
	}
	return &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (c *CategoryService) ListCategories(ctx context.Context, in *pb.Blank) (*pb.CategoryList, error) {
	categories, err := c.CategoryDB.ListAllCategories()
	if err != nil {
		return nil, err
	}
	var Categories []*pb.Category
	for _, category := range categories {
		Categories = append(Categories, &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
	}
	return &pb.CategoryList{Categories: Categories}, nil
}

func (c *CategoryService) GetCategory(ctx context.Context, in *pb.CategoryGetRequest) (*pb.Category, error) {
	category, err := c.CategoryDB.GetCategoryByID(in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (c *CategoryService) CreateCategoryStream(in pb.CategoryService_CreateCategoryStreamServer) error {
	categories := &pb.CategoryList{}
	for {
		category, err := in.Recv()
		if err == io.EOF {
			return in.SendAndClose(categories)
		}
		if err != nil {
			return err
		}
		categoryResult, err := c.CategoryDB.CreateCategory(category.Name, category.Description)
		if err != nil {
			return err
		}
		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryResult.Name,
			Description: categoryResult.Description,
		})

	}
}

func (c *CategoryService) CreateCategoryStreamBidirectional(stream pb.CategoryService_CreateCategoryStreamBidirectionalServer) error {

	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		categoryResult, err := c.CategoryDB.CreateCategory(category.Name, category.Description)
		if err != nil {
			return err
		}
		if err := stream.Send(&pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryResult.Name,
			Description: categoryResult.Description,
		}); err != nil {
			return err
		}
	}
}
