package ecurve

import (
	"gitgit/tools"
	"math/big"
)

/*
	常规来看仿射坐标中中椭圆曲线表达如下：y^2 = x^3 + ax + b
	但在标准射影坐标中，点(x, y, z)对应仿射坐标中的点(x/z, y/z)
	-->故有新的方程组 y^2 * z = x^3 + ax * z^2 + b * z^3
	如果采用加重射影坐标,点(x, y, z)对应仿射坐标中的点(x/z^2, y/z^3),则有
	-->y^2 = x^3 + ax * z^4 + b * z^6

	Ps.下面都是加重射影坐标
*/

// 获取仿射坐标的z值
func zForAffine(x, y *big.Int) *big.Int {
	z := new(big.Int)
	//不会都为0，否则为无穷点
	if x.Sign() != 0 || y.Sign() != 0 {
		z.SetInt64(1)
	}
	return z
}

//将仿射坐标转为射影坐标

// 将射影坐标转为仿射坐标
func (c *CurveWeierstrass) affineFromJacobian(x, y, z *big.Int) (xOut, yOut *big.Int) {
	if z.Sign() == 0 {
		//如果z为0，一般定义为无穷点
		return new(big.Int), new(big.Int)
	}

	zInv := new(big.Int).ModInverse(z, c.P)
	zInvSquare := new(big.Int).Mul(zInv, zInv)

	//xOut = x/z^2
	//yOut = y/z^3
	xOut = new(big.Int).Mul(x, zInvSquare)
	xOut.Mod(xOut, c.P)
	zInvSquare.Mul(zInvSquare, zInv)
	yOut = new(big.Int).Mul(y, zInvSquare)
	yOut.Mod(yOut, c.P)
	return xOut, yOut
}

// 返回的是雅克比表达
func (c *CurveWeierstrass) addJacobian(x1, y1, z1, x2, y2, z2 *big.Int) (*big.Int, *big.Int, *big.Int) {
	//https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-3.html#addition-add-2007-bl
	x3, y3, z3 := new(big.Int), new(big.Int), new(big.Int)
	//判断是否存在无穷点，存在则返回另一个点
	if z1.Sign() == 0 {
		x3.Set(x2)
		y3.Set(y2)
		z3.Set(z2)
		return x3, y3, z3
	}
	if z2.Sign() == 0 {
		x3.Set(x1)
		y3.Set(y1)
		z3.Set(z1)
		return x3, y3, z3
	}

	z1Square := new(big.Int).Mul(z1, z1) // z1^2
	z1Square.Mod(z1Square, c.P)

	z2Square := new(big.Int).Mul(z2, z2) // z2^2
	z2Square.Mod(z2Square, c.P)

	u1 := new(big.Int).Mul(x1, z2Square) // t1 = x1 * z2^2
	u1.Mod(u1, c.P)
	u2 := new(big.Int).Mul(x2, z1Square) // t2 = x2 * z1^2
	u2.Mod(u2, c.P)

	h := new(big.Int).Sub(u2, u1) // t3 = t1 - t2
	xEqual := h.Sign() == 0       // x3是否为0
	if h.Sign() == -1 {
		//如果t3小于0，则调整为正数
		h.Add(h, c.P)
	}

	i := new(big.Int).Lsh(h, 1) // i = 4 * (t1 - t2)^2
	i.Mul(i, i)
	j := new(big.Int).Mul(i, h) // j = 4 * t3^3

	//xEqual := t3.Sign() == 0
	s1 := new(big.Int).Mul(y1, z2) // t4 =  y1 * z2^3
	s1.Mul(s1, z2Square)
	s1.Mod(s1, c.P)
	s2 := new(big.Int).Mul(y2, z1) // t5 = y2 * z1^3
	s2.Mul(s2, z1Square)
	s2.Mod(s2, c.P)

	r := new(big.Int).Sub(s2, s1) // r = 2 * (t5 - t4)

	yEqual := r.Sign() == 0 // 判断y3是否为0
	if r.Sign() == -1 {
		//如果t5小于0，则调整为正数
		r.Add(r, c.P)
	}
	if xEqual && yEqual {
		//如果x3和y3都为0，说明两者为同一个点
		return c.doubleJacobian(x1, y1, z1)
	}

	r.Lsh(r, 1)
	v := new(big.Int).Mul(u1, i) // v = 4 * t1 * (t1 - t2)^2

	x3.Set(r)
	x3.Mul(x3, x3)
	x3.Sub(x3, j)
	x3.Sub(x3, v)
	x3.Sub(x3, v)
	x3.Mod(x3, c.P)

	y3.Set(r)
	v.Sub(v, x3)
	y3.Mul(y3, v)
	s1.Mul(s1, j)
	s1.Lsh(s1, 1)
	y3.Sub(y3, s1)
	y3.Mod(y3, c.P)

	z3.Add(z1, z2)
	z3.Mul(z3, z3)
	z3.Sub(z3, z1Square)
	z3.Sub(z3, z2Square)
	z3.Mul(z3, h)
	z3.Mod(z3, c.P)

	return x3, y3, z3
}

// AddJacobian 返回的是仿射坐标表达
func (c *CurveWeierstrass) AddJacobian(x1, y1, x2, y2 *big.Int) (*big.Int, *big.Int) {
	//初始化
	//基于x y生成z，一般z为1
	z1 := zForAffine(x1, y1)
	z2 := zForAffine(x2, y2)

	return c.affineFromJacobian(c.addJacobian(x1, y1, z1, x2, y2, z2))
}

// 辅助小工具，相同的点相加
// 返回射影坐标结果
func (c *CurveWeierstrass) doubleJacobian(x, y, z *big.Int) (*big.Int, *big.Int, *big.Int) {
	//https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-3.html#doubling-dbl-2001-b
	//fmt.Println("点的倍乘公式执行")
	delta := new(big.Int).Mul(z, z) // delta = z^2
	delta.Mod(delta, c.P)
	gamma := new(big.Int).Mul(y, y) // gamma = y^2
	gamma.Mod(gamma, c.P)
	alpha := new(big.Int).Sub(x, delta) // alpha = 3 * (x - delta) * (x + delta)
	if alpha.Sign() == -1 {
		// 负数
		alpha.Add(alpha, c.P)
	}
	alpha2 := new(big.Int).Add(x, delta)
	alpha.Mul(alpha, alpha2)
	alpha2.Set(alpha)
	alpha.Lsh(alpha, 1)
	alpha.Add(alpha, alpha2)

	beta := alpha2.Mul(x, gamma) //  beta = x * gamma

	x3 := new(big.Int).Mul(alpha, alpha) // x3 = alpha^2 - 8*delta
	beta8 := new(big.Int).Lsh(beta, 3)
	beta8.Mod(beta8, c.P)
	x3.Sub(x3, beta8)
	if x3.Sign() == -1 {
		x3.Add(x3, c.P)
	}
	x3.Mod(x3, c.P)

	z3 := new(big.Int).Add(y, z) // z3 = (y + z)^2 - gamma - delta
	z3.Mul(z3, z3)
	z3.Sub(z3, gamma)
	if z3.Sign() == -1 {
		z3.Add(z3, c.P)
	}
	z3.Sub(z3, delta)
	if z3.Sign() == -1 {
		z3.Add(z3, c.P)
	}
	z3.Mod(z3, c.P)

	beta.Lsh(beta, 2)
	beta.Sub(beta, x3)
	if beta.Sign() == -1 {
		beta.Add(beta, c.P)
	}
	y3 := alpha.Mul(alpha, beta) // y3 = alpha*(4*beta-x3)-8*gamma

	gamma.Mul(gamma, gamma)
	gamma.Lsh(gamma, 3)
	gamma.Mod(gamma, c.P)

	y3.Sub(y3, gamma)
	if y3.Sign() == -1 {
		y3.Add(y3, c.P)
	}
	y3.Mod(y3, c.P)

	return x3, y3, z3
}

// DoubleJacobian 相同的点相加，接口开放
// 返回仿射坐标结果
func (c *CurveWeierstrass) DoubleJacobian(Bx, By *big.Int) (*big.Int, *big.Int) {

	return c.affineFromJacobian(c.doubleJacobian(Bx, By, new(big.Int).SetInt64(1)))
}

// ScalarMult 点乘
func (c *CurveWeierstrass) ScalarMult(Bx, By, k *big.Int) (*big.Int, *big.Int) {
	// 时间消耗上
	// 1w次点乘 标准实现耗时 1.188 个人实现 1.299

	//下面为标准实现
	//Bz := new(big.Int).SetInt64(1)
	//x, y, z := new(big.Int), new(big.Int), new(big.Int)
	//
	//for _, b := range k.Bytes() {
	//	for bitNum := 0; bitNum < 8; bitNum++ {
	//		x, y, z = c.doubleJacobian(x, y, z)
	//		if b&0x80 == 0x80 {
	//			x, y, z = c.addJacobian(Bx, By, Bz, x, y, z)
	//		}
	//		b <<= 1
	//	}
	//}
	//return c.affineFromJacobian(x, y, z)

	//下面为个人实现

	binaryK := tools.BigNumBaseConversion(k, 2)

	x, y, z := new(big.Int).Set(Bx), new(big.Int).Set(By), new(big.Int).SetInt64(1)
	Rx, Ry, Rz := new(big.Int), new(big.Int), new(big.Int) // 初始化为无穷远点

	if string(binaryK[len(binaryK)-1]) == "1" {
		Rx, Ry, Rz = new(big.Int).Set(Bx), new(big.Int).Set(By), new(big.Int).SetInt64(1) // 满足在最小
	}

	for i := 1; i < len(binaryK); i++ {
		x, y, z = c.doubleJacobian(x, y, z)
		if string(binaryK[len(binaryK)-i-1]) == "1" {
			Rx, Ry, Rz = c.addJacobian(x, y, z, Rx, Ry, Rz)
		}
	}
	return c.affineFromJacobian(Rx, Ry, Rz)
}
