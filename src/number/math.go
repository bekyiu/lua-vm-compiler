package number

import "math"

// ------- 算术运算 -------

// 整除 向负无穷取整
func IFloorDiv(a, b int64) int64 {
	if a > 0 && b > 0 || a < 0 && b < 0 || a%b == 0 {
		return a / b
	}
	return a/b - 1
}

func FFloorDiv(a, b float64) float64 {
	return math.Floor(a / b)
}

// 取模
func IMod(a, b int64) int64 {
	return a - IFloorDiv(a, b)*b

}

func FMod(a, b float64) float64 {
	return a - FFloorDiv(a, b)*b
}

// ------ 位运算 --------
// 左移操作, n < 0相当于反向移动n bit
func ShiftLeft(a, n int64) int64 {
	if n >= 0 {
		return a << uint64(n)
	}
	return ShiftRight(a, -n)
}
// 无符号右移, 空位补0
func ShiftRight(a int64, n int64) int64 {
	if n >= 0 {
		return int64(uint64(a) >> uint64(n))
	}
	return ShiftLeft(a, -n)
}
