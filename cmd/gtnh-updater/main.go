package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Caedis/gtnh-updater/internal/asset"
	"github.com/Caedis/gtnh-updater/internal/gatherer"
	"github.com/Caedis/gtnh-updater/internal/maven"
	"github.com/Caedis/gtnh-updater/internal/utils"
)

func main() {
	// Parse command-line flags
	minecraftDir := flag.String("dir", "", "Path to the Minecraft instance directory")
	useLatest := flag.Bool("latest", false, "Use the latest version of mods instead of nightly builds")
	_ = useLatest
	flag.Parse()

	if *minecraftDir == "" {
		fmt.Println("Error: --dir is required")
		os.Exit(1)
	}

	modsDir := filepath.Join(*minecraftDir, "mods")

	fmt.Printf("Updating Minecraft directory: %s\n", *minecraftDir)

	// Fetch assets
	fmt.Println("Fetching latest gtnh-assets.json")
	assets, err := asset.FetchAssets()
	if err != nil {
		log.Fatalf("Failed to fetch assets: %v\n", err)
	}

	// Gather mods in the mod directory
	packMods, err := gatherer.GatherJars(*minecraftDir)
	if err != nil {
		log.Fatalf("Failed to gather mods: %v\n", err)
	}

	// Initialize downloader
	dl, err := maven.NewDownloader()
	if err != nil {
		log.Fatalf("Failed to initialize downloader: %v\n", err)
	}

	// Process mods
	for _, mod := range assets.Mods {
		if mod.Side == "NONE" || len(mod.Versions) == 0 {
			continue
		}

		latestMod := mod.Versions[len(mod.Versions)-1]
		downloadURL := mod.Versions[len(mod.Versions)-1].Download_url
		if mod.Source == "" {
			downloadURL = fmt.Sprintf("https://nexus.gtnewhorizons.com/service/rest/v1/search/assets/download?sort=version&repository=public&group=com.github.GTNewHorizons&name=%s&maven.extension=jar&maven.classifier&version=%s", mod.Name, mod.Latest_version)
		}

		fmt.Printf("Downloading %s - %s\n", mod.Name, downloadURL)

		modPath, err := dl.Download(latestMod.Filename, downloadURL)
		if err != nil {
			log.Printf("Failed to download %s: %v\n", mod.Name, err)
			continue
		}

		// Remove old versions
		for _, modVer := range mod.Versions {
			if modVer.Filename != latestMod.Filename {
				if _, exists := packMods[modVer.Filename]; exists {
					fmt.Printf("Deleting old version: %s\n", modVer.Filename)
					os.Remove(filepath.Join(modsDir, modVer.Filename))
					delete(packMods, modVer.Filename)
				}
			}
		}

		// Move the latest mod to the mods directory
		latestModPath := filepath.Join(modsDir, latestMod.Filename)
		if modPath != latestModPath {
			fmt.Printf("Copying latest version to mods directory: %s\n", latestModPath)
			err := utils.CopyFile(modPath, latestModPath)
			if err != nil {
				log.Printf("Failed to copy file %s to mods directory: %v\n", modPath, err)
			}
		}

		// Update the packMods map with the new file
		packMods[latestMod.Filename] = struct{}{}
	}
}
