package flagger

import "flag"

var (
	Sha256string string
	FileSavePath string
	TotalPart int
	Url string
	Filename string
)

func Init() {
	flag.StringVar(&FileSavePath, "path", "/Users/lee/workspace", "input file save path")
	flag.StringVar(&Sha256string, "sha256", "", "input file sha256 string")
	flag.IntVar(&TotalPart, "total", 10, "input job number")
	flag.StringVar(&Url, "url", "", "input file download url")
	flag.StringVar(&Filename, "name", "", "input file save name")
}
