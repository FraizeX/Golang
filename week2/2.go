package main
import "fmt"
func main() {
	
	cars := []string{"Ferrari", "Honda", "Ford", "BYD"}
	fmt.Println("cars:", cars, "has old length", len(cars), "and capacity", cap(cars)) //capacity of cars is 4
	cars = append(cars, "Toyota")
 	fmt.Println("cars:", cars, "has new length", len(cars), "and capacity", cap(cars)) // What is the len and capacity of cars now? Why?
	a := cars[:cap(cars)]
	fmt.Println(a) 
}