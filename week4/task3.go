package main

import "fmt"

type Counter struct{ n int }
func (c Counter) Get() int { return c.n }
func (c *Counter) Inc() { c.n++ }
func NewCounter() Counter { return Counter{} }

func main() {
	var c Counter
	(c.Inc()) // 1 да
	(&c).Inc() // 2 да
	Counter{}.Get() // 3 да
	Counter{}.Inc() // 4 нет
	NewCounter().Inc() // 5 нетт
	m := map[string]Counter{"a": {}}
	m["a"].Get() // 6 да
	m["a"].Inc() // 7 нет
	s := []Counter{{}}
	s[0].Inc() // 8 да
	pm := map[string]*Counter{"a": {}}
	pm["a"].Inc() // 9 да


	
}