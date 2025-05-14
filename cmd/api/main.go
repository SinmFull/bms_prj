package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"strings"
	"time"

	"github.com/SinmFull/BMS_prj/internal/data"
	"github.com/SinmFull/BMS_prj/internal/jsonlog"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	mqtt struct {
		user     string
		password string
	}
}

type application struct {
	config     config
	logger     *jsonlog.Logger
	models     data.Models
	mqttClient mqtt.Client
}

func main() {
	var cfg config

	// viper.SetConfigName("config")
	// viper.SetConfigType("yaml")
	// viper.AddConfigPath(".")
	// if err := viper.ReadInConfig(); err != nil {
	// 	fmt.Printf("Error reading config file, %s", err)
	// }
	// fmt.Println("viper config:", viper.GetString("database.dsn"))

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "MySQL DSN")
	flag.StringVar(&cfg.mqtt.user, "mqtt-user", "admin", "MQTT user")
	flag.StringVar(&cfg.mqtt.password, "mqtt-password", "qwe123456", "MQTT password")

	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)

	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetUsername(cfg.mqtt.user).SetPassword(cfg.mqtt.password)
	mqttClient := mqtt.NewClient(opts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		logger.PrintFatal(token.Error(), nil)
	}

	app := &application{
		config:     cfg,
		logger:     logger,
		models:     data.NewModels(db),
		mqttClient: mqttClient,
	}

	app.mqttClient.Subscribe("MQTT_RT_DATA", 0, app.mqttMessageHandler)
	app.mqttClient.Subscribe("MQTT_ENY_NOW", 0, app.mqttMinuteMessageHandler)
	app.mqttClient.Subscribe("MQTT_DAY_DATA", 0, app.mqttDayMessageHandler)

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	dsn := strings.TrimPrefix(cfg.db.dsn, "mysql://")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
