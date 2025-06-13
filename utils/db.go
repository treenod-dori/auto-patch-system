package utils

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strings"
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

// ConnectToMySQL - MySQL 데이터베이스 연결
func ConnectToMySQL2(dbConfig *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbConfig.Database.User, dbConfig.Database.Password, dbConfig.Database.Host, dbConfig.Database.DBName)
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

func InitMySQL2(mysqlConfig *Config) error {
	var err error

	// SQLite 연결
	if mysqlDB == nil {
		mysqlDB, err = ConnectToMySQL2(mysqlConfig)
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

func SplitSQLQueries(sql string) ([]string, error) {
	var queries []string
	var currentQuery strings.Builder
	inString := false
	inSingleLineComment := false
	inMultiLineComment := false

	for i, r := range sql {
		// 문자열 리터럴 시작/끝
		if r == '\'' && !inMultiLineComment && !inSingleLineComment {
			inString = !inString
		}

		// 싱글 라인 주석
		if r == '-' && i+1 < len(sql) && sql[i+1] == ' ' && !inString && !inMultiLineComment {
			inSingleLineComment = true
		}

		// 멀티 라인 주석 시작
		if r == '/' && i+1 < len(sql) && sql[i+1] == '*' && !inString && !inSingleLineComment {
			inMultiLineComment = true
		}

		// 멀티 라인 주석 끝
		if r == '*' && i+1 < len(sql) && sql[i+1] == '/' && inMultiLineComment {
			inMultiLineComment = false
			i++ // Skip the '/' character
		}

		// 싱글 라인 주석 끝
		if r == '\n' && inSingleLineComment {
			inSingleLineComment = false
		}

		// 세미콜론이 쿼리의 끝일 때
		if r == ';' && !inString && !inSingleLineComment && !inMultiLineComment {
			queries = append(queries, currentQuery.String())
			currentQuery.Reset() // 쿼리 버퍼 초기화
		} else {
			currentQuery.WriteRune(r) // 쿼리 내용 추가
		}
	}

	// 마지막 쿼리 추가 (세미콜론 없는 경우)
	if currentQuery.Len() > 0 {
		queries = append(queries, currentQuery.String())
	}

	return queries, nil
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
