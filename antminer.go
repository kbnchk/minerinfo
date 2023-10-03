package minerinfo

import (
	"strings"

	"slices"
)

type antminer struct {
	genericMiner
}

func (a antminer) Type() string {
	return "Antminer"
}

func (a antminer) Model() string {
	if len(a.data.Stats) > 0 {
		typemodel := a.data.Stats[0].Type
		modelver := strings.Join(strings.Fields(typemodel)[1:], " ")
		return strings.TrimSpace(strings.Split(modelver, "(")[0])
	} else {
		return a.genericMiner.Model()
	}
}

func (a antminer) Version() string {
	var version string
	if len(a.data.Stats) > 0 {
		typemodel := a.data.Stats[0].Type
		modelversion := strings.Fields(typemodel)
		if len(modelversion) > 1 {
			mvSplitted := strings.Split(strings.Join(modelversion[1:], " "), "(")
			if len(mvSplitted) > 1 {
				version = strings.TrimSuffix(strings.Join(mvSplitted[1:], "("), ")")
			}
		}
	}
	return version
}

func (a antminer) Hashrate() float64 {
	h := a.genericMiner.Hashrate()
	mhsAsGhsModels := []string{"L3+", "L7"} //L3+ and L7 returns MHS as GHS
	currentModel := a.Model()
	if slices.Contains(mhsAsGhsModels, currentModel) {
		return h / 1000
	} else {
		return h
	}
}
