package i18n

type i18nTorrentStates struct {
	Error               string
	FileCheckQueued     string
	CheckingFiles       string
	DownloadingMetadata string
	Queued              string
	Paused              string
	Downloading         string
	Finished            string
	SeedingQueued       string
	Seeding             string
	Unknown             string
}

type i18nKeybinds struct {
	YesKeybind                 string
	NoKeybind                  string
	EscKeybind                 string
	QKeybind                   string
	DoneKeybind                string
	SelectKeybind              string
	ChangePageKeybind          string
	SelectReducedKeybind       string
	AddTorrentKeybind          string
	PauseResumeKeybind         string
	RemoveTorrentKeybind       string
	MoveTorrentKeybind         string
	ToggleMagnetTorrentKeybind string
}

type I18n struct {
	AddTorrent string

	MagnetLink        string
	PathToTorrentFile string
	SavePath          string

	DeletingTorrentName string
	KeepDataQuestion    string

	MovingTorrentName string
	NewSavePath       string

	TorrentsActive string
	Torrent        string

	PageInfo string

	MagnetLinkPlaceHolder string
	TorrentFilePath       string
	SaveDirPlaceHolder    string
	NewSaveDirPlaceHolder string

	Keybinds i18nKeybinds

	TorrentStates i18nTorrentStates

	SeeYouLater              string
	ErrorNonExistantView     string
	ConnectingToPorlaBackend string
}

var English = I18n{
	AddTorrent: "Add torrent",

	MagnetLink:        "Magnet link",
	PathToTorrentFile: "Path to .torrent file",
	SavePath:          "Save path",

	DeletingTorrentName: "Deleting %s",
	KeepDataQuestion:    "Keep data?",

	MovingTorrentName: "Moving %s",
	NewSavePath:       "New save path",

	TorrentsActive: "%s active",
	Torrent:        "torrent",

	PageInfo: "Page %d/%d (max %d results)",

	MagnetLinkPlaceHolder: "magnet:...",
	TorrentFilePath:       "/path/to/torrent/file",
	SaveDirPlaceHolder:    "/path/to/save/dir/",
	NewSaveDirPlaceHolder: "/path/to/new/save/dir/",
	Keybinds: i18nKeybinds{
		YesKeybind:                 "y: yes",
		NoKeybind:                  "n: no",
		EscKeybind:                 "esc: back",
		QKeybind:                   "q: quit",
		DoneKeybind:                "enter: done",
		SelectKeybind:              "j/k, up/down: select",
		ChangePageKeybind:          "g/h, left/right: change page",
		SelectReducedKeybind:       "up/down: select",
		AddTorrentKeybind:          "a: add new torrent",
		PauseResumeKeybind:         "p: pause/resume torrent",
		RemoveTorrentKeybind:       "r: remove torrent",
		MoveTorrentKeybind:         "m: move torrent",
		ToggleMagnetTorrentKeybind: "tab: toggle magnet/.torrent",
	},

	TorrentStates: i18nTorrentStates{
		Error:               "error",
		FileCheckQueued:     "file check queued",
		CheckingFiles:       "checking files",
		DownloadingMetadata: "downloading metadata",
		Queued:              "queued",
		Paused:              "paused",
		Downloading:         "downloading",
		Finished:            "finished",
		SeedingQueued:       "seeding queued",
		Seeding:             "seeding",
		Unknown:             "unknown",
	},

	SeeYouLater:              "See you later!",
	ErrorNonExistantView:     "Error: Non existant view called.",
	ConnectingToPorlaBackend: "Establishing connection to Porla backend.",
}

var French = I18n{
	AddTorrent: "Ajout de torrent",

	MagnetLink:        "Lien aimant",
	PathToTorrentFile: "Chemin du fichier .torrent",
	SavePath:          "Chemin d'enregistrement",

	DeletingTorrentName: "Suppression de %s",
	KeepDataQuestion:    "Garder les données?",

	MovingTorrentName: "Déplacement de %s",
	NewSavePath:       "Nouveau chemin d'enregistrement",

	TorrentsActive: "%s actif(s)",
	Torrent:        "torrent",

	PageInfo: "Page %d/%d (max %d résultats)",

	MagnetLinkPlaceHolder: "magnet:...",
	TorrentFilePath:       "/chemin/du/fichier/torrent",
	SaveDirPlaceHolder:    "/chemin/du/dossier/denregistrement/",
	NewSaveDirPlaceHolder: "/chemin/du/nouveau/dossier/denregistrement/",
	Keybinds: i18nKeybinds{
		YesKeybind:                 "y: oui",
		NoKeybind:                  "n: non",
		EscKeybind:                 "esc: retour",
		QKeybind:                   "q: quitter",
		DoneKeybind:                "enter: finir",
		SelectKeybind:              "j/k, up/down: séléctioner",
		ChangePageKeybind:          "g/h, left/right: changer de page",
		SelectReducedKeybind:       "up/down: séléctioner",
		AddTorrentKeybind:          "a: ajouter un torrent",
		PauseResumeKeybind:         "p: interrompre/relancer le torrent",
		RemoveTorrentKeybind:       "r: supprimer torrent",
		MoveTorrentKeybind:         "m: déplacer torrent",
		ToggleMagnetTorrentKeybind: "tab: basculer de lien magnet à fichier .torrent",
	},

	TorrentStates: i18nTorrentStates{
		Error:               "erreur",
		FileCheckQueued:     "fichiers en file d'attente de vérification",
		CheckingFiles:       "vérification des fichiers",
		DownloadingMetadata: "téléchargement des métadonnées",
		Queued:              "en file d'attente",
		Paused:              "interrompu",
		Downloading:         "téléchargement",
		Finished:            "terminé",
		SeedingQueued:       "en file d'attente pour diffusion",
		Seeding:             "diffusion",
		Unknown:             "inconnu",
	},

	SeeYouLater:              "Au revoir!",
	ErrorNonExistantView:     "Erreur: Vue non existante appelée.",
	ConnectingToPorlaBackend: "Connection au backend de Porla.",
}

func LoadLanguage(I18nLanguage string) I18n {
	switch I18nLanguage {
	case "French":
		return French
	default:
		//Default to English
		return English
	}
}
