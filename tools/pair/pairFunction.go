package pair

import (
	"gitgit/tools"
	"gitgit/tools/ecurve"
	"math/big"
)

func (pair *Pair) Init(G1, G2 *ecurve.CurveWeierstrass) *Pair {
	//p素数域的阶
	//r椭圆曲线子群的阶
	//k嵌入度满足  r|p^k-1
	//e(P,Q) = fr(Qx,Qy)^((p^k -1)/r)
	pair.G1 = new(ecurve.CurveWeierstrass).SetCurveWeierstrass(G1)
	pair.G2 = new(ecurve.CurveWeierstrass).SetCurveWeierstrass(G2)

	pair.P = new(big.Int).Set(G1.P)
	pair.R = new(big.Int).Set(G1.N)

	i := int64(1)
	for t := new(big.Int).Exp(pair.P, tools.Positive1, pair.R); t.String() != "1"; i++ {
		t.Exp(pair.P, new(big.Int).SetInt64(i), pair.R)
		//fmt.Println("i=",i,"求余",new(big.Int).Mod(new(big.Int).Sub(new(big.Int).Exp(pair.P, new(big.Int).SetInt64(i), nil), tools.Positive1), pair.R))
	}
	pair.K = new(big.Int).SetInt64(i - 1)

	return pair
}

// 递推关系如下
// f1 = 1
// f_{i+1} = f_i * l_{iP,P}
// f_{2i} = (f_{i})^2 * l_{iP, -2iP}
func fr(r *big.Int, P, Q *ecurve.Point, pair *Pair, c *ecurve.CurveWeierstrass) *ImaginaryNumber {
	if r.String() == "1" {
		//返回1 + 0i
		//终止递归
		return new(ImaginaryNumber).SetInt64(1, 0)
	}

	//判断r是否为2的倍数
	if tools.IsEven(r) {
		//二的倍数 r=2i
		//i := new(big.Int).SetInt64(int64(len(tools.BigNumBaseConversion(r, 2)) - 1) / 2)
		i := new(big.Int).Set(new(big.Int).Div(r, tools.Positive2))
		//return f_r = (f_i)^2
		//计算表达式l
		p1 := new(ImaginaryNumber).SetPoint(c.Mul(i, P))
		negative2i := new(big.Int).Mul(tools.Negative2, i) //-2i可以模子群的阶 即 生成元的阶
		p2 := new(ImaginaryNumber).SetPoint(c.Mul(new(big.Int).Mod(negative2i, pair.R), P))

		l := new(line).CalculateLineByIms(p1, p2)
		//l.SpeMod(pair.P)
		//Q代入表达式
		rl := l.CalculteImNum(new(ImaginaryNumber).SetPoint(Q))
		//求模
		rl.Mod(rl, c.P)
		//fmt.Println(l.String())
		//fmt.Println(rl.String())
		//fmt.Println("****************************************************************")

		tr := fr(i, P, Q, pair, c) //记录f_{i^2}，用于求平方
		//fmt.Println("tr:", tr.String())
		//fmt.Println(new(ImaginaryNumber).Mul(tr, tr))
		//fmt.Println("************************************************")
		tMod := new(ImaginaryNumber).Mul(new(ImaginaryNumber).Mul(tr, tr), rl)
		return tMod.Mod(tMod, c.P)
	} else {
		//非二次幂
		i := new(big.Int).Sub(r, tools.Positive1)
		//计算表达式l
		p1 := new(ImaginaryNumber).SetPoint(c.Mul(i, P))
		l := new(line).CalculateLineByIms(p1, new(ImaginaryNumber).SetPoint(P))
		//l.SpeMod(pair.P)
		//Q代入表达式
		rl := l.CalculteImNum(new(ImaginaryNumber).SetPoint(Q))
		//求模
		rl.Mod(rl, c.P)
		//fmt.Println(l.String())
		//fmt.Println(rl.String())
		//fmt.Println("****************************************************************")

		tr := fr(i, P, Q, pair, c) //记录f_i
		tMod := new(ImaginaryNumber).Mul(tr, rl)
		return tMod.Mod(tMod, c.P)
	}

}

// Miller 算法
// 计算e(P,Q)
func (pair *Pair) Miller(P, Q *ecurve.Point, c *ecurve.CurveWeierstrass) *ImaginaryNumber {
	//e(P,Q) = fr(Qx,Qy)^((p^k -1)/r)
	//a = (p^k -1)/r
	a := new(big.Int).Exp(pair.P, pair.K, nil)
	a.Div(a, pair.R)
	//fmt.Println(a)
	r := fr(pair.R, P, Q, pair, c)
	r.Exp(r, a, pair.P)
	r.Mod(r, pair.P)
	return r
}

//BLNK算法
