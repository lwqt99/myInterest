package tools
import (
	"math/big"
	"strconv"
)

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
		//result = SwitchNum2Str(modresult.Int64()) + result
		result = strconv.FormatInt(modresult.Int64(), base) + result
		n.Set(divresult)
		//更新余数和商
		modresult.Mod(n,Base)
		divresult.Div(n,Base)
	}
	result = strconv.FormatInt(modresult.Int64(), base) + result
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


//如果a是素数，则(p ^ (a - 1)) % a恒等于1
func fmod(a *big.Int,p int64) bool {
	one,_ := new(big.Int).SetString("1",10)
	a_ := new(big.Int).Sub(a,one)
	result := new(big.Int).Exp(new(big.Int).SetInt64(p),a_,a)
	if result.String()!="1" {
		return false//此时出错 返回false 结果必须要为1
	}
	return true
}


//MillerRabbin 素性检验
func MillerRabbin(a *big.Int) bool {

	p := new(big.Int).Set(a)

	rand.Seed(time.Now().UnixNano())
	//进行1000次检验
	for i := 1; i < 100; i++ {
		//判断失败则退出
		n := rand.Int63()
		if new(big.Int).SetInt64(n).Cmp(p) == 1 {
			n = rand.Int63n(p.Int64() - 1) + 1
		}
		if !fmod(p, n) {
			return false
		}
	}
	return true
}

