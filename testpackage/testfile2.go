package testpackage

import "fmt"

//cligo:command "Adding numbers"
//cligo:option x "First number to add"
//cligo:option y "Second number to add"
func CliAdd(x, y float64) {
	fmt.Println(x + y)
}

//cligo:command "Subtracting numbers"
//cligo:option x "Subtracted number"
//cligo:option y "Subtracting number"
func CliSub(x, y float64) {
	fmt.Println(x - y)
}

//cligo:command "Multiplying numbers"
//cligo:option x2 "Multiplied number"
//cligo:option y "Multiplier"
func CliMul(x, y float64) {
	fmt.Println(x * y)
}

//cligo:command "Dividing numbers"
//cligo:option x "Divided number"
//cligo:option y "Divider"
func CliDiv(x, y float64) {
	if y == 0 {
		fmt.Println("cannot divide number by 0")
	} else {
		fmt.Println(x / y)
	}
}
