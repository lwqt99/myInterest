package tools

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

type googleDriver struct {
	service *selenium.Service
	webDriver *selenium.WebDriver
}

var GoogleDriver googleDriver

func StartChrome() (error) {

	opts := []selenium.ServiceOption{}
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}

	// 禁止加载图片，加快渲染速度
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}

	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  "",
		Args: []string{
			"--headless", // 设置Chrome无头模式
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/604.4.7 (KHTML, like Gecko) " +
				"Version/11.0.2 Safari/604.4.7", // 模拟user-agent，防反爬
		},
	}
	caps.AddChrome(chromeCaps)

	// 启动chromedriver，端口号可自定义
	service, err := selenium.NewChromeDriverService("C:/Goproject/src/gitgit/source/chromedriver.exe", 9516, opts...)
	if err != nil {
		return err
	}
	GoogleDriver.service = service

	// 调起chrome浏览器
	webDriver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 9516))
	if err != nil {
		return err
	}
	GoogleDriver.webDriver = &webDriver
	return nil
}

func SkillChrome()  {
	defer GoogleDriver.service.Stop()  // 停止chromedriver
	defer (*(GoogleDriver.webDriver)).Quit() // 关闭浏览器
}

func GetPageText(targeUrl string) (error, string) {
	err := (*(GoogleDriver.webDriver)).Get(targeUrl)
	if err != nil {
		return err,""
	}

	getpagesource,err := (*(GoogleDriver.webDriver)).PageSource()
	if err != nil{
		return err,""
	}

	time.Sleep(5 * time.Second)
	return nil,getpagesource
}


func dealPageSource(pagesource string) ([]string) {
	r := "<a href=\"//.+\" target=\"_blank\" class=\"cover\">"
	reg := regexp.MustCompile(r)
	if reg == nil { //解释失败，返回nil
		fmt.Println("正则表达式创建失败")
		return nil
	}

	result := reg.FindAllStringSubmatch(pagesource, -1)
	t := make([]string, len(result))
	for i := 0; i < len(result); i++ {
		t[i] = result[i][0][11:46]
	}
	return t
}

func DownloadAllUrlById(id string) error {
	baseUrl := "https://space.bilibili.com/"+id+"/video?tid=0&page="
	var urls []string
	err,pagesource := GetPageText(baseUrl+"1")
	if err != nil {
		return err
	}
	urls = append(urls, dealPageSource(pagesource)...)

	r1 := "<span class=\"be-pager-total\">共 . 页，</span>"
	reg1 := regexp.MustCompile(r1)
	if reg1 == nil { //解释失败，返回nil
		return errors.New("正则表达式创建失败")
	}

	result := reg1.FindStringSubmatch(pagesource)

	totalPages := result[0][33:34+len(result[0])-len(r1)]
	totalPagesInt64,err := strconv.ParseInt(totalPages,10,64)

	for i := 2; i < int(totalPagesInt64) + 1; i++ {
		fmt.Println("第"+strconv.Itoa(i)+"页")
		url := baseUrl + strconv.Itoa(i)
		SkillChrome()
		StartChrome()
		err,pagesource := GetPageText(url)
		if err != nil{
			fmt.Println(err.Error())
		}
		urls = append(urls, dealPageSource(pagesource)...)
	}

	//列表去重
	tempMap := make(map[string]bool)
	for i := 0; i < len(urls); i++ {
		_,ok := tempMap[urls[i]]
		if !ok{
			tempMap[urls[i]] = true
		}
	}

	var finalUrls []string
	for k,_ := range tempMap{
		finalUrls = append(finalUrls, k)
	}

	c := ""
	for i := 0; i < len(finalUrls); i++ {
		c += finalUrls[i] + " "
	}

	fmt.Println(c)

	fmt.Println(len(finalUrls))
	DownVideo(c)

	//for len(finalUrls) != 0 {
	//	fmt.Println("剩余待处理URL:"+strconv.Itoa(len(finalUrls)))
	//	for i := 0; i < len(finalUrls); i++ {
	//		fmt.Println(i)
	//		err = DownVideo(finalUrls[i])
	//		if err != nil{
	//			fmt.Println(err.Error())
	//			continue
	//		}else {
	//			finalUrls = append(finalUrls[:i], finalUrls[i+1:]...)
	//		}
	//	}
	//}
	return nil
}

func DownVideo(url string) error {
	exec.Command("E:").Run()
	exec.Command("cd","E:\\Videos").Run()
	cmd1 := exec.Command("you-get", url)
	//cmd2 := exec.Command("ffmpeg", "-i", "E:/Videos", url)
	err :=cmd1.Run()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	//var out bytes.Buffer
	//cmd1.Stdout = &out
	//fmt.Println(out.String())


	//err =cmd2.Run()
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return err
	//}

	return nil
}