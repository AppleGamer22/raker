package main

type config struct {
	Secret      string
	URI         string
	Database    string
	Storage     string
	Directories bool
	Port        uint
}

var conf = config{
	URI:         "mongodb://localhost:27017/rake",
	Database:    "rake",
	Storage:     ".",
	Directories: false,
	Port:        4100,
}
