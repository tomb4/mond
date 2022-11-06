package client

import (
    bytes2 "bytes"
    "context"
    "encoding/json"
    "fmt"
    "go.uber.org/zap"
    "io"
    "io/ioutil"
    "mond/wind/config"
    "mond/wind/env"
    merr "mond/wind/err"
    "mond/wind/logger"
    "mond/wind/trace"
    "mond/wind/utils/constant"
    "mond/wind/utils/endpoint"
    "net"
    "net/http"
    "net/url"
    "sync"
    "time"
)

var (
    lock      sync.Mutex
    clientMap map[string]*Client
)

func init() {
    clientMap = make(map[string]*Client)
}

type Client struct {
    c    *http.Client
    addr string
    _do  func(ctx context.Context, req *http.Request) (*http.Response, error)
    _log logger.Logger
}

func GetHttpClient(appId string) (*Client, error) {
    lock.Lock()
    defer lock.Unlock()
    if env.GetAppState() != env.Starting {
        panic("GetHttpClient只能在ResourceInit中调用")
    }
    if clientMap[appId] != nil {
        return clientMap[appId], nil
    }
    config, err := config.GetHttpClientOption(appId)
    if err != nil {
        return nil, err
    }
    if config.MaxConn == 0 {
        config.MaxConn = 100
    }
    if config.Timeout == 0 {
        config.Timeout = 5000
    }
    c := Client{
        addr: config.Addr,
        _log: logger.GetLogger(),
    }
    c.initDo()
    transport := &http.Transport{
        DialContext: (&net.Dialer{
            Timeout:   time.Second,      //连接超时时间
            KeepAlive: time.Second * 30, //连接保持超时时间
        }).DialContext,
        MaxConnsPerHost:       config.MaxConn,
        MaxIdleConns:          config.MaxConn, //client对与所有host最大空闲连接数总和
        MaxIdleConnsPerHost:   config.MaxConn,
        IdleConnTimeout:       120 * time.Second, //空闲连接在连接池中的超时时间
        TLSHandshakeTimeout:   3 * time.Second,   //TLS安全连接握手超时时间
        ExpectContinueTimeout: 1 * time.Second,   //发送完请求到接收到响应头的超时时间
    }
    c.c = &http.Client{
        Transport: transport,
        Timeout:   time.Millisecond * time.Duration(config.Timeout),
    }
    clientMap[appId] = &c
    return &c, nil
}

type Method string

const (
    GET    Method = "GET"
    POST   Method = "POST"
    DELETE Method = "DELETE"
    PUT    Method = "PUT"
)

type Option func(req *http.Request)

func WithHeader(key, value string) Option {
    return func(req *http.Request) {
        req.Header.Add(key, value)
    }
}
func (m *Client) Do(ctx context.Context, path string, method Method, req, resp interface{}, opts ...Option) error {
    var input io.Reader
    bytes, err := json.Marshal(req)
    if err != nil {
        return err
    }
    uri := fmt.Sprintf("%s%s", m.addr, path)
    if method == GET {
        reqMap := make(map[string]interface{})
        err = json.Unmarshal(bytes, &reqMap)
        if err != nil {
            return err
        }
        params := url.Values{}
        for k, v := range reqMap {
            params.Add(k, fmt.Sprintf("%v", v))
        }
        uri = fmt.Sprintf("%s?%s", uri, params.Encode())
    } else {
        input = bytes2.NewBuffer(bytes)
    }
    _req, err := http.NewRequest(string(method), uri, input)
    _req.Header.Set("Content-Type", "application/json")
    for _, opt := range opts {
        opt(_req)
    }
    if err != nil {
        return err
    }
    _resp, err := m._do(ctx, _req)
    if err != nil {
        return err
    }
    if _resp.StatusCode >= 400 {
        return merr.HttpStatusErr.SetMsg(_resp.Status)
    }
    respBytes, err := ioutil.ReadAll(_resp.Body)
    if err != nil {
        return err
    }

    err = json.Unmarshal(respBytes, resp)
    if err != nil {
        return err
    }
    return nil
}

func (m *Client) initDo() {
    ep := trace.HttpClientMiddleware(makeEndpoint(m.do))
    m._do = m.makeDo(ep)
}

func (m *Client) makeDo(endpoint endpoint.Endpoint) func(ctx context.Context, req *http.Request) (*http.Response, error) {
    return func(ctx context.Context, req *http.Request) (*http.Response, error) {
        ctx = context.WithValue(ctx, constant.HttpClientHost, m.addr)
        ctx = context.WithValue(ctx, constant.HttpClientPath, req.URL.Path)
        ctx = context.WithValue(ctx, constant.HttpClientMethod, req.Method)
        resp, err := endpoint(ctx, req)
        if err != nil {
            m._log.Error(ctx, "http请求失败", zap.Error(err))
            return nil, err
        }
        r := resp.(*http.Response)
        if r.StatusCode >= 400 {
            m._log.Error(ctx, "http请求失败 status code", zap.Int("status_code", r.StatusCode))
        }
        return r, nil
    }
}

func makeEndpoint(f func(ctx context.Context, req *http.Request) (*http.Response, error)) endpoint.Endpoint {
    return func(ctx context.Context, req interface{}) (interface{}, error) {
        resp, err := f(ctx, req.(*http.Request))
        return resp, err
    }
}

func (m *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
    select {
    case <-ctx.Done():
        return nil, merr.SysErrTimeoutErr
    default:
    }
    resp, err := m.c.Do(req)
    if err != nil {
        return nil, err
    }
    return resp, nil
}
