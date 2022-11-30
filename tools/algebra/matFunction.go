package algebra

import (
	"fmt"
	"gitgit/tools"
	"github.com/pkg/errors"
	"math/big"
)

var negative1 = new(big.Int).SetInt64(-1)
var positive1 = new(big.Int).SetInt64(1)
var zero = new(big.Int).SetInt64(0)

//初始化数组
func (m *Mat) Initialization(width, height int) error {
	if height <= 0 || width <= 0 {
		return errors.New("width or height should be positive!")
	}
	m.width = width
	m.height = height
	for i := 0; i < height; i++ {
		m.mat = append(m.mat, make([]big.Int, m.width))
		m.matInt64 = append(m.matInt64, make([]int64, m.width))
	}
	return nil
}

//深度拷贝mat
func (m *Mat) SetBigInt(mat [][]big.Int) error {
	height, width := len(mat), len(mat[0])
	m.mat = make([][]big.Int, height)
	for i := 0; i < height; i++ {
		if len(mat[i]) != width {
			return errors.New("col size should be equal")
		}
		m.mat[i] = append(m.mat[i], make([]big.Int, len(mat[i]))...)
		copy(m.mat[i], mat[i])
	}
	m.height, m.width = height, width
	return nil
}

//深度拷贝mat
func (m *Mat) SetInt64(mat [][]int64) error {
	height, width := len(mat), len(mat[0])
	m.matInt64 = make([][]int64, height)
	for i := 0; i < height; i++ {
		if len(mat[i]) != width {
			return errors.New("col size should be equal")
		}
		m.matInt64[i] = append(m.matInt64[i], make([]int64, len(mat[i]))...)
		copy(m.matInt64[i], mat[i])
	}
	m.height, m.width = height, width
	return nil
}

//将Int64矩阵copy至BigInt里
func (m *Mat) SetBigIntByInt64() error {
	mat := m.matInt64
	height, width := len(mat), len(mat[0])
	m.mat = make([][]big.Int, height)
	for i := 0; i < height; i++ {
		if len(mat[i]) != width {
			return errors.New("col size should be equal")
		}
		m.mat[i] = append(m.mat[i], make([]big.Int, len(mat[i]))...)
		for j := 0; j < len(mat[i]); j++ {
			m.mat[i][j].SetInt64(mat[i][j])
		}
	}
	m.height, m.width = height, width
	return nil
}

func (m *Mat) GetMat() [][]big.Int {
	return m.mat
}

//行列式计算——递归法 高位矩阵很慢，不可行
func DetRecursion(mat [][]big.Int) *big.Int {
	//2*2时直接返回结果
	length := len(mat)
	if length == 2 {
		return new(big.Int).Sub(
			new(big.Int).Mul(&mat[0][0],&mat[1][1]),
			new(big.Int).Mul(&mat[0][1],&mat[1][0]))
	}

	r := new(big.Int).SetInt64(0)
	indictor := new(big.Int).SetInt64(1)//指示符变量，判断相加正负情况
	t := new(big.Int)//用于存储计算结果

	t.Mul(DetRecursion(MatIndex(mat,1,1,length-1,length-1)), indictor)
	t.Mul(t, &mat[0][0])
	r.Add(r, t)
	//fmt.Println(t.String())
	for i := 1; i < length - 1; i++ {
		indictor.Mul(indictor, negative1) //与负1相乘
		t.Mul(DetRecursion(MatIndexInternal(mat,0,i-1,i+1,length-1,1,length-1)), indictor)
		t.Mul(t, &mat[0][i])
		r.Add(r, t)
		//fmt.Println(t.String())
	}
	indictor.Mul(indictor, negative1) //与负1相乘
	t.Mul(DetRecursion(MatIndex(mat,0,1,length-2,length-1)), indictor)
	t.Mul(t, &mat[0][length-1])
	//fmt.Println(t.String())
	r.Add(r, t)

	return r
}

//递归法计算矩阵行列式
func (m *Mat) DetRecursion() *big.Int {
	//ShowMat(m.mat)
	//fmt.Println(det(m.mat))
	return DetRecursion(m.mat)
}

//行变换计算行列式——复杂度 O(n^3)
//p为模数
//p限制为质数主要因为需要找寻nx - y = mp
//只有在质数前提下才能保证每个元素都为生成元
func DetBigInt(m [][]big.Int, p *big.Int) (error,*big.Int) {
	//对p进行素性检验
	if !tools.MillerRabbin(p){
		return errors.New("p should be prime number"), nil
	}
	//深度拷贝mat
	length := len(m)
	mat := make([][]big.Int, length)
	for i := 0; i < length; i++ {
		mat[i] = append(mat[i], make([]big.Int,length)...)
		copy(mat[i], m[i])
	}
	//预处理mat 先进行求模
	//这一步主要是避免在MatchXY时无法匹配（元素得在域里）
	for i := 0; i < length; i++ {
		for j := 0; j < length; j++ {
			mat[i][j].Mod(&mat[i][j], p)
		}
	}

	indictor := new(big.Int).SetInt64(1)//指示符变量，判断行变换的乘积正负号
	for i := 0; i < length; i++ {
		//行最小计算——不要0
		_,index := MinVectorByCol(mat[:][:], i, i)
		//行交换
		if index != i {
			mat[i], mat[index] = mat[index], mat[i]
			indictor.Mul(indictor, negative1)
		}
		//如果非0的最小值交换完成后仍然为0，则应该跳过（其实值直接为0）
		if mat[i][i].String() == "0" {
			continue
		}
		//遍历下面的每一行
		//i行与j行做差
		for j := i + 1; j < length; j++ {
			//x和y元素之间在模p前提下的倍数关系
			//nx - y = mp
			//计算比率
			n := tools.MatchXY(&mat[i][i],&mat[j][i], p)
			//行间做差细化到每一个元素
			err,vector := VectorSub(mat[j], mat[i], positive1, n)
			if err != nil {
				return err,nil
			}
			vector = VectorMod(vector, p)
			mat[j] = vector
		}
	}
	r := new(big.Int).SetInt64(1)
	for i := 0; i < length; i++ {
		r.Mul(r, &mat[i][i])
	}
	r.Mul(r, indictor)
	r.Mod(r, p)
	return nil,r
}

func (m *Mat)DetBigInt(p *big.Int) (error,*big.Int) {
	err, r := DetBigInt(m.mat, p)
	return err,r
}

func DetInt64(m [][]int64, p int64) (error,*big.Int) {
	//对p进行素性检验
	if !tools.MillerRabbin(new(big.Int).SetInt64(p)){
		return errors.New("p should be prime number"), nil
	}
	//深度拷贝mat
	length := len(m)
	mat := make([][]int64, length)
	for i := 0; i < length; i++ {
		mat[i] = append(mat[i], make([]int64,length)...)
		copy(mat[i], m[i])
	}
	//预处理mat 先进行求模
	//这一步主要是避免在MatchXY时无法匹配（元素得在域里）
	for i := 0; i < length; i++ {
		for j := 0; j < length; j++ {
			mat[i][j] = mat[i][j] % p
		}
	}
	indictor := new(big.Int).SetInt64(1)//指示符变量，判断行变换的乘积正负号
	for i := 0; i < length; i++ {
		//行最小计算——不要0
		_,index := MinVectorByColInt64(mat[:][:], i, i)
		//行交换
		if index != i {
			mat[i], mat[index] = mat[index], mat[i]
			indictor.Mul(indictor, negative1)
		}
		//如果非0的最小值交换完成后仍然为0，则应该跳过（其实值直接为0）
		if mat[i][i] == 0 {
			continue
		}
		//遍历下面的每一行
		//i行与j行做差
		for j := i + 1; j < length; j++ {
			//x和y元素之间在模p前提下的倍数关系
			//nx - y = mp
			//计算比率
			n := tools.MatchXYInt64(mat[i][i],mat[j][i], p)
			//行间做差细化到每一个元素
			err,vector := VectorSubInt64(mat[j], mat[i], 1, n)
			if err != nil {
				return err,nil
			}
			vector = VectorModInt64(vector, p)
			mat[j] = vector
		}
	}
	r := new(big.Int).SetInt64(1)
	for i := 0; i < length; i++ {
		r.Mul(r, new(big.Int).SetInt64(mat[i][i]))
	}
	r.Mul(r, indictor)
	r.Mod(r, new(big.Int).SetInt64(p))
	return nil,r
}

func (m *Mat)DetInt64(p int64) (error,*big.Int) {
	err, r := DetInt64(m.matInt64, p)
	return err,r
}

//在模p的前提下计算逆矩阵
func InverseMatBigInt(m [][]big.Int, p *big.Int) (error, [][]big.Int) {
	//首先判断行列式是否为0
	err, det := DetBigInt(m, p)
	if err != nil {
		return err, nil
	}
	//判断是否存在逆矩阵
	if det.Cmp(zero) == 0 {
		return errors.New("Matrix should be full rank"), nil
	}
	//深度拷贝mat
	length := len(m)
	mat := make([][]big.Int, length)
	for i := 0; i < length; i++ {
		mat[i] = append(mat[i], make([]big.Int,length)...)
		copy(mat[i], m[i])
	}
	//预处理mat 先进行求模
	//这一步主要是避免在MatchXY时无法匹配（元素得在域里）
	for i := 0; i < length; i++ {
		for j := 0; j < length; j++ {
			mat[i][j].Mod(&mat[i][j], p)
		}
	}
	//计算逆矩阵
	//这一步如果和计算det合并可以加速，因为本质上为行变换
	//[mat I]~[I mat^-1]
	//<1>化简为上三角
	//初始化I
	_, I := OnesBigInt(length, length)
	//化简mat成为上三角
	for i := 0; i < length; i++ {
		//行最小计算——不要0
		_,index := MinVectorByCol(mat[:][:], i, i)
		//行交换
		if index != i {
			mat[i], mat[index] = mat[index], mat[i]
			I[i], I[index] = I[index], I[i]
		}
		//如果非0的最小值交换完成后仍然为0，则应该跳过（其实值直接为0）
		if mat[i][i].String() == "0" {
			continue
		}
		//遍历下面的每一行
		//i行与j行做差
		for j := i + 1; j < length; j++ {
			//x和y元素之间在模p前提下的倍数关系
			//nx - y = mp
			//计算比率
			n := tools.MatchXY(&mat[i][i],&mat[j][i], p)
			//行间做差细化到每一个元素
			err,vector1 := VectorSub(mat[j], mat[i], positive1, n)
			if err != nil {
				return err,nil
			}
			err,vector2 := VectorSub(I[j], I[i], positive1, n)
			if err != nil {
				return err,nil
			}
			vector1 = VectorMod(vector1, p)
			vector2 = VectorMod(vector2, p)
			mat[j] = vector1
			I[j] = vector2
		}
	}
	//<2>归一&&做差
	//从最后一行开始逐渐向上
	for i := length-1; i >= 0; i-- {
		//归一化第i行
		n := tools.MatchXY(&mat[i][i], positive1, p)
		vector1 := VectorMulEle(mat[i], n)
		vector1 = VectorMod(vector1, p)
		mat[i] = vector1
		vector2 := VectorMulEle(I[i], n)
		vector2 = VectorMod(vector2, p)
		I[i] = vector2
		//遍历下面的每一行
		//i行与j行做差
		for j := i - 1; j >= 0; j-- {
			//逐渐向上做差
			//行间做差细化到每一个元素
			//这里由于已经进行了行归一化，所以mat[i][i]为1
			err,vector1 := VectorSub(mat[j], mat[i], positive1, &mat[j][i])
			if err != nil {
				return err,nil
			}
			err,vector2 := VectorSub(I[j], I[i], positive1, &mat[j][i])
			if err != nil {
				return err,nil
			}
			vector1 = VectorMod(vector1, p)
			vector2 = VectorMod(vector2, p)
			mat[j] = vector1
			I[j] = vector2
		}
	}
	return nil, I
}

func (m *Mat) InverseMatBigInt(p *big.Int) (error, [][]big.Int) {
	return InverseMatBigInt(m.mat, p)
}

//不模p的情况下如何计算逆矩阵？
//关键在于匹配机制
//利用乘法来保证运行环境均在整数条件下
//计算方式 n*x = m*y直接令 n=y，m=x
func InverseMatBigIntWithoutP(mat [][]big.Int) (error, [][]big.Int) {
	//深度拷贝mat
	length := len(mat)
	A := make([][]big.Int, length)
	for i := 0; i < length; i++ {
		A[i] = append(A[i], make([]big.Int,length)...)
		copy(A[i], mat[i])
	}
	//计算逆矩阵
	//这一步如果和计算det合并可以加速，因为本质上为行变换
	//[mat I]~[I mat^-1]
	//<1>化简为上三角
	//初始化I
	I := make([][]big.Int, length)
	for i := 0; i < length; i++ {
		I[i] = append(I[i], make([]big.Int,1)...)
	}
	I[0][0].SetInt64(1)
	//scale := new(big.Int).SetInt64(1)
	//化简mat成为上三角
	ShowMat(mat)
	for i := 0; i < length; i++ {
		//行最小计算——不要0
		_,index := MinVectorByCol(A[:][:], i, i)
		//行交换
		if index != i {
			A[i], A[index] = A[index], A[i]
			I[i], I[index] = I[index], I[i]
		}
		//如果非0的最小值交换完成后仍然为0，则应该跳过（其实值直接为0）
		if A[i][i].String() == "0" {
			continue
		}
		//遍历下面的每一行
		//i行与j行做差
		for j := i + 1; j < length; j++ {
			//如果已经为0则跳过
			if A[j][i].String() == "0" {
				continue
			}
			//x和y元素之间在模p前提下的倍数关系
			//nx - y = mp
			//计算比率
			//行间做差细化到每一个元素
			r := tools.Gcd(&A[i][i], &A[j][i])
			n, m := new(big.Int).Div(&A[i][i], r), new(big.Int).Div(&A[j][i], r)
			fmt.Println(n.String(), " ", m.String())
			err,vector1 := VectorSub(A[j], A[i], n, m)
			if err != nil {
				return err,nil
			}
			err,vector2 := VectorSub(I[j], I[i], n, m)
			if err != nil {
				return err,nil
			}
			A[j] = vector1
			I[j] = vector2
		}
	}
	ShowMat(A)
	//<2>做差
	//从最后一行开始逐渐向上
	for i := length-1; i >= 0; i-- {
		//归一化
		//不过这里应当仅作为记录？

		//遍历下面的每一行
		//i行与j行做差
		for j := i - 1; j >= 0; j-- {
			//逐渐向上做差
			//行间做差细化到每一个元素
			//这里由于已经进行了行归一化，所以mat[i][i]为1
			err,vector1 := VectorSub(A[j], A[i], &A[i][i], &A[j][i])
			if err != nil {
				return err,nil
			}
			err,vector2 := VectorSub(I[j], I[i], &A[i][i], &A[j][i])
			if err != nil {
				return err,nil
			}
			A[j] = vector1
			I[j] = vector2
		}
	}
	ShowMat(A)
	//_, t := MatMulBigInt(mat, A)
	//ShowMat(t)

	//return nil, FormatMat2VectorByCol(I)
	return nil, nil
}

//Ax = b求解
func CalculateAxb(A_, b [][]big.Int) (error, [][]big.Int) {
	length := len(A_)
	//做行除法
	//A/B的格式可否
	A := make([][]BigFloat, length)
	for i := 0; i < length; i++ {
		A[i] = append(A[i], make([]BigFloat, length)...)
		for j := 0; j < length; j++ {
			A[i][j].A.SetInt64(A_[i][j].Int64())
			A[i][j].B.SetInt64(1)
		}
	}
	//全新算法模式
	//1 0 0 ==> 3 0 1 ==> 1/3 0 1/3 ==>
	//1 2 0		0 2 1	   0  2  1
	//1 0 3		0 0 1	   0  0  1

	//A[0][0].Div()


	return nil, nil
}