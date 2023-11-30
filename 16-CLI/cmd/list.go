/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/christianferraz/goexpert/16-CLI/database"
	"github.com/spf13/cobra"
)

func newListCmd(categoryDB database.Category) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List a category",
		RunE:  runList(categoryDB),
	}
}

func runList(categoryDB database.Category) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {

		categories, err := categoryDB.ListAllCategories()
		if err != nil {
			return err
		}
		for _, category := range categories {
			fmt.Printf("id: %v\nName: %v\nDescription: %v\n\n", category.ID, category.Name, category.Description)
		}
		return nil
	}
}

func init() {
	listCmd := newListCmd(GetCategoryDB(GetDB()))
	categoryCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("all", "a", "", "List all categories")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
