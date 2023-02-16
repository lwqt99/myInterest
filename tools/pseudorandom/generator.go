package pseudorandom

import (
	"crypto/rand"
	"math/big"
)

func LengthDoublingPseudorandom(r *big.Int) (*big.Int, *big.Int) {
	r0, _ := rand.Int(rand.Reader, r)
	r1, _ := rand.Int(rand.Reader, r)
	//fmt.Println("*******************************************")
	//fmt.Println(r0, r1)
	return r0, r1
}
