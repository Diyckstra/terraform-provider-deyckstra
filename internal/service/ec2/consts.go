package ec2

const (
	// https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_CreditSpecificationRequest.html#API_CreditSpecificationRequest_Contents
	CPUCreditsStandard  = "standard"
	CPUCreditsUnlimited = "unlimited"
)

func CPUCredits_Values() []string {
	return []string{
		CPUCreditsStandard,
		CPUCreditsUnlimited,
	}
}

const (
	// https://docs.aws.amazon.com/vpc/latest/privatelink/vpce-interface.html#vpce-interface-lifecycle
	VpcEndpointStateAvailable         = "available"
	VpcEndpointStateDeleted           = "deleted"
	VpcEndpointStateDeleting          = "deleting"
	VpcEndpointStateFailed            = "failed"
	VpcEndpointStatePending           = "pending"
	VpcEndpointStatePendingAcceptance = "pendingAcceptance"
	VpcEndpointStateRejected          = "rejected"
)

const (
	VpnStateModifying = "modifying"
)

// See https://docs.aws.amazon.com/vm-import/latest/userguide/vmimport-image-import.html#check-import-task-status
const (
	EBSSnapshotImportStateActive     = "active"
	EBSSnapshotImportStateDeleting   = "deleting"
	EBSSnapshotImportStateDeleted    = "deleted"
	EBSSnapshotImportStateUpdating   = "updating"
	EBSSnapshotImportStateValidating = "validating"
	EBSSnapshotImportStateValidated  = "validated"
	EBSSnapshotImportStateConverting = "converting"
	EBSSnapshotImportStateCompleted  = "completed"
)

// See https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_CreateNetworkInterface.html#API_CreateNetworkInterface_Example_2_Response.
const (
	NetworkInterfaceStatusPending = "pending"
)

// See https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeInternetGateways.html#API_DescribeInternetGateways_Example_1_Response.
const (
	InternetGatewayAttachmentStateAvailable = "available"
)

// See https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_CustomerGateway.html#API_CustomerGateway_Contents.
const (
	CustomerGatewayStateAvailable = "available"
	CustomerGatewayStateDeleted   = "deleted"
	CustomerGatewayStateDeleting  = "deleting"
	CustomerGatewayStatePending   = "pending"
)

const (
	VpnTunnelOptionsDPDTimeoutActionClear   = "clear"
	VpnTunnelOptionsDPDTimeoutActionNone    = "none"
	VpnTunnelOptionsDPDTimeoutActionRestart = "restart"
)

func VpnTunnelOptionsDPDTimeoutAction_Values() []string {
	return []string{
		VpnTunnelOptionsDPDTimeoutActionClear,
		VpnTunnelOptionsDPDTimeoutActionNone,
		VpnTunnelOptionsDPDTimeoutActionRestart,
	}
}

const (
	VpnTunnelOptionsIKEVersion1 = "ikev1"
	VpnTunnelOptionsIKEVersion2 = "ikev2"
)

func VpnTunnelOptionsIKEVersion_Values() []string {
	return []string{
		VpnTunnelOptionsIKEVersion1,
		VpnTunnelOptionsIKEVersion2,
	}
}

const (
	VpnTunnelOptionsEncryptionAlgorithmAES128           = "aes128"
	VpnTunnelOptionsEncryptionAlgorithmAES256           = "aes256"
	VpnTunnelOptionsEncryptionAlgorithmAESCCM128        = "aes_ccm128"
	VpnTunnelOptionsEncryptionAlgorithmAESCCM256        = "aes_ccm256"
	VpnTunnelOptionsEncryptionAlgorithmAESCTR128        = "aes_ctr128"
	VpnTunnelOptionsEncryptionAlgorithmAESCTR256        = "aes_ctr256"
	VpnTunnelOptionsEncryptionAlgorithmAESGCM128        = "aes_gcm128"
	VpnTunnelOptionsEncryptionAlgorithmAESGCM256        = "aes_gcm256"
	VpnTunnelOptionsEncryptionAlgorithmCamellia128      = "camellia128"
	VpnTunnelOptionsEncryptionAlgorithmCamellia256      = "camellia256"
	VpnTunnelOptionsEncryptionAlgorithmChaCha20Poly1305 = "chacha20poly1305"
)

// The aes_ccm128 and aes_ccm256 algorithms are supported for the second IKE phase only.
func VpnTunnelOptionsPhase1EncryptionAlgorithm_Values() []string {
	return []string{
		VpnTunnelOptionsEncryptionAlgorithmAES128,
		VpnTunnelOptionsEncryptionAlgorithmAES256,
		VpnTunnelOptionsEncryptionAlgorithmAESCTR128,
		VpnTunnelOptionsEncryptionAlgorithmAESCTR256,
		VpnTunnelOptionsEncryptionAlgorithmAESGCM128,
		VpnTunnelOptionsEncryptionAlgorithmAESGCM256,
		VpnTunnelOptionsEncryptionAlgorithmCamellia128,
		VpnTunnelOptionsEncryptionAlgorithmCamellia256,
		VpnTunnelOptionsEncryptionAlgorithmChaCha20Poly1305,
	}
}

func VpnTunnelOptionsPhase2EncryptionAlgorithm_Values() []string {
	return append(VpnTunnelOptionsPhase1EncryptionAlgorithm_Values(),
		VpnTunnelOptionsEncryptionAlgorithmAESCCM128,
		VpnTunnelOptionsEncryptionAlgorithmAESCCM256,
	)
}

const (
	VpnTunnelOptionsIntegrityAlgorithmSHA1   = "sha1"
	VpnTunnelOptionsIntegrityAlgorithmSHA256 = "sha256"
	VpnTunnelOptionsIntegrityAlgorithmSHA384 = "sha384"
	VpnTunnelOptionsIntegrityAlgorithmSHA512 = "sha512"
)

func VpnTunnelOptionsIntegrityAlgorithm_Values() []string {
	return []string{
		VpnTunnelOptionsIntegrityAlgorithmSHA1,
		VpnTunnelOptionsIntegrityAlgorithmSHA256,
		VpnTunnelOptionsIntegrityAlgorithmSHA384,
		VpnTunnelOptionsIntegrityAlgorithmSHA512,
	}
}

// The zero group number is supported for the second IKE phase only, it disables
// the Perfect Forward Secrecy.
func VpnTunnelOptionsPhase1DHGroupNumber_Values() []int {
	return []int{2, 5, 14, 15, 16, 17, 18, 19, 20, 21}
}

func VpnTunnelOptionsPhase2DHGroupNumber_Values() []int {
	return append([]int{0}, VpnTunnelOptionsPhase1DHGroupNumber_Values()...)
}

const (
	VpnConnectionTypeIpsec1      = "ipsec.1"
	VpnConnectionTypeIpsecLegacy = "ipsec.legacy"
)

func VpnConnectionType_Values() []string {
	return []string{
		VpnConnectionTypeIpsec1,
		VpnConnectionTypeIpsecLegacy,
	}
}

const (
	AmazonIPv6PoolID = "Amazon"
)

const (
	DefaultDHCPOptionsID = "default"
)

const (
	DefaultSecurityGroupName = "default"
)

const (
	LaunchTemplateVersionDefault = "$Default"
	LaunchTemplateVersionLatest  = "$Latest"
)
