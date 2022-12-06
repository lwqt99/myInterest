package tools

import (
	"crypto/rand"
	"encoding/pem"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
	"io/ioutil"
	"os"
)


func GenerateSM2KeyPairs() error {
	priv, err := sm2.GenerateKey(rand.Reader) // 生成密钥对
	if err != nil {
		return err
	}
	// 序列化私匙
	priBytes, err := x509.MarshalSm2UnecryptedPrivateKey(priv)
	if err != nil {
		return err
	}
	priBlock := pem.Block{
		Type:  "ECD PRIVATE KEY",
		Bytes: priBytes,
	}
	priFile, _ := os.Create("./gitgit/certificate/sm2-prikey.pem")
	// 编码私匙,写入文件
	if err := pem.Encode(priFile, &priBlock); err != nil {
		return err
	}
	return nil
}

// 加载私匙公匙
func LoadKey() (error, *sm2.PrivateKey) {
	// 读取密匙
	pri, _ := ioutil.ReadFile("./gitgit/certificate/sm2-prikey.pem")
	//pub, _ := ioutil.ReadFile("./gitgit/certificate/sm2-pubkey.pem")
	// 解码私匙
	block, _ := pem.Decode(pri)
	var err error
	// 反序列化私匙

	privKey, err := x509.ParsePKCS8UnecryptedPrivateKey(block.Bytes)
	if err != nil {
		return err,nil
	}


	//// 解码公匙
	//block, _ = pem.Decode(pub)
	//// 反序列化公匙
	//var t interface{}
	//t, err = x509.ParsePKIXPublicKey(block.Bytes)
	//if err != nil {
	//	return err, nil
	//}
	//var ok bool
	//privKey.PublicKey, ok = t.(sm2.PublicKey)
	//if !ok {
	//	return errors.New("公钥interface转换失败"), nil
	//}
	return nil, privKey
}

/*
	使用对方的公钥加密
	content是需要加密的内容
*/
func CreateSm2Encrypt(priv *sm2.PrivateKey, msg []byte) ([]byte,*sm2.PublicKey,error) {
	pub := &priv.PublicKey
	ciphertxt, err := pub.EncryptAsn1(msg,rand.Reader) //sm2加密
	if err != nil {
		return nil, nil, err
	}

	return ciphertxt,pub,nil
}

/*
	私钥进行解密操作
	需要使用匹配的私钥进行解密
*/
func Sm2Decrypt(priv *sm2.PrivateKey,ciphertxt []byte) (string,error) {
	//读取密钥对
	plaintxt,err :=  priv.DecryptAsn1(ciphertxt)  //sm2解密
	if err != nil {
		return "", err
	}
	return string(plaintxt),nil
}

