package main

import "fmt"

func twoSum(nums []int, target int) []int {
	numMap := make(map[int]int)
	for i, num := range nums {
		complement := target - num

		if j, ok := numMap[complement]; ok {
			return []int{j, i}
		}

		numMap[num] = i
	}
	return nil
}

func main() {
	i := twoSum([]int{2, 7, 11, 15}, 24)
	fmt.Println(i)
}
