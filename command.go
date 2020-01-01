package main

import (
	"bytes"
	"github.com/pkg/errors"
	"go/ast"
	"strings"
)

const (
	CLIGO_COMMAND  = "cligo:command"
	CLIGO_ARGUMENT = "cligo:argument"
	CLIGO_OPTION   = "cligo:option"
)

type command struct {
	name        string
	funcName    string
	description string
	options     []*option
	arguments   []*argument
}

type param interface {
	getDescription() string
}

type argument struct {
	name        string
	description string
	_type       string
}

func (a *argument) getDescription() string {
	return a.description
}

type option struct {
	name        string
	description string
	_type       string
}

func (o *option) getDescription() string {
	return o.description
}

func parseCommand(funcDecl *ast.FuncDecl) (*command, error) {
	if funcDecl == nil {
		return nil, errors.New("cannot parse nil function declaration")
	}

	command := &command{
		name:      strings.ToLower(funcDecl.Name.Name[3:]),
		funcName:  funcDecl.Name.Name,
		options:   make([]*option, 0),
		arguments: make([]*argument, 0),
	}

	paramsMap := make(map[string]string)

	for _, paramList := range funcDecl.Type.Params.List {
		if _type, ok := paramList.Type.(*ast.Ident); ok {
			for _, param := range paramList.Names {
				paramsMap[param.Name] = _type.Name
			}
		} else {
			return nil, errors.New("cannot create a sub-command based on a function that " +
				"requires a non-atomic argument. function name: " + funcDecl.Name.Name)
		}
	}

	for _, comment := range funcDecl.Doc.List {
		if isCligo([]byte(comment.Text)) {
			tokens, err := lex(comment.Text)
			if err != nil {
				return nil, errors.Wrap(err, "failed to lex comment: "+comment.Text)
			}

			switch tokens[0].value {
			case CLIGO_COMMAND:
				{
					if len(tokens) != 2 {
						return nil, errors.New("an argument directive must contain exactly one element")
					}

					if tokens[1].kind != STRING {
						return nil, errors.New("a command directive's element must be STRING")
					}

					command.description = tokens[1].value
				}
			case CLIGO_OPTION:
				{
					opt, err := parseOption(tokens[1:])
					if err != nil {
						return nil, errors.Wrap(err, "failed to parse CLIGO option directive: "+comment.Text)
					}

					if v, ok := paramsMap[opt.name]; ok {
						opt._type = v
						command.options = append(command.options, opt)
					} else {
						return nil, errors.New("could not parse option directive with name - " +
							opt.name +
							": param does not exist")
					}
				}
			case CLIGO_ARGUMENT:
				{
					arg, err := parseArgument(tokens[1:])
					if err != nil {
						return nil, errors.Wrap(err, "failed to parse CLIGO argument directive: "+comment.Text)
					}

					if v, ok := paramsMap[arg.name]; ok {
						arg._type = v
						command.arguments = append(command.arguments, arg)
					} else {
						return nil, errors.New("could not parse argument directive with name - " +
							arg.name +
							": param does not exist")
					}
				}
			default:
				{
					return nil, errors.New("invalid beginning of CLIGO directive: " + tokens[0].value)
				}
			}

		}
	}

	return command, nil
}

func parseArgument(tokens []*Token) (*argument, error) {
	if len(tokens) != 2 {
		return nil, errors.New("an argument directive must contain exactly two elements")
	}

	if tokens[0].kind != NAME {
		return nil, errors.New("an argument directive's first element must be NAME")
	}

	if tokens[1].kind != STRING {
		return nil, errors.New("an argument directive's first element must be STRING")
	}

	return &argument{
		name:        tokens[0].value,
		description: tokens[1].value,
	}, nil
}

func parseOption(tokens []*Token) (*option, error) {
	if len(tokens) != 2 {
		return nil, errors.New("an option directive must contain exactly two elements")
	}

	if tokens[0].kind != NAME {
		return nil, errors.New("an option directive's first element must be NAME")
	}

	if tokens[1].kind != STRING {
		return nil, errors.New("an option directive's first element must be STRING")
	}

	return &option{
		name:        tokens[0].value,
		description: tokens[1].value,
	}, nil
}

func isCligo(buf []byte) bool {
	return bytes.HasPrefix(buf, []byte("//cligo:command ")) || bytes.HasPrefix(buf, []byte("//cligo:command \t")) ||
		bytes.HasPrefix(buf, []byte("//cligo:argument ")) || bytes.HasPrefix(buf, []byte("//cligo:argument \t")) ||
		bytes.HasPrefix(buf, []byte("//cligo:option ")) || bytes.HasPrefix(buf, []byte("//cligo:option \t"))
}
