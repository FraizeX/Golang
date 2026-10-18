package main

import (
	"fmt"
)
func main() {
	var matrix[2][4]int
	matrix[0]=[4]int{1,2,3,4}
	matrix[1]=[4]int{5,6,7,8}
	for rowIndex, row := range matrix {
		fmt.Println("Row:", rowIndex, "Values:", row,)
		for colIndex, col := range row {
			fmt.Println("Column:", colIndex, "Value:", col,)
		}
	}
}



