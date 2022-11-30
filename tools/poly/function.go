package poly

import "math/big"

func max(a, b int) int {
	if a > b {
		return a
	}else {
		return b
	}
}

func min(a, b int) int {
	if a > b {
		return b
	}else {
		return a
	}
}

func (p *Poly) Init(f []*big.Int) *Poly {
	p.F = make([]*big.Int,len(f))
	for i := 0; i < len(f); i++ {
		p.F[i] = f[i]
	}
	return p
}

func (p *Poly) InitSigs(args ...*big.Int) *Poly {
	p.F = make([]*big.Int,len(args))
	for i := 0; i < len(args); i++ {
		p.F[i] = args[i]
	}
	return p
}

func (p *Poly) SetBigInt(f *big.Int) *Poly {
	p.F = append(p.F, f)
	return p
}

func (p *Poly) Add(x,y *Poly) *Poly {
	max_ := max(len(x.F),len(y.F))
	f := make([]*big.Int,max_)
	if len(x.F)>len(y.F) {
		for i := 0; i < len(x.F); i++ {
			f[i] = x.F[i]
		}
		for i := 0; i < len(y.F); i++ {
			f[i].Add(f[i],y.F[i])
		}
	}else {
		for i := 0; i < len(y.F); i++ {
			f[i] = y.F[i]
		}
		for i := 0; i < len(x.F); i++ {
			f[i].Add(f[i],x.F[i])
		}
	}
	p.F = f
	return p
}

func (p *Poly) Mul(x, y *Poly) *Poly {
	n := len(x.F) + len(y.F) - 1
	r := make([]*big.Int, n)
	//初始化
	for i := 0; i < n; i++ {
		r[i] = new(big.Int).SetInt64(0)
	}
	//多项式乘法
	for i := 0; i < len(x.F); i++ {
		for j := 0; j < len(y.F); j++ {
			r [(i+j) % n].Add(r [(i+j) % n], new(big.Int).Mul(x.F[i],y.F[j]))
		}
	}
	p.F = r
	return p
}
