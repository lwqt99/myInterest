package tools

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/pem"
	"errors"
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
	ECDSA.randKey = "ljz abc 123456 random choose some helpful id well interest how dare you are"
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



// 基于椭圆曲线的DH密码体质
func Diffie_Hellman(){

}
