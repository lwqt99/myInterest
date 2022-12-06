package ecurve

import (
	"crypto/elliptic"
	crand "crypto/rand"
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

	c := new(CurveWeierstrass).Init(new(big.Int).SetInt64(-3),B, P, N, new(Point).Init(Gx, Gy), 384)

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
	dm := c.Add(c2, c.revPoint(c.Mul(pri, c1)))

	fmt.Println("原文：", m.String())
	fmt.Println("解码结果：",dm.String())
	return

}

//测试密钥生成是否正确
func TestGenKeyCor() {
	//测试用参数
	P, _ := new(big.Int).SetString("39402006196394479212279040100143613805079739270465446667948293404245721771496870329047266088258938001861606973112319", 10)
	N, _ := new(big.Int).SetString("39402006196394479212279040100143613805079739270465446667946905279627659399113263569398956308152294913554433653942643", 10)
	B, _ := new(big.Int).SetString("b3312fa7e23ee7e4988e056be3f82d19181d9c6efe8141120314088f5013875ac656398d8a2ed19d2a85c8edd3ec2aef", 16)
	Gx, _ := new(big.Int).SetString("aa87ca22be8b05378eb1c71ef320ad746e1d3b628ba79b9859f741e082542a385502f25dbf55296c3a545e3872760ab7", 16)
	Gy, _ := new(big.Int).SetString("3617de4a96262c6f5d9e98bf9292dc29f8f41dbd289a147ce9da3113b5f0b8c00a60b1ce1d7e819d7a431d7c90ea0e5f", 16)

	c := new(CurveWeierstrass).Init(new(big.Int).SetInt64(-3),B, P, N, new(Point).Init(Gx, Gy), 384)

	rand := crand.Reader
	k, _ := tools.GenerateBigIntByByte(c.BitSize, rand)

	pri := new(PrivateKey).Init()

	pri.K = new(big.Int).Set(k)
	pri.Pub.P.SetPoint(c.Mul(k, c.G))//自行运算的结果

	curve := elliptic.P384()
	t := new(Point).Init(curve.ScalarBaseMult(k.Bytes()))//真实的结果

	fmt.Println(pri.Pub.P)
	fmt.Println(t)
	fmt.Println(pri.Pub.P.Equal(t))

}

func TestWeierstrassSignature() {
	c := new(CurveWeierstrass).InitP384()
	pri, err := c.GenerateKey(nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	m := "123"
	s1, s2 := c.Signature(m, pri)
	fmt.Println(c.Verify(m, s1, s2, pri.Pub))
}