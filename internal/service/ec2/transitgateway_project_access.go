package ec2

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

const (
	idSeparator = "/"
)

func ResourceTransitGatewayProjectAccess() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: ResourceTransitGatewayProjectAccessCreate,
		ReadWithoutTimeout:   ResourceTransitGatewayProjectAccessRead,
		DeleteWithoutTimeout: ResourceTransitGatewayProjectAccessDelete,

		Importer: &schema.ResourceImporter{
			StateContext: ResourceTransitGatewayProjectAccessImport,
		},

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transit_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func ResourceTransitGatewayProjectAccessCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).EC2Conn

	transitGatewayID := d.Get("transit_gateway_id").(string)
	projectID := d.Get("account_id").(string)

	modifyInput := &ec2.ModifyTransitGatewayInput{
		TransitGatewayId: aws.String(transitGatewayID),
		Options: &ec2.ModifyTransitGatewayOptions{
			AddSharedOwners: []*string{aws.String(projectID)},
		},
	}

	log.Printf("[DEBUG] Creating Transit Gateway Project Access: %s", modifyInput)
	_, err := conn.ModifyTransitGateway(modifyInput)
	if err != nil {
		return diag.Errorf("error sharing Transit Gateway (%s) to project (%s): %s", transitGatewayID, projectID, err)
	}

	id := buildID(transitGatewayID, projectID)
	d.SetId(id)

	return ResourceTransitGatewayProjectAccessRead(ctx, d, meta)
}

func ResourceTransitGatewayProjectAccessRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).EC2Conn

	transitGatewayID := d.Get("transit_gateway_id").(string)
	projectID := d.Get("account_id").(string)

	transitGateway, err := FindTransitGatewayByID(conn, transitGatewayID)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Transit Gateway (%s) not found, removing from state", transitGatewayID)
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("error reading Transit Gateway (%s): %s", transitGatewayID, err)
	}

	d.Set("transit_gateway_id", transitGatewayID)
	d.Set("account_id", projectID)

	if transitGateway.Options != nil && transitGateway.Options.SharedOwners != nil {
		for _, owner := range transitGateway.Options.SharedOwners {
			if aws.StringValue(owner) == projectID {
				return nil
			}
		}
	}

	log.Printf("[WARN] Project (%s) not found in Transit Gateway (%s) shared owners, removing from state", projectID, transitGatewayID)
	d.SetId("")
	return nil
}

func ResourceTransitGatewayProjectAccessDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).EC2Conn

	transitGatewayID := d.Get("transit_gateway_id").(string)
	projectID := d.Get("account_id").(string)

	modifyInput := &ec2.ModifyTransitGatewayInput{
		TransitGatewayId: aws.String(transitGatewayID),
		Options: &ec2.ModifyTransitGatewayOptions{
			RemoveSharedOwners: []*string{aws.String(projectID)},
		},
	}

	log.Printf("[DEBUG] Deleting Transit Gateway Project Access: %s", modifyInput)
	_, err := conn.ModifyTransitGateway(modifyInput)

	if tfresource.NotFound(err) {
		log.Printf("[WARN] Transit Gateway (%s) not found, removing from state", transitGatewayID)
		return nil
	}

	if err != nil {
		return diag.Errorf("error removing project (%s) access from Transit Gateway (%s): %s", projectID, transitGatewayID, err)
	}

	return nil
}

func ResourceTransitGatewayProjectAccessImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	transitGatewayID, projectID, err := parseID(d.Id())
	if err != nil {
		return nil, err
	}

	d.SetId(d.Id())
	d.Set("transit_gateway_id", transitGatewayID)
	d.Set("account_id", projectID)

	return []*schema.ResourceData{d}, nil
}

func buildID(transitGatewayID, projectID string) string {
	return fmt.Sprintf("%s%s%s", transitGatewayID, idSeparator, projectID)
}

func parseID(id string) (string, string, error) {
	parts := strings.Split(id, idSeparator)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format for ID (%s), expected TRANSIT-GATEWAY-ID%sPROJECT-ID", id, idSeparator)
	}

	return parts[0], parts[1], nil
}
