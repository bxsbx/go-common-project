package http

import (
	"StandardProject/common/errorz"
	"StandardProject/common/logz"
	"StandardProject/common/tracer"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	RespType int
	Client   *http.Client
}

type response1 struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const (
	// 结构类型
	CUSTOM_RESPONSE = 0
	RESPONSE1       = 1
	RESPONSE2       = 2

	// 请求体头Content——Type
	FORM = 1
	JSON = 2

	// 请求方法
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"

	//请求头
	ACCEPT       = "Accept"
	CONTENT_TYPE = "Content-Type"

	BODY_FORM = "application/x-www-form-urlencoded"
	BODY_JSON = "application/json"
)

func selectStructToResponse(bytes []byte, respType int, resultData interface{}) error {
	decoder := json.NewDecoder(strings.NewReader(string(bytes)))
	decoder.UseNumber()
	var data interface{}
	switch respType {
	case RESPONSE1:
		var resp response1
		err := decoder.Decode(&resp)
		if err != nil {
			return errorz.CodeError(errorz.ERR_DECODE, err)
		}
		if resp.Code != 10000 {
			return errorz.CodeMsg(resp.Code, resp.Msg)
		}
		data = resp.Data
	case RESPONSE2:
	default:
		err := decoder.Decode(resultData)
		if err != nil {
			return errorz.CodeError(errorz.ERR_DECODE, err)
		}
		return nil
	}
	if resultData == nil || data == nil {
		return nil
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		return errorz.CodeError(errorz.ERR_MARSHAL, err)
	}
	err = json.Unmarshal(marshal, resultData)
	if err != nil {
		return errorz.CodeError(errorz.ERR_UNMARSHAL, err)
	}
	return nil
}

func recordeLog(path, method string, header, query, body interface{}, err error) {
	data := make(map[string]interface{})
	data["url"] = path
	data["method"] = method
	data["header"] = header
	data["query"] = query
	data["body"] = body
	if err != nil {
		data["err"] = err
		logz.Error(data)
	} else {
		logz.Info(data)
	}
}

func (c *Client) Request(path string, method string, header map[string]string, query map[string]string, body io.Reader, resultData interface{}, appCtx context.Context) error {
	if len(query) > 0 {
		params := url.Values{}
		for k, v := range query {
			params.Set(k, v)
		}
		path += "?" + params.Encode()
	}

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return errorz.CodeError(errorz.NEW_REQUEST, err)
	}

	if len(header) > 0 {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("Accept", "application/json")

	// 链路跟踪记录请求头信息
	tracer.InjectTracerSpan(appCtx, req.Header)

	resp, err := c.Client.Do(req)
	if err != nil {
		return errorz.CodeError(errorz.REQUEST_ERR, err)
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return errorz.CodeError(errorz.IO_READ_ERR, err)
	}

	if resp.StatusCode != http.StatusOK {
		return errorz.CodeMsg(resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	err = selectStructToResponse(bytes, c.RespType, resultData)
	if err != nil {
		return err
	}
	recordeLog(path, method, header, query, body, err)
	return nil
}

func (c *Client) Get(path string, query map[string]string, resultData interface{}, appCtx context.Context) error {
	return c.Request(path, GET, nil, query, nil, resultData, appCtx)
}

func (c *Client) Put(path string, query map[string]string, body io.Reader, resultData interface{}, appCtx context.Context) error {
	return c.Request(path, PUT, nil, query, body, resultData, appCtx)
}

func (c *Client) Delete(path string, query map[string]string, resultData interface{}, appCtx context.Context) error {
	return c.Request(path, DELETE, nil, query, nil, resultData, appCtx)
}

func (c *Client) Post(path string, query map[string]string, body io.Reader, resultData interface{}, appCtx context.Context) error {
	return c.Request(path, POST, nil, query, body, resultData, appCtx)
}

func (c *Client) PostQuery(path string, query map[string]string, resultData interface{}, appCtx context.Context) error {
	return c.Request(path, POST, nil, query, nil, resultData, appCtx)
}

func (c *Client) PostForm(path string, form map[string]string, resultData interface{}, appCtx context.Context) error {
	header := make(map[string]string)
	header[CONTENT_TYPE] = BODY_FORM
	bodyForm := url.Values{}
	for k, v := range form {
		bodyForm.Set(k, v)
	}
	return c.Request(path, GET, header, nil, strings.NewReader(bodyForm.Encode()), resultData, appCtx)
}

func (c *Client) PostJson(path string, body interface{}, resultData interface{}, appCtx context.Context) error {
	header := make(map[string]string)
	header[CONTENT_TYPE] = BODY_JSON
	marshal, _ := json.Marshal(body)
	return c.Request(path, GET, header, nil, strings.NewReader(string(marshal)), resultData, appCtx)
}
