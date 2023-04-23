package datalogger

import (
	"fmt"
	"os"
	"path"
)

const MaxSizeFile int = 10 * 1024 * 1024
const BckCount int = 10

//rotatingFileLog writes log a file, if file size exceeds maxBytes,
//it will backup current file and open a new one.
//
//max backup file number is set by backupCount, it will delete oldest if backups too many.
type RotatingFileLog struct {
	Fdescr *os.File

	fileName    string
	maxBytes    int
	backupCount int
}

// NewRotatingFileLog creates dirs and opens the logfile
func NewRotatingFileLog(fileName string, maxBytes int, backupCount int) (*RotatingFileLog, error) {
	dir := path.Dir(fileName)
	os.Mkdir(dir, 0777)

	h := new(RotatingFileLog)
	if maxBytes <= 0 {
		return nil, fmt.Errorf("invalid max bytes")
	}

	h.fileName = fileName
	h.maxBytes = maxBytes
	h.backupCount = backupCount

	var err error
	h.Fdescr, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (h *RotatingFileLog) doRollover() {
	f, err := h.Fdescr.Stat()
	if err != nil {
		return
	}

	if h.maxBytes <= 0 {
		return
	} else if f.Size() < int64(h.maxBytes) {
		return
	}

	if h.backupCount > 0 {
		h.Fdescr.Close()

		for i := h.backupCount - 1; i > 0; i-- {
			sfn := fmt.Sprintf("%s.%d", h.fileName, i)
			dfn := fmt.Sprintf("%s.%d", h.fileName, i+1)

			os.Rename(sfn, dfn)
		}

		dfn := fmt.Sprintf("%s.1", h.fileName)
		os.Rename(h.fileName, dfn)

		h.Fdescr, _ = os.OpenFile(h.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	}
}

func (h *RotatingFileLog) Write(p []byte) (n int, err error) {
	h.doRollover()
	return h.Fdescr.Write(p)
}

// Close simply closes the File
func (h *RotatingFileLog) Close() error {
	if h.Fdescr != nil {
		return h.Fdescr.Close()
	}
	return nil
}
