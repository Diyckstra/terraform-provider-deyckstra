package paas

import (
	"context"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/paas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/service/paas/services"
)

func DataSourceBackup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBackupRead,

		Schema: map[string]*schema.Schema{
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
			"age_days": {
				Type:          schema.TypeInt,
				Optional:      true,
				Default:       ageDaysDefault,
				ConflictsWith: []string{"id"},
				ValidateFunc:  validation.IntAtLeast(0),
			},
			"ready_only": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       true,
				ConflictsWith: []string{"id"},
			},
			"database_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"age_days", "database_name", "ready_only", "service_class", "service_id", "service_type"},
			},
			"protected": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"service_class": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id"},
				ValidateFunc:  validation.StringInSlice(services.ServiceClassValues(), false),
			},
			"service_deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"service_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id"},
			},
			"service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id"},
				ValidateFunc:  validation.StringInSlice(services.ServiceTypeValues(), false),
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

func dataSourceBackupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn

	ageDays := d.Get("age_days").(int)
	databaseName := d.Get("database_name").(string)
	readyOnly := d.Get("ready_only").(bool)

	var filters []func(backup *paas.Backup) bool

	if readyOnly {
		filters = append(filters,
			func(b *paas.Backup) bool {
				return aws.StringValue(b.Status) == BackupStatusCreated
			},
		)
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	if ageDays != ageDaysDefault {
		filters = append(filters,
			func(b *paas.Backup) bool {
				if b.Time == nil {
					return false
				}

				creationTime := time.Unix(aws.Int64Value(b.Time), 0)
				creationDate := time.Date(
					creationTime.Year(), creationTime.Month(), creationTime.Day(),
					0, 0, 0, 0, time.UTC,
				)

				return today.Sub(creationDate).Hours()/DayHours == float64(ageDays)
			},
		)
	}

	var dbFilters []func(backup *paas.DatabaseBackup) bool

	if databaseName != "" {
		dbFilters = append(dbFilters,
			func(db *paas.DatabaseBackup) bool {
				return aws.StringValue(db.Name) == databaseName
			},
		)
	}

	var err error
	var res *paas.Backup
	if v, ok := d.GetOk("id"); ok {
		id := v.(string)

		res, err = FindBackupById(conn, id)

		if err != nil {
			return diag.Errorf("error reading PaaS Service Backup (%s): %s", id, err)
		}
	} else {
		serviceClass := d.Get("service_class").(string)
		serviceId := d.Get("service_id").(string)
		serviceType := d.Get("service_type").(string)

		backups, err := FindBackups(conn, serviceClass, serviceId, serviceType)

		if err != nil {
			return diag.Errorf("error reading PaaS Service Backups: %s", err)
		}

		if len(backups) == 0 {
			return diag.Errorf("no backups found")
		}

		var filtered []*paas.Backup
		for _, backup := range backups {
			if !checkBackup(backup, filters, dbFilters) {
				continue
			}
			filtered = append(filtered, backup)
		}

		if len(filtered) == 0 {
			return diag.Errorf("no backups match the specified criteria")
		}

		// always return the most recent backup
		// FIXME: fix semgrep warning
		sort.Slice(filtered, func(i, j int) bool { return *filtered[i].Time > *filtered[j].Time }) // nosemgrep: prefer-aws-go-sdk-pointer-conversion-conditional
		res = filtered[0]
	}

	d.SetId(aws.StringValue(res.Id))

	d.Set("protected", res.Protected)
	d.Set("service_class", res.ServiceClass)
	d.Set("service_deleted", res.ServiceDeleted)
	d.Set("service_id", res.ServiceId)
	d.Set("service_name", res.ServiceName)
	d.Set("service_type", res.ServiceType)
	d.Set("status", res.Status)

	if res.Time != nil {
		d.Set("time", time.Unix(aws.Int64Value(res.Time), 0).Format(time.RFC3339))
	}

	d.Set("databases", flattenDatabaseBackups(res.Databases))

	return nil
}

func flattenDatabaseBackups(backups []*paas.DatabaseBackup) []map[string]interface{} {
	if backups == nil {
		return []map[string]interface{}{}
	}

	var tfList []map[string]interface{}
	for _, backup := range backups {
		if backup == nil {
			continue
		}

		tfMap := map[string]interface{}{}

		if v := backup.Location; v != nil {
			tfMap["location"] = v
		}

		if v := backup.Logfile; v != nil {
			tfMap["logfile"] = v
		}

		if v := backup.Name; v != nil {
			tfMap["name"] = v
		}

		if v := backup.Size; v != nil {
			tfMap["size"] = v
		}

		if v := backup.Status; v != nil {
			tfMap["status"] = v
		}

		tfList = append(tfList, tfMap)
	}

	return tfList
}

func checkBackup(
	backup *paas.Backup,
	filters []func(*paas.Backup) bool,
	dbFilters []func(databaseBackup *paas.DatabaseBackup) bool,
) bool {
	for _, filter := range filters {
		if !filter(backup) {
			return false
		}
	}

	if len(dbFilters) == 0 {
		return true
	}

	// at least one database should be filtered
	for _, db := range backup.Databases {
		if checkDatabaseBackup(db, dbFilters) {
			return true
		}
	}

	return false
}

func checkDatabaseBackup(db *paas.DatabaseBackup, filters []func(*paas.DatabaseBackup) bool) bool {
	for _, filter := range filters {
		if !filter(db) {
			return false
		}
	}

	return true
}
