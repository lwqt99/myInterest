package algebra

import (
	"gitgit/tools"
	"github.com/pkg/errors"
	"math/big"
)

//多项式求逆
//在模p下的逆计算
func InverseVectorBigInt(vector []big.Int, p *big.Int) (error, []big.Int) {
	//首先验证是否存在逆
	//存在的条件为vector构建的方程组有解
	//具体格式如下，a_i为Vector中的对应元素
	//A = [ a_0			a_(N-1)		...		a3		a2		a1	]
	//		.			.			...		.		.		.	]
	//		a_(N-2)		a_(N-3)		...		a1		a0		a_(N-1) ]
	//		a_(N-1)		a_(N-2)		...		a2		a1		a0		]
	//具体表现为A的行列式不为0
	//<1>构建A
	length := len(vector)
	//声明
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
	//<2>判断A行列式结果，验证是否有解
	err,det := DetBigInt(A, p)
	if err != nil {
		return err, nil
	}
	if det.String() == "0" {
		return errors.New("This vector has not invsere"), nil
	}
	//<3>计算逆 利用増广矩阵行变换解决
	//初始化I = [1 0 0 ... 0 0 0]
	//_, I := OnesBigInt(length, length)
	I := make([][]big.Int, length)
	for i := 0; i < length; i++ {
		I[i] = append(I[i], make([]big.Int,1)...)
	}
	I[0][0].SetInt64(1)
	//化简mat成为上三角
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
			//x和y元素之间在模p前提下的倍数关系
			//nx - y = mp
			//计算比率
			n := tools.MatchXY(&A[i][i],&A[j][i], p)
			//行间做差细化到每一个元素
			err,vector1 := VectorSub(A[j], A[i], positive1, n)
			if err != nil {
				return err,nil
			}
			err,vector2 := VectorSub(I[j], I[i], positive1, n)
			if err != nil {
				return err,nil
			}
			vector1 = VectorMod(vector1, p)
			vector2 = VectorMod(vector2, p)
			A[j] = vector1
			I[j] = vector2
		}
	}
	//<2>归一&&做差
	//从最后一行开始逐渐向上
	for i := length-1; i >= 0; i-- {
		//归一化第i行
		n := tools.MatchXY(&A[i][i], positive1, p)
		vector1 := VectorMulEle(A[i], n)
		vector1 = VectorMod(vector1, p)
		A[i] = vector1
		vector2 := VectorMulEle(I[i], n)
		vector2 = VectorMod(vector2, p)
		I[i] = vector2

		//遍历下面的每一行
		//i行与j行做差
		for j := i - 1; j >= 0; j-- {
			//逐渐向上做差
			//行间做差细化到每一个元素
			//这里由于已经进行了行归一化，所以mat[i][i]为1
			err,vector1 := VectorSub(A[j], A[i], positive1, &A[j][i])
			if err != nil {
				return err,nil
			}
			err,vector2 := VectorSub(I[j], I[i], positive1, &A[j][i])
			if err != nil {
				return err,nil
			}
			vector1 = VectorMod(vector1, p)
			vector2 = VectorMod(vector2, p)
			A[j] = vector1
			I[j] = vector2
		}
	}
	return FormatMat2VectorByCol(I)
}
