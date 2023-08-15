package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mangalorg/mangal/fs"
	"github.com/mangalorg/mangal/path"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
)

var L = lo.Must(newLogger())

func newLogger() (*zerolog.Logger, error) {
	today := time.Now().Format("2006-01-02")

	logPath := filepath.Join(path.LogDir(), fmt.Sprint(today, ".log"))

	file, err := fs.Afero.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return nil, err
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	logger := zerolog.New(file)

	return &logger, nil
}
