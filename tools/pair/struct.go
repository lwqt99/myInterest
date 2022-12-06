package pair

import (
	"gitgit/tools/ecurve"
	"math/big"
)

type Pair struct {
	G1 *ecurve.CurveWeierstrass
	G2 *ecurve.CurveWeierstrass
	P * big.Int//p素数域的阶
	R *big.Int//r椭圆曲线子群的阶
	K *big.Int//k嵌入度
}

//定义直线
//此处为专属表达，与常规直线不同
//Ax + By + C
type line struct {
	A *big.Int
	B *big.Int
	C *big.Int
}

//定义虚数
//a + bi
type ImaginaryNumber struct {
	A *big.Int
	B *big.Int
}

