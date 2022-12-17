package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"osprey/config"
	"osprey/data/torrents"
	"osprey/utils"
)

var client *http.Client

func InitHTTPClient() {
	client = &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	return nil
}

func UpdateTorrentList() torrents.TorrentList {
	requestBody := []byte(`{
		"jsonrpc": "2.0",
		"id": 0,
		"method": "torrents.list",
		"params": {}
	}`)
	req, err := http.NewRequest("POST", config.Config.JSONRPCEndpointURL, bytes.NewBuffer(requestBody))
	utils.CheckError(err)
	req.Header.Add("Authorization", "Bearer "+config.Config.SecretKey)
	resp, err := client.Do(req)
	utils.CheckError(err)
	body, err := ioutil.ReadAll(resp.Body)
	utils.CheckError(err)
	var torrentListRequestResponse torrents.TorrentListRequestResponse
	json.Unmarshal(body, &torrentListRequestResponse)
	return torrentListRequestResponse.Result
}
