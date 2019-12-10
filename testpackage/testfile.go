package testpackage

import (
	"fmt"
	"strconv"
)

// A1 is an exported const
const A1 = 4

// a1 is an unexported const
const a1 = 4

// A is an exported function called A
// Test comment
// @CLIGO_COMMAND(install)
// @CLIGO_OPTION(option1: O1)
// @CLIGO_OPTION(option2: O2)
// @CLIGO_ARGUMENT(arg1: x)
// @CLIGO_ARGUMENT(arg2: y)
func A(option1, option2 bool, arg1, arg2 int) {
	fmt.Println("executing function A() - option1=" +
		strconv.FormatBool(option1) +
		"option2=" +
		strconv.FormatBool(option2))
}

// a is an unexported function called a
// @CLIGO_COMMAND(alias)
func a() {
	fmt.Println("executing function a()")
}

// @CLIGO_ARGUMENT(a: min)
// @CLIGO_ARGUMENT(b: max)
// B is an exported function called B
func B(a, b int) {
	fmt.Println("executing function B()")
}

// c is an unexported function called c
func c() {
	fmt.Println("executing function c()")
}

// Func is an exported function called Func
func Func() {
	fmt.Println("executing function Func()")
}

// fn is an unexported function called fn
func fn() {
	fmt.Println("executing function fn()")
}

/*
Z is an exported function called Z
@CLIGO_COMMAND
*/
func Z() {
	fmt.Println("executing function Z()")
}
