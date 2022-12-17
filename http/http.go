package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"osprey/config"
	"osprey/data/torrents"
	"osprey/utils"
	"strings"
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

func DeleteTorrent(torrent torrents.Torrent, keepData bool) {
	removeDataJSON, err := json.Marshal(!keepData)
	utils.CheckError(err)
	requestBody := []byte(`{
		"jsonrpc": "2.0",
		"method": "torrents.remove",
		"params": {
			"info_hashes":[` + getMarshalledInfoHash(torrent) + `],
			"remove_data":` + string(removeDataJSON) + `
		}
	}`)
	req, err := http.NewRequest("POST", config.Config.JSONRPCEndpointURL, bytes.NewBuffer(requestBody))
	utils.CheckError(err)
	req.Header.Add("Authorization", "Bearer "+config.Config.SecretKey)
	_, err = client.Do(req)
	utils.CheckError(err)
}

func PauseResumeTorrent(torrent torrents.Torrent) {
	method := "torrents.pause"
	if torrents.IsPaused(torrent.Flags) {
		method = "torrents.resume"
	}
	requestBody := []byte(`{
		"jsonrpc": "2.0",
		"method": "` + method + `",
		"params": {
			"info_hash":` + getMarshalledInfoHash(torrent) + `
		}
	}`)
	req, err := http.NewRequest("POST", config.Config.JSONRPCEndpointURL, bytes.NewBuffer(requestBody))
	utils.CheckError(err)
	req.Header.Add("Authorization", "Bearer "+config.Config.SecretKey)
	_, err = client.Do(req)
	utils.CheckError(err)
}

func getMarshalledInfoHash(torrent torrents.Torrent) string {
	marshalledInfoHash, err := json.Marshal(torrent.InfoHash)
	utils.CheckError(err)
	marshalledInfoHashWithNull := strings.Replace(string(marshalledInfoHash), "\"\"", "null", -1)
	utils.CheckError(err)
	return marshalledInfoHashWithNull
}
