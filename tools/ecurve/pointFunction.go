package ecurve

import (
	"math/big"
)

var identityElement = new(Point).SetInt64(-1, -1)//椭圆曲线上的单位元

func (p *Point) Init(x, y *big.Int) *Point {
	p.X = new(big.Int).Set(x)
	p.Y = new(big.Int).Set(y)
	return p
}

func (p *Point) SetPoint(point *Point) *Point {
	p.X = new(big.Int).Set(point.X)
	p.Y = new(big.Int).Set(point.Y)
	return p
}

func (p *Point) SetInt64(x, y int64) *Point {
	p.X = new(big.Int).SetInt64(x)
	p.Y = new(big.Int).SetInt64(y)
	return p
}

func (p *Point) Equal(point *Point) bool {
	if p.X.Cmp(point.X) == 0 && p.Y.Cmp(point.Y) == 0 {
		return true
	}
	return false
}

//判断是否为单位元
func IE(point *Point) bool {
	if point.X.Cmp(identityElement.X) == 0 && point.Y.Cmp(identityElement.Y) == 0 {
		return true
	}
	return false
}

//输出String
func (p *Point) String() string {
	return "x=" + p.X.String() + ", y=" + p.Y.String()
}
