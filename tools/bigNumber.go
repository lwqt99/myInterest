package tools

import (
	"errors"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"strconv"
	"time"
)

var Zero = new(big.Int).SetInt64(0)
var Positive1 = new(big.Int).SetInt64(1)
var Negative1 = new(big.Int).SetInt64(-1)
var Positive2 = new(big.Int).SetInt64(2)
var Negative2 = new(big.Int).SetInt64(-2)

/*
big包下的数字转为其它进制数
适合大位数转化
返回字符串
返回：从大到小 例如 二进制 1000 = 8
*/
func BigNumBaseConversion(n *big.Int, base int) string {
	result := ""

	t := new(big.Int).Set(n)
	Base := new(big.Int).SetInt64(int64(base)) //目标进制
	modresult := new(big.Int).Mod(n, Base)
	divresult := new(big.Int).Div(n, Base)

	for divresult.String() != "0" {
		result = strconv.FormatInt(modresult.Int64(), base) + result
		t.Set(divresult)
		//更新余数和商
		modresult.Mod(t, Base)
		divresult.Div(t, Base)
	}
	result = strconv.FormatInt(modresult.Int64(), base) + result
	return result
}

// 判断是否为2次幂
func JudgePow2(x *big.Int) bool {
	strX := BigNumBaseConversion(x, 2)
	//查看是否有额外的1
	for i := 1; i < len(strX); i++ {
		if string(strX[i]) == "1" {
			return false
		}
	}
	return true
}

// Floor 向下取整 x / y
func Floor(x, y *big.Int) *big.Int {
	return new(big.Int).Div(x, y)
}

// Ceil 向上取整
func Ceil(x, y *big.Int) *big.Int {
	if new(big.Int).Mod(x, y).String() == "0" {
		return new(big.Int).Div(x, y)
	}else {
		return new(big.Int).Add(new(big.Int).Div(x, y), Positive1)
	}
}

// Gcd 求最大公约数
func Gcd(x, y *big.Int) *big.Int {
	modresult := new(big.Int)

	a := new(big.Int).Set(x)
	b := new(big.Int).Set(y)
	modresult.Mod(a, b)
	for modresult.String() != "0" {
		a = new(big.Int).Set(b)
		b = new(big.Int).Set(modresult)
		modresult.Mod(a, b)
	}
	return b
}

// Lcm 最小公倍数
func Lcm(x, y *big.Int) *big.Int {
	r := new(big.Int)
	r.Mul(x, y)
	r.Div(r, Gcd(x, y))
	return r
}

// Exgcd x>y 拓展欧几里得算法
func Exgcd(x *big.Int, y *big.Int) (*big.Int, *big.Int) {
	a := new(big.Int).Set(x)
	b := new(big.Int).Set(y)
	//var temp string
	t1 := new(big.Int).SetInt64(0)
	t2 := new(big.Int).SetInt64(1)

	s1 := new(big.Int).SetInt64(1)
	s2 := new(big.Int).SetInt64(0)

	moderesult := new(big.Int)
	q := new(big.Int)

	for b.String() != "0" {
		moderesult = new(big.Int).Mod(a, b)
		q = new(big.Int).Div(a, b) //(a-modresult)/b
		a, b = new(big.Int).Set(b), new(big.Int).Set(moderesult)
		s1, s2 = new(big.Int).Set(s2), new(big.Int).Sub(s1, new(big.Int).Mul(q, s2))
		t1, t2 = new(big.Int).Set(t2), new(big.Int).Sub(t1, new(big.Int).Mul(q, t2))
	}
	// t1 * y + s1 * x = gcd
	//fmt.Println(new(big.Int).Add(new(big.Int).Mul(t1, y), new(big.Int).Mul(s1, x)))
	return t1, s1
}

// RelativePrime 判断是否互素（relatively prime）
func RelativePrime(x, y *big.Int) bool {
	if Gcd(x, y).String() == "1" {
		return true
	}
	return false
}

// FindSmallRelativePrime 从小到大找出互素
func FindSmallRelativePrime(p *big.Int) *big.Int {
	return nil
}

// 如果a是素数，则(p ^ (a - 1)) % a恒等于1
func fmod(a *big.Int, p int64) bool {
	one, _ := new(big.Int).SetString("1", 10)
	a_ := new(big.Int).Sub(a, one)
	result := new(big.Int).Exp(new(big.Int).SetInt64(p), a_, a)
	if result.String() != "1" {
		return false //此时出错 返回false 结果必须要为1
	}
	return true
}

// MillerRabbin 素性检验
func MillerRabbin(a *big.Int) bool {

	p := new(big.Int).Set(a)

	rand.Seed(time.Now().UnixNano())
	//进行1000次检验
	for i := 1; i < 100; i++ {
		//判断失败则退出
		n := rand.Int63()
		if new(big.Int).SetInt64(n).Cmp(p) == 1 {
			n = rand.Int63n(p.Int64()-1) + 1
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
	re, _ := new(big.Int).SetString("10", 10)
	re.Exp(re, length, nil)
	return re
}

// GenerateBigPrimeP
/*
	用于生成大素数
	n是长度
*/
func GenerateBigPrimeP(n int64) *big.Int {
	numRange := GenerateBigRange(n)
	ran := rand.New(rand.NewSource(time.Now().UnixNano())) //创建的时候需要初始化其中一个值 用于生成随机数
	ran.Seed(time.Now().UnixNano())

	p := new(big.Int).Rand(ran, numRange)
	for !MillerRabbin(p) {
		p.Rand(ran, numRange) //更新p
	}

	return p
}

/*
随机数生成器
*/
func GenerateBigIntByRange(p *big.Int) *big.Int {
	r := new(big.Int)
	ran := rand.New(rand.NewSource(time.Now().UnixNano())) //创建的时候需要初始化其中一个值 用于生成随机数
	time.Sleep(100)                                        //避免重复生成一样的数值
	ran.Seed(time.Now().UnixNano())
	r.Rand(ran, p)
	return r
}

/*
基于大数的中国剩余定理
测试
c(i) = c mod m(i)
*/
func TBigNumCRT(n, primeRange int) {
	/*
		n:	用于测试的方程组数量
		primeRange:	生成测试的素数长度
	*/
	ran := rand.New(rand.NewSource(time.Now().UnixNano())) //创建的时候需要初始化其中一个值 用于生成随机数
	ran.Seed(time.Now().UnixNano())
	t := new(big.Int) //用于间接运算

	var m []*big.Int
	//生成m(i)
	for i := 0; i < n; i++ {
		m = append(m, GenerateBigPrimeP(int64(primeRange)))
	}
	//生成c(i)
	var ci []*big.Int
	for i := 0; i < n; i++ {
		t.Rand(ran, m[i])
		for !MillerRabbin(t) {
			t.Rand(ran, m[i]) //更新p
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
		_, t := Exgcd(Mi[i], m[i])
		Ni = append(Ni, t)
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
		fmt.Println("c % mi =", t.Mod(c, m[i]).Cmp(ci[i]))
	}

}

/*
基于大数的中国剩余定理
c(i) = y(i) mod p(i)
*/
func BigNumCRT(yi, pi []*big.Int) {
	//n:用于测试的方程组数量
	n := len(yi)
	t := new(big.Int)

	//计算同余方程组的解
	//P = p1*p2*...*pn
	P := new(big.Int).SetInt64(1)
	for i := 0; i < n; i++ {
		P.Mul(P, pi[i])
	}
	//P(i) = P / pi
	Pi := make([]*big.Int, n)
	for i := 0; i < n; i++ {
		Pi[i] = new(big.Int).Div(P, pi[i])
	}
	//Pi*Ni mod pi = 1
	Ni := make([]*big.Int, n)
	for i := 0; i < n; i++ {
		_, t := Exgcd(Pi[i], pi[i])
		Ni[i] = t
	}

	//最终计算c
	c := new(big.Int).SetInt64(0)
	for i := 0; i < n; i++ {
		t.Mul(yi[i], Ni[i])
		t = t.Mul(t, Pi[i])
		c.Add(c, t)
	}
	//验证同余方程是否成立 全部为0则成立
	for i := 0; i < n; i++ {
		fmt.Println("c % pi =", t.Mod(c, pi[i]).Cmp(yi[i]))
	}

}

// 计算nx mod p = y
// 需要在合适的时候退出
// 引入error
func MatchXY(x, y, p *big.Int) *big.Int {
	r := new(big.Int).Set(x)
	t := new(big.Int).SetInt64(1)
	for r.Cmp(y) != 0 {
		t.Add(t, Positive1)
		r.Mul(x, t)
		r.Mod(r, p)
	}
	return t
}

// 计算nx mod p = y
func MatchXYInt64(x, y, p int64) int64 {
	r := x
	t := int64(1)
	for r != y {
		t++
		r = x * t
		r = r % p
	}
	return t
}

/*
用于辅助下面的寻找生成元
*/
func HelpFindAllMebComG(P_ []*big.Int, n *big.Int) bool {
	length := len(P_)
	for i := 0; i < length; i++ {
		//不为1说明两者不互质
		if Gcd(P_[i], n).String() != "1" {
			return false
		}
	}
	return true
}

/*
找到各个Pi的共同的生成元g
输入Pi数组
原理为：找到与各个质数互质的值 满足最大公约数为1 即gcd为1
*/
func FindAllMebComG(P []*big.Int) *big.Int {
	n := len(P)
	var b []bool //布尔数组
	for i := 0; i < n; i++ {
		b = append(b, false)
	}
	var P_ []*big.Int
	temp := new(big.Int)
	for i := 0; i < n; i++ {
		P_ = append(P_, new(big.Int).Set(temp.Sub(P[i], Positive1)))
	}
	//找到共同的生成元
	re := new(big.Int).SetInt64(int64(2))
	for !HelpFindAllMebComG(P_, re) {
		re.Add(re, Positive1) //自增1
	}
	return re
}

/*
根据输入的bitsize生成大数
*/
func GenerateBigIntByByte(bitsize int64, rand io.Reader) (*big.Int, error) {
	b := make([]byte, bitsize)
	_, err := io.ReadFull(rand, b)
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(b), nil
}

// LegendreSymbol
/*
	Legendre symbol
	n 整数, p素数
*/
func LegendreSymbol(n, p *big.Int) (*big.Int, error) {
	//判断p是否为素数
	if !MillerRabbin(p) {
		return nil, errors.New("p should be prime number")
	}
	//重置n
	n.Mod(n, p)
	//计算 n^((p-1)/2) mod p
	t := new(big.Int).Sub(p, Positive1)
	t.Div(t, Positive2) //t = (p - 1) / 2

	LS := new(big.Int).Exp(n, t, p)
	//fmt.Println("LS=", LS.String())
	return LS, nil
}

// Cipolla
/*
	Cipolla算法
	解决二次剩余计算问题：求解 x*x = n mod p
*/
func Cipolla(n, p *big.Int) (*big.Int, error) {
	//判断是否存在解
	LS, err := LegendreSymbol(n, p)
	if err != nil || LS.String() != "1" {
		return nil, errors.New("There doesn't exit an solution. ")
	}
	//求解
	//随机生成r
	r := GenerateBigIntByRange(p)
	squareR := new(big.Int).Exp(r, Positive2, p)
	//重新计算LS
	for LS.String() != new(big.Int).Sub(p, Positive1).String() {
		r = GenerateBigIntByRange(p)
		squareR = new(big.Int).Exp(r, Positive2, p)
		//计算r^2 - n的勒让德符号
		t := new(big.Int).Sub(squareR, n)
		t.Mod(t, p) //t = r^2 - n
		LS, _ = LegendreSymbol(t, p)
	}
	//计算二次剩余
	//x := new(big.Int)
	//定义 w^2 = r^2 - n
	squareW := new(big.Int).Sub(squareR, n)
	squareW.Mod(squareW, p)
	//x = (a + w) ^ ((p+1)/2)
	//复数的快速幂
	t := new(big.Int).Add(p, Positive1)
	t.Div(t, Positive2) //t = ((p+1)/2)
	//定义虚数
	im := struct {
		x *big.Int
		y *big.Int
	}{}
	im.x = new(big.Int).SetInt64(1)
	im.y = new(big.Int).SetInt64(0)
	//定义Ax Ay
	Ax, Ay := new(big.Int).Set(r), new(big.Int).SetInt64(1)

	seq := BigNumBaseConversion(t, 2)
	//fmt.Println(seq)
	for i := 0; i < len(seq); i++ {
		if string(seq[len(seq)-i-1]) == "1" {
			//做乘法
			//new(big.Int).Mul(im.y, Ay)要乘上w*w
			im.x, im.y = new(big.Int).Add(new(big.Int).Mul(im.x, Ax), new(big.Int).Mul(new(big.Int).Mul(im.y, Ay), squareW)),
				new(big.Int).Add(new(big.Int).Mul(im.x, Ay), new(big.Int).Mul(im.y, Ax))

			im.x = im.x.Mod(im.x, p)
			im.y = im.y.Mod(im.y, p)
		}
		//Ax Ay 自乘
		Ax, Ay = new(big.Int).Add(new(big.Int).Mul(Ax, Ax), new(big.Int).Mul(new(big.Int).Mul(Ay, Ay), squareW)),
			new(big.Int).Mul(new(big.Int).Mul(Ax, Ay), Positive2)

		Ax = Ax.Mod(Ax, p)
		Ay = Ay.Mod(Ay, p)
	}
	//fmt.Println("x=", im.x, "y=", im.y)
	//fmt.Println("x=", new(big.Int).Sub(p, im.x), "y=", im.y)
	//验证
	//v := new(big.Int).Exp(im.x, Positive2, p)
	//fmt.Println("x^2=", v.String())

	return im.x, nil
}
