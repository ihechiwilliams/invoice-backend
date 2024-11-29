package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
)

type goEnv struct {
	GoMod string `json:"GOMOD"`
}

func LoadConfig(cfg interface{}) error {
	envFilePath, fileExists := initLocal()

	var readCfgErr error

	if fileExists {
		readCfgErr = cleanenv.ReadConfig(envFilePath, cfg)
	} else {
		readCfgErr = cleanenv.ReadEnv(cfg)
	}

	if readCfgErr != nil {
		return readCfgErr
	}

	updateEnvErr := cleanenv.UpdateEnv(cfg)
	if updateEnvErr != nil {
		return updateEnvErr
	}

	return nil
}

func initLocal() (string, bool) {
	if !goExists() {
		return "", false
	}

	modRoot := getModuleRoot()
	envFilePath := fmt.Sprintf("%s/.env", modRoot)

	if envFileExists(envFilePath) {
		return envFilePath, true
	}

	return "", false
}

func envFileExists(envFilePath string) bool {
	_, err := os.Stat(envFilePath)

	return err == nil
}

func goExists() bool {
	_, err := exec.LookPath("go")

	return err == nil
}

func getModuleRoot() string {
	goEnvRaw, err := exec.Command("go", "env", "-json").Output()
	if err != nil {
		log.Fatal().Err(err).Msg("go env command failed")

		return ""
	}

	env := new(goEnv)

	err = json.Unmarshal(goEnvRaw, env)
	if err != nil {
		log.Fatal().Err(err).Msg("go mod unmarshalling failed")

		return ""
	}

	return strings.TrimSuffix(env.GoMod, "/go.mod")
}
