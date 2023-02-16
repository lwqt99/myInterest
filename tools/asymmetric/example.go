package asymmetric

import (
	"fmt"
	"gitgit/tools"
	"math/big"
)

func TestTradition() {
	p := tools.GenerateBigPrimeP(30)
	g := new(big.Int).SetInt64(7)

	r := tools.GenerateBigIntByRange(p)

	sk, pk := KeyGen(r, g, p)

	s := "14894"
	bs := []byte(s)
	s1, s2 := Sign(bs, g, p, sk)
	//fmt.Println("p:", p)
	fmt.Println(Verify(bs, s1, s2, g, p, pk))
}
