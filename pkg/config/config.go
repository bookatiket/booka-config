package config

import (
	"fmt"
	"github.com/gojektech/heimdall/v6"
	"github.com/gojektech/heimdall/v6/httpclient"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"time"
)

// Setup ---
func Setup() (*Configuration, error) {
	env := os.Getenv("SERVICE_ENVIRONMENT")

	configFile := os.Getenv("CONFIG_FILE")
	configPath := os.Getenv("CONFIG_PATH")

	if len(configFile) == 0 {
		panic("config file doesn't exist")
	}

	if len(configPath) == 0 {
		panic("config path doesn't exist")
	}

	logrus.WithFields(logrus.Fields{"config.env": env, "config.file": configFile, "config.path": configPath}).Info("configuration environment")

	viper.SetConfigName(configFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Configuration{
		HttpClient: getHttpClient(),
		DB:         getDBConfig(),
		Server:     getServerConfig(),
	}, nil
}

// InitLog ---
func InitLog() {
	l, err := logrus.ParseLevel(viper.GetString("logger.level"))
	if err != nil {
		l = logrus.DebugLevel
	}
	logrus.SetLevel(l)

	if viper.GetString("SERVICE_ENVIRONMENT") == "prod" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		file, err := os.OpenFile(viper.GetString("logger.output"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logrus.SetOutput(file)
		} else {
			logrus.WithFields(logrus.Fields{"error": err.Error()}).Warning("error creating log file")
			logrus.Info("Failed to log to file, using default stderr")
		}
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

}

func getServerConfig() *ServerConfig {
	tickerDuration, err := time.ParseDuration(viper.GetString("ticker.interval"))
	if err != nil {
		tickerDuration = 300 * time.Second
	}
	return &ServerConfig{
		AirlineHost:    viper.GetString("aggregator.airline.host"),
		AirlinePort:    viper.GetString("aggregator.airline.port"),
		AirlinePath:    viper.GetString("aggregator.airline.path"),
		KaiHost:        viper.GetString("aggregator.kai.host"),
		KaiPort:        viper.GetString("aggregator.kai.port"),
		KaiPath:        viper.GetString("aggregator.kai.path"),
		TravelHost:     viper.GetString("aggregator.travel.host"),
		TravelPort:     viper.GetString("aggregator.travel.port"),
		TravelPath:     viper.GetString("aggregator.travel.path"),
		Port:           viper.GetString("app.search.port"),
		BookingPort:    viper.GetString("app.booking.port"),
		BookingHost:    viper.GetString("app.booking.host"),
		PricePort:      viper.GetString("app.price.port"),
		PriceHost:      viper.GetString("app.price.host"),
		AuthPort:       viper.GetString("app.auth.port"),
		AuthHost:       viper.GetString("app.auth.host"),
		ApiKey:         viper.GetString("api-key"),
		SecretKey:      viper.GetString("secret-key"),
		TickerDuration: tickerDuration,
	}
}

func getHttpClient() heimdall.Client {
	timedOut := 10 * time.Second
	t, err := time.ParseDuration(viper.GetString("http.timedout"))
	if err == nil {
		timedOut = t
	}
	return httpclient.NewClient(httpclient.WithHTTPTimeout(timedOut))
}

func getDBConfig() *DBConfig {
	return &DBConfig{
		Port:     viper.GetString("db.conn.port"),
		Host:     viper.GetString("db.conn.host"),
		Username: viper.GetString("db.conn.username"),
		Password: viper.GetString("db.conn.password"),
		Schema:   viper.GetString("db.conn.schema"),
	}
}

type DBConfig struct {
	Port     string
	Host     string
	Username string
	Password string
	Schema   string
}

type ServerConfig struct {
	AirlineHost    string
	AirlinePort    string
	AirlinePath    string
	KaiHost        string
	KaiPort        string
	KaiPath        string
	TravelHost     string
	TravelPort     string
	TravelPath     string
	Port           string
	BookingPort    string
	BookingHost    string
	PricePort      string
	PriceHost      string
	AuthPort       string
	AuthHost       string
	ApiKey         string
	SecretKey      string
	TickerDuration time.Duration
}

func (t *ServerConfig) GetTickerDuration() time.Duration {
	return t.TickerDuration
}

type Configuration struct {
	HttpClient heimdall.Client
	DB         *DBConfig
	Server     *ServerConfig
}

type StdFormatter struct{}

func (s *StdFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	format := "%s - - [%s] \"%s\" %s\n"
	log := fmt.Sprintf(format, entry.Level, entry.Time.String(), entry.Message, entry.Data)
	return []byte(log), nil
}
