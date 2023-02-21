package main

import forward_secure "gitgit/tools/forward-secure"

func main() {

	//text := ""
	//for true {
	//	t, _ := clipboard.ReadAll()
	//	if text != t {
	//		text = t
	//		fmt.Println(text)
	//	}
	//}

	//curve := new(ecurve.Curve).SetInt64(1,3,23)
	//point1 := new(ecurve.Point).SetInt64(2, 17)
	//point2 := new(ecurve.Point).SetInt64(4, 5)
	//fmt.Println(curve.Add(point1, point2).String())
	//fmt.Println(curve.VerifyPointExit(curve.Add(point1, point2)))
	//fmt.Println(curve.Mul(new(big.Int).SetInt64(6), point1))
	//asymmetric.TestTradition()

	//ecurve.TestMulJacobian()
	forward_secure.TestSumAdd()

	//k, _ := new(big.Int).SetString("9", 10)
	//fmt.Println(k.Lsh())

	//file.TestReadFilesByExt()
	//pair.TestMil()
	//ecurve.TestGenKeyCor()
	//fmt.Println(curve.VerifyPointExit(curve.Add(point1, point2)))
	//curve.ShowPoint()

	//exp.TestExp11(10)

}
