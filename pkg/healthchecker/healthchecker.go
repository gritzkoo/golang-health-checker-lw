package healthchecker

import (
	"sync"
	"time"
)

func New(conf Config) *HealthCheck {
	return &HealthCheck{
		config: conf,
	}
}

type HealthCheck struct {
	config Config
}

func (h *HealthCheck) Liveness() Liveness {
	return Liveness{
		Status:  fullyFunctional,
		Version: h.config.Version,
	}
}

func (h *HealthCheck) Readiness() Readiness {
	var (
		start     = time.Now()
		wg        sync.WaitGroup
		checklist = make(chan integration, len(h.config.Integrations))
		result    = Readiness{
			Name:    h.config.Name,
			Version: h.config.Name,
			Status:  true,
			Date:    start.Format(time.RFC3339),
		}
	)
	wg.Add(len(h.config.Integrations))
	for _, v := range h.config.Integrations {
		go step(v, &result, &wg, checklist)
	}
	go func() {
		wg.Wait()
		close(checklist)
		result.Duration = time.Since(start).Seconds()
	}()
	for chk := range checklist {
		result.Integrations = append(result.Integrations, chk)
	}

	return result
}

func step(c Check, result *Readiness, wg *sync.WaitGroup, checklist chan integration) {
	defer (*wg).Done()
	st := time.Now()
	validation := c.Handle()
	check := integration{
		Name:         c.Name,
		URL:          validation.URL,
		ResponseTime: time.Since(st).Seconds(),
		Status:       validation.Error == nil,
		Error:        validation.Error,
	}
	if !check.Status {
		result.Status = false
	}
	checklist <- check
}
