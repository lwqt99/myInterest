package secretshare

import (
	"gitgit/tools"
	"gitgit/tools/poly"
	"gitgit/tools/set"
	"github.com/pkg/errors"
	"math/big"
)

/*
	q:有限域q中
 */
func Shamir(q, secret *big.Int, n int) ([]*big.Int,[]*big.Int,error) {
	//判断q是否为素数
	if !tools.MillerRabbin(q){
		return nil,nil,errors.New("q should be prime number")
	}
	//判断n是否≥1
	if n < 1{
		return nil,nil,errors.New("n should be greater than 1")
	}
	//生成秘密
	//f(x) = secret + a1x + a2x^2 + ... + a(n-1)x^(n-1)
	//生成n-1个随机数
	a := make([]*big.Int, n)
	a[0] = new(big.Int).Set(secret)
	for i := 1; i < n; i++ {
		a[i] = tools.GenerateBigIntByRange(q)
		for new(set.Set).InitBigInt(a, i+1).Len() != i + 1 {
			a[i] = tools.GenerateBigIntByRange(q)
		}
	}
	//生成n个随机数x
	x := make([]*big.Int, n)
	for i := 0; i < n; i++ {
		x[i] = tools.GenerateBigIntByRange(q)
		for new(set.Set).InitBigInt(a, i+1).Len() != i + 1 {
			x[i] = tools.GenerateBigIntByRange(q)
		}
	}
	//代入方程计算
	fi := make([]*big.Int, n)
	for i := 0; i < n; i++ {
		fi[i] = new(big.Int).SetInt64(0)
		for j := 0; j < n; j++ {
			fi[i].Add(fi[i],new(big.Int).Mul(a[j],new(big.Int).Exp(x[i], new(big.Int).SetInt64(int64(j)), q)))
		}
		fi[i].Mod(fi[i], q)
	}

	return x,fi,nil
}

func SolveShamir(x, y []*big.Int, q *big.Int) *big.Int {
	f := poly.LagrangeInterBigInt(x, y, q)
	return f.F[0].Mod(f.F[0], q)
}