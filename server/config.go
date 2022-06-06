package main

type config struct {
	Secret      string
	URI         string
	Database    string
	Storage     string
	Directories bool
	Users       string
	Port        uint
	// Executable  string
}

var conf = config{
	URI:         "mongodb://localhost:27017/rake",
	Database:    "rake",
	Storage:     ".",
	Directories: false,
	Users:       ".",
	Port:        4100,
}
