package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ZabbixAPI struct {
	authToken  string
	endpoint   string
	requiestId int
	Verbose    bool

	userLogin bool
}

type ZbxAPIRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Id      int         `json:"id"`
	Params  interface{} `json:"params"`
}

type ZbxAPIResponse struct {
	RPCVersion string      `json:"jsonrpc"`
	Result     any         `json:"result,omitempty"`
	Id         int         `json:"id"`
	Error      ZbxAPIError `json:"error,omitempty"`
}

type ZbxAPIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type ZbxLoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ZbxGetParams struct {
	Filter ZbxFilterParams `json:"filter,omitempty"`
	Output []string        `json:"output,omitempty"`
	Host   string          `json:"host,omitempty"`
	HostID string          `json:"hostids,omitempty"`
}

type ZbxFilterParams struct {
	Host []string `json:"host,omitempty"`
	Key  []string `json:"key_,omitempty"`
}

type ZbxHost struct {
	Hostid string `json:"hostid"`
}

type ZbxItem struct {
	ItemID    string   `json:"itemid,omitempty"`
	Name      string   `json:"name"`
	Key       string   `json:"key_"`
	Type      int      `json:"type"`
	ValueType int      `json:"value_type"`
	HostID    string   `json:"hostid,omitempty"`
	Tags      []ZbxTag `json:"tags,omitempty"`
}

type ZbxTag struct {
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

type ZbxItemIdList struct {
	ItemIDs []string `json:"itemids,omitempty"`
}

/*
0 - Zabbix agent;
2 - Zabbix trapper;
7 - Zabbix agent (active);
*/
const ZBX_ITEM_AGENT int = 0
const ZBX_ITEM_TRAPPER int = 2
const ZBX_ITEM_AGENT_ACTIVE int = 7

/*
0 - numeric float;
3 - numeric unsigned;
4 - text.
*/
const ZBX_ITEM_TYPE_FLOAT int = 0
const ZBX_ITEM_TYPE_UNSIGNED_INT int = 3
const ZBX_ITEM_TYPE_TEXT int = 4

func NewZabbixAPI(endpoint, authToken string) *ZabbixAPI {
	return &ZabbixAPI{authToken: authToken, endpoint: endpoint}
}

func (s *ZabbixAPI) Call(method string, params, response any) error {
	requestJSON := ZbxAPIRequest{Jsonrpc: "2.0", Method: method, Id: s.requiestId, Params: params}
	requestBody, err := json.Marshal(requestJSON)
	if err != nil {
		fmt.Println("Request serialization error:", err.Error())
		os.Exit(1)
	}

	if s.Verbose {
		fmt.Println("Request:", string(requestBody))
	}

	request, err := http.NewRequest("POST", s.endpoint, bytes.NewReader(requestBody))
	if err != nil {
		panic(err.Error())
	}

	request.Header.Set("Content-Type", "application/json")

	if s.authToken != "" {
		request.Header.Set("Authorization", "Bearer "+s.authToken)
	}

	client := &http.Client{}
	res, err := client.Do(request)
	s.requiestId++
	if err != nil {
		panic(err.Error())
	}

	resultBody, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	if response != nil {
		data := ZbxAPIResponse{Result: response}

		err = json.Unmarshal(resultBody, &data)
		if err != nil {
			panic(err.Error())
		}

		if data.Error.Code != 0 {
			if s.Verbose {
				fmt.Println(data.Error)
			}
			return errors.New(data.Error.Message + "\n" + data.Error.Data)
		}
	}

	return nil
}

func NewZabbixAPIUsingCreds(endpoint, login, pass string) (*ZabbixAPI, error) {
	loginParams := ZbxLoginParams{login, pass}

	api := ZabbixAPI{endpoint: endpoint}

	err := api.Call("user.login", loginParams, &api.authToken)
	if err != nil {
		return nil, err
	}

	return &api, nil
}

func (s *ZabbixAPI) Logout() {
	s.Call("user.logout", nil, nil)
}

func (s *ZabbixAPI) GetHostID(host string) (string, error) {
	hostID := ""

	filterParams := ZbxFilterParams{Host: []string{host}}
	hostGetParams := ZbxGetParams{Filter: filterParams, Output: []string{"hostid"}}

	var hostList []ZbxHost

	err := s.Call("host.get", hostGetParams, &hostList)
	if err != nil {
		return "", err
	}

	if len(hostList) != 0 {
		hostID = hostList[0].Hostid
	}

	return hostID, nil
}

func (s *ZabbixAPI) GetItemByKey(host, key string) (ZbxItem, error) {
	zbxItem := ZbxItem{}

	filterParams := ZbxFilterParams{Key: []string{key}}
	itemGetParams := ZbxGetParams{Filter: filterParams, Output: []string{"itemid"}, Host: host}

	var itemList []ZbxItem
	err := s.Call("item.get", itemGetParams, &itemList)
	if err != nil {
		return zbxItem, err
	} else if len(itemList) > 0 {
		zbxItem = itemList[0]
	}

	return zbxItem, nil
}

func (s *ZabbixAPI) CreateItems(items []ZbxItem) ([]string, error) {
	var itemIDs []string

	var res ZbxItemIdList
	err := s.Call("item.create", items, &res)
	if err != nil {
		return itemIDs, err
	} else {
		itemIDs = res.ItemIDs
	}

	return itemIDs, nil
}

func (s *ZabbixAPI) GetItemsWithKeys(hostID string, keys []string) ([]ZbxItem, error) {

	filterParams := ZbxFilterParams{Key: keys}
	itemGetParams := ZbxGetParams{Filter: filterParams, Output: []string{"itemid", "key_"}, HostID: hostID}

	var itemList []ZbxItem
	err := s.Call("item.get", itemGetParams, &itemList)
	if err != nil {
		return []ZbxItem{}, err
	} else {
		return itemList, nil
	}
}
