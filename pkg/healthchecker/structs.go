package healthchecker

// Readiness response
type Readiness struct {
	Name         string        `json:"name,omitempty"`
	Status       bool          `json:"status"`
	Version      string        `json:"version,omitempty"`
	Date         string        `json:"date"`
	Duration     float64       `json:"duration"`
	Integrations []integration `json:"integrations"`
}

const fullyFunctional = "fully functional"

// Liveness response
type Liveness struct {
	Status  string `json:"status"`
	Version string `json:"version,omitempty"`
}

// integration default contract
type integration struct {
	Name         string  `json:"name"`
	Status       bool    `json:"status"`
	ResponseTime float64 `json:"response_time"`
	URL          string  `json:"url"`
	Error        error   `json:"error,omitempty"`
}

// Config the setup of this package
type Config struct {
	Name         string  `json:"name"`
	Version      string  `json:"version"`
	Integrations []Check `json:"integrations"`
}

// Check used to inform each integration config
type Check struct {
	Name   string               `json:"name"`
	Handle func() CheckResponse `json:"custom_check,omitempty"`
}

type CheckResponse struct {
	Error error  `json:"error,omitempty"`
	URL   string `json:"url,omitempty"`
}
