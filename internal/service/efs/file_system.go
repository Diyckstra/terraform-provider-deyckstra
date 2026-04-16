package efs

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceFileSystem() *schema.Resource {
	return &schema.Resource{
		Create: resourceFileSystemCreate,
		Read:   resourceFileSystemRead,
		Update: resourceFileSystemUpdate,
		Delete: resourceFileSystemDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		CustomizeDiff: verify.SetTagsDiff,

		Schema: map[string]*schema.Schema{
			"creation_token": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"performance_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(efs.PerformanceMode_Values(), false),
			},
			"number_of_mount_targets": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),

			"size_in_bytes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceFileSystemCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EFSConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	creationToken := ""
	if v, ok := d.GetOk("creation_token"); ok {
		creationToken = v.(string)
	} else {
		creationToken = resource.UniqueId()
	}

	createOpts := &efs.CreateFileSystemInput{
		CreationToken: aws.String(creationToken),
		Tags:          Tags(tags.IgnoreAWS()),
	}

	log.Printf("[DEBUG] Creating EFS file system: %s", createOpts)
	fs, err := conn.CreateFileSystem(createOpts)

	if err != nil {
		return fmt.Errorf("error creating EFS file system: %w", err)
	}

	d.SetId(aws.StringValue(fs.FileSystemId))

	if _, err := waitFileSystemAvailable(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting for EFS file system (%s) to be available: %w", d.Id(), err)
	}

	return resourceFileSystemRead(d, meta)
}

func resourceFileSystemUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EFSConn

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Id(), o, n); err != nil {
			return fmt.Errorf("error updating EFS file system (%s) tags: %w", d.Id(), err)
		}
	}

	return resourceFileSystemRead(d, meta)
}

func resourceFileSystemRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EFSConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	fs, err := FindFileSystemByID(conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] EFS file system (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading EFS file system (%s): %w", d.Id(), err)
	}

	d.Set("creation_token", fs.CreationToken)
	d.Set("performance_mode", fs.PerformanceMode)
	d.Set("owner_id", fs.OwnerId)
	d.Set("number_of_mount_targets", fs.NumberOfMountTargets)

	tags := KeyValueTags(fs.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	if err := d.Set("size_in_bytes", flattenFileSystemSizeInBytes(fs.SizeInBytes)); err != nil {
		return fmt.Errorf("error setting size_in_bytes: %w", err)
	}

	return nil
}

func resourceFileSystemDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EFSConn

	log.Printf("[DEBUG] Deleting EFS file system: %s", d.Id())
	_, err := conn.DeleteFileSystem(&efs.DeleteFileSystemInput{
		FileSystemId: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, efs.ErrCodeFileSystemNotFound) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting EFS file system (%s): %w", d.Id(), err)
	}

	if _, err := waitFileSystemDeleted(conn, d.Id()); err != nil {
		if tfawserr.ErrCodeEquals(err, efs.ErrCodeFileSystemNotFound) {
			return nil
		}
		return fmt.Errorf("error waiting for EFS file system (%s) deletion: %w", d.Id(), err)
	}

	return nil
}

func flattenFileSystemSizeInBytes(sizeInBytes *efs.FileSystemSize) []interface{} {
	if sizeInBytes == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"value": aws.Int64Value(sizeInBytes.Value),
	}

	if sizeInBytes.ValueInIA != nil {
		m["value_in_ia"] = aws.Int64Value(sizeInBytes.ValueInIA)
	}

	if sizeInBytes.ValueInStandard != nil {
		m["value_in_standard"] = aws.Int64Value(sizeInBytes.ValueInStandard)
	}

	return []interface{}{m}
}
