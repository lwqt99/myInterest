package algebra

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"time"
)

//输出长度筛选
func formatFilteLength(str string, length int) {
	switch length {
	case 1:
		fmt.Printf("%1v ", str)
	case 2:
		fmt.Printf("%2v ", str)
	case 3:
		fmt.Printf("%3v ", str)
	case 4:
		fmt.Printf("%4v ", str)
	case 5:
		fmt.Printf("%5v ", str)
	case 6:
		fmt.Printf("%6v ", str)
	case 7:
		fmt.Printf("%7v ", str)
	case 8:
		fmt.Printf("%8v ", str)
	case 9:
		fmt.Printf("%9v ", str)
	case 10:
		fmt.Printf("%10v ", str)
	case 11:
		fmt.Printf("%11v ", str)
	case 12:
		fmt.Printf("%12v ", str)
	case 13:
		fmt.Printf("%13v ", str)
	case 14:
		fmt.Printf("%14v ", str)
	case 15:
		fmt.Printf("%15v ", str)
	default:
		//fmt.Printf("%100v ", str)
		fmt.Print(str+" ")
	}
}

//格式化输出矩阵
func ShowMat(mat [][]big.Int) {
	maxStr := 0
	//获取最长的元素 长度
	for i := 0; i < len(mat); i++ {
		for j := 0; j < len(mat[0]); j++ {
			if len(mat[i][j].String()) > maxStr {
				maxStr = len(mat[i][j].String())
			}
		}
	}
	fmt.Println("********************以下为矩阵展示********************")
	for i := 0; i < len(mat); i++ {
		for j := 0; j < len(mat[0]); j++ {
			formatFilteLength(mat[i][j].String(), maxStr)
		}
		fmt.Println()
	}
	fmt.Println("********************结束********************")
}

//格式化输出矩阵
func (m *Mat)ShowMat() {
	ShowMat(m.mat)
}

//格式化输出矩阵
func ShowMatInt64(mat [][]int64) {
	maxStr := 0
	//获取最长的元素 长度
	for i := 0; i < len(mat); i++ {
		for j := 0; j < len(mat[0]); j++ {
			if len(strconv.FormatInt(mat[i][j], 10)) > maxStr {
				maxStr = len(strconv.FormatInt(mat[i][j], 10))
			}
		}
	}
	for i := 0; i < len(mat); i++ {
		for j := 0; j < len(mat[0]); j++ {
			formatFilteLength(strconv.FormatInt(mat[i][j], 10), maxStr)
		}
		fmt.Println()
	}
}

//格式化输出矩阵
func (m *Mat)ShowMatInt64() {
	//fmt.Println(m.matInt64)
	ShowMatInt64(m.matInt64)
}

//随机数矩阵
func RandomMatBigInt(range_ int64, width, height int) [][]big.Int {
	//设置随机数种子
	rand.Seed(time.Now().UnixNano())

	var mat [][]big.Int
	for i := 0; i < height; i++ {
		mat = append(mat, make([]big.Int, width))
		for j := 0; j < width; j++ {
			mat[i][j].SetInt64(rand.Int63n(range_))
		}
	}

	return mat
}

//随机数矩阵
func RandomMatInt64(range_ int64, width, height int) [][]int64 {
	//设置随机数种子
	rand.Seed(time.Now().UnixNano())
	var mat [][]int64
	for i := 0; i < height; i++ {
		mat = append(mat, make([]int64, width))
		for j := 0; j < width; j++ {
			mat[i][j] = rand.Int63n(range_)
		}
	}
	return mat
}

//单位阵
//返回的是Mat结构
func OnesBigIntMat(width, height int) (error, *Mat) {
	if height <= 0 || width <= 0 {
		return errors.New("width or height should be positive"), nil
	}
	m := new(Mat)
	m.width = width
	m.height = height

	min := width
	if width > height {
		min = height
	}
	for i := 0; i < min; i++ {
		row := make([]big.Int, width)
		row[i].SetInt64(1)
		m.mat[i] = append(m.mat[i], row...)
	}
	return nil, m
}

//单位阵
//返回的是单纯的mat
func OnesBigInt(width, height int) (error, [][]big.Int) {
	if height <= 0 || width <= 0 {
		return errors.New("width or height should be positive"), nil
	}
	m := make([][]big.Int, height)
	min := width
	if width > height {
		min = height
	}
	for i := 0; i < min; i++ {
		row := make([]big.Int, width)
		row[i].SetInt64(1)
		m[i] = append(m[i], row...)
	}
	return nil, m
}

//索引数组
func MatIndex(m [][]big.Int,x1,y1,x2,y2 int) [][]big.Int {
	var mat [][]big.Int
	for i := 0; i < y2 - y1 + 1; i++ {
		mat = append(mat, make([]big.Int, x2 - x1 + 1))
		//fmt.Println(m[y1 + i][x1:x2+1])
		copy(mat[i], m[y1 + i][x1:x2+1])
		//mat[i] = m[y1 + i][x1:x2+1]
	}
	return mat
}

//间隔索引数组
func MatIndexInternal(m [][]big.Int,x1, x2, x3, x4, y1, y2 int) [][]big.Int {
	mat := MatIndex(m,x1,y1,x2,y2)
	for i := 0; i < y2 - y1 + 1; i++ {
		mat[i] = append(mat[i], m[y1 + i][x3:x4+1]...)
	}
	return mat
}

//索引数组
func (m *Mat) MatIndex(x1,y1,x2,y2 int) [][]big.Int {
	var mat [][]big.Int
	for i := 0; i < y2 - y1 + 1; i++ {
		mat = append(mat, make([]big.Int, x2 - x1 + 1))
		mat[y1 + i] = m.mat[y1 + i][x1:x2+1]
	}
	return mat
}

//间隔索引数组
func (m *Mat) MatIndexInternal(x1, x2, x3, x4, y1, y2 int) [][]big.Int {
	mat := m.MatIndex(x1,y1,x2,y2)
	for i := 0; i < y2 - y1 + 1; i++ {
		mat[i] = append(mat[i], m.mat[y1 + i][x3:x4+1]...)
	}
	//for _,v := range mat{
	//	fmt.Println(v)
	//}
	return mat
}

//矩阵乘法
func MatMulBigInt(x, y [][]big.Int) (error,[][]big.Int) {
	if len(x[0]) != len(y) {
		return errors.New("x and y should match their dim"), nil
	}
	height, width := len(x), len(y[0])
	r := make([][]big.Int, height)
	for i := 0; i < height; i++ {
		r[i] = append(r[i], make([]big.Int, width)...)
	}
	//矩阵乘法
	//列索引优化——待完成
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			//计算r[i][j]的值
			for k := 0; k < len(y); k++ {
				r[i][j].Add(&r[i][j], new(big.Int).Mul(&x[i][k], &y[k][j]))
			}
		}
	}
	return nil, r
}

//矩阵乘法
func MatMulInt64(x, y [][]int64) (error,[][]int64) {
	if len(x[0]) != len(y) {
		return errors.New("x and y should match their dim"), nil
	}
	height, width := len(x), len(y[0])
	r := make([][]int64, height)
	for i := 0; i < height; i++ {
		r[i] = append(r[i], make([]int64, width)...)
	}
	//矩阵乘法
	//列索引优化——待完成
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			//计算r[i][j]的值
			for k := 0; k < len(y); k++ {
				r[i][j] += x[i][k]*y[k][j]
			}
		}
	}
	return nil, r
}

func MatMod(x [][]big.Int, p *big.Int) [][]big.Int {
	height,width := len(x), len(x[0])
	r := make([][]big.Int, height)
	for i := 0; i < height; i++ {
		r[i] = append(r[i], make([]big.Int, height)...)
	}
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			r[i][j].Mod(&x[i][j], p)
		}
	}
	return r
}

//转置
func TransposeInt64(mat [][]int64) [][]int64 {
	height,width := len(mat), len(mat[0])
	r := make([][]int64, width)
	for i := 0; i < width; i++ {
		r[i] = append(r[i], make([]int64, height)...)
		for j := 0; j < height; j++ {
			r[i][j] = mat[j][i]
		}
	}
	return r
}

func TransposeBigInt(mat [][]big.Int) [][]big.Int {
	height,width := len(mat), len(mat[0])
	r := make([][]big.Int, width)
	for i := 0; i < width; i++ {
		r[i] = append(r[i], make([]big.Int, height)...)
		for j := 0; j < height; j++ {
			r[i][j].Set(&mat[j][i])
		}
	}
	return r
}