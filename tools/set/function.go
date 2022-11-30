package set

import "math/big"

var full = 1

func (s *Set) InitBigInt(arr []*big.Int, length int) (*Set) {
	s.M = make(map[interface{}]interface{})
	for i := 0; i < length; i++ {
		s.M[arr[i]] = full
	}
	return s
}

func (s *Set) Remove(key interface{}) {
	delete(s.M, key)
}

func (s *Set) Add(key interface{}) {
	s.M[key] = full
}

func (s *Set) Len() int {
	return len(s.M)
}