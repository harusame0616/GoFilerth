package setting

import (
	"log"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/harusame0616/GoFilerth/gofilerth/usecase"
)

// ファイル種類ごとの設定
type fileTypeConfig struct {
	Color string
}

var fileTypeConfigMap map[string]fileTypeConfig

const (
	SETTING_FILE_PATH_USER_HOME = "~/.gofilerth/type.toml"
	SETTING_FILE_CURRENT_FOLDER = "./type.toml"
)

func init() {
	var decodeError error

	if _, err := os.Stat(SETTING_FILE_PATH_USER_HOME); err == nil {
		_, decodeError = toml.DecodeFile(SETTING_FILE_PATH_USER_HOME, &fileTypeConfigMap)
	} else if _, err := os.Stat(SETTING_FILE_CURRENT_FOLDER); err == nil {
		_, decodeError = toml.DecodeFile(SETTING_FILE_CURRENT_FOLDER, &fileTypeConfigMap)
	} else {
		// do nothing
	}

	if decodeError != nil {
		log.Fatal(decodeError)
	}
}

func GetColor(fileDto usecase.FileDto) (int32, error) {
	config, ok := fileTypeConfigMap[fileDto.FileType]

	if ok == false || config.Color == "" {
		return 0xFFFFFF, nil
	}

	color, err := strconv.ParseInt(config.Color[1:], 16, 64)
	if err != nil {
		return 0, err
	}

	return int32(color), nil
}
