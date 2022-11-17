package main

import (
	"errors"
	"io"
	"math"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrFileParse             = errors.New("file parsing error")
	ErrFileOpen              = errors.New("file opening error")
	ErrFileCreate            = errors.New("file creating error")
	ErrFileWrite             = errors.New("file writing error")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func getFile(fileName string, offset, limit int64) (file *os.File, readSize int, err error) {
	file, err = os.Open(fileName)
	if err != nil {
		err = ErrFileOpen
		return
	}

	fileInfo, err := file.Stat()
	if err != nil {
		err = ErrFileParse
		return
	}

	fileSize := int(fileInfo.Size())
	if fileSize == 0 {
		err = ErrUnsupportedFile
		return
	}

	if int(offset) > fileSize {
		err = ErrOffsetExceedsFileSize
		return
	}

	if _, err = file.Seek(offset, 0); err != nil {
		return
	}

	readSize = fileSize
	if offset > 0 {
		readSize -= int(offset)
	}
	if limit > 0 && readSize > int(limit) {
		readSize = int(limit)
	}

	return
}

func transfer(input io.Reader, output io.Writer, readSize int, writeProgress func(value int)) error {
	var counter int
	var currentChunkSize int
	defaultChunkSize := 10
	buffer := make([]byte, defaultChunkSize)
	for {
		currentChunkSize = int(math.Min(float64(readSize-counter), float64(defaultChunkSize)))
		if currentChunkSize != defaultChunkSize {
			buffer = make([]byte, currentChunkSize)
		}

		count, err := input.Read(buffer)
		if writeProgress != nil {
			writeProgress(count)
		}

		if errors.Is(err, io.EOF) {
			break
		}

		_, err = output.Write(buffer)
		if err != nil {
			return ErrFileWrite
		}

		counter += count
		if counter >= readSize {
			break
		}
	}

	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	sourceFile, readSize, err := getFile(fromPath, offset, limit)
	if err != nil {
		return ErrFileOpen
	}

	targetFile, err := os.Create(toPath)
	if err != nil {
		return ErrFileCreate
	}

	defer func() {
		_ = sourceFile.Close()
		_ = targetFile.Close()
	}()

	bar := pb.StartNew(readSize)
	defer bar.Finish()
	writeProgress := func(value int) {
		bar.Add(value)
	}

	err = transfer(sourceFile, targetFile, readSize, writeProgress)
	if err != nil {
		return err
	}

	return nil
}
