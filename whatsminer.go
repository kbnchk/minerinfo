package minerinfo

import "strings"

type whatsminer struct {
	genericMiner
}

// Returns miner
func (w whatsminer) Type() string {
	return "Whatsminer"
}

func (w whatsminer) Model() string {
	if len(w.data.Devdetails) > 0 {
		modelstr := w.data.Devdetails[0].Model
		modelversion := strings.Split(modelstr, "V")
		return modelversion[0]
	} else {
		return w.genericMiner.Model()
	}

}

func (w whatsminer) Version() string {
	var version string
	if len(w.data.Devdetails) > 0 {
		modelstr := w.data.Devdetails[0].Model
		modelversion := strings.Split(modelstr, "V")
		if len(modelversion) > 1 {
			return "V" + strings.Join(modelversion[1:], "V")
		}
	}
	return version
}
