package structs

/* Tasks */
// TID (Task Input Data)
type TaskIn struct {
	ID       string `json:"id"`
	ArriveAt string `json:"arriveAt"`
	Type     string `json:"type"`
}

// TOD (Task Output Data)
type TaskOut struct {
	ID         string  `json:"id"`
	ArriveAt   string  `json:"arriveAt"`
	OffloadTo  string  `json:"offloadTo"`
	Type       string  `json:"type"`
	Deadline   float64 `json:"deadline"`
	LocalDelay float64 `json:"localDelay"`
	LocalEff   float64 `json:"localEff"`
	TODelay    float64 `json:"toDelay"`
	TOEff      float64 `json:"toEff"`
}

// Req-Col Pair
type Pair struct {
	Src  int     `json:"src"`
	Des  int     `json:"des"`
	Cost float64 `json:"cost"`
}

/* Networks */
type ComNet struct {
	ID    string `json:"id"`
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

type Node struct {
	IP            string     `json:"ip"`
	CGs           []CG       `json:"cgs"`
	RestComputing [3]float64 `json:"restComputing"` // kernel, frequency, util
	RestStorage   [2]float64 `json:"restStorage"`   // storage, util
}

type Link struct {
	NodeFrom string  `json:"nodeFrom"`
	NodeTo   string  `json:"nodeTo"`
	Rate     float64 `json:"rate"`    // upload rate
	EsDelay  float64 `json:"esDelay"` // delay to establish the socket
}

type CG struct {
	Type    string `json:"type"`    // available task type
	Utility [3]int `json:"utility"` // [vacant, ]idle, prepared, busy
}

type Tor struct {
	IP    string `json:"ip"`
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Type   string `json:"type"`
	Number int    `json:"number"`
}

/* Simulations */
// S-NPD (Simulated Network Performance Data)
type SiNet struct {
	RestComputing  float64    // computing resource margin
	RestStorage    float64    // storage resource margin
	LimitPerDocker [2]float64 // upper and lower limits of computing resources for each container
}

// S-TFD (Simulated Task Feature Data)
type SiTask struct {
	WorkLoad float64
	DataSize float64
	Deadline float64
}
