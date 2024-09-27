package work_report

import (
	"archive/zip"
	"fmt"
	"io"
)

const (
	MaxDocNameLen              = 100
	MaxWorkReportTopicTitleLen = 300
	ErrorClosingZipFile        = "error closing zip file: %s"
	ErrorClosingFileReader     = "error closing file reader: %s"
)

func UnzipWorkReport(zipFilePath string) ([]byte, error) {
	zipFile, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return nil, err
	}
	defer closeZipFile(zipFile)

	for _, file := range zipFile.File {
		if !file.FileInfo().IsDir() {
			return readFile(file)
		}
	}

	return nil, nil
}

func closeZipFile(zipFile *zip.ReadCloser) {
	if err := zipFile.Close(); err != nil {
		_ = fmt.Errorf(ErrorClosingZipFile, err)
	}
}

func readFile(file *zip.File) ([]byte, error) {
	fileReader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer closeFileReader(fileReader)

	return io.ReadAll(fileReader)
}

func closeFileReader(fileReader io.ReadCloser) {
	if err := fileReader.Close(); err != nil {
		_ = fmt.Errorf(ErrorClosingFileReader, err)
	}
}
