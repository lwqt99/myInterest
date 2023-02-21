package ecurve

import (
	"math/big"
)

// Point 定义椭圆曲线上的点
// 不存在任何特殊运算
// 单纯用于集成表示
type Point struct {
	X *big.Int
	Y *big.Int
}

// CurveWeierstrass 定义椭圆曲线 Weierstrass形式
type CurveWeierstrass struct {
	P *big.Int //域p
	//y^2 = x^3 + ax + b
	//4a^3 + 27b^2 != 0
	A       *big.Int
	B       *big.Int
	G       *Point   //base point
	N       *big.Int //G元素的阶
	BitSize int64
}

type PublicKey struct {
	P *Point
}

type PrivateKey struct {
	K   *big.Int
	Pub *PublicKey // pub = k * G
}
