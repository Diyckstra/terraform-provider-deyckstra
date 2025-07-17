package paas

// TODO: move to sdk
const (
	ServiceStatusPending      = "PENDING"
	ServiceStatusClaimed      = "CLAIMED"
	ServiceStatusCreating     = "CREATING"
	ServiceStatusProvisioning = "PROVISIONING"
	ServiceStatusUpdating     = "UPDATING"
	ServiceStatusDeleting     = "DELETING"
	ServiceStatusDeleted      = "DELETED"
	ServiceStatusReady        = "READY"
	ServiceStatusError        = "ERROR"
	ServiceStatusRecovering   = "RECOVERING"
)

const (
	ServiceNotFoundCode = "PaasServiceNotFound"
)

const (
	BackupStatusCreated = "CREATED"
)

const (
	DayHours = 24
)

const (
	ageDaysDefault = -1
)
