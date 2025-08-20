package main

import (
	"fmt"
	"slices"
)

func main() {
	nums := []int{2, 2, 1, 4, 1}
	fmt.Println("只出现一次的数字:", singleNumber(nums))

	fmt.Println("是否为回文数:", isPalindrome(1221))

	fmt.Println("有效的括号:", isValid("()[]{}"))

	fmt.Println("最长公共前缀:", longestCommonPrefix([]string{"flower", "flow", "flight"}))

	fmt.Println("加一:", plusOne([]int{1, 2, 9}))

	fmt.Println("删除有序数组中的重复项:", removeDuplicates([]int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}))

	fmt.Println("合并区间:", merge([][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}))

	fmt.Println("两数之和:", twoSum([]int{2, 7, 11, 15}, 9))
}

// 只出现一次的数字
func singleNumber(nums []int) int {
	numMap := make(map[int]int)
	for _, num := range nums {
		numMap[num]++
	}
	for num, count := range numMap {
		if count == 1 {
			return num
		}
	}
	return -1
}

// 回文数
func isPalindrome(x int) bool {
	if x < 0 {
		return false
	}
	if x < 10 {
		return true
	}
	reversed := 0
	original := x
	for x > 0 {
		reversed = reversed*10 + x%10
		x /= 10
	}
	return original == reversed
}

// 有效的括号
func isValid(s string) bool {
	stack := []rune{}
	mapping := map[rune]rune{']': '[', '}': '{', ')': '('}
	for _, char := range s {
		if opening, exists := mapping[char]; exists {
			if len(stack) == 0 || stack[len(stack)-1] != opening {
				return false
			}
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, char)
		}
	}
	return len(stack) == 0
}

// 最长公共前缀
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for i := 1; i < len(strs); i++ {
		for j := 0; j < len(prefix) && j < len(strs[i]); j++ {
			if prefix[j] != strs[i][j] {
				prefix = prefix[:j]
				break
			}
		}
	}
	return prefix
}

// 加一
func plusOne(digits []int) []int {
	for i := len(digits) - 1; i >= 0; i-- {
		if digits[i] < 9 {
			digits[i]++
			return digits
		}
		digits[i] = 0
	}
	return append([]int{1}, digits...)
}

// 删除有序数组中的重复项
func removeDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	index := 1
	for i := 1; i < len(nums); i++ {
		if nums[i] != nums[i-1] {
			nums[index] = nums[i]
			index++
		}
	}
	return index
}

// 合并区间
func merge(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return intervals
	}
	slices.SortFunc(intervals, func(i, j []int) int { return i[0] - j[0] })
	merged := [][]int{intervals[0]}
	for _, i := range intervals {
		if i[0] <= merged[len(merged)-1][1] {
			merged[len(merged)-1][1] = max(merged[len(merged)-1][1], i[1])
		} else {
			merged = append(merged, i)
		}
	}
	return merged
}

// 两数之和
func twoSum(nums []int, target int) []int {
	numMap := make(map[int]int)
	for i, num := range nums {
		if j, found := numMap[target-num]; found {
			return []int{j, i}
		}
		numMap[num] = i
	}
	return nil
}
