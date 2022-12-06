package ecurve

import "math/big"

//初始化
func (k *PrivateKey) Init() *PrivateKey {
	k.K = new(big.Int)
	k.Pub = new(PublicKey)
	k.Pub.P = new(Point)
	return k
}

