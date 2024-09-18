package main

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
}

var users = []User{
	{ID: "1", Name: "John", LastName: "Doe"},
	{ID: "2", Name: "Jane", LastName: "Smith"},
	{ID: "3", Name: "Alice", LastName: "Johnson"},
	{ID: "4", Name: "Bob", LastName: "Brown"},
	{ID: "5", Name: "Charlie", LastName: "Davis"},
	{ID: "6", Name: "Diana", LastName: "Miller"},
	{ID: "7", Name: "Eve", LastName: "Wilson"},
	{ID: "8", Name: "Frank", LastName: "Moore"},
	{ID: "9", Name: "Grace", LastName: "Taylor"},
	{ID: "10", Name: "Hank", LastName: "Anderson"},
}
