func sum() int {
	nums := make([]int, 1000)

	for i := range nums {
		nums[i] = i
	}

	total := 0

	for _, n := range nums {
		total += n
	}

	return total
}