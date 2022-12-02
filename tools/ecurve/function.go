package ecurve

import (
	"fmt"
	"gitgit/tools"
	"math/big"
	"strconv"
)



func (c *Curve) Init(a, b, p, n *big.Int, g *Point) *Curve {
	//y^2 = x^3 + ax + b
	c.P = new(big.Int).Set(p)
	c.A = new(big.Int).Set(a)
	c.B = new(big.Int).Set(b)
	c.N = new(big.Int).Set(n)

	c.G = new(Point).SetPoint(g)
	return c
}

func (c *Curve) SetInt64(a, b, p int64) *Curve {
	c.P = new(big.Int).SetInt64(p)
	c.A = new(big.Int).SetInt64(a)
	c.B = new(big.Int).SetInt64(b)
	return c
}

//椭圆曲线上逆元
//A+A' = 0, A'为A的对称点
func (c *Curve) revPoint(p *Point) *Point {
	return new(Point).Init(p.X, new(big.Int).Sub(c.P, p.Y))
}

//比较和p/2的关系
func (c *Curve) cmpMidP(x *big.Int) int64 {
	mid := new(big.Int).Div(c.P, tools.Positive2)
	if x.Cmp(mid) != 1 {
		return -1
	}
	return 1
}

//验证点是否在曲线上
func (c *Curve) VerifyPointExit(point *Point) bool {
	//无穷点直接返回正确
	if IE(point) {
		return true
	}

	//验证点是否在曲线上
	t1 := new(big.Int).Exp(point.X, new(big.Int).SetInt64(3), c.P)
	t1.Add(t1, c.B)
	t1.Add(t1, new(big.Int).Mul(c.A, point.X))
	t1.Mod(t1, c.P)

	t2 := new(big.Int).Exp(point.Y, tools.Positive2, c.P)

	return t1.String() == t2.String()
}

//两点相加的斜率
//需要审慎考虑mod的使用
//应当定义 小于 p/2 的y为负数
func (c *Curve) getSlope(point1, point2 *Point) *big.Int {
	//fmt.Println(point2.Y,point1.Y)
	//判断两点是否相等
	if point1.Equal(point2) {
		//相等
		//判断正负
		//求 (3x1^2 + a) / 2y1
		if point1.Y.String() == "0" {
			return nil
		}
		//求2y1的逆
		t := new(big.Int).Mul(point1.Y, tools.Positive2)
		t.Mod(t, c.P)
		tRev,_ := tools.Exgcd(c.P, t)

		tRev.Mul(tRev, new(big.Int).Add(new(big.Int).Mul(new(big.Int).SetInt64(3),new(big.Int).Exp(point1.X,tools.Positive2,c.P)),c.A))
		//tRev.Mod(tRev, c.P)
		return tRev.Mod(tRev, c.P)
	}else {
		//不相等
		//求(y2-y1)/(x2-x1)
		if point2.X.Cmp(point1.X) == 0 {
			return nil
		}
		//求x2-x1的逆
		t := new(big.Int).Sub(point2.X, point1.X)
		tRev,_ := tools.Exgcd(c.P, t)

		//fmt.Println("tRev=",tRev)
		tRev.Mul(tRev, new(big.Int).Sub(point2.Y,point1.Y))
		//fmt.Println("tRev=",tRev)

		return tRev.Mod(tRev, c.P)
	}
}

//椭圆曲线上的点相加
//做直线，找交点，选择对称点
//重合则是切线
func (c *Curve) Add(point1, point2 *Point) *Point {
	//fmt.Println(point1)
	//fmt.Println(point2)
	//fmt.Println("******************************************")
	if !c.VerifyPointExit(point1) || !c.VerifyPointExit(point2) {
		fmt.Println("no such point exit in this curve")
		return nil
	}
	//如果存在无穷点则应该直接返回结果
	if IE(point1) {
		//fmt.Println("point1 is Identity Element")
		return point2
	}else if IE(point2) {
		//fmt.Println("point2 is Identity Element")
		return point1
	}
	//求切线斜率
	s := c.getSlope(point1, point2)
	//s = nil 则在无穷点
	if s == nil{
		//fmt.Println("result is Identity Element")
		return new(Point).SetPoint(identityElement)
	}
	//计算结果
	//x3 = s^2 - 2x1
	x := new(big.Int).Exp(s, tools.Positive2, c.P)
	x.Sub(x, point1.X)
	x.Sub(x, point2.X)
	x.Mod(x, c.P)
	//y3 = m(x1-x3)-y1
	y := new(big.Int).Sub(point1.X, x)
	y.Mul(y, s)
	y.Sub(y, point1.Y)
	y.Mod(y, c.P)

	r := new(Point).Init(x, y)

	return r
}

//多倍点
func (c *Curve) Mul(k *big.Int, p *Point) *Point {
	//k转二进制
	binaryK := tools.BigNumBaseConversion(k, 2)
	storeV := make([]*Point, len(binaryK))//存储 (2^i)p
	r := new(Point).SetPoint(identityElement)//结果

	//计算多倍点的二次幂
	storeV[0] = new(Point).SetPoint(p)
	for i := 1; i < len(binaryK); i++ {
		storeV[i] = c.Add(storeV[i - 1], storeV[i - 1])
	}
	//求和
	for i := 0; i < len(binaryK); i++ {
		if string(binaryK[len(binaryK) - i - 1]) == "1" {
			r = c.Add(r, storeV[i])
		}
	}
	//fmt.Println()

	return r
}


func (c *Curve) ShowPoint()  {
	point := make(map[string]string)
	//值得注意，可以将范围缩小至 p/2
	// y^2 = n mod p
	//-->y^2 + p*p - 2*y*p = n mod p
	//-->(p - y)^2 = n mod p
	for y := int64(0); y < c.P.Int64(); y++ {
		//计算y^2
		t := new(big.Int).Exp(new(big.Int).SetInt64(y), tools.Positive2, c.P).String()
		if _,ok := point[t];ok{
			point[new(big.Int).Exp(new(big.Int).SetInt64(y), tools.Positive2, c.P).String()] += "/" + 	strconv.Itoa(int(y))
		}else {
			point[new(big.Int).Exp(new(big.Int).SetInt64(y), tools.Positive2, c.P).String()] = strconv.Itoa(int(y))
		}

	}
	n := 0
	for x := int64(0); x < c.P.Int64(); x++ {
		//计算x^3 + ax + b
		t := new(big.Int).Exp(new(big.Int).SetInt64(x), new(big.Int).SetInt64(3), c.P)
		t.Add(t, new(big.Int).Mul(c.A, new(big.Int).SetInt64(x)))
		t.Add(t, c.B)
		t.Mod(t, c.P)

		_, ok := point[t.String()]
		if ok{
			n++
			fmt.Println("x=",x,"y=",point[t.String()])
		}
	}
	fmt.Println("元素的个数（含无穷点）=",n*2+1)
}
