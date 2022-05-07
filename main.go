package main

import (
	"fmt"
	"gitgit/tools"
	"math/big"
)

func main()  {
	x,y := new(big.Int).SetInt64(5),new(big.Int).SetInt64(15)
	fmt.Println(tools.Gcd(x, y))
}


