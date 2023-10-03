package minerinfo

import "strconv"

type genericMiner struct {
	data minerInfo
}

func (g genericMiner) Type() string {
	return "Unknown"
}

func (g genericMiner) Model() string {
	return "Unknown"
}

func (g genericMiner) Version() string {
	return "Unknown"
}

func (g genericMiner) Hashrate() float64 {
	if len(g.data.Summary) > 0 {
		if g.data.Summary[0].MHS != nil {

			return parseHR(g.data.Summary[0].MHS)
		} else {
			return parseHR(g.data.Summary[0].GHS) * 1000
		}
	} else {
		return 0
	}
}

func (g genericMiner) Pools() []Pools {
	return g.data.Pools
}

// asserts type of returned hashrate value and tries to parse it to float
func parseHR(value any) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case string:
		ret, _ := strconv.ParseFloat(v, 64)
		return ret
	default:
		return 0
	}
}
