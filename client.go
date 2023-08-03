package yandex_direct_sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mg-realcom/yandex-direct-sdk/common"
	"github.com/mg-realcom/yandex-direct-sdk/statistics"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Client struct {
	Tr              *http.Client
	Login           string
	Token           *string
	App             *App
	host            environment
	statisticsLimit statisticsLimits
}

type App struct {
	ID     string
	Secret string
}

type statisticsLimits struct {
	retryInterval  int32
	reportsInQueue int8
}

type environment string

const (
	LIVE    environment = "api.direct.yandex.com"
	SANDBOX environment = "api-sandbox.direct.yandex.com"
)

func NewClient(tr *http.Client, login string, token *string, app *App, sandbox bool) *Client {
	if sandbox {
		return &Client{
			Login: login,
			Token: token,
			host:  SANDBOX,
		}
	}

	return &Client{
		Tr:    tr,
		Login: login,
		Token: token,
		App:   app,
		host:  LIVE,
		statisticsLimit: statisticsLimits{
			retryInterval:  0,
			reportsInQueue: 0,
		},
	}
}

func (c *Client) buildHeader(req *http.Request) {
	req.Header.Add("Authorization", "Bearer "+*c.Token)
	req.Header.Add("Client-Login", c.Login)
	req.Header.Add("Accept-Language", "ru")
	req.Header.Add("skipReportHeader", "true")
	req.Header.Add("skipReportSummary", "true")
}

type Payload struct {
	Method string `json:"method"`
	Params struct {
		Ads []struct {
			AdGroupID int `json:"AdGroupId"`
			TextAd    struct {
				Text   string `json:"Text"`
				Title  string `json:"Title"`
				Href   string `json:"Href,omitempty"`
				Mobile string `json:"Mobile"`
			} `json:"TextAd"`
		} `json:"Ads"`
	} `json:"params"`
}

func (c *Client) GetReport(titleRequest, dir string, typeReport statistics.ReportType, fields []string, filter []statistics.Filter, dateRange statistics.DateRange, dtRangeType statistics.DateRangeType) (string, error) {
	t := time.Now().Format("2006-01-02")
	reportName := fmt.Sprintf("%s_%s_%s_%s_%s_%s", c.Login, typeReport, t, titleRequest, dateRange.From, dateRange.To)
	params := statistics.ReportDefinition{
		Selection: &statistics.SelectionCriteria{
			DateFrom: dateRange.From,
			DateTo:   dateRange.To,
			Filter:   filter,
		},
		FieldNames:    fields,
		Page:          &common.Page{Limit: 10000000, Offset: 0},
		ReportName:    reportName,
		ReportType:    typeReport,
		DateRangeType: dtRangeType,
		Format:        common.FormatTSV,
		IncludeVAT:    common.YES,
	}

	ctx := context.Background()
	for {
		req, err := c.createGetReportRequest(ctx, params)
		if err != nil {
			return "", fmt.Errorf("createGetReportRequest: %w", err)
		}

		c.waitInfo(reportName)
		time.Sleep(time.Duration(c.statisticsLimit.retryInterval) * time.Second)

		resp, err := c.Tr.Do(req)
		if err != nil {
			return "", fmt.Errorf("do request: %w", err)
		}

		switch resp.StatusCode {
		case http.StatusOK:
			file, err := createTSVFile(dir, reportName, resp)
			if err != nil {
				return "", fmt.Errorf("createTSVFile: %w", err)
			}

			return file, nil
		case http.StatusCreated, http.StatusAccepted:
			err := c.waitInit(resp)
			if err != nil {
				return "", fmt.Errorf("waitInit: %w", err)
			}
		case http.StatusInternalServerError:
			return "", fmt.Errorf("internal server error")
		case http.StatusBadRequest:
			fmt.Println(resp.Status)

			data, err := c.badRequestPrepare(resp)
			if err != nil {
				return "", fmt.Errorf("cannot prepare bad request: %w", err)
			}

			return "", fmt.Errorf("ошибка отчета %s", data.Error.ErrorDetail)
		default:
			return "", fmt.Errorf("cтатус код сервера при получении отчета %v", resp.StatusCode)
		}
	}
}

type Request struct {
	Params statistics.ReportDefinition `json:"params"`
}

type Response struct {
	Error struct {
		ErrorDetail string `json:"error_detail"`
		RequestID   string `json:"request_id"`
		ErrorCode   string `json:"error_code"`
		ErrorString string `json:"error_string"`
	} `json:"error"`
}

func (c *Client) createGetReportRequest(ctx context.Context, params statistics.ReportDefinition) (*http.Request, error) {
	reqContent := Request{Params: params}

	body, err := json.Marshal(reqContent)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.direct.yandex.com/json/v5/reports", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	c.buildHeader(req)

	return req, nil
}

func (c *Client) badRequestPrepare(resp *http.Response) (Response, error) {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("cant read response body: %w", err)
	}

	var data Response

	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return Response{}, fmt.Errorf("cant unmarshal response body: %w", err)
	}

	return data, nil
}

func (c *Client) waitInit(resp *http.Response) error {
	if resp == nil {
		return fmt.Errorf("response is nil")
	}

	retryIn, err := strconv.Atoi(resp.Header.Get("retryIn"))
	if err != nil {
		return fmt.Errorf("retryIn: %w", err)
	}

	c.statisticsLimit.retryInterval = int32(retryIn)

	reportsInQueue, err := strconv.Atoi(resp.Header.Get("reportsInQueue"))
	if err != nil {
		return fmt.Errorf("reportsInQueue: %w", err)
	}

	c.statisticsLimit.reportsInQueue = int8(reportsInQueue)

	return nil
}

func (c *Client) waitInfo(reportName string) {
	if c.statisticsLimit.retryInterval > 1 {
		fmt.Printf("Повтор запроса на отчет %s через %v\n", reportName, c.statisticsLimit.retryInterval)
	}

	if c.statisticsLimit.reportsInQueue > 1 {
		fmt.Printf("Количество отчетов в очереди %v\n", c.statisticsLimit.reportsInQueue)
	}
}

func createTSVFile(dir string, filename string, resp *http.Response) (string, error) {
	if resp == nil {
		return "", fmt.Errorf("response is nil")
	}

	f, err := os.CreateTemp(dir, fmt.Sprintf("%s_*.tsv", filename))
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}

	defer f.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	err = os.WriteFile(f.Name(), responseBody, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write to file: %w", err)
	}

	return f.Name(), nil
}
