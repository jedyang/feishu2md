package main

import (
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

var version = "yunsheng-v3-1"

func main() {
	app := &cli.App{
		Name:    "feishu2md",
		Version: strings.TrimSpace(string(version)),
		Usage:   "Download feishu/larksuite document to markdown file",
		Action: func(ctx *cli.Context) error {
			cli.ShowAppHelp(ctx)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "config",
				Usage: "Read config file or set field(s) if provided",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "appId",
						Value:       "",
						Usage:       "Set app id for the OPEN API",
						Destination: &configOpts.appId,
					},
					&cli.StringFlag{
						Name:        "appSecret",
						Value:       "",
						Usage:       "Set app secret for the OPEN API",
						Destination: &configOpts.appSecret,
					},
					&cli.StringFlag{
						Name:        "ossAccessKeyId",
						Value:       "",
						Usage:       "Set OSS access key id",
						Destination: &configOpts.ossAccessKeyId,
					},
					&cli.StringFlag{
						Name:        "ossAccessKeySecret",
						Value:       "",
						Usage:       "Set OSS access key secret",
						Destination: &configOpts.ossAccessKeySecret,
					},
					&cli.StringFlag{
						Name:        "ossBucketName",
						Value:       "",
						Usage:       "Set OSS bucket name",
						Destination: &configOpts.ossBucketName,
					},
					&cli.StringFlag{
						Name:        "ossEndpoint",
						Value:       "",
						Usage:       "Set OSS endpoint",
						Destination: &configOpts.ossEndpoint,
					},
					&cli.StringFlag{
						Name:        "ossRegion",
						Value:       "",
						Usage:       "Set OSS region",
						Destination: &configOpts.ossRegion,
					},
					&cli.StringFlag{
						Name:        "ossPrefix",
						Value:       "",
						Usage:       "Set OSS prefix",
						Destination: &configOpts.ossPrefix,
					},
				},
				Action: func(ctx *cli.Context) error {
					return handleConfigCommand()
				},
			},
			{
				Name:    "download",
				Aliases: []string{"dl"},
				Usage:   "Download feishu/larksuite document to markdown file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "output",
						Aliases:     []string{"o"},
						Value:       "./",
						Usage:       "Specify the output directory for the markdown files",
						Destination: &dlOpts.outputDir,
					},
					&cli.BoolFlag{
						Name:        "dump",
						Value:       false,
						Usage:       "Dump json response of the OPEN API",
						Destination: &dlOpts.dump,
					},
					&cli.BoolFlag{
						Name:        "batch",
						Value:       false,
						Usage:       "Download all documents under a folder",
						Destination: &dlOpts.batch,
					},
					&cli.BoolFlag{
						Name:        "wiki",
						Value:       false,
						Usage:       "Download all documents within the wiki.",
						Destination: &dlOpts.wiki,
					},
					&cli.BoolFlag{
						Name:        "uploadpic",
						Value:       false,
						Usage:       "Upload images to Alibaba Cloud OSS instead of downloading to local",
						Destination: &dlOpts.uploadPic,
					},
					&cli.StringFlag{
						Name:        "name",
						Value:       "",
						Usage:       "Specify the markdown file name",
						Destination: &dlOpts.name,
					},
				},
				ArgsUsage: "<url>",
				Action: func(ctx *cli.Context) error {
					if ctx.NArg() == 0 {
						return cli.Exit("Please specify the document/folder/wiki url", 1)
					} else {
						url := ctx.Args().First()
						return handleDownloadCommand(url)
					}
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
