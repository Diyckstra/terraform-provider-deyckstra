package paas

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/paas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func waitServiceCreated(ctx context.Context, conn *paas.PaaS, id string, timeout time.Duration) (*paas.Service, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{ServiceStatusPending, ServiceStatusClaimed, ServiceStatusCreating, ServiceStatusProvisioning},
		Target:  []string{ServiceStatusReady},
		Refresh: statusService(conn, id),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*paas.Service); ok {
		if err != nil {
			setServiceErrorToResourceLastError(output, err)
		}

		return output, err
	}

	return nil, err
}

func waitServiceUpdated(ctx context.Context, conn *paas.PaaS, id string, timeout time.Duration) (*paas.Service, error) { //nolint:unparam
	stateConf := &resource.StateChangeConf{
		Pending: []string{ServiceStatusClaimed, ServiceStatusUpdating, ServiceStatusRecovering},
		Target:  []string{ServiceStatusReady},
		Refresh: statusService(conn, id),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*paas.Service); ok {
		if err != nil {
			setServiceErrorToResourceLastError(output, err)
		}

		if aws.StringValue(output.Status) == ServiceStatusReady && aws.BoolValue(output.IsRolledBack) {
			return output, errors.New("an error occurred while updating the service and " +
				"it was rolled back to the previous version. " +
				"Please check the updated parameters and apply the changes again",
			)
		}

		return output, err
	}

	return nil, err
}

func waitServiceReady(ctx context.Context, conn *paas.PaaS, id string, timeout time.Duration) (*paas.Service, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			ServiceStatusPending,
			ServiceStatusClaimed,
			ServiceStatusCreating,
			ServiceStatusProvisioning,
			ServiceStatusUpdating,
			ServiceStatusRecovering,
		},
		Target:  []string{ServiceStatusReady},
		Refresh: statusService(conn, id),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if output, ok := outputRaw.(*paas.Service); ok {
		if err != nil {
			setServiceErrorToResourceLastError(output, err)
		}

		return output, err
	}

	return nil, err
}

const (
	prometheusOperationNotAllowedErrorCode = "OperationNotAllowed"
	prometheusTaskInProgressErrorMessage   = "another task is in progress"
)

var errPrometheusChildDesiredStateNotObservable = errors.New("desired Prometheus child resource state is not yet observable")

func retryPrometheusMutationWhenTaskInProgress(ctx context.Context, mutate func() error) error {
	timeout, err := remainingPrometheusMutationTimeout(ctx)
	if err != nil {
		return err
	}

	_, err = tfresource.RetryWhenAWSErrMessageContainsContext(
		ctx,
		timeout,
		func() (interface{}, error) {
			return nil, mutate()
		},
		prometheusOperationNotAllowedErrorCode,
		prometheusTaskInProgressErrorMessage,
	)

	return err
}

func waitPrometheusScrapeJobCreatedOrUpdated(
	ctx context.Context,
	conn *paas.PaaS,
	serviceID string,
	jobID string,
	desired prometheusChildDesiredState,
) error {
	return waitPrometheusChildCreatedOrUpdated(
		ctx,
		func() (interface{}, error) {
			return FindPrometheusScrapeJobByID(ctx, conn, serviceID, jobID)
		},
		func(output interface{}) bool {
			job, ok := output.(*paas.PrometheusScrapeJob)
			return ok && prometheusScrapeJobMatchesDesired(job, desired)
		},
	)
}

func waitPrometheusNotificationChannelCreatedOrUpdated(
	ctx context.Context,
	conn *paas.PaaS,
	serviceID string,
	channelID string,
	desired prometheusChildDesiredState,
) error {
	return waitPrometheusChildCreatedOrUpdated(
		ctx,
		func() (interface{}, error) {
			return FindPrometheusNotificationChannelByID(ctx, conn, serviceID, channelID)
		},
		func(output interface{}) bool {
			channel, ok := output.(*paas.NotificationChannel)
			return ok && prometheusNotificationChannelMatchesDesired(channel, desired)
		},
	)
}

func waitPrometheusRouteCreatedOrUpdated(
	ctx context.Context,
	conn *paas.PaaS,
	serviceID string,
	routeID string,
	desired prometheusChildDesiredState,
) error {
	return waitPrometheusChildCreatedOrUpdated(
		ctx,
		func() (interface{}, error) {
			return FindPrometheusRouteByID(ctx, conn, serviceID, routeID)
		},
		func(output interface{}) bool {
			route, ok := output.(*paas.PrometheusRoute)
			return ok && prometheusRouteMatchesDesired(route, desired)
		},
	)
}

func waitPrometheusChildCreatedOrUpdated(
	ctx context.Context,
	find func() (interface{}, error),
	matches func(interface{}) bool,
) error {
	timeout, err := remainingPrometheusMutationTimeout(ctx)
	if err != nil {
		return err
	}

	findAndMatch := func() (interface{}, error) {
		output, err := find()
		if err != nil {
			return nil, err
		}
		if !matches(output) {
			return nil, errPrometheusChildDesiredStateNotObservable
		}

		return output, nil
	}

	_, err = tfresource.RetryWhenContext(
		ctx,
		timeout,
		findAndMatch,
		func(err error) (bool, error) {
			if tfresource.NotFound(err) || errors.Is(err, errPrometheusChildDesiredStateNotObservable) {
				return true, err
			}

			return false, err
		},
	)

	return err
}

func waitPrometheusScrapeJobDeleted(
	ctx context.Context,
	conn *paas.PaaS,
	serviceID string,
	jobID string,
) error {
	return waitPrometheusChildDeleted(ctx, func() (interface{}, error) {
		return FindPrometheusScrapeJobByID(ctx, conn, serviceID, jobID)
	})
}

func waitPrometheusNotificationChannelDeleted(
	ctx context.Context,
	conn *paas.PaaS,
	serviceID string,
	channelID string,
) error {
	return waitPrometheusChildDeleted(ctx, func() (interface{}, error) {
		return FindPrometheusNotificationChannelByID(ctx, conn, serviceID, channelID)
	})
}

func waitPrometheusRouteDeleted(
	ctx context.Context,
	conn *paas.PaaS,
	serviceID string,
	routeID string,
) error {
	return waitPrometheusChildDeleted(ctx, func() (interface{}, error) {
		return FindPrometheusRouteByID(ctx, conn, serviceID, routeID)
	})
}

func waitPrometheusChildDeleted(ctx context.Context, find func() (interface{}, error)) error {
	timeout, err := remainingPrometheusMutationTimeout(ctx)
	if err != nil {
		return err
	}

	_, err = tfresource.RetryUntilNotFoundContext(ctx, timeout, find)

	return err
}

func waitPrometheusServiceUpdated(
	ctx context.Context,
	conn *paas.PaaS,
	serviceID string,
) error {
	timeout, err := remainingPrometheusMutationTimeout(ctx)
	if err != nil {
		return err
	}

	_, err = waitServiceUpdated(ctx, conn, serviceID, timeout)

	return err
}

func waitPrometheusServiceReadyBeforeMutation(
	ctx context.Context,
	conn *paas.PaaS,
	serviceID string,
) error {
	timeout, err := remainingPrometheusMutationTimeout(ctx)
	if err != nil {
		return err
	}

	// IsRolledBack describes the result of a previous service update. It must
	// not permanently block a later child mutation that can repair the service.
	_, err = waitServiceReady(ctx, conn, serviceID, timeout)
	return err
}

func remainingPrometheusMutationTimeout(ctx context.Context) (time.Duration, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return 0, errors.New("Prometheus child resource mutation context has no deadline")
	}

	timeout := time.Until(deadline)
	if timeout <= 0 {
		return 0, context.DeadlineExceeded
	}

	return timeout, nil
}

func waitServiceDeleted(ctx context.Context, conn *paas.PaaS, id string, timeout time.Duration) (*paas.Service, error) {
	stateConf := &resource.StateChangeConf{
		Pending:        []string{ServiceStatusPending, ServiceStatusClaimed, ServiceStatusDeleting},
		Target:         []string{ServiceStatusDeleted},
		Refresh:        statusService(conn, id),
		Timeout:        timeout,
		NotFoundChecks: 1,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if tfresource.NotFound(err) {
		return nil, nil
	}

	if output, ok := outputRaw.(*paas.Service); ok {
		if err != nil {
			setServiceErrorToResourceLastError(output, err)
		}

		return output, err
	}

	return nil, err
}

func setServiceErrorToResourceLastError(service *paas.Service, err error) {
	if status := aws.StringValue(service.Status); status != ServiceStatusError {
		return
	}

	errCode := aws.StringValue(service.ErrorCode)
	errDesc := aws.StringValue(service.ErrorDescription)
	tfresource.SetLastError(err, fmt.Errorf("code: %s, description: %s", errCode, errDesc))
}
