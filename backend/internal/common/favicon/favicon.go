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
	"sun-panel/internal/common"
	"sun-panel/internal/global"
	"sun-panel/internal/infra/storage"
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
	resp, err := http.Head(url)
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
	return size, nil
}

// 下载图片
func DownloadImage(ctx context.Context, url string, storage storage.RcloneStorage) (string, error) {
	// 获取远程文件大小
	fileSize, err := GetRemoteFileSize(url)
	if err != nil {
		return "", err
	}

	// 判断文件大小是否在阈值内（设置为 10MB）
	maxSize := int64(10 * 1024 * 1024)
	if fileSize > maxSize {
		return "", fmt.Errorf("文件太大，不下载。大小：%d字节", fileSize)
	}

	// 发送HTTP GET请求获取图片数据
	response, err := http.Get(url)
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

	urlFileName := path.Base(url)
	fileExt := path.Ext(url)
	if fileExt == "" {
		fileExt = ".ico"
	}
	fileName := common.Md5(fmt.Sprintf("%s%s", urlFileName, time.Now().String())) + fileExt

	// 使用 rclone 存储接口上传文件
	filepath, err := storage.Upload(ctx, response.Body, fileName)
	if err != nil {
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
