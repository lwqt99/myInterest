package ecurve

import (
	"fmt"
	"gitgit/tools"
	"math/big"
)

func Tr() {

}

//测试Weierstrass表达的椭圆曲线
func TestWeierstrassEcc() {
	//测试用参数
	P, _ := new(big.Int).SetString("39402006196394479212279040100143613805079739270465446667948293404245721771496870329047266088258938001861606973112319", 10)
	N, _ := new(big.Int).SetString("39402006196394479212279040100143613805079739270465446667946905279627659399113263569398956308152294913554433653942643", 10)
	B, _ := new(big.Int).SetString("b3312fa7e23ee7e4988e056be3f82d19181d9c6efe8141120314088f5013875ac656398d8a2ed19d2a85c8edd3ec2aef", 16)
	Gx, _ := new(big.Int).SetString("aa87ca22be8b05378eb1c71ef320ad746e1d3b628ba79b9859f741e082542a385502f25dbf55296c3a545e3872760ab7", 16)
	Gy, _ := new(big.Int).SetString("3617de4a96262c6f5d9e98bf9292dc29f8f41dbd289a147ce9da3113b5f0b8c00a60b1ce1d7e819d7a431d7c90ea0e5f", 16)

	c := new(Curve).Init(new(big.Int).SetInt64(-3),B, P, N, new(Point).Init(Gx, Gy))

	//curve := elliptic.P384()
	//fmt.Println(c.Mul(new(big.Int).SetInt64(3), c.G).String())
	//fmt.Println(c.Add(c.G, new(Point).Init(curve.ScalarBaseMult(new(big.Int).SetInt64(2).Bytes()))))
	//fmt.Println(curve.ScalarBaseMult(new(big.Int).SetInt64(3).Bytes()))

	/**/
	//生成私钥
	pri := tools.GenerateBigIntByRange(N)
	pub := c.Mul(pri, c.G)//pri * g

	//待加密信息点
	m := c.Mul(new(big.Int).SetInt64(4), c.G)
	//加密
	//生成加密用的随机数r
	r := tools.GenerateBigIntByRange(N)
	rPub := c.Mul(r, pub)//r * pub

	c1 := c.Mul(r, c.G)//r * g
	c2 := c.Add(rPub, m)//r * pub + m

	//解密
	//c2 - pri * c1
	//=r * pub + m - pri * r * g
	//=r * pri * g - pri * r * g + m
	//fmt.Println(c.Mul(new(big.Int).SetInt64(12), c.G))
	//fmt.Println(c.Mul(new(big.Int).SetInt64(3),c.Mul(new(big.Int).SetInt64(4), c.G)))

	dm := c.Add(c2, c.revPoint(c.Mul(pri, c1)))

	//fmt.Println("r * pub =", rPub.String())
	//fmt.Println("pri * r * g = ", c.Mul(pri, c1).String())
	fmt.Println("原文：", m.String())
	fmt.Println("解码结果：",dm.String())
	return

}
