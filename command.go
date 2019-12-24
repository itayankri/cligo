package main

type command struct {
	name      string
	funcName  string
	alias     string
	arguments []*argument
}

type argument struct {
	name    string
	aliases []string
	_type   string
}

type option struct {
	name    string
	aliases []string
	_type   string
}
