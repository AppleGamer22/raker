package main

type config struct {
	secret   string
	database string
	storage  string
	users    string
	port     uint
}

var conf = config{
	database: "mongodb://localhost:27017/rake",
	storage:  ".",
	users:    ".",
	port:     4200,
}
