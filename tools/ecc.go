package tools

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"github.com/tjfoc/gmsm/sm3"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
)

type Ecdsa struct {
	randKey string
	priKey *ecdsa.PrivateKey
	pubKey *ecdsa.PublicKey
}

type Cipher struct {
	XCoordinate *big.Int
	YCoordinate *big.Int
	HASH        []byte
	CipherText  []byte
}

var ECDSA Ecdsa

//用于生成自定义的参数
func initPX() elliptic.Curve {
	var p *elliptic.CurveParams
	p = &elliptic.CurveParams{Name: "P-X"}
	p.P, _ = new(big.Int).SetString("39402006196394479212279040100143613805079739270465446667948293404245721771496870329047266088258938001861606973112319", 10)
	p.N, _ = new(big.Int).SetString("39402006196394479212279040100143613805079739270465446667946905279627659399113263569398956308152294913554433653942643", 10)
	p.B, _ = new(big.Int).SetString("b3312fa7e23ee7e4988e056be3f82d19181d9c6efe8141120314088f5013875ac656398d8a2ed19d2a85c8edd3ec2aef", 16)
	p.Gx, _ = new(big.Int).SetString("aa87ca22be8b05378eb1c71ef320ad746e1d3b628ba79b9859f741e082542a385502f25dbf55296c3a545e3872760ab7", 16)
	p.Gy, _ = new(big.Int).SetString("3617de4a96262c6f5d9e98bf9292dc29f8f41dbd289a147ce9da3113b5f0b8c00a60b1ce1d7e819d7a431d7c90ea0e5f", 16)
	p.BitSize = 384
	return p
}

//也可以使用自定义的曲线 具体可修改上面的函数
func (*Ecdsa)generateKey(priFile, pubFile *os.File) error {
	lenth := len(ECDSA.randKey)
	if lenth < 224/8 + 8 {
		return errors.New("randKey长度太短，至少为36位！")
	}
	// 根据随机密匙的长度创建私匙
	var curve elliptic.Curve
	if lenth > 521/8+8 {
		curve = elliptic.P521()
	} else if lenth > 384/8+8 {
		curve = elliptic.P384()
	} else if lenth > 256/8+8 {
		curve = elliptic.P256()
	} else if lenth > 224/8+8 {
		curve = elliptic.P224()
	}
	// 生成私匙
	priKey, err := ecdsa.GenerateKey(curve, strings.NewReader(ECDSA.randKey))
	if err != nil {
		return err
	}
	// 序列化私匙
	priBytes, err := x509.MarshalECPrivateKey(priKey)
	if err != nil {
		return err
	}
	priBlock := pem.Block{
		Type:  "ECD PRIVATE KEY",
		Bytes: priBytes,
	}
	// 编码私匙,写入文件
	if err := pem.Encode(priFile, &priBlock); err != nil {
		return err
	}
	// 序列化公匙
	pubBytes, err := x509.MarshalPKIXPublicKey(&priKey.PublicKey)
	if err != nil {
		return err
	}
	pubBlock := pem.Block{
		Type:  "ECD PUBLIC KEY",
		Bytes: pubBytes,
	}
	// 编码公匙,写入文件
	if err := pem.Encode(pubFile, &pubBlock); err != nil {
		return err
	}
	return nil
}

// y^2  = x^3 + ax + b
func (*Ecdsa)initECDSA() error {
	ECDSA.randKey = "ljz abc 123456 random choose some helpful id"
	// 初始化生成私匙公匙
	priFile, _ := os.Create("./gitgit/certificate/ecdsa-prikey.pem")
	pubFile, _ := os.Create("./gitgit/certificate/ecdsa-pubkey.pem")
	if err := ECDSA.generateKey(priFile, pubFile); err != nil {
		return err
	}
	return nil
}

func (*Ecdsa)GenerateKeys() error {
	if err := ECDSA.initECDSA(); err != nil{
		return err
	}
	return nil
}


// 加载私匙公匙
func (*Ecdsa)LoadKey() error {
	// 读取密匙
	pri, _ := ioutil.ReadFile("./gitgit/certificate/ecdsa-prikey.pem")
	pub, _ := ioutil.ReadFile("./gitgit/certificate/ecdsa-pubkey.pem")
	// 解码私匙
	block, _ := pem.Decode(pri)
	var err error
	// 反序列化私匙
	ECDSA.priKey, err = x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	// 解码公匙
	block, _ = pem.Decode(pub)
	// 反序列化公匙
	var t interface{}
	t, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	var ok bool
	ECDSA.pubKey, ok = t.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("公钥interface转换失败")
	}
	return nil
}

//生成签名
func (*Ecdsa) Sign(message string) (*big.Int, *big.Int, error) {
	hashMessage := sha256.Sum256([]byte(message))
	r,s,err := ecdsa.Sign(rand.Reader, ECDSA.priKey, hashMessage[:])
	if err!=nil {
		return nil,nil, err
	}
	return r, s, nil
}

//签名校验
func (*Ecdsa)Verify(message string,r, s *big.Int) bool {
	hashMessage := sha256.Sum256([]byte(message))
	right := ecdsa.Verify(ECDSA.pubKey, hashMessage[:], r, s)
	return right
}

/*
	用于生成随机数
 */
func randFieldElement(c elliptic.Curve, random io.Reader) (r *big.Int, err error) {
	one := new(big.Int).SetInt64(1)
	if random == nil {
		random = rand.Reader //If there is no external trusted random source,please use rand.Reader to instead of it.
	}
	params := c.Params()
	b := make([]byte, params.BitSize/8+8)
	_, err = io.ReadFull(random, b)
	if err != nil {
		return
	}
	r = new(big.Int).SetBytes(b)
	n := new(big.Int).Sub(params.N, one)
	r.Mod(r, n)
	r.Add(r, one)
	return
}

func intToBytes(x int) []byte {
	var buf = make([]byte, 4)

	binary.BigEndian.PutUint32(buf, uint32(x))
	return buf
}

// 填充
func zeroByteSlice(n int) []byte {
	return make([]byte,n)
}

/*
	加密
 */
func kdf(length, bitsize int, x ...[]byte) ([]byte, bool) {
	var c []byte

	ct := 1
	h := sm3.New()
	for i, j := 0, (length+bitsize-1)/bitsize; i < j; i++ {
		h.Reset()
		for _, xx := range x {
			h.Write(xx)
		}
		h.Write(intToBytes(ct))
		hash := h.Sum(nil)
		if i+1 == j && length%bitsize != 0 {
			c = append(c, hash[:length%bitsize]...)
		} else {
			c = append(c, hash...)
		}
		ct++
	}
	for i := 0; i < length; i++ {
		if c[i] != 0 {
			return c, true
		}
	}
	return c, false
}

func encrypt(message string, random io.Reader) ([]byte, error) {
	data := []byte(message)
	length := len(data)

	for {
		curve := ECDSA.pubKey.Curve
		r, err := randFieldElement(curve, random)
		if err != nil {
			return nil, err
		}
		x1, y1 := curve.ScalarBaseMult(r.Bytes()) //计算rG(kG便于理解)
		x2, y2 := curve.ScalarMult(ECDSA.pubKey.X, ECDSA.pubKey.Y, r.Bytes())

		x1Buf := x1.Bytes()
		y1Buf := y1.Bytes()
		x2Buf := x2.Bytes()
		y2Buf := y2.Bytes()

		//fmt.Println(len(x1Buf))
		bitsize := curve.Params().BitSize
		bitsize_8 := bitsize / 8

		if n := len(x1Buf); n < bitsize_8 {
			x1Buf = append(zeroByteSlice(bitsize_8)[:bitsize_8-n], x1Buf...)
		}
		if n := len(y1Buf); n < bitsize_8 {
			y1Buf = append(zeroByteSlice(bitsize_8)[:bitsize_8-n], y1Buf...)
		}
		if n := len(x2Buf); n < bitsize_8 {
			x2Buf = append(zeroByteSlice(bitsize_8)[:bitsize_8-n], x2Buf...)
		}
		if n := len(y2Buf); n < bitsize_8 {
			y2Buf = append(zeroByteSlice(bitsize_8)[:bitsize_8-n], y2Buf...)
		}

		// 计算加密结果
		c := []byte{}
		c = append(c, x1Buf...) // x分量
		c = append(c, y1Buf...) // y分量
		tm := []byte{}
		tm = append(tm, x2Buf...)
		tm = append(tm, data...)
		tm = append(tm, y2Buf...)

		h := sm3.Sm3Sum(tm)
		c = append(c, h...)
		ct, ok := kdf(length, bitsize_8, x2Buf, y2Buf) // 密文
		if !ok {
			continue
		}
		c = append(c, ct...)
		//for i := 0; i < len(c); i++ {
		//	fmt.Println(c[i])
		//}
		//fmt.Println(len(c))
		for i := 0; i < length; i++ {
			c[bitsize_8*3+i] ^= data[i]
		}

		return append([]byte{0x04}, c...), nil
	}

}

func CipherMarshal(bitsize int, data []byte) ([]byte, error) {
	data = data[1:]
	x := new(big.Int).SetBytes(data[:bitsize])
	y := new(big.Int).SetBytes(data[bitsize:bitsize*2])
	hash := data[bitsize*2:bitsize*3]
	cipherText := data[bitsize*3:]
	return asn1.Marshal(Cipher{x, y, hash, cipherText})
}

func (*Ecdsa) EncryptAsn1(message string, random io.Reader) ([]byte, error) {
	cipher, err := encrypt(message, random)
	if err != nil {
		return nil, err
	}
	return CipherMarshal(ECDSA.pubKey.Curve.Params().BitSize/8,cipher)
}

/*
sm2密文asn.1编码格式转C1|C3|C2拼接格式
*/
func CipherUnmarshal(bitsize int, data []byte) ([]byte, error) {
	var cipher Cipher
	_, err := asn1.Unmarshal(data, &cipher)
	if err != nil {
		return nil, err
	}
	x := cipher.XCoordinate.Bytes()
	y := cipher.YCoordinate.Bytes()
	hash := cipher.HASH
	if err != nil {
		return nil, err
	}
	cipherText := cipher.CipherText
	if err != nil {
		return nil, err
	}
	if n := len(x); n < bitsize {
		x = append(zeroByteSlice(bitsize)[:bitsize-n], x...)
	}
	if n := len(y); n < bitsize {
		y = append(zeroByteSlice(bitsize)[:bitsize-n], y...)
	}
	c := []byte{}
	c = append(c, x...)          // x分量
	c = append(c, y...)          // y分
	c = append(c, hash...)       // x分量
	c = append(c, cipherText...) // y分
	return append([]byte{0x04}, c...), nil
}


func decrypt(priv *ecdsa.PrivateKey, data []byte) ([]byte, error) {
	data = data[1:]
	curve := priv.Curve
	bitsize_8 := curve.Params().BitSize / 8
	length := len(data) - bitsize_8 * 3

	x := new(big.Int).SetBytes(data[:bitsize_8])
	y := new(big.Int).SetBytes(data[bitsize_8:bitsize_8*2])
	x2, y2 := curve.ScalarMult(x, y, priv.D.Bytes())
	x2Buf := x2.Bytes()
	y2Buf := y2.Bytes()
	if n := len(x2Buf); n < bitsize_8 {
		x2Buf = append(zeroByteSlice(bitsize_8)[:bitsize_8-n], x2Buf...)
	}
	if n := len(y2Buf); n < bitsize_8 {
		y2Buf = append(zeroByteSlice(bitsize_8)[:bitsize_8-n], y2Buf...)
	}
	c, ok := kdf(length,bitsize_8, x2Buf, y2Buf)
	if !ok {
		return nil, errors.New("Decrypt: failed to decrypt")
	}
	for i := 0; i < length; i++ {
		c[i] ^= data[i+bitsize_8*3]
	}
	tm := []byte{}
	tm = append(tm, x2Buf...)
	tm = append(tm, c...)
	tm = append(tm, y2Buf...)
	h := sm3.Sm3Sum(tm)
	if bytes.Compare(h, data[bitsize_8*2:bitsize_8*3]) != 0 {
		return c, errors.New("Decrypt: failed to decrypt")
	}
	return c, nil
}


func (*Ecdsa) DecryptAsn1(ciphertxt []byte) ([]byte, error) {
	bitsize_8 := ECDSA.pubKey.Curve.Params().BitSize / 8
	cipher, err := CipherUnmarshal(bitsize_8,ciphertxt)
	if err != nil {
		return nil, err
	}
	return decrypt(ECDSA.priKey, cipher)
}

// 基于椭圆曲线的DH密码体质
func Diffie_Hellman(){

}
