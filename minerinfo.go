package minerinfo

// Aggregated miner information
type minerInfo struct {
	Stats      []stats      `json:"STATS"`
	Summary    []summary    `json:"SUMMARY"`
	Devdetails []devdetails `json:"DEVDETAILS"`
	Pools      []Pools      `json:"POOLS"`
}

// Result of stats reponse from cgminer api
// only neccassary fields for now
type stats struct {
	Type string `json:"Type"`
}

// Result of summary reponse from cgminer api
// only neccassary fields for now
type summary struct {
	MHS any `json:"MHS 5s"`
	GHS any `json:"GHS 5s"`
}

// Result of devdetails reponse from cgminer api
// only neccassary fields for now
type devdetails struct {
	Name  string `json:"Name"`
	Model string `json:"Model"`
}

// Result of Pools reponse from cgminer api
// only neccassary fields for now
type Pools struct {
	Url    string `json:"URL"`
	User   string `json:"User"`
	Status string `json:"Status"`
}

// Generic cgminer api request
type minerRequest struct {
	Command string `json:"command"`
	Args    []any  `json:"args,omitempty"`
}
