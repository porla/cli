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
	TorrentSettingsKeybind     string
	ToggleOptionKeybind        string
	NinjaModeKeybind           string
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

	TorrentSettingsForTorrentName string
	AutomaticallyManaged          string
	SequentialDownload            string
	DownloadLimit                 string
	DownloadLimitHint             string
	MaxConnections                string
	MaxUploads                    string
	UploadLimit                   string

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

	TorrentSettingsForTorrentName: "Torrent settings for %s",
	AutomaticallyManaged:          "Automatically managed",
	SequentialDownload:            "Sequential download",
	DownloadLimit:                 "Download limit",
	DownloadLimitHint:             "The download limit for this torrent. -1 means unlimited.",
	MaxConnections:                "Max connections",
	MaxUploads:                    "Max uploads",
	UploadLimit:                   "Upload limit",

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
		TorrentSettingsKeybind:     "s: torrent settings",
		ToggleMagnetTorrentKeybind: "tab: toggle magnet/.torrent",
		ToggleOptionKeybind:        "space: toggle option",
		NinjaModeKeybind:           "n: toggle ninja mode",
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

	TorrentSettingsForTorrentName: "Réglages du torrent %s",
	AutomaticallyManaged:          "Géré automatiquement",
	SequentialDownload:            "Téléchargement séquentiel",
	DownloadLimit:                 "Limite de la vitesse de téléchargement",
	DownloadLimitHint:             "La limite de la vitesse de téléchargement pour ce torrent. -1 veut dire illimité.",
	MaxConnections:                "Nombre max. de connections",
	MaxUploads:                    "Nombre max. de partages",
	UploadLimit:                   "Upload limit",

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
		TorrentSettingsKeybind:     "s: réglages du torrent",
		ToggleMagnetTorrentKeybind: "tab: basculer de lien magnet à fichier .torrent",
		ToggleOptionKeybind:        "espace: faire basculer l'option",
		NinjaModeKeybind:           "n : activer/désactiver le mode ninja",
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
