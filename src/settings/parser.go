package settings

import (
	"errors"
	"gopkg.in/yaml.v2"
	_ "gopkg.in/yaml.v2"
	"io/ioutil"
)

var s *settings

type settings struct {
	NotificationsChannelURL string `yaml:"notifications_channel_url"`
}

func NotificationsChannelUrl() string {
	return s.NotificationsChannelURL
}

func Init(path string) error {
	if s != nil {
		return errors.New("settings already initialized")
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	data := &settings{}
	err = yaml.Unmarshal(b, data)
	if err != nil {
		return err
	}
	s = data

	return nil
}
