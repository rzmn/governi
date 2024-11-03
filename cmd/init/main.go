package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"verni/internal/common"
	"verni/internal/db"
	"verni/internal/services/logging"
)

func main() {
	logger := logging.DefaultService()
	root, present := os.LookupEnv("VERNI_PROJECT_ROOT")
	if present {
		common.RegisterRelativePathRoot(root)
	}
	configFile, err := os.Open(common.AbsolutePath("./config/prod/verni.json"))
	if err != nil {
		log.Fatalf("failed to open config file: %s", err)
	}
	defer configFile.Close()
	configData, err := io.ReadAll(configFile)
	if err != nil {
		log.Fatalf("failed to read config file: %s", err)
	}
	type Module struct {
		Type   string                 `json:"type"`
		Config map[string]interface{} `json:"config"`
	}
	type Config struct {
		Storage Module `json:"storage"`
	}
	var config Config
	json.Unmarshal([]byte(configData), &config)
	logger.Log("initializing with config %v", config)
	database := func() db.DB {
		switch config.Storage.Type {
		case "postgres":
			data, err := json.Marshal(config.Storage.Config)
			if err != nil {
				log.Fatalf("failed to serialize ydb config err: %v", err)
			}
			var postgresConfig db.PostgresConfig
			json.Unmarshal(data, &postgresConfig)
			logger.Log("creating postgres with config %v", postgresConfig)
			db, err := db.Postgres(postgresConfig, logger)
			if err != nil {
				log.Fatalf("failed to initialize postgres err: %v", err)
			}
			logger.Log("initialized postgres")
			return db
		default:
			log.Fatalf("unknown storage type %s", config.Storage.Type)
			return nil
		}
	}()
	defer database.Close()

	createDatabaseActions(database, logger).setup()
}
