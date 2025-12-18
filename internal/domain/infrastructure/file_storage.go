package infrastructure

import (
	"io"
)

type FileStorage interface {
	SaveFile(userID int64, filename string, fileContent io.Reader) (string, error)
}
