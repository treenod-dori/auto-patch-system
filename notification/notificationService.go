package notification

import (
	"fmt"
	"github.com/levigross/grequests"
	"github.com/spf13/viper"
	"os"
)

type Service interface {
	SendOKNotification() error
	SendFailNotification() error
}

type NotiConfig struct {
	Notification *NotificaionCase `yaml:"notification"`
}

type NotificaionCase struct {
	Ok   string `yaml:"ok, omitempty"`
	Fail string `yaml:"fail, omitempty"`
}

type SlackNotificationService struct {
}

func NewSlackNotificationService() Service {
	return &SlackNotificationService{}
}

func (s SlackNotificationService) SendOKNotification() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		os.Exit(1)
	}

	// 환경에 맞는 설정 가져오기
	config := &NotiConfig{}
	if err := viper.UnmarshalKey("slack", config); err != nil {
		return err
	}

	// HTTP POST 요청
	_, err := grequests.Post(config.Notification.Ok, nil)
	if err != nil {
		//resp, err := grequests.Post(config.Notification.Fail, nil)
		return err
	}

	return nil
}

func (s SlackNotificationService) SendFailNotification() error {
	//TODO implement me
	panic("implement me")
}
