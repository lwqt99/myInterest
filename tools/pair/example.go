package pair

import (
	"fmt"
	"gitgit/tools/ecurve"
)

func TestMil() {
	c := new(ecurve.CurveWeierstrass).SetInt64(0, 3, 17, 101)
	c.G = new(ecurve.Point).SetInt64(1, 1)

	P := new(ecurve.Point).SetInt64(12, 32)
	Q := new(ecurve.Point).SetInt64(10, 16)

	pair := new(Pair).Init(c, c)
	//计算 (p^k - 1 )/ r
	//t := new(big.Int).Div(new(big.Int).Sub(new(big.Int).Exp(pair.P, pair.K, nil), tools.Positive1), pair.R)
	//fmt.Println("除法",new(big.Int).Div(new(big.Int).Sub(new(big.Int).Exp(pair.P, pair.K, nil), tools.Positive1), pair.R),
	//	"求余",new(big.Int).Mod(new(big.Int).Sub(new(big.Int).Exp(pair.P, pair.K, nil), tools.Positive1), pair.R))
	fmt.Println(pair.Miller(P, Q, c).String())

	//p1 := new(ImaginaryNumber).SetInt64(55, 96)
	//tp := new(ImaginaryNumber).Mul(p1, p1)
	//tp.Mod(tp, pair.P)
	//fmt.Println(tp)

	//c.ShowPoint()
}
