package testpackage

// A1 is an exported const
const A1 = 4

// a1 is an unexported const
const a1 = 4

// A is an exported function called A
// Test comment
func A() {

}

// a is an unexported function called a
// @CLIGO_COMMAND
func a() {

}

// @CLIGO_ARGUMENT(a: min)
// @CLIGO_ARGUMENT(b: max)
// B is an exported function called B
func B(a, b int) {

}

// c is an unexported function called c
func c() {

}

// Func is an exported function called Func
func Func() {

}

// fn is an unexported function called fn
func fn() {

}

/*
Z is an exported function called Z
@CLIGO_COMMAND
*/
func Z() {

}
