package main

type command struct {
	name      string
	alias     string
	arguments []*argument
	options   []*option
}

type argument struct {
	name  string
	alias string
	_type string
}

type option struct {
	name  string
	alias string
	_type string
}
