package poly

import (
	"gitgit/tools"
	"math/big"
)

//有限域q下的拉格朗日插值
func LagrangeInterBigInt(x, y []*big.Int, q *big.Int) *Poly {
	resultF := new(Poly).SetBigInt(new(big.Int).SetInt64(0))
	n := len(x)
	for i := 0; i < n; i++ {
		//初始化结果
		roundResult := new(Poly).SetBigInt(y[i])
		//构建(x-xi)(xi-xj)
		for j := 0; j < n; j++ {
			if j != i{
				//计算xi-xj的逆
				t := new(big.Int).Sub(x[i],x[j])
				t.Mod(t, q)
				tRev,_ := tools.Exgcd(q, t)//逆
				//构建
				tp := new(Poly).InitSigs(new(big.Int).Mul(tools.Negative1,new(big.Int).Mul(tRev,x[j])), tRev)
				roundResult.Mul(roundResult, tp)
			}
		}
		resultF.Add(resultF,roundResult)
	}
	return resultF
}
