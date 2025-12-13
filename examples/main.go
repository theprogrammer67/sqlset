package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/theprogrammer67/sqlset"
)

//go:embed queries
var queriesFS embed.FS

func main() {
	// Create a new SQLSet from the embedded filesystem.
	sqlSet, err := sqlset.New(queriesFS)
	if err != nil {
		log.Fatalf("Failed to create SQL set: %v", err)
	}

	// Get a specific query
	query, err := sqlSet.Get("users", "GetUserByID")
	if err != nil {
		log.Fatalf("Failed to get query: %v", err)
	}
	fmt.Println("GetUserByID query:", query)

	// Or, panic if the query is not found
	query = sqlSet.MustGet("users", "CreateUser")
	fmt.Println("CreateUser query:", query)

	// You can also retrieve metadata for all query sets
	metas := sqlSet.GetAllMetas()
	for _, meta := range metas {
		fmt.Printf("Set ID: %s, Name: %s, Description: %s\n", meta.ID, meta.Name, meta.Description)
	}
}
