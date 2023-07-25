package env

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type LoadRemoteEnvParams struct {
	Host    string
	Token   string
	AppName string
	Env     string
}

var ticker *time.Ticker
var tickerMutex sync.Mutex

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Erro carregando arquivo .env")
	}

	logrus.Info("Carregamento do arquivo .env realizado")
}
func LoadRemoteEnv(params LoadRemoteEnvParams) {
	logrus.Debug("Carregando variaveis de ambiente remotamente.")

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", params.Host, params.AppName, params.Env), nil)
	if err != nil {
		logrus.Fatalf("%s", err)
	}

	token := GetEnv("CLOUD_PROPERTIES_TOKEN", "")
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Fatalf("%s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Fatalf("não foi possível obter a configuração: %s", resp.Status)
	}

	var config map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		logrus.Fatalf("%s", err)
	}

	for k, v := range config {
		err := os.Setenv(strings.ToUpper(k), fmt.Sprintf("%v", v))
		if err != nil {
			return
		}
	}

	autoRefresh(params)
}

func GetEnv(variable string, default_ ...string) string {
	value := os.Getenv(variable)

	if value == "" {
		if len(default_) > 0 {
			return default_[0]
		}
		logrus.Fatalf("A variavel de ambiente %s não foi definida", variable)
	}

	return value
}

func autoRefresh(params LoadRemoteEnvParams) {
	refreshTime := GetEnv("ENV_REFRESH_TIME", "1")
	refreshInterval, err := strconv.Atoi(refreshTime)

	if err != nil {
		logrus.Fatal("ENV_REFRESH_TIME deve ser um número inteiro")
	}

	tickerMutex.Lock()
	defer tickerMutex.Unlock()

	ticker = time.NewTicker(time.Duration(refreshInterval) * time.Minute)
	go func() {
		for range ticker.C {
			LoadRemoteEnv(params)
		}
	}()
}
