package testpackage

import "fmt"

func CliAdd(x, y float64) {
	fmt.Println(x + y)
}

//func CliAddMany(nums ...float64) {
//	var sum float64
//
//	for _, num := range nums {
//		sum += num
//	}
//
//	fmt.Println(sum)
//}

func CliSub(x, y float64) {
	fmt.Println(x - y)
}

func CliMul(x, y float64) {
	fmt.Println(x * y)
}

func CliDiv(x, y float64) {
	if y == 0 {
		fmt.Println("cannot divide number by 0")
	} else {
		fmt.Println(x / y)
	}
}
