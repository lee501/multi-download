package filedownloader

import (
	"net/http"
	"os"
	logger "multi-download/lib"
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
		logger.Logger
	}
}