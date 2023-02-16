package forward_secure

import "math/big"

func (s *Sum) getSk0() *big.Int {
	return s.sk0
}

func (s *Sum) getPk0() *big.Int {
	return s.pk0
}
