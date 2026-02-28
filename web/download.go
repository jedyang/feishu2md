package main

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/88250/lute"
	"github.com/gin-gonic/gin"
	"github.com/jedyang/feishu2md/core"
	"github.com/jedyang/feishu2md/utils"
)

func downloadHandler(c *gin.Context) {
	// get parameters
	feishu_docx_url, err := url.QueryUnescape(c.Query("url"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid encoded feishu/larksuite URL")
		return
	}

	// Get additional parameters
	fileName := c.Query("fileName")
	// URL decode fileName to handle Chinese characters
	if fileName != "" {
		decodedFileName, err := url.QueryUnescape(fileName)
		if err == nil {
			fileName = decodedFileName
		}
	}
	enableCloudStorage := c.Query("enableCloudStorage") == "true"

	// Validate the url to download
	docType, docToken, err := utils.ValidateDocumentURL(feishu_docx_url)
	fmt.Println("Captured document token:", docToken)

	// Create client with context
	ctx := context.Background()
	config := core.NewConfig(
		os.Getenv("FEISHU_APP_ID"),
		os.Getenv("FEISHU_APP_SECRET"),
	)
	// Load OSS config from environment variables
	config.OSS.AccessKeyId = os.Getenv("OSS_ACCESS_KEY_ID")
	config.OSS.AccessKeySecret = os.Getenv("OSS_ACCESS_KEY_SECRET")
	config.OSS.BucketName = os.Getenv("OSS_BUCKET_NAME")
	config.OSS.Endpoint = os.Getenv("OSS_ENDPOINT")
	config.OSS.Region = os.Getenv("OSS_REGION")
	config.OSS.Prefix = os.Getenv("OSS_PREFIX")

	client := core.NewClient(
		config.Feishu.AppId, config.Feishu.AppSecret,
	)

	// Process the download
	parser := core.NewParser(config.Output)
	markdown := ""

	// for a wiki page, we need to renew docType and docToken first
	if docType == "wiki" {
		node, err := client.GetWikiNodeInfo(ctx, docToken)
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal error: client.GetWikiNodeInfo: "+err.Error())
			log.Printf("error: %s", err)
			return
		}
		docType = node.ObjType
		docToken = node.ObjToken
	}
	if docType == "docs" {
		c.String(http.StatusBadRequest, "Unsupported docs document type")
		return
	}

	docx, blocks, err := client.GetDocxContent(ctx, docToken)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal error: client.GetDocxContent: "+err.Error())
		log.Printf("error: %s", err)
		return
	}
	markdown = parser.ParseDocxContent(docx, blocks)

	zipBuffer := new(bytes.Buffer)
	writer := zip.NewWriter(zipBuffer)
	for _, imgToken := range parser.ImgTokens {
		if enableCloudStorage {
			// Upload to OSS
			ossURL, err := client.UploadImageToOSS(ctx, imgToken, config.OSS)
			if err != nil {
				c.String(http.StatusInternalServerError, "Internal error: client.UploadImageToOSS: "+err.Error())
				log.Printf("error: %s", err)
				return
			}
			markdown = strings.Replace(markdown, imgToken, ossURL, 1)
		} else {
			// Download to local
			localLink, rawImage, err := client.DownloadImageRaw(ctx, imgToken, config.Output.ImageDir)
			if err != nil {
				c.String(http.StatusInternalServerError, "Internal error: client.DownloadImageRaw: "+err.Error())
				log.Printf("error: %s", err)
				return
			}
			markdown = strings.Replace(markdown, imgToken, localLink, 1)
			f, err := writer.Create(localLink)
			if err != nil {
				c.String(http.StatusInternalServerError, "Internal error: zipWriter.Create")
				log.Printf("error: %s", err)
				return
			}
			_, err = f.Write(rawImage)
			if err != nil {
				c.String(http.StatusInternalServerError, "Internal error: zipWriter.Create.Write")
				log.Printf("error: %s", err)
				return
			}
		}
	}

	engine := lute.New(func(l *lute.Lute) {
		l.RenderOptions.AutoSpace = true
	})
	result := engine.FormatStr("md", markdown)

	// Set response
	if len(parser.ImgTokens) > 0 {
		// Determine markdown filename
		mdName := fmt.Sprintf("%s.md", docToken)
		if fileName != "" {
			mdName = fmt.Sprintf("%s.md", fileName)
		}

		f, err := writer.Create(mdName)
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal error: zipWriter.Create")
			log.Printf("error: %s", err)
			return
		}
		_, err = f.Write([]byte(result))
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal error: zipWriter.Create.Write")
			log.Printf("error: %s", err)
			return
		}

		err = writer.Close()
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal error: zipWriter.Close")
			log.Printf("error: %s", err)
			return
		}

		// Determine zip filename
		zipName := docToken
		if fileName != "" {
			zipName = fileName
		}
		c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.zip"`, zipName))
		c.Data(http.StatusOK, "application/octet-stream", zipBuffer.Bytes())
	} else {
		// Determine markdown filename
		mdName := docToken
		if fileName != "" {
			mdName = fileName
		}
		c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.md"`, mdName))
		c.Data(http.StatusOK, "application/octet-stream", []byte(result))
	}
}
