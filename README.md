# feishu2md
[![Golang - feishu2md](https://img.shields.io/github/go-mod/go-version/jedyang/feishu2md?color=%2376e1fe&logo=go)](https://go.dev/) [![Release](https://img.shields.io/github/v/release/jedyang/feishu2md?color=orange&logo=github)](https://github.com/jedyang/feishu2md/releases)

这是一个下载飞书文档为 Markdown 文件的工具，使用 Go 语言实现。

**这是我基于另一个开源项目进行的二次开发，原项目地址：[Wsine/feishu2md: 一键命令下载飞书文档为 Markdown（寻找维护者）](https://github.com/Wsine/feishu2md)**



## 增加的功能

主要是增加了我自己需要的功能

1. 将飞书文档中的图片直接上传到阿里云oss，并在md文件中直接使用图片url

2. 通过参数指定md文件名，原先是一串随机字母

   

## 如何使用

### 1. 获取 API Token

配置文件需要填写 APP ID 和 APP SECRET 信息，请参考 [飞书官方文档](https://open.feishu.cn/document/ukTMukTMukTM/ukDNz4SO0MjL5QzM/get-) 获取。推荐设置为

- 进入飞书[开发者后台](https://open.feishu.cn/app)
- 创建企业自建应用（个人版），信息随意填写
- （重要）打开权限管理，开通以下必要的权限（可点击以下链接参考 API 调试台->权限配置字段）
  - [获取文档基本信息](https://open.feishu.cn/document/server-docs/docs/docs/docx-v1/document/get)，「查看新版文档」权限 `docx:document:readonly`
  - [获取文档所有块](https://open.feishu.cn/document/server-docs/docs/docs/docx-v1/document/list)，「查看新版文档」权限 `docx:document:readonly`
  - [下载素材](https://open.feishu.cn/document/server-docs/docs/drive-v1/media/download)，「下载云文档中的图片和附件」权限 `docs:document.media:download`
  - [获取文件夹中的文件清单](https://open.feishu.cn/document/server-docs/docs/drive-v1/folder/list)，「查看、评论、编辑和管理云空间中所有文件」权限 `drive:file:readonly`
  - [获取知识空间节点信息](https://open.feishu.cn/document/server-docs/docs/wiki-v2/space-node/get_node)，「查看知识库」权限 `wiki:wiki:readonly`
- 打开凭证与基础信息，获取 App ID 和 App Secret

### 2. 配置

我这个版本启用新版本，避免与原来的feishu2md混淆

```
$ ./feishu2md.exe -h
NAME:
   feishu2md - Download feishu/larksuite document to markdown file

USAGE:
   feishu2md [global options] command [command options] [arguments...]

VERSION:
   yunsheng-v3-1

```



先配置飞书的appId和appSecret。

如果要传到阿里云oss，同时要配置oss相关参数

```
$ ./feishu2md.exe config -h
NAME:
   feishu2md config - Read config file or set field(s) if provided

USAGE:
   feishu2md config [command options] [arguments...]

OPTIONS:
   --appId value               Set app id for the OPEN API
   --appSecret value           Set app secret for the OPEN API
   --ossAccessKeyId value      Set OSS access key id
   --ossAccessKeySecret value  Set OSS access key secret
   --ossBucketName value       Set OSS bucket name
   --ossEndpoint value         Set OSS endpoint
   --ossRegion value           Set OSS region
   --ossPrefix value           Set OSS prefix
   --help, -h                  show help (default: false)

```

   通过 `feishu2md config` 命令可以查看配置文件路径以及是否成功配置。

### 3. 使用

<details>
  <summary>命令行版本</summary>

  借助 Go 语言跨平台的特性，已编译好了主要平台的可执行文件，可以在 Release中下载，并将相应平台的 feishu2md 可执行文件放置在 PATH 路径中即可。

   **查阅帮助文档**

   ```bash
$ ./feishu2md.exe dl -h
NAME:
   feishu2md download - Download feishu/larksuite document to markdown file

USAGE:
   feishu2md download [command options] <url>

OPTIONS:
   --output value, -o value  Specify the output directory for the markdown files (default: "./")
   --dump                    Dump json response of the OPEN API (default: false)
   --batch                   Download all documents under a folder (default: false)
   --wiki                    Download all documents within the wiki. (default: false)
   --uploadpic               Upload images to Alibaba Cloud OSS instead of downloading to local (default: false)
   --name value              Specify the markdown file name
   --help, -h                show help (default: false)

   ```



   **下载单个文档为 Markdown**

   通过 `feishu2md dl <your feishu docx url>` 直接下载，文档链接可以通过 **分享 > 开启链接分享 > 互联网上获得链接的人可阅读 > 复制链接** 获得。

   示例：

   ```bash
   $ feishu2md dl "https://domain.feishu.cn/docx/docxtoken"
   ```

  **批量下载某文件夹内的全部文档为 Markdown**

  此功能暂时不支持Docker版本

  通过`feishu2md dl --batch <your feishu folder url>` 直接下载，文件夹链接可以通过 **分享 > 开启链接分享 > 互联网上获得链接的人可阅读 > 复制链接** 获得。

  示例：

  ```bash
  $ feishu2md dl --batch -o output_directory "https://domain.feishu.cn/drive/folder/foldertoken"
  ```

  **批量下载某知识库的全部文档为 Markdown**

  通过`feishu2md dl --wiki <your feishu wiki setting url>` 直接下载，wiki settings链接可以通过 打开知识库设置获得。

  示例：

  ```bash
  $ feishu2md dl --wiki -o output_directory "https://domain.feishu.cn/wiki/settings/123456789101112"
  ```

</details>

## 感谢

- [feishu2md](https://github.com/Wsine/feishu2md)

- [chyroc/lark](https://github.com/chyroc/lark)

- [chyroc/lark_docs_md](https://github.com/chyroc/lark_docs_md)

  
