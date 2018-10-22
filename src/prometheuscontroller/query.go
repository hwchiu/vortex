package prometheuscontroller

import (
	"fmt"
	"time"

	"github.com/hwchiu/vortex/src/serviceprovider"
	pv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"golang.org/x/net/context"
)

type RangeSetting struct {
	Interval   time.Duration // minutes
	Resolution time.Duration // seconds
	Rate       time.Duration // minutes
}

// Instant query at a single point in time
func query(sp *serviceprovider.Container, expression string) (model.Vector, error) {
	api := sp.Prometheus.API

	testTime := time.Now()
	result, err := api.Query(context.Background(), expression, testTime)

	// https://github.com/prometheus/client_golang/blob/d6a9817c4afc94d51115e4a30d449056a3fbf547/api/prometheus/v1/api.go#L316
	// this api always return the err no matter what
	// so we should use result==nil to determine whether it is a true error
	if result == nil {
		return nil, err
	}

	if result.Type() == model.ValVector {
		return result.(model.Vector), nil
	}
	return nil, fmt.Errorf("the type of the return result can not be identify")
}

// Query over a range of time
func queryRange(sp *serviceprovider.Container, expression string, rs RangeSetting) (model.Matrix, error) {
	api := sp.Prometheus.API

	rangeSet := pv1.Range{Start: time.Now().Add(-time.Minute * rs.Interval), End: time.Now(), Step: time.Second * rs.Resolution}
	result, err := api.QueryRange(context.Background(), expression, rangeSet)

	// https://github.com/prometheus/client_golang/blob/d6a9817c4afc94d51115e4a30d449056a3fbf547/api/prometheus/v1/api.go#L316
	// this api always return the err no matter what
	// so we should use result==nil to determine whether it is a true error
	if result == nil {
		return nil, err
	}

	if result.Type() == model.ValMatrix {
		return result.(model.Matrix), nil
	}
	return nil, fmt.Errorf("the type of the return result can not be identify")
}
