package forward_secure

import (
	"crypto/sha256"
	"fmt"
	"gitgit/tools"
	"gitgit/tools/asymmetric"
	"gitgit/tools/pseudorandom"
	"math/big"
)

func (c *Sum) Init(r1, sk0, pk0, pk1, pk *big.Int) *Sum {
	c.r1 = new(big.Int).Set(r1)
	c.sk0 = new(big.Int).Set(sk0)
	c.pk0 = new(big.Int).Set(pk0)

	c.pk1 = new(big.Int).Set(pk1)
	//计算pk

	c.pk = new(big.Int).Set(pk)

	return c
}

// KeyGen 默认生成时间周期T=2的方案
func (s *Sum) KeyGen() *Sum {
	p := tools.GenerateBigPrimeP(30)
	g := new(big.Int).SetInt64(7)
	r := tools.GenerateBigIntByRange(p)

	//r --> r0, r1
	r0, r1 := pseudorandom.LengthDoublingPseudorandom(r)
	//fmt.Println(r0, r1)

	sk0, pk0 := asymmetric.KeyGen(r0, g, p)
	pk1 := asymmetric.PubKeyGen(r1, g, p)

	t := pk0.Bytes()
	t = append(t, pk1.Bytes()...)
	byteT := sha256.Sum256(t)

	pk := new(big.Int).SetBytes(byteT[:])

	s.Init(r1, sk0, pk0, pk1, pk)

	s.G = new(big.Int).Set(g)
	s.P = new(big.Int).Set(p)

	s.T0 = new(big.Int).SetInt64(1)
	s.T1 = new(big.Int).SetInt64(1)

	return s
}

func (ss *SumSigma) Init(pk0, pk1 *big.Int, sigma [2]*big.Int) *SumSigma {
	ss.Pk0 = new(big.Int).Set(pk0)
	ss.Pk1 = new(big.Int).Set(pk1)
	ss.Sigma[0], ss.Sigma[1] = new(big.Int).Set(sigma[0]), new(big.Int).Set(sigma[1])
	return ss
}

func (s *Sum) Sign(t, sk *big.Int, m []byte) (*SumSigma, *big.Int) {
	var sigma [2]*big.Int
	sigma[0], sigma[1] = asymmetric.Sign(m, s.G, s.P, sk)
	return new(SumSigma).Init(s.pk0, s.pk1, sigma), t
}

func (s *Sum) Verify(pk, t *big.Int, ss *SumSigma, m []byte) bool {
	//验证哈希值的情况
	temp := ss.Pk0.Bytes()
	temp = append(temp, ss.Pk1.Bytes()...)
	byteT := sha256.Sum256(temp)
	//验证哈希正确性 Ps.我觉得没什么用
	if new(big.Int).SetBytes(byteT[:]).String() != pk.String() {
		return false
	}
	//判断时间T
	if t.Cmp(s.T1) == 1 {
		//采用pk0验证
		return asymmetric.Verify(m, ss.Sigma[0], ss.Sigma[1], s.G, s.P, ss.Pk0)
	} else {
		return asymmetric.Verify(m, ss.Sigma[0], ss.Sigma[1], s.G, s.P, ss.Pk1)
	}
}

func (c *Sum) Update(t *big.Int) *Sum {
	if t.Cmp(c.T0) != -1 {
		//擦除sk0
		c.sk0.SetInt64(0)
		c.sk1.Set(c.r1)
		c.r1.SetInt64(0) //擦除随机数
	}

	return c
}

func (ps *ProductSigma) Init(pk *big.Int, sigma0, sigma1 [2]*big.Int) *ProductSigma {
	ps.Pk1 = new(big.Int).Set(pk)
	ps.Sigma0[0], ps.Sigma0[1] = new(big.Int).Set(sigma0[0]), new(big.Int).Set(sigma0[1])
	ps.Sigma1[0], ps.Sigma1[1] = new(big.Int).Set(sigma1[0]), new(big.Int).Set(sigma1[1])
	return ps
}

//签名

func (product *Product) Sign(t, sk *big.Int, m []byte) (*ProductSigma, *big.Int) {
	var sigma1 [2]*big.Int
	sigma1[0], sigma1[1] = asymmetric.Sign(m, product.G, product.P, sk)
	return new(ProductSigma).Init(product.Pk1, product.Sigma, sigma1), new(big.Int).Mod(t, product.T1)
}

// Verify 验证签名
/*

 */
func (product *Product) Verify(M []byte, ps *ProductSigma, t *big.Int) bool {
	//<1>：用pk验证 M：pk1 签名：ps 时间：t/T1 保证pk1合法
	v0 := asymmetric.Verify(product.Pk1.Bytes(), ps.Sigma0[0], ps.Sigma0[1], product.G, product.P, product.Pk)
	//fmt.Println("v0=", v0)
	//<2>：用pk1验证M
	v1 := asymmetric.Verify(M, ps.Sigma1[0], ps.Sigma1[1], product.G, product.P, product.Pk1)
	//fmt.Println("v1=", v1)
	return v0 && v1
}

//密钥生成

func (product *Product) KeyGen() {
	p := tools.GenerateBigPrimeP(30)
	g := new(big.Int).SetInt64(7)
	r := tools.GenerateBigIntByRange(p)

	//r --> r0, r1
	r0, r1 := pseudorandom.LengthDoublingPseudorandom(r)
	//r1 --> r1', r1''
	r1_, r1__ := pseudorandom.LengthDoublingPseudorandom(r1)

	sk0, pk := asymmetric.KeyGen(r0, g, p)
	sk1, pk1 := asymmetric.KeyGen(r1_, g, p)

	var sigma [2]*big.Int
	//sigma := product.Sign(tools.Zero, sk1, pk1.Bytes())
	sigma[0], sigma[1] = asymmetric.Sign(pk1.Bytes(), g, p, sk0) //用sk0对pk1进行签名，pk和sigma用于验证签名合法性
	//fmt.Println(asymmetric.Verify(pk1.Bytes(), sigma[0], sigma[1], g, p, pk))

	product.G = new(big.Int).Set(g)
	product.P = new(big.Int).Set(p)

	product.Sk0 = new(big.Int).Set(sk0)
	product.Sigma[0], product.Sigma[1] = new(big.Int).Set(sigma[0]), new(big.Int).Set(sigma[1])
	product.Sk1 = new(big.Int).Set(sk1)
	product.Pk1 = new(big.Int).Set(pk1)
	product.R1__ = new(big.Int).Set(r1__)
	product.Pk = new(big.Int).Set(pk)

	product.T0 = new(big.Int).SetInt64(2)
	product.T1 = new(big.Int).SetInt64(2)
	return
}

//密钥更新

func (product *Product) Update(t *big.Int) {
	//判断是否需要更新
	if new(big.Int).Mod(t, product.T1).String() == "0" {
		r_, r := pseudorandom.LengthDoublingPseudorandom(product.R1__)
		sk1, pk1 := asymmetric.KeyGen(r_, product.G, product.P) //生成了新的pk1
		var sigma [2]*big.Int
		sigma[0], sigma[1] = asymmetric.Sign(pk1.Bytes(), product.G, product.P, product.Sk0) //更新对pk1的签名

		product.R1__.Set(r)
		product.Sk1.Set(sk1)
		product.Pk1.Set(pk1)
		product.Sigma[0].Set(sigma[0])
		product.Sigma[1].Set(sigma[1])
	}

	return
}

func (product *Product) constructSumSum(sum1, sum2 *Sum) *Product {
	//将两个Sum方案合成，采用乘积的形式
	//嵌入两个通用的方案
	//然后执行方案自带的内容

	product.Universal1 = &Sum{
		G:   new(big.Int).Set(sum1.G),
		P:   new(big.Int).Set(sum1.P),
		r1:  new(big.Int).Set(sum1.r1),
		sk0: new(big.Int).Set(sum1.sk0),
		pk0: new(big.Int).Set(sum1.pk0),
		//sk1: new(big.Int).Set(sum1.sk1),
		pk1: new(big.Int).Set(sum1.pk1),
		pk:  new(big.Int).Set(sum1.pk),
		T0:  new(big.Int).Set(sum1.T0),
		T1:  new(big.Int).Set(sum1.T1),
	}

	return product
}

func (product *Product) constructProductSum(product_ *Product, sum *Sum) *Product {

	return product
}

// MakeTree 构建MMM的树
func MakeTree() {
	//利用 Sum 构造T=2的方案
	c := new(Sum).KeyGen()
	//利用 Product 构造 T = 2*2
	t := new(Sum).KeyGen() //t临时存储T=2的Sum构造方案，用于之后的乘积
	//r用于存储结果
	//初始输入 Sum * Sum
	//后续输入 Product * Sum
	M := "I'm the superintendent of this junior school"
	bM := []byte(M)
	//fmt.Println(c.r0)
	r := new(Product).constructSumSum(c, t)
	ss, _ := r.Universal1.Sign(tools.Zero, r.Universal1.getSk0(), bM)
	fmt.Println(r.Universal1.Verify(r.Universal1.getPk0(), tools.Zero, ss, bM))
	//fmt.Println(ss.Sigma)
	//fmt.Println(ss.Pk0)
	//fmt.Println(ss.Pk1)

	//t = new(Sum).KeyGen() //t临时存储T=2的Sum构造方案，用于之后的乘积
	//r.constructProductSum(r, t)

	return
}

func (product Product) Show() {
	fmt.Println("G:", product.G, " ")
	fmt.Println("P:", product.P, " ")
	fmt.Println("r'':", product.R1__, " ")
	fmt.Println("sk0:", product.Sk0, " ")
	fmt.Println("pk:", product.Pk, " ")
	fmt.Println("sk1:", product.Sk1, " ")
	fmt.Println("pk1:", product.Pk1, " ")
	fmt.Println("sigma:", product.Sigma, " ")
}
