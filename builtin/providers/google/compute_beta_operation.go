package google

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	computeBeta "google.golang.org/api/compute/v0.beta"
)

// OperationWaitType is an enum specifying what type of operation
// we're waiting on.
type ComputeBetaOperationWaitType byte

const (
	ComputeBetaOperationWaitInvalid ComputeBetaOperationWaitType = iota
	ComputeBetaOperationWaitGlobal
	ComputeBetaOperationWaitRegion
	ComputeBetaOperationWaitZone
)

type ComputeBetaOperationWaiter struct {
	Service *computeBeta.Service
	Op      *computeBeta.Operation
	Project string
	Region  string
	Type    ComputeBetaOperationWaitType
	Zone    string
}

func (w *ComputeBetaOperationWaiter) RefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var op *computeBeta.Operation
		var err error

		switch w.Type {
		case ComputeBetaOperationWaitGlobal:
			op, err = w.Service.GlobalOperations.Get(
				w.Project, w.Op.Name).Do()
		case ComputeBetaOperationWaitRegion:
			op, err = w.Service.RegionOperations.Get(
				w.Project, w.Region, w.Op.Name).Do()
		case ComputeBetaOperationWaitZone:
			op, err = w.Service.ZoneOperations.Get(
				w.Project, w.Zone, w.Op.Name).Do()
		default:
			return nil, "bad-type", fmt.Errorf(
				"Invalid wait type: %#v", w.Type)
		}

		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] Got %q when asking for operation %q", op.Status, w.Op.Name)

		return op, op.Status, nil
	}
}

func (w *ComputeBetaOperationWaiter) Conf() *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING"},
		Target:  "DONE",
		Refresh: w.RefreshFunc(),
	}
}

// ComputeBetaOperationError wraps computeBeta.OperationError and implements the
// error interface so it can be returned.
type ComputeBetaOperationError computeBeta.OperationError

func (e ComputeBetaOperationError) Error() string {
	var buf bytes.Buffer

	for _, err := range e.Errors {
		buf.WriteString(err.Message + "\n")
	}

	return buf.String()
}

func computeBetaOperationWaitGlobal(config *Config, op *computeBeta.Operation, activity string) error {
	w := &ComputeBetaOperationWaiter{
		Service: config.clientComputeBeta,
		Op:      op,
		Project: config.Project,
		Type:    ComputeBetaOperationWaitGlobal,
	}

	state := w.Conf()
	state.Delay = 10 * time.Second
	state.Timeout = 4 * time.Minute
	state.MinTimeout = 2 * time.Second
	opRaw, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for %s: %s", activity, err)
	}

	op = opRaw.(*computeBeta.Operation)
	if op.Error != nil {
		return ComputeBetaOperationError(*op.Error)
	}

	return nil
}

func computeBetaOperationWaitRegion(config *Config, op *computeBeta.Operation, region, activity string) error {
	w := &ComputeBetaOperationWaiter{
		Service: config.clientComputeBeta,
		Op:      op,
		Project: config.Project,
		Type:    ComputeBetaOperationWaitRegion,
		Region:  region,
	}

	state := w.Conf()
	state.Delay = 10 * time.Second
	state.Timeout = 4 * time.Minute
	state.MinTimeout = 2 * time.Second
	opRaw, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for %s: %s", activity, err)
	}

	op = opRaw.(*computeBeta.Operation)
	if op.Error != nil {
		return ComputeBetaOperationError(*op.Error)
	}

	return nil
}

func computeBetaOperationWaitZone(config *Config, op *computeBeta.Operation, zone, activity string) error {
	return computeBetaOperationWaitZoneTime(config, op, zone, 4, activity)
}

func computeBetaOperationWaitZoneTime(config *Config, op *computeBeta.Operation, zone string, minutes int, activity string) error {
	w := &ComputeBetaOperationWaiter{
		Service: config.clientComputeBeta,
		Op:      op,
		Project: config.Project,
		Zone:    zone,
		Type:    ComputeBetaOperationWaitZone,
	}
	state := w.Conf()
	state.Delay = 10 * time.Second
	state.Timeout = time.Duration(minutes) * time.Minute
	state.MinTimeout = 2 * time.Second
	opRaw, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for %s: %s", activity, err)
	}
	op = opRaw.(*computeBeta.Operation)
	if op.Error != nil {
		// Return the error
		return ComputeBetaOperationError(*op.Error)
	}
	return nil
}