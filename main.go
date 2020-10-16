package main

import (
	"flag"
	"github.com/sevenelevenlee/multi-download/filedownloader"
	"github.com/sevenelevenlee/multi-download/lib/flagger"
	"github.com/sevenelevenlee/multi-download/lib/logger"
	"log"
	"os"
	"time"
	"fmt"
)

func main() {
	// 初始化变量
	flagger.Init()
	// 把用户传递的命令行参数解析为对应变量的值
	flag.Parse()

	if flagger.Url == "" {
		logger.Logger.Info("please input the file download site")
		os.Exit(-1)
	}
	startTime := time.Now()
	//url, filename, Sha256string, fileSavePath string, totalPart int
	downloader := filedownloader.NewFileDownloader(flagger.Url, flagger.Filename, flagger.Sha256string, flagger.FileSavePath, flagger.TotalPart)
	if err := downloader.Run(); err != nil {
		// fmt.Printf("\n%s", err)
		log.Fatal(err)
	}
	fmt.Printf("\n 文件下载完成耗时: %f second\n", time.Now().Sub(startTime).Seconds())
}
