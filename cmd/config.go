package main

import (
	"fmt"
	"os"

	"github.com/Wsine/feishu2md/core"
	"github.com/Wsine/feishu2md/utils"
)

type ConfigOpts struct {
	appId              string
	appSecret          string
	ossAccessKeyId     string
	ossAccessKeySecret string
	ossBucketName      string
	ossEndpoint        string
	ossRegion          string
	ossPrefix          string
}

var configOpts = ConfigOpts{}

func handleConfigCommand() error {
	configPath, err := core.GetConfigFilePath()
	if err != nil {
		return err
	}

	fmt.Println("Configuration file on: " + configPath)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := core.NewConfig(configOpts.appId, configOpts.appSecret)
		config.OSS.AccessKeyId = configOpts.ossAccessKeyId
		config.OSS.AccessKeySecret = configOpts.ossAccessKeySecret
		config.OSS.BucketName = configOpts.ossBucketName
		config.OSS.Endpoint = configOpts.ossEndpoint
		config.OSS.Region = configOpts.ossRegion
		config.OSS.Prefix = configOpts.ossPrefix
		if err = config.WriteConfig2File(configPath); err != nil {
			return err
		}
		fmt.Println(utils.PrettyPrint(config))
	} else {
		config, err := core.ReadConfigFromFile(configPath)
		if err != nil {
			return err
		}
		if configOpts.appId != "" {
			config.Feishu.AppId = configOpts.appId
		}
		if configOpts.appSecret != "" {
			config.Feishu.AppSecret = configOpts.appSecret
		}
		if configOpts.ossAccessKeyId != "" {
			config.OSS.AccessKeyId = configOpts.ossAccessKeyId
		}
		if configOpts.ossAccessKeySecret != "" {
			config.OSS.AccessKeySecret = configOpts.ossAccessKeySecret
		}
		if configOpts.ossBucketName != "" {
			config.OSS.BucketName = configOpts.ossBucketName
		}
		if configOpts.ossEndpoint != "" {
			config.OSS.Endpoint = configOpts.ossEndpoint
		}
		if configOpts.ossRegion != "" {
			config.OSS.Region = configOpts.ossRegion
		}
		if configOpts.ossPrefix != "" {
			config.OSS.Prefix = configOpts.ossPrefix
		}
		if configOpts.appId != "" || configOpts.appSecret != "" || configOpts.ossAccessKeyId != "" || configOpts.ossAccessKeySecret != "" || configOpts.ossBucketName != "" || configOpts.ossEndpoint != "" || configOpts.ossRegion != "" || configOpts.ossPrefix != "" {
			if err = config.WriteConfig2File(configPath); err != nil {
				return err
			}
		}
		fmt.Println(utils.PrettyPrint(config))
	}
	return nil
}
