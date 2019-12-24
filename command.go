package main

type command struct {
	name      string
	alias     string
	arguments []*argument
	options   []*option
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
