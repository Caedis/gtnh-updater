package asset

import (
	"encoding/json"
	"github.com/Caedis/gtnh-updater/internal/models"
	"io"
	"net/http"
)

func FetchAssets() (*models.GTNHAsset, error) {
	resp, err := http.Get("https://github.com/GTNewHorizons/DreamAssemblerXXL/raw/refs/heads/master/gtnh-assets.json")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Decode the JSON data
	var data models.GTNHAsset
	err = json.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
