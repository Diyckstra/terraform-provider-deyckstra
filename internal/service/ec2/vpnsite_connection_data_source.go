package ec2

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func DataSourceVPNConnection() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPNConnectionRead,

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
				Optional: true,
				Computed: true,
			},
			"filter": DataSourceFiltersSchema(),
			"high_availability": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"local_ipv4_network_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remote_ipv4_network_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
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
			"tunnel1_inside_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tunnel1_preshared_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
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
			"tunnel2_inside_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tunnel2_preshared_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"tunnel2_vgw_inside_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceVPNConnectionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	input := &ec2.DescribeVpnConnectionsInput{}

	if v, ok := d.GetOk("filter"); ok {
		input.Filters = BuildFiltersDataSource(v.(*schema.Set))
	}

	input.Filters = append(input.Filters, BuildAttributeFilterList(
		map[string]string{
			"customer-gateway-id": d.Get("customer_gateway_id").(string),
			"type":                d.Get("type").(string),
			"vpn-gateway-id":      d.Get("vpn_gateway_id").(string),
		},
	)...)

	if len(input.Filters) == 0 {
		input.Filters = nil
	}

	if v, ok := d.GetOk("id"); ok {
		input.VpnConnectionIds = aws.StringSlice([]string{v.(string)})
	}

	vpnConnection, err := FindVPNConnection(conn, input)

	if err != nil {
		return tfresource.SingularDataSourceFindError("EC2 VPN Connection", err)
	}

	d.SetId(aws.StringValue(vpnConnection.VpnConnectionId))

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   ec2.ServiceName,
		Region:    meta.(*conns.AWSClient).Region,
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("vpn-connection/%s", d.Id()),
	}.String()
	d.Set("arn", arn)
	d.Set("customer_gateway_configuration", vpnConnection.CustomerGatewayConfiguration)
	d.Set("customer_gateway_id", vpnConnection.CustomerGatewayId)
	d.Set("state", vpnConnection.State)
	d.Set("type", vpnConnection.Type)
	d.Set("vpn_gateway_id", vpnConnection.VpnGatewayId)

	if err := d.Set("vgw_telemetry", flattenVgwTelemetries(vpnConnection.VgwTelemetry)); err != nil {
		return fmt.Errorf("error setting vgw_telemetry: %w", err)
	}

	if err := d.Set("tags", KeyValueTags(vpnConnection.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	// The platform reports neither the connection options nor the tunnel ones, so the
	// values are taken from the customer gateway configuration. It is present in the
	// response only if the connection is in the pending or available state.
	tunnelInfo, err := CustomerGatewayConfigurationToTunnelInfo(aws.StringValue(vpnConnection.CustomerGatewayConfiguration), "", "")

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
	} else if vpnConnection.CustomerGatewayConfiguration != nil {
		log.Printf("[ERROR] Error unmarshaling Customer Gateway XML configuration for (%s): %s", d.Id(), err)
	}

	return nil
}
