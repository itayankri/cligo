package main

const (
	CLIGO_COMMAND  = "@CLIGO_COMMAND"
	CLIGO_ARGUMENT = "@CLIGO_ARGUMENT"
	CLIGO_OPTION   = "@CLIGO_OPTION"
)

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
