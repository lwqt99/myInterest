package algebra

import (
	"github.com/pkg/errors"
	"math/big"
)

//支持多项式计算

//取值出来
func (p polynomial) GetVectorInt64() []int64 {
	return p.vectorInt64
}

func (p polynomial) GetVectorBigInt() []big.Int {
	return p.vectorBigInt
}

//构建A矩阵
func (p *polynomial) generateAInt64() {
	length := len(p.vectorInt64)
	A := make([][]int64, length)
	for i := 0; i < length; i++ {
		A[i] = append(A[i], make([]int64, length)...)
	}
	//赋值
	//规律：每一列开始为 (0-j+length)
	for i := 0; i < length; i++ {
		for j := 0; j < length; j++ {
			A[i][j] = p.vectorInt64[(length-j+i)%length]
		}
	}
}

//生成特殊矩阵A
func GenerateAInt64(vector []int64) [][]int64 {
	length := len(vector)
	A := make([][]int64, length)
	for i := 0; i < length; i++ {
		A[i] = append(A[i], make([]int64, length)...)
	}
	//赋值
	//规律：每一列开始为 (0-j+length)
	for i := 0; i < length; i++ {
		for j := 0; j < length; j++ {
			A[i][j] = vector[(length-j+i)%length]
		}
	}
	return A
}

func GenerateABigInt(vector []big.Int) [][]big.Int {
	length := len(vector)
	A := make([][]big.Int, length)
	for i := 0; i < length; i++ {
		A[i] = append(A[i], make([]big.Int, length)...)
	}
	//赋值
	//规律：每一列开始为 (0-j+length)
	for i := 0; i < length; i++ {
		for j := 0; j < length; j++ {
			A[i][j].Set(&vector[(length-j+i)%length])
		}
	}
	return A
}

func (p *polynomial) generateABigInt() {
	length := len(p.vectorBigInt)
	A := make([][]big.Int, length)
	for i := 0; i < length; i++ {
		A[i] = append(A[i], make([]big.Int, length)...)
	}
	//赋值
	//规律：每一列开始为 (0-j+length)
	for i := 0; i < length; i++ {
		for j := 0; j < length; j++ {
			A[i][j].Set(&p.vectorBigInt[(length-j+i)%length])
		}
	}
}

//多项式乘法
//在环上
//限制为N-1次
//x、y长度一致
func PolynomialMulBigInt(x, y []big.Int) (error,[]big.Int) {
	if len(x) != len(y) {
		return errors.New("x y should in the same ring"), nil
	}
	A := GenerateABigInt(x)
	if A == nil {
		return errors.New("you should initialize A first"), nil
	}
	//目标：A*v = r (N*N) * (N*1) --> N*1
	//多项式乘法规律：
	//A乘一下即可
	err,v := MatMulBigInt(A, TransposeBigInt(FormatVector2MatBigInt(y)))
	if err != nil {
		return err, nil
	}
	return FormatMat2VectorByCol(v)
}

func PolynomialMulInt64(x, y []int64) (error,[]int64) {
	if len(x) != len(y) {
		return errors.New("x y should in the same ring"), nil
	}
	A := GenerateAInt64(x)
	if A == nil {
		return errors.New("you should initialize A first"), nil
	}
	//目标：A*v = r (N*N) * (N*1) --> N*1
	//多项式乘法规律：
	//A乘一下即可
	err,v := MatMulInt64(A, TransposeInt64(FormatVector2MatInt64(y)))
	if err != nil {
		return err, nil
	}
	return nil, FormatMat2VectorByColInt64(v)
}

func (*polynomial)PolynomialMul(x, y *polynomial) (error,[]int64) {
	if len(x.vectorInt64) != len(y.vectorInt64) {
		return errors.New("x y should in the same ring"), nil
	}
	if x.AInt64 == nil {
		return errors.New("you should initialize A first"), nil
	}
	//目标：A*v = r (N*N) * (N*1) --> N*1
	//多项式乘法规律：
	//A乘一下即可
	err,v := MatMulInt64(x.AInt64, FormatVector2MatInt64(y.vectorInt64))
	if err != nil {
		return err, nil
	}
	return nil, FormatMat2VectorByRowInt64(v)
}


