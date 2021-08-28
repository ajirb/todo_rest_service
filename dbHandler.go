package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DRIVER string = "mysql"
var URL string = "root:toor@tcp(127.0.0.1:3306)/?parseTime=true"
var DB_NAME string = "testdb"
var CREATE_TABLE_QUERY string = "create table if not exists TODO( " +
	"  id integer AUTO_INCREMENT PRIMARY KEY, " +
	"  name varchar(50), " +
	"  description varchar(256)," +
	"  priority integer, " +
	"  due date default (CURRENT_DATE), " +
	"  completed bool, " +
	"  completion_date date " +
	");"

func initialDbSetup() {
	fmt.Println("Connnecting to DB ...")
	db, err := sql.Open(DRIVER, URL)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	fmt.Println("...Connected")

	execQuery(db, "CREATE DATABASE IF NOT EXISTS "+DB_NAME)
	execQuery(db, "USE "+DB_NAME)
	execQuery(db, CREATE_TABLE_QUERY)
}

func execQuery(d *sql.DB, query string) {
	_, err := d.Exec(query)
	if err != nil {
		panic(err.Error())
	}
}

func insertData(todo Todo) Todo {
	fmt.Println("Connnecting to DB ...")
	db, err := sql.Open(DRIVER, URL)
	if err != nil {
		panic(err.Error())
	}
	execQuery(db, "USE "+DB_NAME)
	defer db.Close()
	fmt.Println("...Connected")

	var due_date = "null"
	if todo.Due_Date.Valid {
		due_date = fmt.Sprintf("str_to_date(\"%v\",\"%%m-%%d-%%Y\")", todo.Due_Date.Time.Format("01-02-2006"))
	}

	var completion_date = "null"
	if todo.Completion_Date.Valid {
		completion_date = fmt.Sprintf("str_to_date(\"%v\",\"%%m-%%d-%%Y\")", todo.Completion_Date.Time.Format("01-02-2006"))
	}

	insert_query := fmt.Sprintf("insert into todo (name,"+
		"description,priority,due,completed,completion_date)"+
		"value(\"%v\",\"%v\",%v,%v,%v,%v)", todo.Name.String, todo.Description.String,
		todo.Priority.Int64, due_date, todo.Completed.Bool, completion_date)
	fmt.Println(insert_query)

	result, err := db.Query(insert_query)

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&todo.Id, &todo.Name, &todo.Description, &todo.Priority, &todo.Due_Date, &todo.Completed, &todo.Completion_Date)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(todo.Due_Date)
		break
	}

	return todo
}

func readTable(query string) []Todo {
	var todos []Todo
	fmt.Println("Connnecting to DB ...")
	db, err := sql.Open(DRIVER, URL)
	if err != nil {
		panic(err.Error())
	}
	execQuery(db, "USE "+DB_NAME)
	defer db.Close()
	fmt.Println("...Connected")

	fmt.Println(query)
	result, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		var todo Todo
		err = result.Scan(&todo.Id, &todo.Name, &todo.Description, &todo.Priority, &todo.Due_Date, &todo.Completed, &todo.Completion_Date)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(todo.Due_Date)
		todos = append(todos, todo)
	}
	return todos
}

func updateTodo(todo Todo) Todo {
	fmt.Println("Connnecting to DB ...")
	db, err := sql.Open(DRIVER, URL)
	if err != nil {
		panic(err.Error())
	}
	execQuery(db, "USE "+DB_NAME)
	defer db.Close()
	fmt.Println("...Connected")

	var due_date = "null"
	if todo.Due_Date.Valid {
		due_date = fmt.Sprintf("str_to_date(\"%v\",\"%%m-%%d-%%Y\")", todo.Due_Date.Time.Format("01-02-2006"))
	}
	var completion_date = "null"
	if todo.Completion_Date.Valid {
		completion_date = fmt.Sprintf("str_to_date(\"%v\",\"%%m-%%d-%%Y\")", todo.Completion_Date.Time.Format("01-02-2006"))
	}

	update_query := fmt.Sprintf("update todo set name=\"%v\",description=\"%v\",priority=%v,due=%v,completed=%v,"+
		"completion_date=%v where id=%v", todo.Name.String, todo.Description.String,
		todo.Priority.Int64, due_date, todo.Completed.Bool, completion_date, todo.Id)
	fmt.Println(update_query)

	execQuery(db, update_query)
	return todo
}

func deleteTodo(id string) bool {
	fmt.Println("Connnecting to DB ...")
	db, err := sql.Open(DRIVER, URL)
	if err != nil {
		panic(err.Error())
	}
	execQuery(db, "USE "+DB_NAME)
	defer db.Close()
	fmt.Println("...Connected")

	result, err := db.Query("delete from todo where id=" + id)

	if err != nil {
		panic(err.Error())
	}
	return result.Next()
}
