package pair

import (
	"gitgit/tools"
	"math/big"
)

// 初始化表达式
func (l *line) Init(a, b, c *big.Int) *line {
	l.A.Set(a)
	l.B.Set(b)
	l.C.Set(c)

	return l
}

func (l *line) SetInt64(a, b, c int64) *line {
	l.A.SetInt64(a)
	l.B.SetInt64(b)
	l.C.SetInt64(c)

	return l
}

func (l *line) String() string {
	return l.A.String() + "x +" + l.B.String() + " y + " + l.C.String()
}

func (l *line) CalculateLineByIms(point1, point2 *ImaginaryNumber) *line {
	//fmt.Println("p1:", point1.String())
	//fmt.Println("p2:", point2.String())

	if point1.A.Cmp(point2.A) == 0 {
		l.A = new(big.Int).SetInt64(1)
		l.B = new(big.Int).SetInt64(0)
		l.C = new(big.Int).Set(new(big.Int).Mul(tools.Negative1, point1.A))
		return l
	}
	//(y - y1)*(x1 - x2) = (x - x1)(y1 - y2)
	//Ax + By + C
	l.A = new(big.Int).Sub(point2.B, point1.B) //A = y2 - y1
	l.B = new(big.Int).Sub(point1.A, point2.A) //B = x1 - x2
	l.C = new(big.Int).Sub(new(big.Int).Mul(new(big.Int).Sub(point2.A, point1.A), point1.B),
		new(big.Int).Mul(new(big.Int).Sub(point2.B, point1.B), point1.A)) //C = (x2 - x1)*y1 + (y2 - y1)*x1

	//若B小于0
	//调整方向
	if l.B.Cmp(new(big.Int).SetInt64(0)) == -1 {
		l.A.Mul(l.A, tools.Negative1)
		l.B.Mul(l.B, tools.Negative1)
		l.C.Mul(l.C, tools.Negative1)
	}

	//fmt.Println(l.String())
	return l
}

func (l *line) SpeMod(p *big.Int) *line {
	//这里采用特殊Mod
	l.A.Mod(l.A, p)
	l.B.Mod(l.B, p)
	l.C.Mod(l.C, p)

	//若大于 p / 2则反转
	if tools.CmpMidP(l.A, p) == 1 {
		l.A.Sub(p, l.A)
	}
	if tools.CmpMidP(l.B, p) == 1 {
		l.B.Sub(p, l.B)
	}
	if tools.CmpMidP(l.C, p) == 1 {
		l.B.Sub(p, l.C)
	}
	return l
}

// 代入虚数计算
func (l *line) CalculteImNum(point *ImaginaryNumber) *ImaginaryNumber {
	rPoint := new(ImaginaryNumber)

	//Ax + By + C
	rPoint.A = new(big.Int).Add(l.C, new(big.Int).Mul(point.A, l.A))
	rPoint.B = new(big.Int).Mul(l.B, point.B)

	return rPoint
}
