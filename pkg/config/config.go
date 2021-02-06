package config

import (
	"fmt"
	"os"
	"time"

	"github.com/gojektech/heimdall/v6"
	"github.com/gojektech/heimdall/v6/httpclient"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
		HTTPClient: getHTTPClient(),
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
		PriceGrpcHost:  viper.GetString("app.price.grpc.host"),
		PriceGrpcPort:  viper.GetString("app.price.grpc.port"),
		AuthPort:       viper.GetString("app.auth.port"),
		AuthHost:       viper.GetString("app.auth.host"),
		APIKey:         viper.GetString("api-key"),
		SecretKey:      viper.GetString("secret-key"),
		TickerDuration: tickerDuration,
	}
}

func getHTTPClient() heimdall.Client {
	timedOut := 20 * time.Second
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

//DBConfig - struct for database configuration
type DBConfig struct {
	Port     string
	Host     string
	Username string
	Password string
	Schema   string
}

//ServerConfig - struct for server configuration
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
	PriceGrpcHost  string
	PriceGrpcPort  string
	AuthPort       string
	AuthHost       string
	APIKey         string
	SecretKey      string
	TickerDuration time.Duration
}

//GetTickerDuration - function to get ticker time
func (t *ServerConfig) GetTickerDuration() time.Duration {
	return t.TickerDuration
}

//Configuration - struct for main configuration
type Configuration struct {
	HTTPClient heimdall.Client
	DB         *DBConfig
	Server     *ServerConfig
}

//StdFormatter - standard log formatter
type StdFormatter struct{}

//Format - function for formatting log
func (s *StdFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	format := "%s - - [%s] \"%s\" %s\n"
	log := fmt.Sprintf(format, entry.Level, entry.Time.String(), entry.Message, entry.Data)
	return []byte(log), nil
}
