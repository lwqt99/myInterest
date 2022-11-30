package algebra

import "math/big"

//定义基础的类——抽象的起点

type Mat struct {
	height int
	width int
	mat [][]big.Int
	matInt64 [][]int64
}

//多项式
type polynomial struct {
	vectorInt64 []int64
	vectorBigInt []big.Int
	AInt64 [][]int64//构建矩阵——用于支持特定计算
	ABigInt [][]big.Int
}

//封装BigFloat
type BigFloat struct {
	A big.Int
	B big.Int
}