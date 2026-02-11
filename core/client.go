package core

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/chyroc/lark"
	"github.com/chyroc/lark_rate_limiter"
)

type Client struct {
	larkClient *lark.Lark
}

func NewClient(appID, appSecret string) *Client {
	return &Client{
		larkClient: lark.New(
			lark.WithAppCredential(appID, appSecret),
			lark.WithTimeout(60*time.Second),
			lark.WithApiMiddleware(lark_rate_limiter.Wait(4, 4)),
		),
	}
}

func (c *Client) DownloadImage(ctx context.Context, imgToken, outDir string) (string, error) {
	resp, _, err := c.larkClient.Drive.DownloadDriveMedia(ctx, &lark.DownloadDriveMediaReq{
		FileToken: imgToken,
	})
	if err != nil {
		return imgToken, err
	}
	fileext := filepath.Ext(resp.Filename)
	filename := fmt.Sprintf("%s/%s%s", outDir, imgToken, fileext)
	err = os.MkdirAll(filepath.Dir(filename), 0o755)
	if err != nil {
		return imgToken, err
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		return imgToken, err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.File)
	if err != nil {
		return imgToken, err
	}
	return filename, nil
}

func (c *Client) DownloadImageRaw(ctx context.Context, imgToken, imgDir string) (string, []byte, error) {
	resp, _, err := c.larkClient.Drive.DownloadDriveMedia(ctx, &lark.DownloadDriveMediaReq{
		FileToken: imgToken,
	})
	if err != nil {
		return imgToken, nil, err
	}
	fileext := filepath.Ext(resp.Filename)
	filename := fmt.Sprintf("%s/%s%s", imgDir, imgToken, fileext)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.File)
	return filename, buf.Bytes(), nil
}

func (c *Client) UploadImageToOSS(ctx context.Context, imgToken string, ossConfig OSSConfig) (string, error) {
	// Download image from Feishu
	_, imgData, err := c.DownloadImageRaw(ctx, imgToken, "")
	if err != nil {
		return imgToken, err
	}

	// Generate object key
	// Create OSS client configuration
	provider := credentials.NewStaticCredentialsProvider(ossConfig.AccessKeyId, ossConfig.AccessKeySecret, "")
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(provider).
		WithRegion(ossConfig.Region)

	// Create OSS client
	client := oss.NewClient(cfg)

	fileext := ".png" // Default extension
	if len(imgData) > 0 {
		// Simple file type detection
		if len(imgData) >= 4 && imgData[0] == 0xFF && imgData[1] == 0xD8 && imgData[2] == 0xFF {
			fileext = ".jpg"
		} else if len(imgData) >= 8 && imgData[0] == 0x89 && imgData[1] == 0x50 && imgData[2] == 0x4E && imgData[3] == 0x47 {
			fileext = ".png"
		} else if len(imgData) >= 2 && imgData[0] == 0x47 && imgData[1] == 0x49 {
			fileext = ".gif"
		} else if len(imgData) >= 4 && imgData[0] == 0x52 && imgData[1] == 0x49 && imgData[2] == 0x46 && imgData[3] == 0x46 {
			fileext = ".webp"
		}
	}

	objectKey := fmt.Sprintf("%s%s%s", ossConfig.Prefix, imgToken, fileext)

	// Upload to OSS
	result, err := client.PutObject(ctx, &oss.PutObjectRequest{
		Bucket: oss.Ptr(ossConfig.BucketName),
		Key:    oss.Ptr(objectKey),
		Body:   bytes.NewReader(imgData),
	})

	if err != nil {
		return imgToken, err
	}

	fmt.Printf("Upload image to OSS successfully, ETag: %v\n", result.ETag)

	// Generate OSS URL
	ossURL := fmt.Sprintf("https://%s.%s/%s", ossConfig.BucketName, ossConfig.Endpoint, objectKey)

	return ossURL, nil
}

func (c *Client) GetDocxContent(ctx context.Context, docToken string) (*lark.DocxDocument, []*lark.DocxBlock, error) {
	resp, _, err := c.larkClient.Drive.GetDocxDocument(ctx, &lark.GetDocxDocumentReq{
		DocumentID: docToken,
	})
	if err != nil {
		return nil, nil, err
	}
	docx := &lark.DocxDocument{
		DocumentID: resp.Document.DocumentID,
		RevisionID: resp.Document.RevisionID,
		Title:      resp.Document.Title,
	}
	var blocks []*lark.DocxBlock
	var pageToken *string
	for {
		resp2, _, err := c.larkClient.Drive.GetDocxBlockListOfDocument(ctx, &lark.GetDocxBlockListOfDocumentReq{
			DocumentID: docx.DocumentID,
			PageToken:  pageToken,
		})
		if err != nil {
			return docx, nil, err
		}
		blocks = append(blocks, resp2.Items...)
		pageToken = &resp2.PageToken
		if !resp2.HasMore {
			break
		}
	}
	return docx, blocks, nil
}

func (c *Client) GetWikiNodeInfo(ctx context.Context, token string) (*lark.GetWikiNodeRespNode, error) {
	resp, _, err := c.larkClient.Drive.GetWikiNode(ctx, &lark.GetWikiNodeReq{
		Token: token,
	})
	if err != nil {
		return nil, err
	}
	return resp.Node, nil
}

func (c *Client) GetDriveFolderFileList(ctx context.Context, pageToken *string, folderToken *string) ([]*lark.GetDriveFileListRespFile, error) {
	resp, _, err := c.larkClient.Drive.GetDriveFileList(ctx, &lark.GetDriveFileListReq{
		PageSize:    nil,
		PageToken:   pageToken,
		FolderToken: folderToken,
	})
	if err != nil {
		return nil, err
	}
	files := resp.Files
	for resp.HasMore {
		resp, _, err = c.larkClient.Drive.GetDriveFileList(ctx, &lark.GetDriveFileListReq{
			PageSize:    nil,
			PageToken:   &resp.NextPageToken,
			FolderToken: folderToken,
		})
		if err != nil {
			return nil, err
		}
		files = append(files, resp.Files...)
	}
	return files, nil
}

func (c *Client) GetWikiName(ctx context.Context, spaceID string) (string, error) {
	resp, _, err := c.larkClient.Drive.GetWikiSpace(ctx, &lark.GetWikiSpaceReq{
		SpaceID: spaceID,
	})

	if err != nil {
		return "", err
	}

	return resp.Space.Name, nil
}

func (c *Client) GetWikiNodeList(ctx context.Context, spaceID string, parentNodeToken *string) ([]*lark.GetWikiNodeListRespItem, error) {
	resp, _, err := c.larkClient.Drive.GetWikiNodeList(ctx, &lark.GetWikiNodeListReq{
		SpaceID:         spaceID,
		PageSize:        nil,
		PageToken:       nil,
		ParentNodeToken: parentNodeToken,
	})

	if err != nil {
		return nil, err
	}

	nodes := resp.Items
	previousPageToken := ""

	for resp.HasMore && previousPageToken != resp.PageToken {
		previousPageToken = resp.PageToken
		resp, _, err := c.larkClient.Drive.GetWikiNodeList(ctx, &lark.GetWikiNodeListReq{
			SpaceID:         spaceID,
			PageSize:        nil,
			PageToken:       &resp.PageToken,
			ParentNodeToken: parentNodeToken,
		})

		if err != nil {
			return nil, err
		}

		nodes = append(nodes, resp.Items...)
	}

	return nodes, nil
}
