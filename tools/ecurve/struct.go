package ecurve

import (
	"math/big"
)

// 定义椭圆曲线上的点
type Point struct {
	X *big.Int
	Y *big.Int
}

// 定义椭圆曲线 Weierstrass形式
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
