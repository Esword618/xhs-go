package xhs

import (
	"errors"
	"fmt"
	"github.com/Esword618/xhs-go/consts"
	"github.com/Esword618/xhs-go/utils"
	"github.com/bitly/go-simplejson"
	"github.com/imroc/req/v3"
	"github.com/playwright-community/playwright-go"
	"net/http"
	"net/url"
	"strings"
)

type IXhsClient interface {
	Initialize(cookieStr string, headless ...bool) error
	SendCloseSignal()
	WaitForCloseSignal()
	Close()
	GetNoteById(noteId string) (*simplejson.Json, error)
	GetUserInfo(userId string) (*simplejson.Json, error)
	GetUserNotes(userId string, cursor ...string) (*simplejson.Json, error)
	GetNoteByKeyword(keyword string, args ...interface{}) (*simplejson.Json, error)
	GetHomeFeed(feedType consts.FeedType) (*simplejson.Json, error)
}

var xhsclient = xhsClient{}

func NewXhs() IXhsClient {
	return &xhsclient
}

type xhsClient struct {
	client         *req.Client
	PW             *playwright.Playwright
	Browser        playwright.Browser
	BrowserContext playwright.BrowserContext
	ContextPage    playwright.Page
	CloseCh        chan struct{}
	cookieMap      map[string]string
}

// Initialize 初始化
func (x *xhsClient) Initialize(cookieStr string, headless ...bool) error {
	headless1 := true
	if len(headless) > 0 {
		headless1 = headless[0]
	}
	runOption := &playwright.RunOptions{
		SkipInstallBrowsers: true,
	}
	err := playwright.Install(runOption)
	if err != nil {
		return fmt.Errorf("could not install playwright dependencies: %v", err)
	}

	// 初始化一个pw
	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("failed to start Playwright: %v", err)
	}
	// browser
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Channel:  playwright.String("chrome"),
		Headless: playwright.Bool(headless1),
	})
	if err != nil {
		return fmt.Errorf("failed to launch Chromium: %v", err)
	}

	// browserContext
	browserContext, err := browser.NewContext(playwright.BrowserNewContextOptions{
		UserAgent: playwright.String("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"),
		Viewport: &playwright.ViewportSize{
			Width:  1520,
			Height: 1080,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create browser context: %v", err)
	}
	err = browserContext.AddInitScript(playwright.BrowserContextAddInitScriptOptions{
		Script: playwright.String(consts.StealthMinJs),
	})
	if err != nil {
		return fmt.Errorf("执行js报错: %v", err)
	}
	var cookiesP []playwright.OptionalCookie
	var cookies []*http.Cookie
	cookiemap := utils.ConvertStrCookieToDict(cookieStr)
	for key, value := range cookiemap {
		cookiesP = append(cookiesP, playwright.OptionalCookie{
			Name:   playwright.String(key),
			Value:  playwright.String(value),
			Domain: playwright.String(".xiaohongshu.com"),
			Path:   playwright.String("/"),
		})
		cookies = append(cookies, &http.Cookie{
			Name:   key,
			Value:  value,
			Domain: ".xiaohongshu.com",
			Path:   "/",
		})
	}

	err = browserContext.AddCookies(cookiesP...)
	if err != nil {
		return fmt.Errorf("add cookie 失败: %v", err)
	}

	// contextPage
	contextPage, err := browserContext.NewPage()
	if err != nil {
		return fmt.Errorf("failed to create page: %v", err)
	}

	// 执行页面导航操作
	_, err = contextPage.Goto("https://www.xiaohongshu.com")
	if err != nil {
		return fmt.Errorf("Failed to navigate to URL:", err)
	}

	// client
	contextPage.WaitForTimeout(3000)
	client := req.C()
	client.SetTLSFingerprintChrome()
	client.SetCommonHeaders(map[string]string{
		"user-agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
		"content-type": "application/json;charset=UTF-8",
	})
	client.SetBaseURL("https://edith.xiaohongshu.com")
	client.SetCommonCookies(cookies...)
	client.OnBeforeRequest(func(client *req.Client, req *req.Request) error {

		var rawURL string
		var data string
		switch req.Method {
		case http.MethodPost:
			rawURL = req.RawURL
			data = string(req.Body)
		case http.MethodGet:
			queryParams := url.Values{}
			for k, v := range req.QueryParams {
				queryParams.Set(k, v[0])
			}

			queryString := queryParams.Encode()
			rawURL = fmt.Sprintf("%s?%s", req.RawURL, queryString)
			data = ""
		default:
			return errors.New("未知请求")
		}
		encryptParams, err := x.sign(rawURL, data, req.Method)
		if err != nil {
			fmt.Printf("before request,%v\n", err)
		}
		for key, value := range encryptParams {
			req.Headers.Add(key, value)
		}
		return nil
	})
	x.PW = pw
	x.Browser = browser
	x.BrowserContext = browserContext
	x.ContextPage = contextPage
	x.CloseCh = make(chan struct{})
	x.client = client
	x.cookieMap = cookiemap
	return nil
}

// SendCloseSignal 发送关闭信号
func (x *xhsClient) SendCloseSignal() {
	x.CloseCh <- struct{}{}
}

// WaitForCloseSignal 等待关闭信号
func (x *xhsClient) WaitForCloseSignal() {
	<-x.CloseCh
}

// Close 关闭
func (x *xhsClient) Close() {
	if x.CloseCh != nil {
		close(x.CloseCh)
	}
	if x.ContextPage != nil {
		x.ContextPage.Close()
	}
	if x.BrowserContext != nil {
		x.BrowserContext.Close()
	}
	if x.Browser != nil {
		x.Browser.Close()
	}
	if x.PW != nil {
		x.PW.Stop()
	}
}

// sign值获取
func (x *xhsClient) sign(rawUrl, data, method string) (map[string]string, error) {
	var jsStr string
	switch method {
	case http.MethodPost:
		jsStr = fmt.Sprintf(`
		const url = "%s";
		const data = JSON.parse('%s');
		window._webmsxyw(url, data);
		`, rawUrl, data)
	case http.MethodGet:
		jsStr = fmt.Sprintf(`
		const url = "%s";
		const data = "%s";
		window._webmsxyw(url);
		`, rawUrl, data)
	}
	encryptParams, err := x.ContextPage.Evaluate(jsStr)
	convertedParams := make(map[string]string)
	if err != nil {
		fmt.Println(err)
		return convertedParams, err
	}
	encryptParamsMap, ok := encryptParams.(map[string]interface{})
	if !ok {
		// 处理类型不匹配的情况
		return convertedParams, fmt.Errorf("encryptParams is not a map[string]interface{}")
	}

	for key, value := range encryptParamsMap {
		strValue := fmt.Sprintf("%v", value)
		convertedParams[strings.ToLower(key)] = strValue
	}
	return convertedParams, nil
}
