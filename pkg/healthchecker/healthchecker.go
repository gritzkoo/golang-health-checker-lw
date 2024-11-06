package healthchecker

import (
	"sync"
	"time"
)

// Create a new pointer to HealthCheck ready to use
func New(conf Config) *HealthCheck {
	return &HealthCheck{
		config: conf,
	}
}

func (h *HealthCheck) getConcurrence() int {
	concurrence := h.config.Concurrence
	if concurrence == 0 {
		concurrence = 10
	}
	return concurrence
}

// The main object where the Liveness and Readiness actions reside!
type HealthCheck struct {
	config Config
}

/*
Liveness function will return a status and version fields.

Used to endpoint /health-check/liveness <- optional, just a convention
is used only to display if you application is up and running without verify if
any of its integrations is OK
*/
func (h *HealthCheck) Liveness() Liveness {
	return Liveness{
		Status:  fullyFunctional,
		Version: h.config.Version,
	}
}

/*
Readiness action

This function will execute all checks passed in
healthchecker.Config.Integrations[*].Handle functions
and return a detailed response
*/
func (h *HealthCheck) Readiness() Readiness {
	var (
		start     = time.Now()
		wg        sync.WaitGroup
		checklist = make(chan Integration, len(h.config.Integrations))
		semaphore = make(chan struct{}, h.getConcurrence())
		result    = Readiness{
			Name:    h.config.Name,
			Version: h.config.Version,
			Status:  true,
			Date:    start.Format(time.RFC3339),
		}
	)
	wg.Add(len(h.config.Integrations))
	for _, v := range h.config.Integrations {
		go step(v, &wg, checklist, semaphore)
	}
	go func() {
		wg.Wait()
		close(checklist)
	}()
	result.Duration = time.Since(start).Seconds()
	for chk := range checklist {
		if !chk.Status {
			result.Status = false
		}
		result.Integrations = append(result.Integrations, chk)
	}

	return result
}

// internal function to only execute the Check.Handle function async
func step(c Check, wg *sync.WaitGroup, checklist chan Integration, semaphore chan struct{}) {
	semaphore <- struct{}{} // reserve a spot on semaphore
	defer func() {
		wg.Done()
		<-semaphore // release semaphore spot
	}()
	st := time.Now()
	validation := c.Handle()
	check := Integration{
		Name:         c.Name,
		URL:          validation.URL,
		ResponseTime: time.Since(st).Seconds(),
		Status:       validation.Error == nil,
	}
	if !check.Status {
		check.Error = validation.Error.Error()
	}
	checklist <- check
}
