package tools

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type Spider struct {
	response *http.Response
}

var SpiderMan Spider

func genIpaddr() string {
	rand.Seed(time.Now().Unix())
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}

func GetOutBoundIP()(ip string, err error)  {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}


func (s Spider) parse(response *http.Response) (string, error) {
	var htmlDocument string
	buf := make([]byte, 4096)
	for {
		n, err := response.Body.Read(buf)
		if n == 0 {
			fmt.Println("网页读取成功")
			return htmlDocument, nil
		}
		if err != nil && err != io.EOF {
			return "", err
		}
		htmlDocument += string(buf[:n])
	}
}

// 将字符串存储为HTML格式的文件
func (s Spider) saveAsHtml(htmlDocument string) error {
	file, err := os.Create("csdn.html")
	if err != nil {
		return err
	}
	_, err = file.WriteString(htmlDocument)
	if err != nil {
		return err
	}
	//fmt.Println(n)
	return nil
}


func HttpGet(url,mode string) (error, string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return err, ""
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New("status code error"),""
	}
	htmlDocument,err := SpiderMan.parse(res)
	if err != nil {
		return err,""
	}
	if mode == "save"{
		err = SpiderMan.saveAsHtml(htmlDocument)
		if err != nil {
			return err,""
		}
		return nil,""
	}else if mode == "get" {
		return  nil, htmlDocument
	}

	return nil,""
}

//可以使用代理来伪造ip
func (s *Spider) NewHttpGet(url string) error {
	//proxyIP := "//" + genIpaddr() + ":8080"
	//proxyIP := "175.150.110.188:8080"
	//proxy := func(_ *http.Request) (*u.URL, error) {
	//	return u.Parse(proxyIP)
	//}
	//transport := &http.Transport{Proxy: proxy}
	//client := &http.Client{Transport: transport}

	//声明client 参数为默认
	client := &http.Client{}

	//method:GET/POST等 url:顾名思义 io.Reader:
	var r io.Reader

	request, err := http.NewRequest("GET", url, r)
	//request.Header.Add("X-Forwarded-For", genIpaddr() + ":8080")
	//request.Header.Set()

	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = client.Do(request)
	if err != nil {
		return err
	}
	//response, _ := client.Do(request)
	//
	//if response.StatusCode != 200 {
	//	return errors.New("status code error")
	//}

	//htmlDocument,err := SpiderMan.parse(response)
	//if err != nil {
	//	return err
	//}
	//err = SpiderMan.saveAsHtml(htmlDocument)
	//if err != nil {
	//	return err
	//}
	return nil
}

//udp访问，且伪造ip
func (s Spider) UdpAccessByRandomIp(url string) error {
	ServerAddr,err := net.ResolveUDPAddr("udp4",":8080")
	if err != nil {
		return err
	}
	LocalAddr, err := net.ResolveUDPAddr("udp4", genIpaddr()+":8080")
	if err != nil {
		return err
	}

	RemoteAddr, err := net.ResolveUDPAddr("udp4", "192.168.164.255:10002")
	if err != nil {
		return err
	}
	Conn, err := net.DialUDP("udp4", LocalAddr, RemoteAddr)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	defer Conn.Close()

	ServerConn, err := net.ListenUDP("udp4", ServerAddr)


	defer ServerConn.Close()



	return nil
}
