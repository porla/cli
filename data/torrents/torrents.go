package torrents

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

type TorrentListRequestResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  TorrentList `json:"result"`
}

func checkBit(flags uint64, bit uint64) bool {
	return (flags & (1 << bit)) == 1<<bit
}

func isAutoManaged(flags uint64) bool {
	return checkBit(flags, 5)
}

func isPaused(flags uint64) bool {
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
		return "error"
	}

	switch torrent.State {
	case 1:
		{
			if isPaused(torrent.Flags) {
				return "file check queued"
			}
			return "checking files"
		}
	case 2:
		return "downloading metadata"
	case 3:
		{
			if isPaused(torrent.Flags) {
				if isAutoManaged(torrent.Flags) {
					return "queued"
				}
				return "paused"
			}
			return "downloading"
		}
	case 4:
		return "finished"
	case 5:
		{
			if isPaused(torrent.Flags) {
				if isAutoManaged(torrent.Flags) {
					return "seeding queued"
				}
				return "finished"
			}
			return "seeding"
		}
	}
	return "unknown"
}
