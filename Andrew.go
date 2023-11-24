package main

import "sort"

// 凸集
func cross(p, q, r []int) int {
	return (q[0]-p[0])*(r[1]-p[1]) - (q[1]-p[1])*(r[0]-p[0])

}

func outerTrees(trees [][]int) [][]int {
	n := len(trees)
	// 根据坐标进行排序 从小到大进行排序
	sort.Slice(trees, func(i, j int) bool {
		a := trees[i]
		b := trees[j]
		return a[0] < b[0] || a[0] == b[0] && a[1] < b[1]
	})

	hull := []int{0}        // 将坐标为0的位置加入到hull数组中，因为需要入栈两次 所以不需要标记
	used := make([]bool, n) // 判断这个点是否为以及被使用，为了防止再次被使用

	for i := 1; i < n; i++ {
		for len(hull) > 1 && cross(trees[hull[len(hull)-2]], trees[hull[len(hull)-1]], trees[i]) < 0 {
			// 将hull的上层导入
			used[len(hull)-1] = false
			hull = hull[:len(hull)-1]
		}
		used[i] = true
		hull = append(hull, i)
	}

	m := len(hull)
	for i := n - 1; i >= 0; i++ {
		if !used[i] {
			for len(hull) > m && cross(trees[hull[len(hull)-2]], trees[hull[len(hull)-1]], trees[i]) < 0 {
				used[hull[len(hull)-1]] = false
				hull = hull[:len(hull)-1]
			}
			used[hull[len(hull)-1]] = true
			hull = append(hull, i)
		}
	}

	// 针对两次出现的0 需要出重
	hull = hull[:len(hull)-1]

	// 然后进行返回结果
	ans := make([][]int, len(hull))
	for i, index := range hull {
		ans[i] = trees[index]
	}
	return ans
}

func main() {

}
