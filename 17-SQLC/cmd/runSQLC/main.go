package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/christianferraz/goexpert/17-SQLC/internal/db"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ctx := context.Background()
	dbConn, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/courses")
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	queries := db.New(dbConn)

	// courseArgs := CourseParams{
	// 	ID:          uuid.New().String(),
	// 	Name:        "Go",
	// 	Description: sql.NullString{String: "Go course", Valid: true},
	// }
	// categoryArgs := CategoryParams{
	// 	ID:          uuid.New().String(),
	// 	Name:        "Backend",
	// 	Description: sql.NullString{String: "Backend Course", Valid: true},
	// }
	// courseDB := NewCourseDB(dbConn)
	// err = courseDB.CreateCourseAndCategory(ctx, categoryArgs, courseArgs)
	// if err != nil {
	// 	panic(err)
	// }

	courses, err := queries.ListCourses(ctx)
	if err != nil {
		panic(err)
	}
	for _, c := range courses {
		fmt.Printf("Category: %s, Course ID: %s, Course Name: %s, Course Description: %s, Course Price: %f", c.CategoryID, c.ID, c.Name, c.Description.String, c.Price)
	}

}
