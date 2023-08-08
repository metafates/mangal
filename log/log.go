package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	"github.com/mangalorg/mangal/fs"
	"github.com/mangalorg/mangal/path"
	"github.com/samber/lo"
)

var L = lo.Must(newLogger())

func newLogger() (logger *log.Logger, err error) {
	today := time.Now().Format("2006-01-02")

	logPath := filepath.Join(path.LogDir(), fmt.Sprint(today, ".log"))

	file, err := fs.Afero.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return nil, err
	}

	logger = log.New(file)

	logger.SetOutput(file)
	logger.SetFormatter(log.TextFormatter)

	return
}
