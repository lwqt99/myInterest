package tools

import "math/big"


/*
	对应数字转为进制
	11 --> a
	支持到16进制
*/
func SwitchNum2Str(n int64) string {
	switch n {
	case 0:
		return "0"
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "3"
	case 4:
		return "4"
	case 5:
		return "5"
	case 6:
		return "6"
	case 7:
		return "7"
	case 8:
		return "8"
	case 9:
		return "9"
	case 10:
		return "a"
	case 11:
		return "b"
	case 12:
		return "c"
	case 13:
		return "d"
	case 14:
		return "e"
	case 15:
		return "f"
	default:
		return ""
	}
}


/*
	big包下的数字转为其它进制数
	适合大位数转化
	返回字符串
*/
func BigNumBaseConversion(n *big.Int,base int) string {
	result := ""

	modresult := new(big.Int)
	divresult := new(big.Int)

	Base := new(big.Int).SetInt64(int64(base))//目标进制

	modresult.Mod(n,Base)
	divresult.Div(n,Base)

	for divresult.String()!="0" {
		result = SwitchNum2Str(modresult.Int64()) + result
		n.Set(divresult)

		//更新余数和商
		modresult.Mod(n,Base)
		divresult.Div(n,Base)
	}
	result = modresult.String() + result

	return result
}


//求最大公约数
func Gcd(x *big.Int,y *big.Int) *big.Int {
	modresult := new(big.Int)

	a := new(big.Int).Set(x)
	b := new(big.Int).Set(y)

	modresult.Mod(a,b)

	for modresult.String()!="0" {
		a = new(big.Int).Set(b)
		b = new(big.Int).Set(modresult)
		modresult.Mod(a,b)
	}

	return b
}


//x>y 拓展欧几里得算法
func Exgcd(x *big.Int,y *big.Int) *big.Int {
	one,_ := new(big.Int).SetString("1",10)
	two,_ := new(big.Int).SetString("1",10)

	a := new(big.Int).Set(x)
	b := new(big.Int).Set(y)
	//var temp string
	t1,_ := new(big.Int).SetString("0",10)
	t2,_ := new(big.Int).SetString("1",10)

	s1,_ := new(big.Int).SetString("1",10)
	s2,_ := new(big.Int).SetString("0",10)

	moderesult,_ := new(big.Int).SetString("1",10)//初始化
	q,_ := new(big.Int).SetString("1",10)//初始化

	for b.String() != "0"  {
		moderesult = new(big.Int).Set(one.Mod(a,b))
		q = new(big.Int).Set(one.Div(a,b))//(a-modresult)/b
		a, b = new(big.Int).Set(b), new(big.Int).Set(moderesult)
		s1,s2 = new(big.Int).Set(s2),new(big.Int).Set(one.Sub(s1,two.Mul(q,s2)))
		t1,t2 = new(big.Int).Set(t2),new(big.Int).Set(one.Sub(t1,two.Mul(q,t2)))
	}
	// t1 * b + s1 * a = q
	return t1
}

