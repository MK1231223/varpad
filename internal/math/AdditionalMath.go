package app_math

type Vector1 struct {
	X1 int
	X2 int
}

func (v1 *Vector1) ResetTo(num int) {
	v1.X1 = num
	v1.X2 = num
}

func PercentageI32(x, y int32) int32 {
	return x * y / 100
}
