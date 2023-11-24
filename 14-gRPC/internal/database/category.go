package database

import (
	"database/sql"

	"github.com/google/uuid"
)

type Category struct {
	ID          string
	Name        string
	Description string
	db          *sql.DB
}

func NewCategory(db *sql.DB) *Category {
	return &Category{db: db}
}

func (c *Category) CreateCategory(name, description string) (*Category, error) {
	uuid := uuid.New().String()
	_, err := c.db.Exec("INSERT INTO categories (id, name, description) VALUES ($1, $2, $3)", uuid, name, description)
	if err != nil {
		return nil, err
	}
	return &Category{
		ID:          uuid,
		Name:        name,
		Description: description,
	}, err
}

func (c *Category) ListAllCategories() ([]*Category, error) {
	rows, err := c.db.Query("SELECT id, name, description FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []*Category
	for rows.Next() {
		var id, name, description string
		if err := rows.Scan(&id, &name, &description); err != nil {
			return nil, err
		}
		category := &Category{
			ID:          id,
			Name:        name,
			Description: description,
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (c *Category) GetCategoryByID(id string) (*Category, error) {
	var category Category
	err := c.db.QueryRow("SELECT id, name, description FROM categories WHERE id = $1", id).Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		return nil, err
	}
	return &category, nil
}
