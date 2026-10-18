package main
import "fmt"

func main() {
	numbers := [5]float64{1.1, 2.5, 3.1, 4.2, 5.6}
	// other variant is var numbers [5]float64 = [5]float64{1.1, 2.5, 3.1, 4.2, 5.6}
	var numberss [5]float64 = [5]float64{1.1, 2.5, 3.1, 4.2, 5.6}
	numbers1 := [5]float64{1.1, 2.5, 3.1}
	var a [3]float64
	a[0] = 1.1
	a[1] = 2.5
	a[2] = 3.1	
	fmt.Println(numbers)
	fmt.Println(numberss)
	fmt.Println(numbers1)
	fmt.Println(a)
}