package configs

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// O mapstructure vai ler os arquivos .env e jogar nas variaveis do struct

type Config struct {
	RedisSrc      string `mapstructure:"REDIS_SRC"`
	RedisPass     string `mapstructure:"REDIS_PASS"`
	AllowedTokens map[string][]int
}

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}
	// Aqui, você está definindo o nome do arquivo de configuração como "app_config".
	// Isso indica que o Viper procurará por um arquivo de configuração chamado "app_config"
	//  quando você tentar ler a configuração.
	viper.SetConfigName("app_config")
	// Aqui, você está definindo o tipo de configuração como "env". Isso sugere que o
	// Viper espera que o arquivo de configuração esteja no formato de variáveis de ambiente,
	// comum em arquivos .env.
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(path + "/.env")
	// Esta linha indica ao Viper para automaticamente ler variáveis de ambiente com nomes correspondentes
	//  às chaves no arquivo de configuração. Por exemplo, se o arquivo de configuração tiver uma
	//  chave "DATABASE_URL", o Viper procurará automaticamente por uma variável de ambiente chamada
	//  "DATABASE_URL" para definir o valor correspondente.
	viper.AutomaticEnv()
	// Aqui, você está tentando ler o arquivo de configuração com base nas configurações anteriores.
	// Se for bem-sucedido, o conteúdo do arquivo de configuração será carregado nas configurações do
	// Viper. Se ocorrer algum erro durante a leitura, o erro será retornado na variável err
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(cfg)
	if err != nil {
		panic(err)
	}
	// cfg.AllowedTokens = map[string][]int{
	// 	// "token": {req/s, tempo de bloqueio em segundos}
	// 	"token1":      {100, 0},
	// 	"token2":      {30, 1000},
	// 	"192.168.0.1": {100, 0},
	// }
	cfg.AllowedTokens = getTokens()
	return cfg, nil
}

func getTokens() map[string][]int {
	AllowedTokens := make(map[string][]int)
	i := 1
	for {
		tokenKey := fmt.Sprintf("TOKEN_%d", i)
		if !viper.IsSet(tokenKey) {
			break
		}
		tokenValue := viper.GetString(tokenKey)
		parseAndAddToMap(tokenValue, AllowedTokens)
		i++
	}
	j := 1
	for {
		ipKey := fmt.Sprintf("IP_LIMIT_%d", j)
		if !viper.IsSet(ipKey) {
			break
		}
		ipValue := viper.GetString(ipKey)
		parseAndAddToMap(ipValue, AllowedTokens)
		j++
	}
	return AllowedTokens
}

func parseAndAddToMap(value string, tokenMap map[string][]int) {
	parts := strings.Split(value, ":")
	if len(parts) != 3 {
		log.Printf("Invalid format for value: %s", value)
		return
	}

	tokenOrIP := parts[0]
	reqPerSec, err1 := strconv.Atoi(parts[1])
	blockTime, err2 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil {
		log.Printf("Error parsing limits for %s: reqPerSec %v, blockTime %v", tokenOrIP, err1, err2)
		return
	}

	tokenMap[tokenOrIP] = []int{reqPerSec, blockTime}
}
