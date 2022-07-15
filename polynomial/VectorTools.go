package polynomial

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/pkg/errors"
	"math/big"
	"math/rand"
	"time"
)

//格式化输出Vector
func ShowVector(v []big.Int) {
	for i := 0; i < len(v); i++ {
		fmt.Print(v[i].String()+" ")
	}
	fmt.Println()
}

//随机Vector
func RandomVectorBigInt(range_ int64, length int) []big.Int {
	rand.Seed(time.Now().UnixNano())

	v := make([]big.Int, length)
	for i := 0; i < length; i++ {
		v[i].SetInt64(rand.Int63n(range_))
	}
	return v
}

//行最小值计算
func MinVectorByRow(v []big.Int) big.Int {
	r := new(big.Int)
	r.Set(&v[0])
	for i := 1; i < len(v); i++ {
		if r.Cmp(&v[i]) == 1{
			r.Set(&v[i])
		}
	}
	return *r
}

//列最小值计算
//col为所在的列
//start为开始计算的列数（前面的不用比）
func MinVectorByCol(mat [][]big.Int, col, start int) (big.Int, int) {
	r := new(big.Int)
	r.SetInt64(math.MaxInt64)
	//r.Set(&mat[start][col])
	p := start//记录最小值所在的行数
	for i := start + 1; i < len(mat); i++ {
		if r.Cmp(&mat[i][col]) == 1 && mat[i][col].String() != "0"{
			r.Set(&mat[i][col])
			p = i
		}
	}
	return *r, p
}

//列最小值计算
//col为所在的列
//start为开始计算的列数（前面的不用比）
func MinVectorByColInt64(mat [][]int64, col, start int) (int64, int) {
	var r int64 = math.MaxInt64
	//r.Set(&mat[start][col])
	p := start//记录最小值所在的行数
	for i := start + 1; i < len(mat); i++ {
		if r > mat[i][col] && mat[i][col] != 0{
			r= mat[i][col]
			p = i
		}
	}
	return r, p
}

//两个Vector做差
//计算方式如下
//p*vector1 - q*vector2
func VectorSub(vector1, vector2 []big.Int, p, q *big.Int) (error,[]big.Int) {
	if len(vector1) != len(vector2) {
		return errors.New("dim does not match"),nil
	}
	r := make([]big.Int, len(vector1))

	for i := 0; i < len(r); i++ {
		r[i].Sub(
			new(big.Int).Mul(&vector1[i], p),
			new(big.Int).Mul(&vector2[i], q))
	}

	return nil,r[:]
}

//两个Vector做差
//计算方式如下
//p*vector1 - q*vector2
func VectorSubInt64(vector1, vector2 []int64, p, q int64) (error,[]int64) {
	if len(vector1) != len(vector2) {
		return errors.New("dim does not match"),nil
	}
	r := make([]int64, len(vector1))

	for i := 0; i < len(r); i++ {
		r[i] = vector1[i]*p - vector2[i]*q
	}

	return nil,r[:]
}


//两个Vector相加
//计算方式如下
//p*vector1 + q*vector2
func VectorAdd(vector1, vector2 []big.Int, p, q *big.Int) (error,[]big.Int) {
	if len(vector1) != len(vector2) {
		return errors.New("dim does not match"),nil
	}
	r := make([]big.Int, len(vector1))

	for i := 0; i < len(r); i++ {
		r[i].Add(
			new(big.Int).Mul(&vector1[i], p),
			new(big.Int).Mul(&vector2[i], q))
	}

	return nil,r[:]
}

//Vector乘积
func VectorMulEle(vector []big.Int, p *big.Int) []big.Int {
	r := make([]big.Int, len(vector))
	for i := 0; i < len(vector); i++ {
		r[i].Mul(&vector[i], p)
	}
	return r
}

//Vector求模
func VectorMod(vector []big.Int, p *big.Int) []big.Int {
	r := make([]big.Int, len(vector))
	for i := 0; i < len(vector); i++ {
		r[i].Mod(&vector[i], p)
	}
	return r
}

//Vector求模
func VectorModInt64(vector []int64, p int64) []int64 {
	r := make([]int64, len(vector))
	for i := 0; i < len(vector); i++ {
		//如果r[i]为负数则求模结果异常
		r[i] = (vector[i] % p + p) % p
	}
	return r
}