package utils

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	mysql2 "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"os"
)

var mysqlDB *sql.DB

// ConnectToMySQLGorm Connect to a database handle from a connection string.
func ConnectToMySQLGorm() (*gorm.DB, error) {
	mySQLConfig, _ := NewMySQLConfig("sandbox")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?multiStatements=true", mySQLConfig.Database.User, mySQLConfig.Database.Password, mySQLConfig.Database.Host, mySQLConfig.Database.DBName)
	db, err := gorm.Open(mysql2.Open(dsn))

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

func GetMySqlDB() *sql.DB {
	return mysqlDB
}

func GetMySqlGormDB() *gorm.DB {
	gormMySQL, err := ConnectToMySQLGorm()
	if err != nil {
		fmt.Printf("Error connecting to MySQL: %s\n", err)
	}
	return gormMySQL
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
