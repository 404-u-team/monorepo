package system

import (
	consts "DevSpace/Consts"
	"fmt"
	"io/fs"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Port uint
}

func NullConfig() Config {
	return Config{Port: 0}
}

func DefaultConfig() Config {
	return Config{Port: consts.APIPortDefault}
}

func GetConfig() (Config, error) {
	user := os.Getenv("USER")
	pathToDir := fmt.Sprintf("/home/%s/.config/DevSpace", user)
	pathToConf := fmt.Sprintf("%s/%s.toml", pathToDir, consts.ConfigName)

	if user == "" {
		return NullConfig(), EnvError{VarName: "USER"}
	}

	_, err := os.Stat(pathToDir)

	//то бишь такой папки нет
	if err != nil {
		mkdirError := os.Mkdir(pathToDir, fs.FileMode(0755)) // владелец может все, остальные читают
		//технически это невозхможно, но на всякий
		if mkdirError != nil {
			return NullConfig(), OSError{Err: fmt.Sprintf("Ошибка создания папки для конфига по пути: %s : %s", pathToDir, mkdirError.Error())}
		}
	}

	info, err := os.Stat(pathToConf)
	if err != nil {
		config, mkconfErr := CreateConfig(pathToConf)
		if mkconfErr != nil {
			return NullConfig(), OSError{Err: fmt.Sprintf("Ошибка создания конфига по пути %s : %s", pathToConf, mkconfErr.Error())}
		}

		return config, nil
	}

	//чем черт не шутит
	if info.IsDir() {
		rmdirErr := os.RemoveAll(pathToDir)
		if rmdirErr != nil {
			return NullConfig(), OSError{Err: fmt.Sprintf("Ошибка удаления непредвиденного каталога по пути %s : %s", pathToConf, rmdirErr.Error())}
		}
		_, mkconfErr := os.Create(pathToConf)
		if mkconfErr != nil {
			return NullConfig(), OSError{Err: fmt.Sprintf("Ошибка создания конфига по пути %s : %s", pathToConf, mkconfErr.Error())}
		}
	}

	var config Config

	_, parceError := toml.DecodeFile(pathToConf, &config)
	if parceError != nil {
		return NullConfig(), ParseError{}
	}

	return config, nil
}

func CreateConfig(path string) (Config, error) {
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		return NullConfig(), OSError{Err: fmt.Sprintf("Ошибка создания конфига по пути %s : %s", path, err.Error())}
	}

	encoder := toml.NewEncoder(file)
	encErr := encoder.Encode(DefaultConfig())
	if encErr != nil {
		return NullConfig(), encErr
	}

	return DefaultConfig(), nil
}

func DebugPrintConfigFile() {
	user := os.Getenv("USER")
	pathToConf := fmt.Sprintf("/home/%s/.config/DevSpace/%s.toml", user, consts.ConfigName)

	fmt.Println("Читаю файл:", pathToConf)
	data, err := os.ReadFile(pathToConf)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return
	}

	fmt.Println("Содержимое файла:")
	fmt.Println(string(data))
}
