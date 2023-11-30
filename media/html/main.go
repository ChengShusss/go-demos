package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/gohouse/golib/random"
)

var (
	url = "https://www.zhihu.com/question/554409392/answer/3304018835"
	// url = "https://www.zhihu.com/question/632248312/answer/3305395429"

	// "https://www.zhihu.com/question/632248312/answer/3305575082"
)

func main() {
	// 请求 HTML 页面
	// res, err := http.Get(url)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer res.Body.Close()

	// 从本地获取 html 作为模拟输入
	if len(os.Args) < 2 {
		fmt.Println("please specific input file")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("failed to open file, err: ", err)
		os.Exit(1)
	}
	defer f.Close()

	// 解析 HTML 文档
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}

	ReplaceImgPaths(doc)

	s, err := Html2Md(doc)
	if err != nil {
		fmt.Printf("Failed to trans, err: %v\n", err)
		return
	}

	fmt.Println(s)

}

func Html2Md(doc *goquery.Document) (string, error) {
	h, err := doc.Find("div.css-376mun").Html()

	if err != nil {
		return "", err
	}

	converter := md.NewConverter("", true, nil)
	return converter.ConvertString(h)
}

func ExtractImgPaths(doc *goquery.Document) (imgUrls []string) {

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("src")
		if ok {
			imgUrls = append(imgUrls, url)
		}
	})

	return
}

func ReplaceImgPaths(doc *goquery.Document) {

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("src")
		if !ok {
			return
		}

		new, err := DownloadImgs(url, "./data/")
		if err != nil {
			fmt.Printf("download [%v] failed, err: %v\n", url, err)
			return
		}
		s.SetAttr("src", new)
	})
}

func DownloadImgs(url, localDir string) (string, error) {
	// better use a new name

	// Do http request
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	contentType := res.Header.Get("Content-Type")
	exts, err := mime.ExtensionsByType(contentType)
	if err != nil || len(exts) == 0 {
		exts = []string{"jpg"}
	} else {
		fmt.Printf("exts: %v\n", exts)
	}
	newName := fmt.Sprintf("%d-%s%s", time.Now().Unix(), random.RandString(10), exts[len(exts)-1])

	fmt.Println("new name: ", newName)

	f, err := os.OpenFile(filepath.Join(localDir, newName), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return "", err
	}

	return newName, nil
}

func PrintImgs(doc *goquery.Document) {

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		a, ok := s.Attr("src")
		fmt.Printf("s.attr: %v, %v\n", a, ok)
	})

}
