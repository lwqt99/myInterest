package forward_secure

import (
	"fmt"
	"gitgit/tools"
)

func TestMMMSmall() {
	product := new(Product)
	product.KeyGen()
	//product.Show()
	M := "I'm the superintendent of this junior school"
	bM := []byte(M)

	ps, _ := product.Sign(tools.Zero, product.Sk1, bM) //sk1给消息签名
	//验证签名
	fmt.Println(product.Verify(bM, ps, tools.Zero))

	//product.Show()
	product.Update(product.T1)
	//product.Show()

	ps, _ = product.Sign(tools.Zero, product.Sk1, bM) //sk1给消息签名
	//验证签名
	fmt.Println(product.Verify(bM, ps, tools.Zero))
	return
}

func TestMMM() {
	MakeTree()
	return
}
