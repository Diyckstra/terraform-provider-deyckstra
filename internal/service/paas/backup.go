package paas

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/paas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/experimental/nullable"
)

func ResourceBackup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBackupCreate,
		ReadContext:   resourceBackupRead,
		UpdateContext: resourceBackupUpdate,
		DeleteContext: resourceBackupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"backup_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"databases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"logfile": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enable_deletion_protection": {
				Type:         nullable.TypeNullableBool,
				Optional:     true,
				ValidateFunc: nullable.ValidateTypeStringNullableBool,
			},
			// TODO: temporarily remove this parameter from the schema until it is approved.
			//"force_delete": {
			//	Type:     schema.TypeBool,
			//	Optional: true,
			//	Default:  false,
			//},
			"protected": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"service_class": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceBackupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Get("backup_id").(string)
	d.SetId(id)

	return resourceBackupUpdate(ctx, d, meta)
}

func resourceBackupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	id := d.Id()

	backup, err := FindBackupById(conn, id)

	if err != nil {
		return diag.Errorf("error reading PaaS Backup (%s): %s", id, err)
	}

	d.Set("enable_deletion_protection", aws.String(strconv.FormatBool(aws.BoolValue(backup.Protected))))
	d.Set("protected", backup.Protected)
	d.Set("service_class", backup.ServiceClass)
	d.Set("service_deleted", backup.ServiceDeleted)
	d.Set("service_id", backup.ServiceId)
	d.Set("service_name", backup.ServiceName)
	d.Set("service_type", backup.ServiceType)
	d.Set("status", backup.Status)

	if backup.Time != nil {
		d.Set("time", time.Unix(aws.Int64Value(backup.Time), 0).Format(time.RFC3339))
	}

	d.Set("databases", flattenDatabaseBackups(backup.Databases))

	return nil
}

func resourceBackupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	id := d.Id()

	if v, null, _ := nullable.Bool(d.Get("enable_deletion_protection").(string)).Value(); !null {
		input := &paas.ModifyBackupInput{
			BackupId:  aws.String(id),
			Protected: aws.Bool(v),
		}
		log.Printf("[DEBUG] Modifying PaaS Service Backup: %s", input)

		_, err := conn.ModifyBackup(input)

		if err != nil {
			return diag.Errorf("error modifying PaaS Service Backup (%s): %s", id, err)
		}
	}

	return resourceBackupRead(ctx, d, meta)
}

func resourceBackupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	id := d.Id()

	if v, ok := d.Get("force_delete").(bool); ok && v {
		input := &paas.DeleteBackupsInput{
			ServiceId: aws.String(d.Get("service_id").(string)),
			BackupIds: []*string{aws.String(id)},
		}

		log.Printf("[DEBUG] Deleting PaaS Service Backup: %s", input)
		_, err := conn.DeleteBackups(input)

		if err != nil {
			return diag.Errorf("error deleting PaaS Service Backup (%s): %s", id, err)
		}
	}

	log.Printf("[WARN] PaaS Service Backup (%s) was not deleted, removing from state", id)

	return nil
}
