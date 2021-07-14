package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func validate(fromPath string, offset int64) error {
	fi, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	if !fi.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset > 0 && offset >= fi.Size() {
		return ErrOffsetExceedsFileSize
	}

	return nil
}

type progressWriter struct {
	limit   int64
	printed int64
}

func (pw *progressWriter) Write(b []byte) (n int, err error) {
	if len(b) == 0 {
		fmt.Printf("\r[%s>] %.2f%%", strings.Repeat("=", 100), 100.0)
	}

	for k := range b {
		pCur := float32(pw.printed+int64(k)+1) / float32(pw.limit) * 100
		pLeft := 100 - pCur
		fmt.Printf("\r[%s>%s] %.2f%%", strings.Repeat("=", int(pCur)),
			strings.Repeat("_", int(pLeft)), pCur)
	}

	pw.printed += int64(len(b))

	return len(b), nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	err := validate(fromPath, offset)
	if err != nil {
		return err
	}

	f, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	if limit == 0 || limit+offset > fi.Size() {
		limit = fi.Size() - offset
	}

	wF, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer wF.Close()
	mw := io.MultiWriter(wF, &progressWriter{limit: limit})

	buf := make([]byte, 512)
	var doneCnt int64

	for {
		cnt, err := f.ReadAt(buf, offset+doneCnt)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}

		if cnt > int(limit-doneCnt) {
			cnt = int(limit - doneCnt)
		}

		amount, err := mw.Write(buf[:cnt])
		if err != nil {
			return err
		}

		doneCnt += int64(amount)

		if limit <= doneCnt {
			break
		}
	}

	fmt.Println()

	err = f.Sync()
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	err = wF.Close()
	if err != nil {
		return err
	}

	return nil
}
