package favicon

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sun-panel/internal/global"
	"sun-panel/internal/infra/storage"
	"sun-panel/internal/util"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func IsHTTPURL(url string) bool {
	httpPattern := `^(http://|https://|//)`
	match, err := regexp.MatchString(httpPattern, url)
	if err != nil {
		return false
	}
	return match
}

func GetOneFaviconURL(urlStr string) (string, error) {
	iconURLs, err := getFaviconURL(urlStr)
	if err != nil {
		return "", err
	}

	for _, v := range iconURLs {
		// 标准的路径地址
		if IsHTTPURL(v) {
			return v, nil
		} else {
			urlInfo, _ := url.Parse(urlStr)
			fullUrl := urlInfo.Scheme + "://" + urlInfo.Host + "/" + strings.TrimPrefix(v, "/")
			return fullUrl, nil
		}
	}
	return "", fmt.Errorf("not found ico")
}

// 获取远程文件的大小
func GetRemoteFileSize(url string) (int64, error) {
	// 创建一个自定义的 HTTP 客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
		// 启用自动重定向
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	// 添加浏览器模拟头信息
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			global.Logger.Errorf("failed to close resp.Body. error : %v", err)
		}
	}()

	// 检查HTTP响应状态
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("HTTP request failed, status code: %d", resp.StatusCode)
	}

	// 获取Content-Length字段，即文件大小
	size := resp.ContentLength

	// 如果服务器没有提供 Content-Length，尝试读取响应体来确定大小
	if size <= 0 {
		global.Logger.Infof("Content-Length not provided for %s, using alternative method", url)
		// 不实际读取整个响应体，因为我们已经有了连接，可以直接关闭
		return 0, fmt.Errorf("Content-Length not provided by server")
	}

	return size, nil
}

// 下载图片
func DownloadImage(ctx context.Context, url string, storage storage.RcloneStorage) (string, error) {
	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	
	// 发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			global.Logger.Errorf("failed to close resp.Body. error : %v", err)
		}
	}()

	// 检查HTTP响应状态
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed, status code: %d", response.StatusCode)
	}

	// 限制最大下载大小为 10MB
	limitedReader := http.MaxBytesReader(nil, response.Body, 10*1024*1024)

	// 生成文件名
	urlFileName := path.Base(url)
	fileExt := path.Ext(url)
	if fileExt == "" {
		fileExt = ".ico"
	}
	fileName := util.Md5(fmt.Sprintf("%s%s", urlFileName, time.Now().String())) + fileExt

	// 上传文件
	filepath, err := storage.Upload(ctx, limitedReader, fileName)
	if err != nil {
		if strings.Contains(err.Error(), "request body too large") {
			return "", fmt.Errorf("文件太大，不下载")
		}
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	return filepath, nil
}

func getFaviconURL(url string) ([]string, error) {
	var icons []string
	icons = make([]string, 0)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return icons, err
	}

	// 设置User-Agent头字段，模拟浏览器请求
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	resp, err := client.Do(req)
	if err != nil {
		return icons, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			global.Logger.Errorf("failed to close resp.Body. error : %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return icons, errors.New("HTTP request failed with status code " + strconv.Itoa(resp.StatusCode))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return icons, err
	}

	// 查找所有link标签，筛选包含rel属性为"icon"的标签
	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		rel, _ := s.Attr("rel")
		href, _ := s.Attr("href")

		if strings.Contains(rel, "icon") && href != "" {
			// fmt.Println(href)
			icons = append(icons, href)
		}
	})

	if len(icons) == 0 {
		return icons, errors.New("favicon not found on the page")
	}

	return icons, nil
}
