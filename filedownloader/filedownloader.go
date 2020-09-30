package filedownloader

import (
	"errors"
	"github.com/sevenelevenlee/multi-download/lib/logger"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strconv"
)

//文件下载器
type FileDownloader struct {
	FileName 	 string
	Url   	 	 string
	FileSize 	 int
	TotalPart 	 int //分片数量
	FileSavePath string //文件保存路径
	DoneFilePart []filePart //数据切片
	Sha256string string //sha256校验
}

type filePart struct {
	Index 	int     //编号
	From	int		//分片起始位置
	To 		int		//结束位置
	Data    []byte  //保存字节数据
}

func NewFileDownloader(url, filename, Sha256string, fileSavePath string, totalPart int) *FileDownloader {
	if fileSavePath == "" {
		wd, _ := os.Getwd()
		fileSavePath = wd
	}
	return &FileDownloader{
		FileName: filename,
		Url: url,
		TotalPart: totalPart,
		FileSavePath: fileSavePath,
		DoneFilePart: make([]filePart, totalPart),
		Sha256string: Sha256string,
	}
}

func (fd *FileDownloader) Run() error {
	
}

//创建http request
func (fd *FileDownloader) getNewRequest(method string) (*http.Request, error)  {
	r, err := http.NewRequest(method, fd.Url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36")
	return r, nil
}

//head请求 检查是否支持断点续传，文件大小
func (fd *FileDownloader) head() (int, error){
	request, err := fd.getNewRequest("HEAD")
	if err != nil {
		logger.Logger.Error("create http request fail", zap.String("method", "head"), zap.Error(err))
		return 0, err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Logger.Error("head request fail", zap.String("method", "head"), zap.Error(err))
		return 0, err
	}
	if response.Header.Get("Accept-Ranges") != "bytes" {
		logger.Logger.Info("服务器不支持断点续传")
		return 0, errors.New("服务器不支持断点续传")
	}
	return strconv.Atoi(response.Header.Get("Content-Length"))
}