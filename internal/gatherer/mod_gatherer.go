package gatherer

import "path/filepath"

func GatherJars(minecraftDir string) (map[string]struct{}, error) {
	modsDir := filepath.Join(minecraftDir, "mods", "*.jar")
	matches, err := filepath.Glob(modsDir)
	if err != nil {
		return nil, err
	}

	packModMap := make(map[string]struct{})
	for _, packMod := range matches {
		packModMap[filepath.Base(packMod)] = struct{}{}
	}

	return packModMap, nil
}
