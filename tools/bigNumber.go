package tools

import (
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"time"
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
	return s1
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

/*
	用于提供长度为n的数，用于提取大素数
*/
func GenerateBigRange(n int64) *big.Int {
	length := new(big.Int).SetInt64(n)
	re,_ := new(big.Int).SetString("10",10)
	re.Exp(re,length,nil)
	return re
}

/*
	用于生成大素数
 */
func GenerateBigPrimeP(n int64) *big.Int {
	numRange := GenerateBigRange(n)
	ran := rand.New(rand.NewSource(time.Now().UnixNano()))//创建的时候需要初始化其中一个值 用于生成随机数
	ran.Seed(time.Now().UnixNano())

	p := new(big.Int).Rand(ran,numRange)
	for !MillerRabbin(p) {
		p.Rand(ran,numRange)//更新p
	}

	return p
}

/*
	基于大数的中国剩余定理
	c(i) = c mod m(i)
 */
func BigNumCRT(n, primeRange int)  {
	/*
		n:	用于测试的方程组数量
		primeRange:	生成测试的素数长度
	 */
	ran := rand.New(rand.NewSource(time.Now().UnixNano()))//创建的时候需要初始化其中一个值 用于生成随机数
	ran.Seed(time.Now().UnixNano())
	t := new(big.Int) //用于间接运算

	var m []*big.Int
	//生成m(i)
	for i:= 0; i<n; i++{
		m = append(m, GenerateBigPrimeP(int64(primeRange)))
	}
	//生成c(i)
	var ci []*big.Int
	for i:=0;i<n;i++{
		t.Rand(ran,m[i])
		for !MillerRabbin(t) {
			t.Rand(ran,m[i])//更新p
		}
		ci = append(ci, t)
	}
	//计算同余方程组的解
	//M = m(1)*m(2)*...*m(n)
	M := new(big.Int).SetInt64(1)
	for i := 0; i < n; i++ {
		M.Mul(M, m[i])
	}
	//M(i) = M / m(i)
	var Mi []*big.Int
	for i := 0; i < n; i++ {
		Mi = append(Mi, t.Div(M, m[i]))
	}
	//Mi*Ni mod mi = 1
	var Ni []*big.Int
	for i := 0; i < n; i++ {
		Ni = append(Ni, Exgcd(Mi[i], m[i]))
	}
	//验证Mi*Ni mod mi = 1是否成立
	//for i := 0; i < n; i++ {
	//	t := new(big.Int).Mul(Mi[i], Ni[i])
	//	t = t.Mod(t, m[i])
	//	fmt.Println(t.String())
	//}
	//最终计算c
	c := new(big.Int).SetInt64(0)
	for i := 0; i < n; i++ {
		t.Mul(ci[i], Ni[i])
		t = t.Mul(t, Mi[i])
		c.Add(c, t)
	}
	//验证同余方程是否成立 全部为0则成立
	for i := 0; i < n; i++ {
		fmt.Println("c % mi =",t.Mod(c, m[i]).Cmp(ci[i]))
	}

}
