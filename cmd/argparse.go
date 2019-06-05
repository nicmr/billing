package cmd

import (
	"time"

	"github.com/Altemista/altemista-billing/pkg/costs"
)

const (
	iso8601 = "2006-01-02"
)

func parseCostProvider(s string) (costapi costs.Provider) {
	costapi = costs.Default()
	switch s {
	case "aws":
		costapi = costs.AWS()
	case "azure":
		costapi = costs.Azure()
	case "on-premise":
		costapi = costs.OnPremise()
	default:
		// stays costs.Default()
	}
	return
}

func parseMonth(s string) (time.Time, error) {
	var parsedMonth time.Time
	switch s {
	case "current":
		parsedMonth = time.Now()
	case "last":
		y, m, _ := time.Now().Date()
		parsedMonth = time.Date(y, m, 1, 0, 0, 0, 0, time.UTC).AddDate(0, -1, 0)
	default:
		// try to parse as iso
		s += "-01"
		var err error
		parsedMonth, err = time.Parse(iso8601, s)
		if err != nil {
			return time.Time{}, err
		}
	}
	return parsedMonth, nil
}
