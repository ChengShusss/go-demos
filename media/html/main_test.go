package main

import (
	"fmt"
	"testing"
)

func TestDownloadImg(t *testing.T) {
	_, err := DownloadImgs(
		"https://picx.zhimg.com/v2-27c9ef9c47f6f2025f76cd64f8fac535_r.jpg?source=2c26e567",
		"/home/cheng/workSpace/codeSpace/tinyGoProjects/go-demos/media/html/data/")

	fmt.Printf("err: %v\n", err)
}
