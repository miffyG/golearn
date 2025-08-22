package main

import (
	"github.com/miffyG/golearn/task3/db"
)

func main() {
	dbConfig, err := db.LoadConfig()
	if err != nil {
		panic(err)
	}
	db.InitSqlxDb(&dbConfig)
	db.InitGormDb(&dbConfig)

	createEmployeesTable()
	insertEmployeesIfNeeded()
	queryTechnicalEmployees()
	queryHighestSalaryEmployee()

	createBooksTable()
	insertBooksIfNeeded()
	queryBooksAbovePrice(50)

	createBlogTables()
	insertBlogTestData()
	getUserPostsAndComments(1)
	getMostCommentedPost()

	db.CloseDBConnections()
}
