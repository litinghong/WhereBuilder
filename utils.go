package whereBuilder

import (
	"fmt"
	"strconv"
	"strings"
)

// 将字符串分割为int数组
func SplitToInt(str, sep string) ([]int, error) {
	s := strings.Split(str, sep)
	r := make([]int, len(s))
	for i, v := range s {
		t, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		r[i] = t
	}

	return r, nil
}

// 将int数组合并为一个字符串
func JoinInt(i []int, sep string) string {
	s := make([]string, len(i))
	for k, v := range i {
		s[k] = strconv.Itoa(v)
	}

	return strings.Join(s, sep)
}

// 将float数组合并为一个字符串
func JoinFloat32(i []float32, sep string) string {
	s := make([]string, len(i))
	for k, v := range i {
		s[k] = fmt.Sprintf("%f", v)
	}

	return strings.Join(s, sep)
}

// 将float数组合并为一个字符串
func JoinFloat64(i []float64, sep string) string {
	return JoinFloat64(i, sep)
}
