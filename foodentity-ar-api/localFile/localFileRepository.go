package localFile

import (
	"encoding/base64"
	"fmt"
	"foodentity-ar-api/model"
	"os"
	"strings"
)

func NewLocalFileRepositoryImpl() LocalFileRepository {
	return &localFileRepositoryImpl{}
}

type LocalFileRepository interface {
	Add(req *model.Request, imageName string) error
}

type localFileRepositoryImpl struct {
	Request *model.Request
}

func (impl *localFileRepositoryImpl) Add(req *model.Request, imageName string) error {
	var base64Data = req.Image[strings.IndexByte(req.Image, ',')+1:]
	fmt.Println("base64Data: " + base64Data)

	data, decodeErr := base64.StdEncoding.DecodeString(base64Data)
	if decodeErr != nil {
		return decodeErr
	}

	file, err := os.Create("/tmp/" + imageName)
	if err != nil {
		return err
	}

	defer file.Close()
	file.Write(data)
	fmt.Printf("Saved to '%v' tempolarily.\n", file.Name())

	return nil
}
