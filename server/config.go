package main

type config struct {
	Secret   string
	URI      string
	Database string
	Storage  string
	Users    string
	Port     uint
}

var conf = config{
	URI:      "mongodb://localhost:27017/rake",
	Database: "rake",
	Storage:  ".",
	Users:    ".",
	Port:     4100,
}
