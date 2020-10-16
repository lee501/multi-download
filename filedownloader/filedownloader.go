package filedownloader

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/sevenelevenlee/multi-download/lib/logger"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
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
	fileSize, err := fd.head()
	if err != nil {
		return err
	}
	fd.FileSize = fileSize

	//设置分片job
	jobs := make([]filePart, fd.TotalPart)
	eachSize := fileSize / fd.TotalPart
	log.Println(fileSize, eachSize)
	//指定filePart编号，起始位置
	for i := range jobs {
		jobs[i].Index = i
		if i == 0 {
			jobs[i].From = 0
		} else {
			jobs[i].From = jobs[i-1].To + 1
		}
		if i < fd.TotalPart - 1 {
			jobs[i].To = jobs[i].From + eachSize
		} else {
			jobs[i].To = fileSize - 1
		}
	}

	var wg sync.WaitGroup
	for jobId, job := range jobs {
		wg.Add(1)
		go fd.worker(&wg, jobId, job)
	}
	wg.Wait()
	return fd.mergeFilePart()
}

//worker
func (fd *FileDownloader) worker(wg *sync.WaitGroup, jobId int, job filePart) {
	defer wg.Done()
	err := fd.download(job)
	if err != nil {
		logger.Logger.Error("file download fail", zap.String("jobId", strconv.Itoa(jobId)), zap.Error(err))
	}
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
		return 0, errors.New("server do not support break point download")
	}
	//设置文件名
	if fd.FileName == "" {
		fd.FileName = parseFileInfo(response)
	}
	return strconv.Atoi(response.Header.Get("Content-Length"))
}

//下载
func (fd *FileDownloader) download(part filePart) error {
	request, err := fd.getNewRequest("GET")
	if err != nil {
		return err
	}
	log.Printf("开始[%d]下载from:%d to:%d\n", part.Index, part.From, part.To)
	request.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", part.From, part.To))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode > 299 {
		return errors.New(fmt.Sprintf("server code error: %v", response.StatusCode))
	}
	defer response.Body.Close()
	//处理下载数据
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if len(bytes) != (part.To - part.From + 1){
		return errors.New("file part size error")
	}
	part.Data = bytes
	fd.DoneFilePart[part.Index] = part
	return nil
}

func parseFileInfo(res *http.Response) string {
	contentDisposition := res.Header.Get("Content-Disposition")
	if contentDisposition != "" {
		_, params, err := mime.ParseMediaType(contentDisposition)
		if err != nil {
			panic(err)
		}
		return params["filename"]
	}
	filename := filepath.Base(res.Request.URL.Path)
	return filename
}

//merge file part
func (fd *FileDownloader) mergeFilePart() error {
	log.Println("start merge file")
	//文件保存路径
	path := filepath.Join(fd.FileSavePath, fd.FileName)
	file, err := os.Create(path)
	if err != nil {
		return nil
	}
	defer file.Close()
	//校验文件
	hash := sha256.New()
	totalSize := 0

	//处理逻辑
	for _, filePart := range fd.DoneFilePart {
		file.Write(filePart.Data)
		if len(fd.Sha256string) != 0 {
			hash.Write(filePart.Data)
		}
		totalSize += len(filePart.Data)
	}

	if totalSize != fd.FileSize {
		return errors.New("file size not whole")
	}
	if fd.Sha256string != "" && hex.EncodeToString(hash.Sum(nil)) != fd.Sha256string {
		return errors.New("sha256 validate fail")
	}
	return nil
}
