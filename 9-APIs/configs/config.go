package configs

import (
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

// O mapstructure vai ler os arquivos .env e jogar nas variaveis do struct

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBHost        string `mapstructure:"DB_HOST"`
	DBPort        string `mapstructure:"DB_PORT"`
	DBUser        string `mapstructure:"DB_USER"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBName        string `mapstructure:"DB_NAME"`
	WebServerPort string `mapstructure:"WEB_SERVER_PORT"`
	JWTSecret     string `mapstructure:"JWT_SECRET"`
	JWTExpiresIn  int    `mapstructure:"JWT_EXPIRESIN"`
	TokenAuth     *jwtauth.JWTAuth
}

func LoadConfig(path string) (*Config, error) {
	var cfg *Config
	// Aqui, você está definindo o nome do arquivo de configuração como "app_config".
	// Isso indica que o Viper procurará por um arquivo de configuração chamado "app_config"
	//  quando você tentar ler a configuração.
	viper.SetConfigName("app_config")
	// Aqui, você está definindo o tipo de configuração como "env". Isso sugere que o
	// Viper espera que o arquivo de configuração esteja no formato de variáveis de ambiente,
	// comum em arquivos .env.
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
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
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	cfg.TokenAuth = jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)
	return cfg, nil
}

// será a primeira funçao a ser executada
// func init() {

// }
