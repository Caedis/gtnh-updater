package models

type GTNHAsset struct {
	Config                    GTNHConfig
	Translations              GTNHConfig
	Mods                      []GTNHMod
	Latest_nightly            int
	Latest_successful_nightly int
}

type GTNHConfig struct {
	Name            string
	Latest_version  string
	Needs_attention bool
	Versions        []GTNHVersion
	Config_type     string `json:"type"`
	Repo_url        string
}

type GTNHMod struct {
	Name           string
	Latest_version string
	Private        bool
	Versions       []GTNHVersion
	License        string
	Repo_url       string
	Side           string
	// Curse
	Source     string
	Slug       string
	Project_id string
}

type GTNHVersion struct {
	Version_tag string
	//Changelog            string
	Prelease             bool
	Tagged_at            string
	Filename             string
	Download_url         string
	Browser_download_url string
}
