package main

import (
	"database/sql"
	"fmt"

	"github.com/miffyG/golearn/task3/db"
)

// 创建employees 表，包含字段 id 、 name 、 department 、 salary
func createEmployeesTable() sql.Result {
	res := db.SqlxDb.MustExec(`
		CREATE TABLE IF NOT EXISTS employees (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			department VARCHAR(100) NOT NULL,
			salary DECIMAL(15, 2) NOT NULL
		)
	`)
	fmt.Println("创建 employees 表成功")
	return res
}

// 如果employees表没有数据则往employees表插入一些数据供测试
func insertEmployeesIfNeeded() {
	var count int
	err := db.SqlxDb.Get(&count, `SELECT COUNT(*) FROM employees`)
	if err != nil {
		fmt.Println("查询员工数量失败:", err)
		return
	}
	if count == 0 {
		db.SqlxDb.MustExec(`
			INSERT INTO employees (name, department, salary) VALUES
			('张三', '技术部', 12000.00),
			('李四', '技术部', 15000.00),
			('王五', '市场部', 9000.00),
			('赵六', '技术部', 18000.00),
			('钱七', '人事部', 8000.00),
			('孙八', '技术部', 16000.00),
			('周九', '市场部', 9500.00),
			('吴十', '技术部', 17000.00),
			('郑十一', '人事部', 8500.00),
			('冯十二', '技术部', 14000.00),
			('陈十三', '市场部', 10000.00),
			('褚十四', '技术部', 15500.00),
			('卫十五', '人事部', 8200.00),
			('蒋十六', '技术部', 16500.00),
			('沈十七', '市场部', 9800.00),
			('韩十八', '技术部', 17500.00),
			('杨十九', '人事部', 8300.00),
			('朱二十', '技术部', 14500.00),
			('秦二十一', '市场部', 10200.00),
			('尤二十二', '技术部', 15800.00)
		`)
		fmt.Println("插入20条员工数据成功")
		return
	}
	fmt.Println("已有员工数据，无须插入")
}

// 使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中
type Employee struct {
	ID         int     `db:"id"`
	Name       string  `db:"name"`
	Department string  `db:"department"`
	Salary     float64 `db:"salary"`
}

func queryTechnicalEmployees() {
	var employees []Employee
	err := db.SqlxDb.Select(&employees, `
		SELECT * FROM employees WHERE department = '技术部'
	`)
	if err != nil {
		fmt.Println("查询技术部员工信息失败:", err)
		return
	}
	fmt.Println("查询到的技术部员工信息:")
	for _, emp := range employees {
		fmt.Printf("ID: %d, Name: %s, Department: %s, Salary: %.2f\n",
			emp.ID, emp.Name, emp.Department, emp.Salary)
	}
}

// 使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中
func queryHighestSalaryEmployee() {
	var employee Employee
	err := db.SqlxDb.Get(&employee, `
		SELECT * FROM employees ORDER BY salary DESC LIMIT 1
	`)
	if err != nil {
		fmt.Println("查询最高工资员工信息失败:", err)
		return
	}
	fmt.Println("查询到的最高工资员工信息:", employee)
}

// 创建一个 books 表，包含字段 id 、 title 、 author 、 price
func createBooksTable() sql.Result {
	res := db.SqlxDb.MustExec(`
		CREATE TABLE IF NOT EXISTS books (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(100) NOT NULL,
			author VARCHAR(100) NOT NULL,
			price DECIMAL(15, 2) NOT NULL
		)
	`)
	fmt.Println("创建 books 表成功")
	return res
}

type Book struct {
	ID     int     `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"`
}

// 往books表插入20条数据供测试
func insertBooksIfNeeded() {
	var count int
	err := db.SqlxDb.Get(&count, `SELECT COUNT(*) FROM books`)
	if err != nil {
		fmt.Println("查询书籍数量失败:", err)
		return
	}
	if count == 0 {
		db.SqlxDb.MustExec(`
			INSERT INTO books (title, author, price) VALUES
			('Go语言入门', '张三', 59.00),
			('Python编程', '李四', 68.00),
			('Java核心技术', '王五', 88.00),
			('数据库原理', '赵六', 72.00),
			('算法导论', '钱七', 99.00),
			('操作系统', '孙八', 65.00),
			('网络基础', '周九', 54.00),
			('前端开发', '吴十', 49.00),
			('后端架构', '郑十一', 79.00),
			('人工智能', '冯十二', 120.00),
			('机器学习', '陈十三', 110.00),
			('深度学习', '褚十四', 15.00),
			('数据结构', '卫十五', 60.00),
			('C语言精粹', '蒋十六', 45.00),
			('Rust实战', '沈十七', 70.00),
			('Kotlin开发', '韩十八', 58.00),
			('Swift编程', '杨十九', 62.00),
			('PHP项目实战', '朱二十', 53.00),
			('Ruby基础', '秦二十一', 47.00),
			('Scala进阶', '尤二十二', 66.00)
		`)
		fmt.Println("插入书籍数据成功")
		return
	}
	fmt.Println("已有书籍数据，无须插入")
}

// 使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全
func queryBooksAbovePrice(price float64) {
	var books []Book
	err := db.SqlxDb.Select(&books, `
		SELECT * FROM books WHERE price > ?
	`, price)
	if err != nil {
		fmt.Println("查询价格大于", price, "的书籍失败:", err)
		return
	}
	fmt.Println("查询到的价格大于", price, "的书籍信息:")
	for _, book := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s, Price: %.2f\n",
			book.ID, book.Title, book.Author, book.Price)
	}
}
