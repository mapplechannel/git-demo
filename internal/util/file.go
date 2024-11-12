package util

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"fmt"
	"hsm-scheduling-back-end/pkg/logger"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// 读取文件数据
func ReadFile(path string) (data string) {
	//打开文件
	open, err := os.Open(path)
	if err != nil {
		logger.Error("打开文件失败 %s", err.Error())
	}
	defer open.Close()
	file, err := ioutil.ReadAll(open)
	if err != nil {
		return
	}
	return string(file)
}

// 读取最后一行文件数据
func ReadLineRawFile(path string) (data string) {
	//打开文件
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("打开文件失败 %s", err.Error())
	}
	defer file.Close()
	//创建一个Scanner对象
	scanner := bufio.NewScanner(file)

	var lastLine string
	for scanner.Scan() {
		lastLine = scanner.Text()

	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}
	return lastLine
}

func WriteContentAppend(path string, content string) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	_, err = file.Write([]byte(content + "\n"))
	if err != nil {
		fmt.Println(err)
	}
}

func WriteFile(filePath, data string) bool {
	//写文件
	file, error := os.Create(filePath)
	if error != nil {
		fmt.Println(error)
		return false
	}
	defer file.Close()
	//写入byte的slice数据
	_, err := file.Write([]byte(data))
	if err != nil {
		log.Fatal(err)
		return false
	}
	// log.Printf("Wrote %d bytes.", bytesWritten)
	return true
}

// 创建文件
func CreateFile(path string) {
	_, err1 := os.Stat(path)
	logger.Info("创建文件：%v", path)
	if err1 != nil {
		file, error := os.Create(path)
		if error != nil {
			fmt.Println(error)
		}
		defer file.Close()
	}
}

// 创建文件目录
func CreateDir(path string) bool {
	_, err1 := os.Stat(path)
	//存在返回错误
	if err1 != nil {
		if os.IsNotExist(err1) {
			//文件夹不存在就创建
			logger.Info("数据文件不存在,正在创建文件目录: %s", path)
			return Mkdir(path)
		} else {
			return false
		}
	} else {
		return true
	}
}

func Mkdir(path string) bool {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		logger.Error("创建文件失败: %s" + err.Error())
		return false
	} else {
		return true
	}

}

// 判断目录文件是否存在
// 存在 true
// 不存在 false
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func PathExists1(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	return false
}

// 判断是否是GBK 格式的文件
func IsGBK(data []byte) bool {
	length := len(data)
	var i int = 0
	for i < length {
		if data[i] <= 0x7f {
			//编码0~127,只有一个字节的编码，兼容ASCII码
			i++
			continue
		} else {
			//大于127的使用双字节编码，落在gbk编码范围内的字符
			if data[i] >= 0x81 &&
				data[i] <= 0xfe &&
				data[i+1] >= 0x40 &&
				data[i+1] <= 0xfe &&
				data[i+1] != 0xf7 {
				i += 2
				continue
			} else {
				return false
			}
		}
	}
	return true
}

func IsUtf8(data []byte) bool {
	// 遍历字节切片中的每个字节
	for i := 0; i < len(data); {
		// 如果字节的最高位为 0，说明是单字节字符，直接跳过
		if data[i]&0x80 == 0x00 {
			i++
			continue
		} else {
			// 否则，根据字节的最高位连续的 1 的个数，判断该字符占用的字节数
			num := preNUm(data[i])
			// 如果字节数大于 2，说明是多字节字符
			if num > 2 {
				// 跳过第一个字节
				i++
				// 遍历后面的 num - 1 个字节，判断是否都以 10 开头
				for j := 0; j < num-1; j++ {
					// 如果不是以 10 开头，说明不是 UTF-8 编码，返回 false
					if data[i]&0xc0 != 0x80 {
						return false
					}
					i++
				}
			} else {
				// 否则，说明不是 UTF-8 编码，返回 false
				return false
			}
		}
	}
	// 如果遍历完所有字节都没有发现不符合 UTF-8 编码规则的情况，返回 true
	return true
}

// preNUm 函数用来计算一个字节的最高位连续的 1 的个数
// 参数 data 为要计算的字节
// 返回值为一个整数，表示最高位连续的 1 的个数
func preNUm(data byte) int {
	// 将字节转换为二进制字符串
	str := fmt.Sprintf("%b", data)
	var i int = 0
	// 遍历字符串中的每一位，如果是 1 就计数，如果是 0 就停止
	for i < len(str) {
		if str[i] != '1' {
			break
		}
		i++
	}
	return i
}

// 将文件打成zip包
func ZipFiles(filePaths []string, zipFileName string) error {
	newZipFile, err := os.Create(zipFileName)
	fmt.Println("创建文件zip文件夹")
	if err != nil {
		return err
	}
	defer newZipFile.Close()
	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	for _, filePath := range filePaths {
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return err
		}
		header.Name = filepath.Base(filePath)
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func CompressFolderToZip(folderPath, zipPath string) error {
	// 创建一个新的ZIP文件
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// 创建一个Writer来写入ZIP文件
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 遍历文件夹
	err = filepath.Walk(folderPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 获取文件在文件夹中的相对路径
		relPath, err := filepath.Rel(folderPath, filePath)
		if err != nil {
			return err
		}

		// 如果是文件夹，则在ZIP文件中创建相应的文件夹
		if info.IsDir() {
			_, err = zipWriter.Create(relPath + "/")
			if err != nil {
				return err
			}
		} else {
			// 如果是文件，则将文件内容写入到ZIP文件中
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			// 创建一个在ZIP文件中的新文件
			zipFile, err := zipWriter.Create(relPath)
			if err != nil {
				return err
			}

			// 将文件内容复制到ZIP文件中的新文件中
			_, err = io.Copy(zipFile, file)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	fmt.Println("文件夹压缩为ZIP成功：", zipPath)
	return nil
}

// 从 ZIP 文件中读取文件内容
func ReadFileContentsFromZIP(zipPath string) (map[string]string, error) {
	fileContents := make(map[string]string)
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()
		content, err := io.ReadAll(rc)
		if err != nil {
			return nil, err
		}
		fileContents[f.Name] = string(content)
	}
	return fileContents, nil
}

// 保存文件内容到指定路径
func SaveFileContent(savePath string, content string) error {
	outputFile, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = outputFile.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

// 删除文件目录
func RemoveFolder(folderPath string) bool {
	err := os.RemoveAll(folderPath)
	if err != nil {
		logger.Error("删除文件失败：%v", folderPath)
		return false
	}
	return true
}

// 将文件添加到tar包中
func addFileToTar(filePath string, tarWriter *tar.Writer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	// 创建tar头部信息
	header, err := tar.FileInfoHeader(fileInfo, "")
	if err != nil {
		return err
	}

	// 写入tar头部信息
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	// 写入文件内容
	if _, err := io.Copy(tarWriter, file); err != nil {
		return err
	}

	return nil
}

// 复制文件到指定文件夹下，若文件夹不存在则创建
func CopyFileToFloder(sourceFile, destinationFolder string) bool {
	// 判断目标文件夹是否存在，不存在则创建
	_, err := os.Stat(destinationFolder)
	if os.IsNotExist(err) {
		err = os.MkdirAll(destinationFolder, os.ModePerm)
		if err != nil {
			return false
		}
	}
	// 打开源文件
	src, err := os.Open(sourceFile)
	if err != nil {
		return false
	}
	defer src.Close()
	// 获取源文件名
	fileName := filepath.Base(sourceFile)
	// 创建目标文件
	dstPath := filepath.Join(destinationFolder, fileName)
	dst, err := os.Create(dstPath)
	if err != nil {
		return false
	}
	defer dst.Close()
	// 复制文件内容
	_, err = io.Copy(dst, src)
	if err != nil {
		return false
	}
	return true
}

func CopyFile(src, destFolder string) {
	// 源文件路径
	// src := "/path/to/source/file.txt"

	// 目标文件夹路径
	// destFolder := "/path/to/destination/folder"

	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()

	// 创建目标文件
	destFile, err := os.Create(destFolder)
	if err != nil {
		panic(err)
	}
	defer destFile.Close()

	// 复制文件内容
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		panic(err)
	}
}

func UnzipFileToPath(zipFilePath, extractDir string) {

	// zipFilePath := "/home/hollysys/testdata/ceshi/job_peah3f30_test789.zip"
	// extractDir := "/home/hollysys/testdata/ceshi/gaosheng"
	zipFile, err := zip.OpenReader(zipFilePath)
	if err != nil {
		fmt.Println("打开 ZIP 文件出错:", err)
		return
	}
	defer zipFile.Close()

	// 创建目标目录
	err = os.MkdirAll(extractDir, os.ModePerm)
	if err != nil {
		fmt.Println("创建目标目录出错:", err)
		return
	}

	// 解压 ZIP 文件
	for _, file := range zipFile.File {
		// 构建解压后的文件路径
		extractedFilePath := filepath.Join(extractDir, file.Name)

		// 如果是目录，则创建对应的目录
		if file.FileInfo().IsDir() {
			os.MkdirAll(extractedFilePath, os.ModePerm)
			continue
		}

		// 创建解压后的文件
		extractedFile, err := os.Create(extractedFilePath)
		if err != nil {
			fmt.Println("创建解压后的文件出错:", err)
			return
		}
		defer extractedFile.Close()

		// 打开 ZIP 文件中的文件
		zipFile, err := file.Open()
		if err != nil {
			fmt.Println("打开 ZIP 文件中的文件出错:", err)
			return
		}
		defer zipFile.Close()

		// 将 ZIP 文件中的内容复制到解压后的文件中
		_, err = io.Copy(extractedFile, zipFile)
		if err != nil {
			fmt.Println("解压文件出错:", err)
			return
		}
	}

	fmt.Println("ZIP 文件解压成功")
}

// 复制文件目录下的所有文件
func CopyFilesToFolder(srcDir, dstDir string) error {
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 获取相对于源目录的路径
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		// 创建目标文件路径
		dstPath := filepath.Join(dstDir, relPath)
		if info.IsDir() {
			// 如果是目录，创建目标目录
			err = os.MkdirAll(dstPath, info.Mode())
			if err != nil {
				return err
			}
		} else {
			// 如果是文件，复制文件内容
			err = copyFile(path, dstPath)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// 复制文件
func copyFile(srcPath, dstPath string) error {
	// 打开源文件
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("打开源文件出错: %v", err)
	}
	defer srcFile.Close()

	// 创建或打开目标文件
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("创建或打开目标文件出错: %v", err)
	}
	defer dstFile.Close()

	// 复制文件内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("复制文件内容出错: %v", err)
	}

	return nil
}
