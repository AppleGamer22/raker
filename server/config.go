package main

type config struct {
	Secret   string
	Database string
	Storage  string
	Users    string
	Port     uint
}

var conf = config{
	Database: "mongodb://localhost:27017/rake",
	Storage:  ".",
	Users:    ".",
	Port:     4200,
}
