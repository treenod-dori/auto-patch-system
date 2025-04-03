package utils

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"os"
)

var sqliteDB *gorm.DB
var mysqlDB *sql.DB

// ConnectToSQLite Connect to a database handle from a connection string.
func ConnectToSQLite(configuration *Config) (*gorm.DB, error) {
	dbName := configuration.Database.DBName
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}

// ConnectToMySQL - MySQL 데이터베이스 연결
func ConnectToMySQL(dbConfig *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?multiStatements=true", dbConfig.Database.User, dbConfig.Database.Password, dbConfig.Database.Host, dbConfig.Database.DBName)
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitSQLite(sqliteConfig *Config) error {
	var err error

	// SQLite 연결
	if sqliteDB == nil {
		sqliteDB, err = ConnectToSQLite(sqliteConfig)
		if err != nil {
			return fmt.Errorf("failed to connect to SQLite: %w", err)
		}
	}

	return nil
}

func InitMySQL(mysqlConfig *Config) error {
	var err error

	// SQLite 연결
	if mysqlDB == nil {
		mysqlDB, err = ConnectToMySQL(mysqlConfig)
		if err != nil {
			return fmt.Errorf("failed to connect to SQLite: %w", err)
		}
	}

	return nil
}

func GetSqliteDB() *gorm.DB {
	return sqliteDB
}

func GetMySqlDB() *sql.DB {
	return mysqlDB
}

type Config struct {
	Database *Database `yaml:"database"` // 환경별로 데이터베이스 설정을 관리
}

type Database struct {
	Host     string `yaml:"host,omitempty"`
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password,omitempty"`
	DBName   string `yaml:"dbName,omitempty"`
}

func NewMySQLConfig(environment string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 설정 파일 읽기
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		os.Exit(1)
	}

	// 환경에 맞는 설정 가져오기
	db := &Config{}
	if err := viper.UnmarshalKey(environment, db); err != nil {
		return nil, fmt.Errorf("error unmarshalling config for environment %s: %s", environment, err)
	}

	return db, nil
}

func NewSQLiteConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 설정 파일 읽기
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		os.Exit(1)
	}

	// 환경에 맞는 설정 가져오기
	config := &Config{}
	if err := viper.UnmarshalKey("sqlite", config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config for environment %s: %s", "sqlite", err)
	}

	return config, nil
}

// 슬랙 메시지를 전송하는 함수
func SendSlackMessage(webhookURL, channel, username, message string) error {
	// 슬랙 페이로드 생성
	payload := `{
		"channel": "` + channel + `",
		"username": "` + username + `",
		"text": "` + message + `"
	}`

	// HTTP POST 요청
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 응답 코드 확인
	if resp.StatusCode != http.StatusOK {
		return errors.New("슬랙 Webhook 전송 실패: 상태 코드 " + resp.Status)
	}

	return nil
}
