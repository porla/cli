package torrents

import "osprey/config"

type Torrent struct {
	DownloadRate  uint64    `json:"download_rate"`
	UploadRate    uint64    `json:"upload_rate"`
	Error         bool      `json:"error"`
	Flags         uint64    `json:"flags"`
	InfoHash      [2]string `json:"info_hash"`
	ListPeers     uint64    `json:"list_peers"`
	ListSeeds     uint64    `json:"list_seeds"`
	Name          string    `json:"name"`
	NumPeers      uint64    `json:"num_peers"`
	NumSeeds      uint64    `json:"num_seeds"`
	Progress      float64   `json:"progress"`
	QueuePosition int64     `json:"queue_position"`
	SavePath      string    `json:"save_path"`
	Size          uint64    `json:"size"`
	State         uint      `json:"state"`
	Total         uint64    `json:"total"`
	TotalDone     uint64    `json:"total_done"`
}

type TorrentList struct {
	Page          int       `json:"page"`
	PageSize      int       `json:"page_size"`
	Torrents      []Torrent `json:"torrents"`
	TorrentsTotal int       `json:"torrents_total"`
}

type TorrentListRequestError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type TorrentListRequestResponse struct {
	JSONRPC string                  `json:"jsonrpc"`
	Result  TorrentList             `json:"result"`
	Error   TorrentListRequestError `json:"error"`
}

func checkBit(flags uint64, bit uint64) bool {
	return (flags & (1 << bit)) == 1<<bit
}

func isAutoManaged(flags uint64) bool {
	return checkBit(flags, 5)
}

func IsPaused(flags uint64) bool {
	return (flags & (1 << 4)) == 1<<4
}

func StateColor(torrent Torrent) string {
	if torrent.Error {
		return "1"
	}

	switch torrent.State {
	case 3:
		return "12"
	case 5:
		return "2"
	}
	return "15"
}

func StateString(torrent Torrent) string {
	if torrent.Error {
		return config.Currenti18n.TorrentStates.Error
	}

	switch torrent.State {
	case 1:
		{
			if IsPaused(torrent.Flags) {
				return config.Currenti18n.TorrentStates.FileCheckQueued
			}
			return config.Currenti18n.TorrentStates.CheckingFiles
		}
	case 2:
		return config.Currenti18n.TorrentStates.DownloadingMetadata
	case 3:
		{
			if IsPaused(torrent.Flags) {
				if isAutoManaged(torrent.Flags) {
					return config.Currenti18n.TorrentStates.Queued
				}
				return config.Currenti18n.TorrentStates.Paused
			}
			return config.Currenti18n.TorrentStates.Downloading
		}
	case 4:
		return config.Currenti18n.TorrentStates.Finished
	case 5:
		{
			if IsPaused(torrent.Flags) {
				if isAutoManaged(torrent.Flags) {
					return config.Currenti18n.TorrentStates.SeedingQueued
				}
				return config.Currenti18n.TorrentStates.Finished
			}
			return config.Currenti18n.TorrentStates.Seeding
		}
	}
	return config.Currenti18n.TorrentStates.Unknown
}
