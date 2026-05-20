package main

import (
	"ai-study-admin/model"
	"ai-study-admin/plugins"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	templateService "github.com/collect-ui/collect/src/collect/service_imp"
	collect "github.com/collect-ui/collect/src/collect/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func getContentType(filePath string) string {
	// 获取文件扩展名
	ext := strings.ToLower(filepath.Ext(filePath))

	// 根据扩展名返回对应的 MIME 类型
	switch ext {
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js", ".mjs":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".webp":
		return "image/webp"
	case ".txt":
		return "text/plain"
	case ".pdf":
		return "application/pdf"
	case ".zip":
		return "application/zip"
	case ".gz":
		return "application/gzip"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".ttf":
		return "font/ttf"
	case ".otf":
		return "font/otf"
	default:
		// 默认返回二进制流类型
		return "application/octet-stream"
	}
}
func serveStatic(urlPrefix, root string, cache bool) gin.HandlerFunc {
	fs := http.Dir(root)
	//fileServer := http.FileServer(fs)

	return func(c *gin.Context) {

		p := c.Request.URL.Path

		// 设置Cache-Control头，缓存60天

		if !strings.HasPrefix(p, urlPrefix) {
			c.Next()
			return
		}

		if cache {
			c.Header("Cache-Control", "public, max-age=5184000") // 60天 = 60 * 24 * 60 * 60 秒
			// 设置Expires头，缓存60天
			expires := time.Now().Add(60 * 24 * time.Hour)
			c.Header("Expires", expires.Format(http.TimeFormat))
		} else {
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
			c.Header("Pragma", "no-cache") // 兼容 HTTP/1.0
			c.Header("Expires", "0")       // 立即过期

		}

		// 去掉 URL 前缀
		p = strings.TrimPrefix(p, urlPrefix)

		// 检查客户端是否支持 gzip
		if strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			gzFilePath := p + ".gz"
			if f, err := fs.Open(gzFilePath); err == nil {
				// 返回 .gz 文件
				c.Header("Content-Encoding", "gzip")
				c.Header("Content-Type", getContentType(p))
				info, _ := f.Stat()
				http.ServeContent(c.Writer, c.Request, info.Name(), info.ModTime(), f)
				return
			}
		}

		// 尝试查找文件
		f, err := fs.Open(p)
		if err != nil {
			// 如果文件不存在，尝试查找 index.html
			f, err = fs.Open("index.html")
			if err != nil {
				c.Next()
				return
			}
		}
		defer f.Close()

		info, err := f.Stat()
		if err != nil || info.IsDir() {
			// 如果文件是目录或无法获取文件信息，尝试查找 index.html
			f, err = fs.Open("index.html")
			if err != nil {
				c.Next()
				return
			}
			defer f.Close()
			info, _ = f.Stat()
		}

		// 提供文件服务
		http.ServeContent(c.Writer, c.Request, info.Name(), info.ModTime(), f)
		c.Abort()
	}
}

func isTrustedOrigin(origin string) bool {
	if origin == "" || strings.HasPrefix(origin, "http://localhost") {
		return true
	}
	// 添加你的域名白名单
	allowedDomains := []string{"https://yourdomain.com", "https://www.yourdomain.com"}
	for _, domain := range allowedDomains {
		if origin == domain {
			return true
		}
	}
	return false
}

func DynamicSessionOptions() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 动态安全配置
		isSecure := c.Request.TLS != nil || // 自动检测HTTPS
			c.GetHeader("X-Forwarded-Proto") == "https" // 支持代理场景

		// 2. SameSite策略 (所有HTTPS域名用None，否则Lax)
		sameSite := http.SameSiteNoneMode
		if !isSecure || isLocalhostOrIP(c.Request.Host) {
			sameSite = http.SameSiteLaxMode
		}

		// 3. 动态设置Domain (仅对非IP和非localhost的域名生效)
		domain := ""
		if !isLocalhostOrIP(c.Request.Host) {
			domain = extractRootDomain(c.Request.Host)
		}

		// 4. 应用配置
		session := sessions.Default(c)
		session.Options(sessions.Options{
			Path:     "/",
			Domain:   domain,
			MaxAge:   86400 * 30,
			Secure:   isSecure,
			HttpOnly: true,
			SameSite: sameSite,
		})

		// 5. 显式保存配置
		if err := session.Save(); err != nil {
			fmt.Printf("Session save error: %v", err)
		}

		c.Next()
	}
}

// 辅助函数：判断是否本地或IP访问
func isLocalhostOrIP(host string) bool {
	host = strings.Split(host, ":")[0] // 去除端口
	return host == "localhost" || net.ParseIP(host) != nil
}

// 辅助函数：提取根域名
func extractRootDomain(host string) string {
	parts := strings.Split(host, ":")
	domainParts := strings.Split(parts[0], ".")
	if len(domainParts) >= 2 {
		return "." + strings.Join(domainParts[len(domainParts)-2:], ".") // 如 ".iqiaoqi.com"
	}
	return parts[0]
}

func initRuntimeLogger() *os.File {
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Printf("create log dir failed: %v", err)
		return nil
	}

	file, err := os.OpenFile("logs/server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("open log file failed: %v", err)
		return nil
	}

	writer := io.MultiWriter(os.Stdout, file)
	gin.DefaultWriter = writer
	gin.DefaultErrorWriter = writer
	log.SetOutput(writer)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	return file
}

func normalizeTemplateDataRequest(c *gin.Context) {
	if c.Request.Method != http.MethodPost || c.Request.URL.Path != "/template_data/data" {
		return
	}

	service := c.Query("service")
	if service == "" {
		return
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("read template_data body failed: %v", err)
		c.Request.Body = io.NopCloser(bytes.NewReader(nil))
		return
	}
	_ = c.Request.Body.Close()

	payload := map[string]interface{}{}
	if len(bytes.TrimSpace(bodyBytes)) > 0 {
		if err := json.Unmarshal(bodyBytes, &payload); err != nil {
			log.Printf("parse template_data body failed: %v", err)
			c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			return
		}
	}

	if _, ok := payload["service"]; !ok {
		payload["service"] = service
	}

	nextBody, err := json.Marshal(payload)
	if err != nil {
		log.Printf("encode template_data body failed: %v", err)
		c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewReader(nextBody))
	c.Request.ContentLength = int64(len(nextBody))
	c.Request.Header.Set("Content-Type", "application/json")
}

func main() {
	logFile := initRuntimeLogger()
	if logFile != nil {
		defer logFile.Close()
	}
	// todo go profile 使用
	//gin.SetMode(gin.ReleaseMode)
	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	// 全局设置跨域头
	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Vary", "Origin")
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
	r.Use(func(c *gin.Context) {
		normalizeTemplateDataRequest(c)
		start := time.Now()
		c.Next()
		if strings.HasPrefix(c.Request.URL.Path, "/template_data/data") {
			log.Printf(
				"lowcode service=%s method=%s status=%d latency=%s client=%s",
				c.Query("service"),
				c.Request.Method,
				c.Writer.Status(),
				time.Since(start),
				c.ClientIP(),
			)
		}
	})
	// 生成cookies
	store := cookie.NewStore([]byte("secret"))

	r.Use(sessions.Sessions("session_id", store))
	r.Use(DynamicSessionOptions()) // 添加动态选项中间件
	r.Static("/static", "./static")

	dirStr := collect.GetAppKey("dirList")
	dirList := strings.Split(dirStr, ";")
	for _, file := range dirList {
		if collect.IsValueEmpty(file) {
			continue
		}
		fileInfo := strings.Split(file, ",")
		if len(fileInfo) < 3 {
			log.Printf("skip invalid dirList item: %s", file)
			continue
		}
		cache := false
		if fileInfo[2] == "true" {
			cache = true
		}
		r.Use(serveStatic(fileInfo[0], fileInfo[1], cache))
	}

	fileStr := collect.GetAppKey("fileList")
	fileList := strings.Split(fileStr, ";")
	for _, file := range fileList {
		if collect.IsValueEmpty(file) {
			continue
		}
		fileInfo := strings.Split(file, ",")
		if len(fileInfo) < 2 {
			log.Printf("skip invalid fileList item: %s", file)
			continue
		}
		r.StaticFile(fileInfo[0], fileInfo[1])
	}
	// 设置数据库
	templateService.SetDatabaseModel(&model.TableData{})
	// 设置外部处理器
	templateService.SetOuterModuleRegister(plugins.GetRegisterList())
	// 添加定时任务
	templateService.RunScheduleService()
	// 添加启动服务
	templateService.RunStartupService()
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/collect-ui")
	})
	r.POST("/template_data/data", func(c *gin.Context) {
		templateService.HandlerRequest(c)
	})

	r.GET("/template_data/ws/:token", func(context *gin.Context) {
		templateService.HandlerWsRequest(context)
	})

	serverPort := collect.GetAppKey("server_port")
	isHttps := collect.GetAppKey("is_https")
	domains := strings.Split(collect.GetAppKey("domain"), ",")
	log.Printf("AI Study admin backend starting on port %s, https=%s", serverPort, isHttps)

	if isHttps == "true" {
		tlsConfig := &tls.Config{}
		for _, domain := range domains {
			certFile := "/etc/letsencrypt/live/" + domain + "/fullchain.pem"
			keyFile := "/etc/letsencrypt/live/" + domain + "/privkey.pem"
			cert, err := tls.LoadX509KeyPair(certFile, keyFile)
			if err != nil {
				panic(err)
			}
			tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
		}
		server := &http.Server{
			Addr:      ":" + serverPort,
			Handler:   r,
			TLSConfig: tlsConfig,
		}
		if err := server.ListenAndServeTLS("", ""); err != nil {
			panic(err)
		}
	} else {
		r.Run(":" + serverPort) // listen and serve on 0.0.0.0:8080
	}

}
