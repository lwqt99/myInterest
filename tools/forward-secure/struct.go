package forward_secure

import (
	"math/big"
)

// PrivateKey 定义通用private key
type PrivateKey interface {
}

// UniversalScheme 定义通用方案
type UniversalScheme interface {
	getSk0() *big.Int
	getPk0() *big.Int

	Sign(t, sk *big.Int, m []byte) (*SumSigma, *big.Int) //执行签名
	Verify(pk, t *big.Int, ss *SumSigma, m []byte) bool
}

// 加法

type Sum struct {
	//常规签名参数
	//R *big.Int
	G *big.Int
	P *big.Int
	//引入通用描述方案
	Universal1 UniversalScheme
	Universal2 UniversalScheme

	r0  *big.Int //将被及时擦除
	r1  *big.Int
	sk0 *big.Int
	pk0 *big.Int
	sk1 *big.Int
	pk1 *big.Int
	pk  *big.Int //pk = Hash(pk0, pk1)

	//时间周期
	T0 *big.Int
	T1 *big.Int
}

//乘法

type Product struct {
	//常规签名参数
	//R *big.Int
	G *big.Int
	P *big.Int
	//引入通用描述方案
	Universal1 UniversalScheme
	Universal2 UniversalScheme

	//前向安全拓展参数
	R1__  *big.Int
	Sk0   *big.Int
	Sigma [2]*big.Int
	Sk1   *big.Int
	Pk1   *big.Int
	Pk    *big.Int //pk = pk0

	//时间周期
	T0 *big.Int
	T1 *big.Int
	//T *big.Int//T = T0 + T1
}

// SumSigma 求和Sigma
type SumSigma struct {
	Sigma [2]*big.Int
	Pk0   *big.Int
	Pk1   *big.Int
}

// ProductSigma 乘法Sigma
type ProductSigma struct {
	Pk1    *big.Int
	Sigma0 [2]*big.Int //σ0
	Sigma1 [2]*big.Int //σ1
}
