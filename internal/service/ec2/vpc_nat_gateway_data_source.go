package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func DataSourceNATGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNATGatewayRead,

		Schema: map[string]*schema.Schema{
			"auto_provision_zones": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connectivity_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filter": CustomFiltersSchema(),
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"nat_gateway_addresses": natGatewayAddressesSchema(),
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceNATGatewayRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	input := &ec2.DescribeNatGatewaysInput{
		Filter: BuildAttributeFilterList(
			map[string]string{
				"state":  d.Get("state").(string),
				"vpc-id": d.Get("vpc_id").(string),
			},
		),
	}

	if v, ok := d.GetOk("id"); ok {
		input.NatGatewayIds = aws.StringSlice([]string{v.(string)})
	}

	if tags, ok := d.GetOk("tags"); ok {
		input.Filter = append(input.Filter, BuildTagFilterList(
			Tags(tftags.New(tags.(map[string]interface{}))),
		)...)
	}

	input.Filter = append(input.Filter, BuildCustomFilterList(
		d.Get("filter").(*schema.Set),
	)...)
	if len(input.Filter) == 0 {
		// Don't send an empty filters list; the EC2 API won't accept it.
		input.Filter = nil
	}

	ngw, err := FindNATGateway(conn, input)

	if err != nil {
		return tfresource.SingularDataSourceFindError("EC2 NAT Gateway", err)
	}

	d.SetId(aws.StringValue(ngw.NatGatewayId))
	d.Set("auto_provision_zones", ngw.AutoProvisionZones)
	d.Set("availability_mode", ngw.AvailabilityMode)
	d.Set("connectivity_type", ngw.ConnectivityType)
	d.Set("state", ngw.State)
	d.Set("vpc_id", ngw.VpcId)

	if err := d.Set("nat_gateway_addresses", flattenNATGatewayAddresses(ngw.NatGatewayAddresses)); err != nil {
		return fmt.Errorf("error setting nat_gateway_addresses: %w", err)
	}

	if err := d.Set("tags", KeyValueTags(ngw.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	return nil
}
