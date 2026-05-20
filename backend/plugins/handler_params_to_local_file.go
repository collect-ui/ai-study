package plugins

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	common "github.com/collect-ui/collect/src/collect/common"
	config "github.com/collect-ui/collect/src/collect/config"
	collect "github.com/collect-ui/collect/src/collect/filters"
	templateService "github.com/collect-ui/collect/src/collect/service_imp"
	utils "github.com/collect-ui/collect/src/collect/utils"
)

type ToLocalFile struct {
	templateService.BaseHandler
}

func getTargetPath(field string, params map[string]interface{}) (string, error) {
	tpl, err := config.CastTemplate(field)
	if err != nil {
		return "", err
	}
	targetPath := utils.RenderTplData(tpl, params).(string)
	if targetPath == "" || targetPath == "~" {
		targetPath = "./"
	}
	return targetPath, nil
}

func shortenFileName(filename string, maxLength int) string {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	timestamp := time.Now().Format("15-04-05")
	suffix := "_" + timestamp

	if containsChinese(name) {
		return name + suffix + ext
	}

	maxNameLength := maxLength - len(suffix) - len(ext)
	if len(name) > maxNameLength {
		name = name[:maxNameLength]
	}

	return name + suffix + ext
}

func containsChinese(s string) bool {
	for _, r := range s {
		if utf8.RuneLen(r) > 1 {
			return true
		}
	}
	return false
}

func (si *ToLocalFile) HandlerData(template *config.Template, handlerParam *config.HandlerParam, ts *templateService.TemplateService) *common.Result {
	params := template.GetParams()
	file := ts.File

	if file == nil {
		return common.NotOk("上传文件不能为空")
	}

	targetDir, err := getTargetPath(handlerParam.Field, params)
	if err != nil {
		return common.NotOk(err.Error())
	}
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		err = os.MkdirAll(targetDir, os.ModePerm)
		if err != nil {
			return common.NotOk(fmt.Sprintf("创建文件夹失败: %s", err.Error()))
		}
	}
	shortenedFileName := shortenFileName(ts.FileHeader.Filename, 50)
	realFilePath := targetDir + shortenedFileName
	httpFilePath := calculateHttpPath(realFilePath, collect.GetKey("file_prefix"), collect.GetKey("local_file_dir"))

	out, err := os.Create(realFilePath)
	if err != nil {
		return common.NotOk(fmt.Sprintf("创建文件失败: %s", err.Error()))
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return common.NotOk(fmt.Sprintf("写入文件失败: %s", err.Error()))
	}
	fileInfo, err := os.Stat(realFilePath)
	if err != nil {
		return common.NotOk(fmt.Sprintf("获取文件大小失败: %s", err.Error()))
	}

	fileData := make(map[string]interface{})
	fileData["path"] = httpFilePath
	fileData["real_path"] = realFilePath
	fileData["size"] = formatFileSize(fileInfo.Size())
	fileData["filename"] = shortenedFileName
	fileData["filetype"] = getFileType(shortenedFileName)
	return common.Ok(fileData, "处理参数成功")
}

func getFileType(shortenedFileName string) string {
	lastDotIndex := strings.LastIndex(shortenedFileName, ".")
	if lastDotIndex == -1 {
		return "unknown"
	}

	fileExtension := shortenedFileName[lastDotIndex+1:]
	switch strings.ToLower(fileExtension) {
	case "jpg", "jpeg", "png", "gif", "svg", "bmp", "tiff", "tif", "webp", "ico", "heic", "heif":
		return "image"
	case "mp4", "avi", "mov", "mkv", "flv", "wmv", "webm":
		return "video"
	case "mp3", "wav", "ogg", "flac", "aac", "m4a":
		return "audio"
	case "txt", "md", "rtf", "log", "csv", "tsv":
		return "text"
	case "pdf":
		return "pdf"
	case "doc", "docx":
		return "word"
	case "xls", "xlsx":
		return "excel"
	case "ppt", "pptx":
		return "powerpoint"
	default:
		return fileExtension
	}
}

func formatFileSize(size int64) string {
	const (
		KB = 1 << 10
		MB = 1 << 20
		GB = 1 << 30
	)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}

func calculateHttpPath(realFilePath, filePrefix, localFileDir string) string {
	if strings.HasPrefix(realFilePath, localFileDir) {
		realFilePath = realFilePath[len(localFileDir):]
	}
	realFilePath = filepath.ToSlash(realFilePath)

	httpPath := filePrefix + realFilePath
	if !strings.HasPrefix(httpPath, "/") {
		httpPath = "/" + httpPath
	}
	return httpPath
}
