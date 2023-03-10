package http

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"osprey/config"
	"osprey/data/torrents"
	"osprey/utils"
	"strconv"
	"strings"
	"time"
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

func UpdateTorrentList(page int) (torrents.TorrentList, int) {
	requestBody := []byte(`{
		"jsonrpc": "2.0",
		"method": "torrents.list",
		"params": {
			"page": ` + strconv.Itoa(page) + `,
  			"page_size": ` + strconv.Itoa(config.Config.PageSize) + `
		}
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
	updatedPageIndex := page
	if len(torrentListRequestResponse.Result.Torrents) == 0 && page != 0 {
		if page > 0 {
			updatedPageIndex--
		}
	}
	if torrentListRequestResponse.Error.Code == -2 {
		if page > 0 {
			time.Sleep(time.Millisecond * 5) // Prevent the server from interrupting the connection by adding a 25ms timeout
			return UpdateTorrentList(page - 1)
		}
	}
	return torrentListRequestResponse.Result, updatedPageIndex
}

func AddTorrent(magnetURI, savePath string, addingMagnetLink bool) {
	requestBody := []byte{}
	if addingMagnetLink {
		requestBody = []byte(`{
			"jsonrpc": "2.0",
			"method": "torrents.add",
			"params": {
				"magnet_uri":"` + magnetURI + `",
				"save_path":"` + savePath + `",
				"metadata": {
					"source": "osprey"
				}
			}
		}`)
	} else {
		f, err := os.Open(magnetURI)
		utils.CheckError(err)

		// Read entire JPG into byte slice.
		reader := bufio.NewReader(f)
		content, err := ioutil.ReadAll(reader)
		utils.CheckError(err)

		// Encode as base64.
		encoded := base64.StdEncoding.EncodeToString(content)
		requestBody = []byte(`{
			"jsonrpc": "2.0",
			"method": "torrents.add",
			"params": {
				"ti":"` + encoded + `",
				"save_path":"` + savePath + `",
				"metadata": {
					"source": "osprey"
				}
			}
		}`)
	}
	req, err := http.NewRequest("POST", config.Config.JSONRPCEndpointURL, bytes.NewBuffer(requestBody))
	utils.CheckError(err)
	req.Header.Add("Authorization", "Bearer "+config.Config.SecretKey)
	_, err = client.Do(req)
	utils.CheckError(err)
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

func MoveTorrent(torrent torrents.Torrent, newPath string) {
	requestBody := []byte(`{
		"jsonrpc": "2.0",
		"method": "torrents.move",
		"params": {
			"info_hash":` + getMarshalledInfoHash(torrent) + `,
			"path": "` + newPath + `"
		}
	}`)
	req, err := http.NewRequest("POST", config.Config.JSONRPCEndpointURL, bytes.NewBuffer(requestBody))
	utils.CheckError(err)
	req.Header.Add("Authorization", "Bearer "+config.Config.SecretKey)
	_, err = client.Do(req)
	utils.CheckError(err)
}

func GetTorrentProperties(torrent torrents.Torrent) torrents.TorrentProperties {
	requestBody := []byte(`{
		"jsonrpc": "2.0",
		"method": "torrents.properties.get",
		"params": {
			"info_hash":` + getMarshalledInfoHash(torrent) + `
		}
	}`)
	req, err := http.NewRequest("POST", config.Config.JSONRPCEndpointURL, bytes.NewBuffer(requestBody))
	utils.CheckError(err)
	req.Header.Add("Authorization", "Bearer "+config.Config.SecretKey)
	resp, err := client.Do(req)
	utils.CheckError(err)
	body, err := ioutil.ReadAll(resp.Body)
	utils.CheckError(err)
	var torrentPropertiesRequestResponse torrents.TorrentPropertiesRequestResponse
	json.Unmarshal(body, &torrentPropertiesRequestResponse)
	return torrentPropertiesRequestResponse.Result
}

func SetTorrentProperties(torrent torrents.Torrent, torrentPropertiesSetData torrents.TorrentPropertiesSetData) {
	set_flags := 0
	if torrentPropertiesSetData.IsAutomaticallyManaged {
		set_flags |= 1 << 5
	}
	if torrentPropertiesSetData.IsSequenciallyDownloading {
		set_flags |= 1 << 9
	}

	unset_flags := 0
	if !torrentPropertiesSetData.IsAutomaticallyManaged {
		unset_flags |= 1 << 5
	}
	if !torrentPropertiesSetData.IsSequenciallyDownloading {
		unset_flags |= 1 << 9
	}
	requestBody := []byte(`{
		"jsonrpc": "2.0",
		"method": "torrents.properties.set",
		"params": {
			"info_hash":` + getMarshalledInfoHash(torrent) + `,
			"auto_managed": ` + marshallBool(torrentPropertiesSetData.IsAutomaticallyManaged) + `,
			"download_limit": ` + torrentPropertiesSetData.DownloadLimit + `,
			"max_connections": ` + torrentPropertiesSetData.MaxConnections + `,
			"max_uploads": ` + torrentPropertiesSetData.MaxUploads + `,
			"sequential_download": ` + marshallBool(torrentPropertiesSetData.IsSequenciallyDownloading) + `,
			"set_flags": ` + strconv.Itoa(set_flags) + `,
			"unset_flags": ` + strconv.Itoa(unset_flags) + `,
			"upload_limit": ` + torrentPropertiesSetData.UploadLimit + `
		}
	}`)
	req, err := http.NewRequest("POST", config.Config.JSONRPCEndpointURL, bytes.NewBuffer(requestBody))
	utils.CheckError(err)
	req.Header.Add("Authorization", "Bearer "+config.Config.SecretKey)
	_, err = client.Do(req)
	utils.CheckError(err)
}

func marshallBool(i bool) string {
	o, err := json.Marshal(i)
	utils.CheckError(err)
	return string(o)
}

func getMarshalledInfoHash(torrent torrents.Torrent) string {
	marshalledInfoHash, err := json.Marshal(torrent.InfoHash)
	utils.CheckError(err)
	marshalledInfoHashWithNull := strings.Replace(string(marshalledInfoHash), "\"\"", "null", -1)
	utils.CheckError(err)
	return marshalledInfoHashWithNull
}
