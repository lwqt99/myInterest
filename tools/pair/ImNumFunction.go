package pair

import (
	"gitgit/tools"
	"gitgit/tools/ecurve"
	"math/big"
)

func (im *ImaginaryNumber) Init(a, b *big.Int) *ImaginaryNumber {
	return im
}

func (im *ImaginaryNumber) SetInt64(a, b int64) *ImaginaryNumber {
	im.A = new(big.Int).SetInt64(a)
	im.B = new(big.Int).SetInt64(b)

	return im
}

func (im *ImaginaryNumber) SetPoint(p *ecurve.Point) *ImaginaryNumber {
	im.A = new(big.Int).Set(p.X)
	im.B = new(big.Int).Set(p.Y)

	return im
}

func (im *ImaginaryNumber) SetIM(p *ImaginaryNumber) *ImaginaryNumber {
	im.A = new(big.Int).Set(p.A)
	im.B = new(big.Int).Set(p.B)

	return im
}

func (im *ImaginaryNumber) Mod(point *ImaginaryNumber, p *big.Int) *ImaginaryNumber {

	im.A = new(big.Int).Mod(point.A, p)
	im.B = new(big.Int).Mod(point.B, p)

	return im
}

func (im *ImaginaryNumber) Mul(im1, im2 *ImaginaryNumber) *ImaginaryNumber {
	im.A, im.B = new(big.Int).Sub(new(big.Int).Mul(im1.A, im2.A), new(big.Int).Mul(new(big.Int).Mul(im1.B, im2.B), tools.Positive2)),
		new(big.Int).Add(new(big.Int).Mul(im1.A, im2.B), new(big.Int).Mul(im1.B, im2.A))

	return im
}

func (im *ImaginaryNumber) Exp(im1 *ImaginaryNumber, a, p *big.Int) *ImaginaryNumber {
	//im = im1^a mod p
	//采用二进制转化
	binaryA := tools.BigNumBaseConversion(a, 2)
	storeV := make([]*ImaginaryNumber, len(binaryA)) //存储 im1^(2^i)
	r := new(ImaginaryNumber).SetInt64(1, 0)

	//计算多倍点的二次幂
	storeV[0] = new(ImaginaryNumber).SetIM(im1)
	for i := 1; i < len(binaryA); i++ {
		storeV[i] = new(ImaginaryNumber).Mul(storeV[i-1], storeV[i-1])
		storeV[i].Mod(storeV[i], p)
	}
	//求乘积
	for i := 0; i < len(binaryA); i++ {
		if string(binaryA[len(binaryA)-i-1]) == "1" {
			r.Mul(r, storeV[i])
		}
	}
	im.SetIM(r)
	return im
}

func (im *ImaginaryNumber) String() string {
	return im.A.String() + "+" + im.B.String() + "i"
}
