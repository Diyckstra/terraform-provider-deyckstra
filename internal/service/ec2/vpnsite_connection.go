package ec2

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceVPNConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPNConnectionCreate,
		Read:   resourceVPNConnectionRead,
		Update: resourceVPNConnectionUpdate,
		Delete: resourceVPNConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"customer_gateway_configuration": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"customer_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"high_availability": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
			"local_ipv4_network_cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDRNetwork(0, 32),
			},
			"remote_ipv4_network_cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDRNetwork(0, 32),
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
			"tunnel1_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tunnel1_bgp_asn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tunnel1_bgp_holdtime": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tunnel1_cgw_inside_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			//lintignore:S018 // the platform accepts a single IKE version
			"tunnel1_ike_versions": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(VpnTunnelOptionsIKEVersion_Values(), false),
				},
			},
			"tunnel1_inside_cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validateVpnConnectionTunnelInsideCIDR(),
			},
			"tunnel1_phase1_dh_group_numbers": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 6,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntInSlice(VpnTunnelOptionsPhase1DHGroupNumber_Values()),
				},
			},
			"tunnel1_phase1_encryption_algorithms": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(VpnTunnelOptionsPhase1EncryptionAlgorithm_Values(), false),
				},
			},
			"tunnel1_phase1_integrity_algorithms": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(VpnTunnelOptionsIntegrityAlgorithm_Values(), false),
				},
			},
			"tunnel1_phase1_lifetime_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(900, 28800),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old == strconv.Itoa(defaultVpnTunnelOptionsPhase1LifetimeSeconds) && new == "0" {
						return true
					}
					return false
				},
			},
			//lintignore:S018 // the platform accepts a single group number
			"tunnel1_phase2_dh_group_numbers": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntInSlice(VpnTunnelOptionsPhase2DHGroupNumber_Values()),
				},
			},
			"tunnel1_phase2_encryption_algorithms": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(VpnTunnelOptionsPhase2EncryptionAlgorithm_Values(), false),
				},
			},
			"tunnel1_phase2_integrity_algorithms": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(VpnTunnelOptionsIntegrityAlgorithm_Values(), false),
				},
			},
			"tunnel1_phase2_lifetime_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(900, 3600),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old == strconv.Itoa(defaultVpnTunnelOptionsPhase2LifetimeSeconds) && new == "0" {
						return true
					}
					return false
				},
			},
			"tunnel1_preshared_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validateVpnConnectionTunnelPreSharedKey(),
			},
			"tunnel1_replay_window_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(32, 2048),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old == strconv.Itoa(defaultVpnTunnelOptionsReplayWindowSize) && new == "0" {
						return true
					}
					return false
				},
			},
			"tunnel1_vgw_inside_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tunnel2_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tunnel2_bgp_asn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tunnel2_bgp_holdtime": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tunnel2_cgw_inside_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			//lintignore:S018 // the platform accepts a single IKE version
			"tunnel2_ike_versions": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(VpnTunnelOptionsIKEVersion_Values(), false),
				},
			},
			"tunnel2_inside_cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validateVpnConnectionTunnelInsideCIDR(),
			},
			"tunnel2_phase1_dh_group_numbers": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 6,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntInSlice(VpnTunnelOptionsPhase1DHGroupNumber_Values()),
				},
			},
			"tunnel2_phase1_encryption_algorithms": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(VpnTunnelOptionsPhase1EncryptionAlgorithm_Values(), false),
				},
			},
			"tunnel2_phase1_integrity_algorithms": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(VpnTunnelOptionsIntegrityAlgorithm_Values(), false),
				},
			},
			"tunnel2_phase1_lifetime_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(900, 28800),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old == strconv.Itoa(defaultVpnTunnelOptionsPhase1LifetimeSeconds) && new == "0" {
						return true
					}
					return false
				},
			},
			//lintignore:S018 // the platform accepts a single group number
			"tunnel2_phase2_dh_group_numbers": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntInSlice(VpnTunnelOptionsPhase2DHGroupNumber_Values()),
				},
			},
			"tunnel2_phase2_encryption_algorithms": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(VpnTunnelOptionsPhase2EncryptionAlgorithm_Values(), false),
				},
			},
			"tunnel2_phase2_integrity_algorithms": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(VpnTunnelOptionsIntegrityAlgorithm_Values(), false),
				},
			},
			"tunnel2_phase2_lifetime_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(900, 3600),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old == strconv.Itoa(defaultVpnTunnelOptionsPhase2LifetimeSeconds) && new == "0" {
						return true
					}
					return false
				},
			},
			"tunnel2_preshared_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validateVpnConnectionTunnelPreSharedKey(),
			},
			"tunnel2_replay_window_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(32, 2048),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old == strconv.Itoa(defaultVpnTunnelOptionsReplayWindowSize) && new == "0" {
						return true
					}
					return false
				},
			},
			"tunnel2_vgw_inside_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(VpnConnectionType_Values(), false),
			},
			"vgw_telemetry": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accepted_route_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"certificate_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_status_change": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"outside_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},

		CustomizeDiff: customdiff.All(
			resourceVPNConnectionCustomizeDiff,
			verify.SetTagsDiff,
		),
	}
}

// resourceVPNConnectionCustomizeDiff rejects the options of the second tunnel for a
// connection without fault tolerance, because such a connection has a single tunnel.
// Only the configuration is checked: the tunnel2 attributes are computed and are
// present in the state of any connection.
func resourceVPNConnectionCustomizeDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	if diff.Get("high_availability").(bool) {
		return nil
	}

	config := diff.GetRawConfig()

	if config.IsNull() || !config.IsKnown() {
		return nil
	}

	for _, key := range []string{
		"tunnel2_ike_versions",
		"tunnel2_inside_cidr",
		"tunnel2_phase1_dh_group_numbers",
		"tunnel2_phase1_encryption_algorithms",
		"tunnel2_phase1_integrity_algorithms",
		"tunnel2_phase1_lifetime_seconds",
		"tunnel2_phase2_dh_group_numbers",
		"tunnel2_phase2_encryption_algorithms",
		"tunnel2_phase2_integrity_algorithms",
		"tunnel2_phase2_lifetime_seconds",
		"tunnel2_preshared_key",
		"tunnel2_replay_window_size",
	} {
		value := config.GetAttr(key)

		if value.IsNull() || !value.IsKnown() {
			continue
		}

		if value.Type().IsCollectionType() && value.LengthInt() == 0 {
			continue
		}

		return fmt.Errorf("%s cannot be set when high_availability is false: such a connection has a single tunnel", key)
	}

	return nil
}

// Tunnel option values applied by the platform when they are not specified.
var (
	defaultVpnTunnelOptionsPhase1LifetimeSeconds = 28800
	defaultVpnTunnelOptionsPhase2LifetimeSeconds = 3600
	defaultVpnTunnelOptionsReplayWindowSize      = 1024
)

func resourceVPNConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	input := &ec2.CreateVpnConnectionInput{
		CustomerGatewayId: aws.String(d.Get("customer_gateway_id").(string)),
		Options:           expandVpnConnectionOptionsSpecification(d),
		TagSpecifications: ec2TagSpecificationsFromKeyValueTags(tags, ec2.ResourceTypeVpnConnection),
		Type:              aws.String(d.Get("type").(string)),
		VpnGatewayId:      aws.String(d.Get("vpn_gateway_id").(string)),
	}

	log.Printf("[DEBUG] Creating EC2 VPN Connection: %s", input)
	output, err := conn.CreateVpnConnection(input)

	if err != nil {
		return fmt.Errorf("error creating EC2 VPN Connection: %w", err)
	}

	d.SetId(aws.StringValue(output.VpnConnection.VpnConnectionId))

	if _, err := WaitVPNConnectionCreated(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting for EC2 VPN Connection (%s) create: %w", d.Id(), err)
	}

	// Read off the API to populate our RO fields.
	return resourceVPNConnectionRead(d, meta)
}

func resourceVPNConnectionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	vpnConnection, err := FindVPNConnectionByID(conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] EC2 VPN Connection (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading EC2 VPN Connection (%s): %w", d.Id(), err)
	}

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   ec2.ServiceName,
		Region:    meta.(*conns.AWSClient).Region,
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("vpn-connection/%s", d.Id()),
	}.String()
	d.Set("arn", arn)
	d.Set("customer_gateway_id", vpnConnection.CustomerGatewayId)
	d.Set("type", vpnConnection.Type)
	d.Set("vpn_gateway_id", vpnConnection.VpnGatewayId)

	if err := d.Set("vgw_telemetry", flattenVgwTelemetries(vpnConnection.VgwTelemetry)); err != nil {
		return fmt.Errorf("error setting vgw_telemetry: %w", err)
	}

	tags := KeyValueTags(vpnConnection.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	d.Set("customer_gateway_configuration", vpnConnection.CustomerGatewayConfiguration)

	// The platform reports neither the connection options nor the tunnel ones, so the
	// values are taken from the customer gateway configuration. It is present in the
	// response only if the connection is in the pending or available state.
	tunnelInfo, err := CustomerGatewayConfigurationToTunnelInfo(
		aws.StringValue(vpnConnection.CustomerGatewayConfiguration),
		d.Get("tunnel1_preshared_key").(string), // Not currently available during import
		d.Get("tunnel1_inside_cidr").(string),
	)

	if err == nil {
		d.Set("high_availability", tunnelInfo.Tunnel2Address != "")
		d.Set("local_ipv4_network_cidr", tunnelInfo.LocalIpv4NetworkCidr)
		d.Set("remote_ipv4_network_cidr", tunnelInfo.RemoteIpv4NetworkCidr)
		d.Set("tunnel1_address", tunnelInfo.Tunnel1Address)
		d.Set("tunnel1_bgp_asn", tunnelInfo.Tunnel1BGPASN)
		d.Set("tunnel1_bgp_holdtime", tunnelInfo.Tunnel1BGPHoldTime)
		d.Set("tunnel1_cgw_inside_address", tunnelInfo.Tunnel1CgwInsideAddress)
		d.Set("tunnel1_inside_cidr", tunnelInfo.Tunnel1InsideCidr)
		d.Set("tunnel1_preshared_key", tunnelInfo.Tunnel1PreSharedKey)
		d.Set("tunnel1_vgw_inside_address", tunnelInfo.Tunnel1VgwInsideAddress)
		d.Set("tunnel2_address", tunnelInfo.Tunnel2Address)
		d.Set("tunnel2_bgp_asn", tunnelInfo.Tunnel2BGPASN)
		d.Set("tunnel2_bgp_holdtime", tunnelInfo.Tunnel2BGPHoldTime)
		d.Set("tunnel2_cgw_inside_address", tunnelInfo.Tunnel2CgwInsideAddress)
		d.Set("tunnel2_inside_cidr", tunnelInfo.Tunnel2InsideCidr)
		d.Set("tunnel2_preshared_key", tunnelInfo.Tunnel2PreSharedKey)
		d.Set("tunnel2_vgw_inside_address", tunnelInfo.Tunnel2VgwInsideAddress)
	} else {
		if vpnConnection.CustomerGatewayConfiguration != nil {
			log.Printf("[ERROR] Error unmarshaling Customer Gateway XML configuration for (%s): %s", d.Id(), err)
		}

		d.Set("tunnel1_address", nil)
		d.Set("tunnel1_bgp_asn", nil)
		d.Set("tunnel1_bgp_holdtime", nil)
		d.Set("tunnel1_cgw_inside_address", nil)
		d.Set("tunnel1_preshared_key", nil)
		d.Set("tunnel1_vgw_inside_address", nil)
		d.Set("tunnel2_address", nil)
		d.Set("tunnel2_bgp_asn", nil)
		d.Set("tunnel2_bgp_holdtime", nil)
		d.Set("tunnel2_cgw_inside_address", nil)
		d.Set("tunnel2_preshared_key", nil)
		d.Set("tunnel2_vgw_inside_address", nil)
	}

	return nil
}

func resourceVPNConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	// The platform supports no modifications of a VPN connection, so every argument
	// except the tags forces a new resource.
	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Id(), o, n); err != nil {
			return fmt.Errorf("error updating EC2 VPN Connection (%s) tags: %w", d.Id(), err)
		}
	}

	return resourceVPNConnectionRead(d, meta)
}

func resourceVPNConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	log.Printf("[INFO] Deleting EC2 VPN Connection: %s", d.Id())
	_, err := conn.DeleteVpnConnection(&ec2.DeleteVpnConnectionInput{
		VpnConnectionId: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, ErrCodeInvalidVpnConnectionIDNotFound) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting EC2 VPN Connection (%s): %w", d.Id(), err)
	}

	if _, err := WaitVPNConnectionDeleted(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting for EC2 VPN Connection (%s) delete: %w", d.Id(), err)
	}

	return nil
}

func expandVpnConnectionOptionsSpecification(d *schema.ResourceData) *ec2.VpnConnectionOptionsSpecification {
	apiObject := &ec2.VpnConnectionOptionsSpecification{}

	if v, ok := d.GetOk("local_ipv4_network_cidr"); ok {
		apiObject.LocalIpv4NetworkCidr = aws.String(v.(string))
	}

	if v, ok := d.GetOk("remote_ipv4_network_cidr"); ok {
		apiObject.RemoteIpv4NetworkCidr = aws.String(v.(string))
	}

	tunnel1Options, tunnel1Configured := expandVpnTunnelOptionsSpecification(d, "tunnel1_")

	// The platform creates a connection with two tunnels unless the options are given
	// for a single tunnel only, so the number of the tunnel options in the request
	// defines whether the connection is fault-tolerant. Tunnel options without any
	// value are not serialized, hence the platform default is sent explicitly to keep
	// the tunnel in the request.
	if !d.Get("high_availability").(bool) {
		if !tunnel1Configured {
			tunnel1Options.ReplayWindowSize = aws.Int64(int64(defaultVpnTunnelOptionsReplayWindowSize))
		}

		apiObject.TunnelOptions = []*ec2.VpnTunnelOptionsSpecification{tunnel1Options}

		return apiObject
	}

	tunnel2Options, tunnel2Configured := expandVpnTunnelOptionsSpecification(d, "tunnel2_")

	if tunnel1Configured != tunnel2Configured {
		if tunnel1Configured {
			tunnel2Options.ReplayWindowSize = aws.Int64(int64(defaultVpnTunnelOptionsReplayWindowSize))
		} else {
			tunnel1Options.ReplayWindowSize = aws.Int64(int64(defaultVpnTunnelOptionsReplayWindowSize))
		}
	}

	apiObject.TunnelOptions = []*ec2.VpnTunnelOptionsSpecification{tunnel1Options, tunnel2Options}

	return apiObject
}

// expandVpnTunnelOptionsSpecification returns the tunnel options and reports
// whether any of them is set in the configuration.
func expandVpnTunnelOptionsSpecification(d *schema.ResourceData, prefix string) (*ec2.VpnTunnelOptionsSpecification, bool) {
	apiObject := &ec2.VpnTunnelOptionsSpecification{}
	configured := false

	if v, ok := d.GetOk(prefix + "ike_versions"); ok {
		for _, v := range v.(*schema.Set).List() {
			apiObject.IKEVersions = append(apiObject.IKEVersions, &ec2.IKEVersionsRequestListValue{Value: aws.String(v.(string))})
			configured = true
		}
	}

	if v, ok := d.GetOk(prefix + "phase1_dh_group_numbers"); ok {
		for _, v := range v.(*schema.Set).List() {
			apiObject.Phase1DHGroupNumbers = append(apiObject.Phase1DHGroupNumbers, &ec2.Phase1DHGroupNumbersRequestListValue{Value: aws.Int64(int64(v.(int)))})
			configured = true
		}
	}

	if v, ok := d.GetOk(prefix + "phase1_encryption_algorithms"); ok {
		for _, v := range v.(*schema.Set).List() {
			apiObject.Phase1EncryptionAlgorithms = append(apiObject.Phase1EncryptionAlgorithms, &ec2.Phase1EncryptionAlgorithmsRequestListValue{Value: aws.String(v.(string))})
			configured = true
		}
	}

	if v, ok := d.GetOk(prefix + "phase1_integrity_algorithms"); ok {
		for _, v := range v.(*schema.Set).List() {
			apiObject.Phase1IntegrityAlgorithms = append(apiObject.Phase1IntegrityAlgorithms, &ec2.Phase1IntegrityAlgorithmsRequestListValue{Value: aws.String(v.(string))})
			configured = true
		}
	}

	if v, ok := d.GetOk(prefix + "phase2_dh_group_numbers"); ok {
		for _, v := range v.(*schema.Set).List() {
			apiObject.Phase2DHGroupNumbers = append(apiObject.Phase2DHGroupNumbers, &ec2.Phase2DHGroupNumbersRequestListValue{Value: aws.Int64(int64(v.(int)))})
			configured = true
		}
	}

	if v, ok := d.GetOk(prefix + "phase2_encryption_algorithms"); ok {
		for _, v := range v.(*schema.Set).List() {
			apiObject.Phase2EncryptionAlgorithms = append(apiObject.Phase2EncryptionAlgorithms, &ec2.Phase2EncryptionAlgorithmsRequestListValue{Value: aws.String(v.(string))})
			configured = true
		}
	}

	if v, ok := d.GetOk(prefix + "phase2_integrity_algorithms"); ok {
		for _, v := range v.(*schema.Set).List() {
			apiObject.Phase2IntegrityAlgorithms = append(apiObject.Phase2IntegrityAlgorithms, &ec2.Phase2IntegrityAlgorithmsRequestListValue{Value: aws.String(v.(string))})
			configured = true
		}
	}

	if v, ok := d.GetOk(prefix + "phase1_lifetime_seconds"); ok {
		apiObject.Phase1LifetimeSeconds = aws.Int64(int64(v.(int)))
		configured = true
	}

	if v, ok := d.GetOk(prefix + "phase2_lifetime_seconds"); ok {
		apiObject.Phase2LifetimeSeconds = aws.Int64(int64(v.(int)))
		configured = true
	}

	if v, ok := d.GetOk(prefix + "preshared_key"); ok {
		apiObject.PreSharedKey = aws.String(v.(string))
		configured = true
	}

	if v, ok := d.GetOk(prefix + "replay_window_size"); ok {
		apiObject.ReplayWindowSize = aws.Int64(int64(v.(int)))
		configured = true
	}

	if v, ok := d.GetOk(prefix + "inside_cidr"); ok {
		apiObject.TunnelInsideCidr = aws.String(v.(string))
		configured = true
	}

	return apiObject, configured
}

func flattenVgwTelemetry(apiObject *ec2.VgwTelemetry) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.AcceptedRouteCount; v != nil {
		tfMap["accepted_route_count"] = aws.Int64Value(v)
	}

	if v := apiObject.CertificateArn; v != nil {
		tfMap["certificate_arn"] = aws.StringValue(v)
	}

	if v := apiObject.LastStatusChange; v != nil {
		tfMap["last_status_change"] = aws.TimeValue(v).Format(time.RFC3339)
	}

	if v := apiObject.OutsideIpAddress; v != nil {
		tfMap["outside_ip_address"] = aws.StringValue(v)
	}

	if v := apiObject.Status; v != nil {
		tfMap["status"] = aws.StringValue(v)
	}

	if v := apiObject.StatusMessage; v != nil {
		tfMap["status_message"] = aws.StringValue(v)
	}

	return tfMap
}

func flattenVgwTelemetries(apiObjects []*ec2.VgwTelemetry) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}

		tfList = append(tfList, flattenVgwTelemetry(apiObject))
	}

	return tfList
}

type XmlVpnConnectionConfig struct {
	LocalIpv4NetworkCidr  string           `xml:"ipv4_network_cidr_local"`
	RemoteIpv4NetworkCidr string           `xml:"ipv4_network_cidr_remote"`
	Tunnels               []XmlIpsecTunnel `xml:"ipsec_tunnel"`
}

type XmlIpsecTunnel struct {
	BGPASN           string `xml:"vpn_gateway>bgp>asn"`
	BGPHoldTime      int    `xml:"vpn_gateway>bgp>hold_time"`
	CgwInsideAddress string `xml:"customer_gateway>tunnel_inside_address>ip_address"`
	InsideCidr       string `xml:"tunnel_cidr"`
	OutsideAddress   string `xml:"vpn_gateway>tunnel_outside_address>ip_address"`
	PreSharedKey     string `xml:"ike>pre_shared_key"`
	VgwInsideAddress string `xml:"vpn_gateway>tunnel_inside_address>ip_address"`
}

type TunnelInfo struct {
	LocalIpv4NetworkCidr    string
	RemoteIpv4NetworkCidr   string
	Tunnel1Address          string
	Tunnel1BGPASN           string
	Tunnel1BGPHoldTime      int
	Tunnel1CgwInsideAddress string
	Tunnel1InsideCidr       string
	Tunnel1PreSharedKey     string
	Tunnel1VgwInsideAddress string
	Tunnel2Address          string
	Tunnel2BGPASN           string
	Tunnel2BGPHoldTime      int
	Tunnel2CgwInsideAddress string
	Tunnel2InsideCidr       string
	Tunnel2PreSharedKey     string
	Tunnel2VgwInsideAddress string
}

func (slice XmlVpnConnectionConfig) Len() int {
	return len(slice.Tunnels)
}

func (slice XmlVpnConnectionConfig) Less(i, j int) bool {
	return slice.Tunnels[i].OutsideAddress < slice.Tunnels[j].OutsideAddress
}

func (slice XmlVpnConnectionConfig) Swap(i, j int) {
	slice.Tunnels[i], slice.Tunnels[j] = slice.Tunnels[j], slice.Tunnels[i]
}

// CustomerGatewayConfigurationToTunnelInfo turns the XML customer gateway configuration
// into a TunnelInfo. A connection without fault tolerance has a single tunnel, so the
// second one is left empty.
func CustomerGatewayConfigurationToTunnelInfo(xmlConfig string, tunnel1PreSharedKey string, tunnel1InsideCidr string) (*TunnelInfo, error) {
	var vpnConfig XmlVpnConnectionConfig

	if err := xml.Unmarshal([]byte(xmlConfig), &vpnConfig); err != nil {
		return nil, err
	}

	if len(vpnConfig.Tunnels) == 0 {
		return nil, fmt.Errorf("no tunnels in the customer gateway configuration")
	}

	// XML tunnel ordering was commented here as being inconsistent since
	// this logic was originally added. The original sorting is based on
	// outside address. Given potential tunnel identifying configuration,
	// we try to correctly align the tunnel ordering before preserving the
	// original outside address sorting fallback for backwards compatibility
	// as to not inadvertently flip existing configurations.
	if len(vpnConfig.Tunnels) > 1 {
		if tunnel1PreSharedKey != "" {
			if tunnel1PreSharedKey != vpnConfig.Tunnels[0].PreSharedKey && tunnel1PreSharedKey == vpnConfig.Tunnels[1].PreSharedKey {
				vpnConfig.Tunnels[0], vpnConfig.Tunnels[1] = vpnConfig.Tunnels[1], vpnConfig.Tunnels[0]
			}
		} else if cidr := tunnel1InsideCidr; cidr != "" {
			if _, ipNet, err := net.ParseCIDR(cidr); err == nil && ipNet != nil {
				vgwInsideAddressIP1 := net.ParseIP(vpnConfig.Tunnels[0].VgwInsideAddress)
				vgwInsideAddressIP2 := net.ParseIP(vpnConfig.Tunnels[1].VgwInsideAddress)

				if !ipNet.Contains(vgwInsideAddressIP1) && ipNet.Contains(vgwInsideAddressIP2) {
					vpnConfig.Tunnels[0], vpnConfig.Tunnels[1] = vpnConfig.Tunnels[1], vpnConfig.Tunnels[0]
				}
			}
		} else {
			sort.Sort(vpnConfig)
		}
	}

	tunnelInfo := &TunnelInfo{
		LocalIpv4NetworkCidr:    vpnConfig.LocalIpv4NetworkCidr,
		RemoteIpv4NetworkCidr:   vpnConfig.RemoteIpv4NetworkCidr,
		Tunnel1Address:          vpnConfig.Tunnels[0].OutsideAddress,
		Tunnel1BGPASN:           vpnConfig.Tunnels[0].BGPASN,
		Tunnel1BGPHoldTime:      vpnConfig.Tunnels[0].BGPHoldTime,
		Tunnel1CgwInsideAddress: vpnConfig.Tunnels[0].CgwInsideAddress,
		Tunnel1InsideCidr:       vpnConfig.Tunnels[0].InsideCidr,
		Tunnel1PreSharedKey:     vpnConfig.Tunnels[0].PreSharedKey,
		Tunnel1VgwInsideAddress: vpnConfig.Tunnels[0].VgwInsideAddress,
	}

	if len(vpnConfig.Tunnels) > 1 {
		tunnelInfo.Tunnel2Address = vpnConfig.Tunnels[1].OutsideAddress
		tunnelInfo.Tunnel2BGPASN = vpnConfig.Tunnels[1].BGPASN
		tunnelInfo.Tunnel2BGPHoldTime = vpnConfig.Tunnels[1].BGPHoldTime
		tunnelInfo.Tunnel2CgwInsideAddress = vpnConfig.Tunnels[1].CgwInsideAddress
		tunnelInfo.Tunnel2InsideCidr = vpnConfig.Tunnels[1].InsideCidr
		tunnelInfo.Tunnel2PreSharedKey = vpnConfig.Tunnels[1].PreSharedKey
		tunnelInfo.Tunnel2VgwInsideAddress = vpnConfig.Tunnels[1].VgwInsideAddress
	}

	return tunnelInfo, nil
}

func validateVpnConnectionTunnelPreSharedKey() schema.SchemaValidateFunc {
	return validation.All(
		validation.StringLenBetween(8, 64),
		validation.StringDoesNotMatch(regexp.MustCompile(`^0`), "cannot start with zero character"),
		validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z_.]+$`), "can only contain alphanumeric, period and underscore characters"),
	)
}

// https://docs.k2.cloud/en/services/interconnect/vpn_connections/operations.html
func validateVpnConnectionTunnelInsideCIDR() schema.SchemaValidateFunc {
	return validation.All(
		validation.IsCIDRNetwork(30, 30),
		validation.StringMatch(regexp.MustCompile(`^169\.254\.25[2-5]\.`), "must be within 169.254.252.0/22"),
	)
}
