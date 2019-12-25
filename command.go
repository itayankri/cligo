package main

type command struct {
	name      string
	funcName  string
	arguments []*argument
}

type argument struct {
	name  string
	_type string
}

type option struct {
	name  string
	_type string
}
