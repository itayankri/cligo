package main

type annotationType string

const (
	CLIGO_COMMAND  annotationType = "@CLIGO_COMMAND"
	CLIGO_ARGUMENT annotationType = "@CLIGO_ARGUMENT"
	CLIGO_OPTION   annotationType = "@CLIGO_OPTION"
)

type annotation struct {
	annotationType annotationType
	name           string
	alias          string
}
