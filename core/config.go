package core

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
)

type Config struct {
	Feishu FeishuConfig `json:"feishu"`
	Output OutputConfig `json:"output"`
	OSS    OSSConfig    `json:"oss"`
}

type FeishuConfig struct {
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type OutputConfig struct {
	ImageDir        string `json:"image_dir"`
	TitleAsFilename bool   `json:"title_as_filename"`
	UseHTMLTags     bool   `json:"use_html_tags"`
	SkipImgDownload bool   `json:"skip_img_download"`
}

type OSSConfig struct {
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	BucketName      string `json:"bucket_name"`
	Endpoint        string `json:"endpoint"`
	Region          string `json:"region"`
	Prefix          string `json:"prefix"`
}

func NewConfig(appId, appSecret string) *Config {
	return &Config{
		Feishu: FeishuConfig{
			AppId:     appId,
			AppSecret: appSecret,
		},
		Output: OutputConfig{
			ImageDir:        "static",
			TitleAsFilename: false,
			UseHTMLTags:     false,
			SkipImgDownload: false,
		},
		OSS: OSSConfig{
			AccessKeyId:     "",
			AccessKeySecret: "",
			BucketName:      "",
			Endpoint:        "",
			Region:          "",
			Prefix:          "",
		},
	}
}

func GetConfigFilePath() (string, error) {
	configPath, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	configFilePath := path.Join(configPath, "feishu2md", "config.json")
	return configFilePath, nil
}

func ReadConfigFromFile(configPath string) (*Config, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	config := NewConfig("", "")
	err = json.Unmarshal([]byte(file), &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (conf *Config) WriteConfig2File(configPath string) error {
	err := os.MkdirAll(filepath.Dir(configPath), 0o755)
	if err != nil {
		return err
	}
	file, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(configPath, file, 0o644)
	return err
}
