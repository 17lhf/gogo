package basics

import (
	"fmt"
	"net/http"

	"gogo/internal/lesson"
)

// SlicesHandler GET /api/basics/slices
// 演示切片 — Go最重要的数据结构之一
func SlicesHandler(w http.ResponseWriter, r *http.Request) {
	// ── 创建切片 ──────────────────────────────────────────────
	// 切片是对底层数组的"视图"，包含：指针、长度、容量
	// Java: new ArrayList<>() 或 new int[]{}

	nums := []int{10, 20, 30, 40, 50}                    // 字面量
	made := make([]int, 3, 10)                            // make(类型, 长度, 容量)
	made[0], made[1], made[2] = 1, 2, 3
	var nilSlice []int                                    // nil切片

	// ── append ────────────────────────────────────────────────
	// 容量不足时自动扩容，但会返回新切片！必须接收返回值
	// Java: list.add(elem)   Go: slice = append(slice, elem)
	nilSlice = append(nilSlice, 1, 2, 3)
	nums = append(nums, 60, 70)

	// ── 切片操作 (slice of slice) ─────────────────────────────
	// Java: list.subList(1, 4)
	sub := nums[1:4]  // [20,30,40]，包含index1，不含index4
	head := nums[:3]  // [10,20,30]
	tail := nums[3:]  // [40,50,60,70]

	// !! sub与nums共享底层数组，修改sub会影响nums !!

	// ── range遍历 ─────────────────────────────────────────────
	// Java: for (int x : list)  Go: for _, v := range slice
	var doubled []int
	for _, v := range nums {   // range返回(index, value)，_丢弃不需要的值
		doubled = append(doubled, v*2)
	}

	// ── 二维切片 ──────────────────────────────────────────────
	matrix := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}

	// ── copy：深拷贝，不共享底层数组 ──────────────────────────
	dst := make([]int, len(nums))
	n := copy(dst, nums)

	lesson.WriteJSON(w, http.StatusOK, lesson.Response{
		Code:  http.StatusOK,
		Topic: "Go切片 (Slice)",
		Java:  "Java ArrayList<T> 或 T[]",
		Summary: "切片是Go最常用的数据结构，本质是(指针+长度+容量)三元组。" +
			"与Java ArrayList不同，append可能返回新切片，必须接收返回值。",
		Points: []string{
			"切片三要素: 底层数组指针、len(长度)、cap(容量)",
			"append(slice, elem) 容量不足时扩容并返回新切片，必须用 slice = append(slice, elem)",
			"nums[low:high] 与原切片共享内存，修改会互相影响",
			"make([]T, len, cap) 预分配容量，避免频繁扩容",
			"range遍历返回(index, value)，_丢弃不需要的值",
			"copy(dst, src) 深拷贝，不共享底层数组",
		},
		Data: map[string]interface{}{
			"literal_slice":    nums,
			"make_slice":       made,
			"appended_nil":     nilSlice,
			"sub_slice":        sub,
			"head_slice":       head,
			"tail_slice":       tail,
			"doubled_range":    doubled,
			"matrix":           matrix,
			"copy_count":       n,
			"copied_slice":     dst,
			"len_cap":          fmt.Sprintf("len=%d, cap=%d", len(nums), cap(nums)),
		},
	})
}
