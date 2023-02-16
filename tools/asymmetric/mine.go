package asymmetric

import (
	"crypto/sha256"
	"gitgit/tools"
	"math/big"
)

func KeyGen(r, g, p *big.Int) (*big.Int, *big.Int) {
	sk := new(big.Int).Set(r)
	pk := new(big.Int).Exp(g, sk, p)

	return sk, pk
}

func PubKeyGen(r, g, p *big.Int) *big.Int {
	_, pk := KeyGen(r, g, p)
	return pk
}

func SecKeyGen(r, g, p *big.Int) *big.Int {
	return r
}

func Sign(m []byte, g, p, sk *big.Int) (*big.Int, *big.Int) {
	hashM := sha256.Sum256(m)
	bigIntM := new(big.Int).SetBytes(hashM[:])
	//fmt.Println("H=",bigIntM)

	r := tools.GenerateBigIntByRange(p)

	s1 := new(big.Int).Exp(g, r, p)
	s2 := new(big.Int).Sub(new(big.Int).Mul(bigIntM, sk), r)
	//s2.Mod(s2, new(big.Int).Sub(p, tools.Positive1))

	return s1, s2
}

func Verify(m []byte, s1, s2, g, p, pk *big.Int) bool {
	hashM := sha256.Sum256(m)
	bigIntM := new(big.Int).SetBytes(hashM[:])
	//fmt.Println("H=",bigIntM)

	waitVerify := new(big.Int).Mul(new(big.Int).Exp(g, s2, p), s1)
	waitVerify.Mod(waitVerify, p)

	shouldResult := new(big.Int).Exp(pk, bigIntM, p)

	//fmt.Println(waitVerify.String())
	//fmt.Println(shouldResult.String())

	return waitVerify.String() == shouldResult.String()
}
