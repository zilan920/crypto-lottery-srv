package consumer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	NumOfWorkers = 12 // 协程数量
	//AttemptsPerKey = 1000 // 每个私钥尝试次数
	EtherscanAPI = "https://api.etherscan.io/api"
)

type AccountResult []struct {
	Account string `json:"account"`
	Balance string `json:"balance"`
}

func GetBalance(addressList []string, apiKey string) (AccountResult, error) {
	url := fmt.Sprintf("%s?module=account&action=balancemulti&address=%s&tag=latest&apikey=%s", EtherscanAPI, strings.Join(addressList[:], ","), apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析API响应
	response := struct {
		Status  string        `json:"status"`
		Message string        `json:"message"`
		Result  AccountResult `json:"result"`
	}{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("查询地址余额失败：%s, resp: %s", err.Error(), response)
	}

	if response.Status != "1" || response.Message != "OK" {
		return nil, fmt.Errorf("查询地址余额失败：%s, %s", response.Message, string(body))
	}

	if err != nil {
		return nil, err
	}

	return response.Result, nil
}
