package tools

import (
	"fmt"
	"math/big"
)

func TestCipolla() {
	fmt.Println(Cipolla(new(big.Int).SetInt64(5), new(big.Int).SetInt64(84906529)))
}
