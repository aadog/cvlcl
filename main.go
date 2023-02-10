package main

import (
	"archive/zip"
	"bytes"
	"compress/zlib"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
)

var (
	gPath string
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("需要3个参数")
		return
	}
	// macOS不支持这种，也不需要支持这种
	genresFile(os.Args[1], os.Args[2])
}

func readZipData(ff *zip.File) []byte {
	if rr, err := ff.Open(); err == nil {
		defer rr.Close()
		bs, err := ioutil.ReadAll(rr)
		if err != nil {
			return nil
		}
		return bs
	}
	return nil
}

// zlib压缩
func ZlibCompress(input []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := zlib.NewWriterLevel(&in, zlib.BestCompression)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(input)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return in.Bytes(), nil
}

func genresByte(input []byte, newFileName string) {
	fmt.Println("genFile: ", newFileName)
	if len(input) == 0 {
		fmt.Println("000000")
		return
	}

	crc32Val := crc32.ChecksumIEEE(input)

	//压缩
	bs, err := ZlibCompress(input)
	if err != nil {
		panic(err)
	}
	code := bytes.NewBuffer(nil)
	code.WriteString("package liblclbinres")
	code.WriteString("\r\n\r\n")
	code.WriteString(fmt.Sprintf("const CRC32Value uint32 = 0x%x\r\n\r\n", crc32Val))

	code.WriteString("var LCLBinRes = []byte(\"")
	for _, b := range bs {
		code.WriteString("\\x" + fmt.Sprintf("%.2x", b))
	}
	code.WriteString("\")\r\n")
	ioutil.WriteFile(newFileName, code.Bytes(), 0666)
}

// 生成字节的单元
func genresFile(fileName, newFileName string) {
	bs, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	genresByte(bs, newFileName)
}
