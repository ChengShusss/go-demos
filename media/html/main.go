package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/gohouse/golib/random"
)

var (
// url = "https://www.zhihu.com/question/554409392/answer/3304018835"
// url = "https://www.zhihu.com/question/629815487/answer/3303525305"
// url = "https://www.zhihu.com/question/632248312/answer/3305395429"

// "https://www.zhihu.com/question/632248312/answer/3305575082"
)

func main() {

	// 从本地获取 html 作为模拟输入
	if len(os.Args) < 3 {
		fmt.Println("please specific input file")
		os.Exit(1)
	}

	input := os.Args[1]
	output := os.Args[2]

	var r io.Reader

	if strings.HasPrefix(input, "http") {
		// 请求 HTML 页面
		client := &http.Client{}
		req, err := http.NewRequest("GET", input, nil)
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Set("authority", "www.zhihu.com")
		req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
		req.Header.Set("cache-control", "no-cache")
		// Need to set cookies
		req.Header.Set("sec-ch-ua", `Google Chrome";v="119", "Chromium";v="119", "Not?A_Brand";v="24"'`)
		req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		r = resp.Body
	} else {
		f, err := os.Open(input)
		if err != nil {
			fmt.Println("failed to open file, err: ", err)
			os.Exit(1)
		}
		defer f.Close()
		r = f
	}

	// 解析 HTML 文档
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Fatal(err)
	}

	s, err := Html2Md(doc)
	if err != nil {
		fmt.Printf("Failed to trans, err: %v\n", err)
		return
	}

	// fmt.Println(s)

	ff, err := os.OpenFile(filepath.Join(output, doc.Find("title").Text()+".md"), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to output, err: %v\n", err)
		return
	}
	defer ff.Close()

	ff.WriteString(s)
	fmt.Println("Succeed to trans website to markdown file")
	// fmt.Println(doc.Html())

}

func Html2Md(doc *goquery.Document) (string, error) {

	main := doc.Find("div.css-376mun")
	ReplaceImgPaths(main)

	// fmt.Println(main.Html())

	h, err := doc.Find("div.css-376mun").Html()

	if err != nil {
		return "", err
	}

	converter := md.NewConverter("", true, nil)
	return converter.ConvertString(h)
}

func ReplaceImgPaths(doc *goquery.Selection) {

	doc.Find("figure").Each(func(i int, s *goquery.Selection) {
		s.ReplaceWithHtml(s.Find("noscript").Text())
		// fmt.Printf("Figure: %v\n", s.Find("noscript").First())
	})

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("src")
		if !ok {
			fmt.Println("Failed to get src")
			s.Remove()
			return
		}

		new, err := DownloadImgs(url, "./data/")
		if err != nil {
			fmt.Printf("download [%v] failed, err: %v\n", url, err)
			s.Remove()
			return
		}
		s.SetAttr("src", new)
		// s.ReplaceWithHtml(strings.TrimSpace(img.Text()))
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
	}
	newName := fmt.Sprintf("%d-%s%s", time.Now().Unix(), random.RandString(10), exts[len(exts)-1])

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
