package iam

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

const (
	RandomPasswordLength = 16
)

func ResourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		CustomizeDiff: verify.SetTagsDiff,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"force_destroy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Delete IAM user even if it has non-Terraform-managed IAM access keys, login profile or MFA devices",
			},
			"last_login_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"login": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 40),
					validation.StringMatch(
						regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_.-]*$`),
						"name must start with a Latin letter "+
							"and can only contain Latin letters, numbers, underscores (_), periods (.) and hyphens (-)",
					),
				),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == strings.ToLower(new)
				},
			},
			"otp_required": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"permissions_boundary": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 2048),
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"phone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secret_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
			"update_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	name := d.Get("name").(string)

	input := &iam.CreateUserInput{
		UserName: aws.String(name),
	}

	if v, ok := d.GetOk("display_name"); ok {
		input.DisplayName = aws.String(v.(string))
	} else {
		input.DisplayName = aws.String(name)
	}

	if v, ok := d.GetOk("email"); ok {
		input.Email = aws.String(v.(string))
	}

	if v, ok := d.GetOk("otp_required"); ok {
		input.OtpRequired = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("path"); ok {
		input.Path = aws.String(v.(string))
	}

	var password string
	if v, ok := d.GetOk("password"); ok {
		password = v.(string)
		input.Password = aws.String(password)
	} else {
		generatedPassword, err := GeneratePassword(RandomPasswordLength)

		if err != nil {
			return fmt.Errorf("error generating password: %w", err)
		}

		password = generatedPassword
		input.Password = aws.String(generatedPassword)
	}

	if v, ok := d.GetOk("permissions_boundary"); ok {
		input.PermissionsBoundary = aws.String(v.(string))
	}

	if len(tags) > 0 {
		input.Tags = Tags(tags.IgnoreAWS())
	}

	log.Printf("[DEBUG] Creating IAM user: %s", input)
	output, err := conn.CreateUser(input)

	if err != nil {
		return fmt.Errorf("failed creating IAM user (%s): %w", name, err)
	}

	d.SetId(aws.StringValue(output.User.UserName))
	d.Set("password", password)
	d.Set("secret_key", output.User.SecretKey)

	return resourceUserUpdate(d, meta)
}

func resourceUserRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	name := d.Id()

	user, err := FindUserByName(conn, name)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] IAM user (%s) not found, removing from state", name)
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading IAM user (%s): %w", name, err)
	}

	d.Set("arn", user.UserArn)
	d.Set("display_name", user.DisplayName)
	d.Set("email", user.Email)
	d.Set("enabled", user.Enabled)

	if user.LastLoginDate != nil {
		d.Set("last_login_date", aws.TimeValue(user.LastLoginDate).Format(time.RFC3339))
	} else {
		d.Set("last_login_date", nil)
	}

	d.Set("login", user.Login)
	d.Set("name", user.UserName)
	d.Set("otp_required", user.OtpRequired)
	d.Set("path", user.Path)

	if user.PermissionsBoundary != nil {
		d.Set("permissions_boundary", user.PermissionsBoundary.PermissionsBoundaryArn)
	}

	d.Set("phone", user.Phone)

	tags := KeyValueTags(user.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	if user.UpdateDate != nil {
		d.Set("update_date", aws.TimeValue(user.UpdateDate).Format(time.RFC3339))
	} else {
		d.Set("update_date", nil)
	}

	d.Set("user_id", user.UserId)

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn
	name := d.Id()

	updatableKeys := []string{
		"display_name",
		"email",
		"otp_required",
		"path",
		"password",
		"phone",
	}

	if d.HasChanges(updatableKeys...) {
		input := &iam.UpdateUserInput{
			UserName: aws.String(name),
		}

		if d.HasChange("display_name") {
			input.DisplayName = aws.String(d.Get("display_name").(string))
		}

		if d.HasChange("email") {
			input.Email = aws.String(d.Get("email").(string))
		}

		if d.HasChange("otp_required") {
			input.OtpRequired = aws.Bool(d.Get("otp_required").(bool))
		}

		if d.HasChange("path") {
			input.NewPath = aws.String(d.Get("path").(string))
		}

		if d.HasChange("password") {
			input.Password = aws.String(d.Get("password").(string))
		}

		if d.HasChange("phone") {
			input.Phone = aws.String(d.Get("phone").(string))
		}

		log.Printf("[DEBUG] Modifying IAM user: %s", input)
		_, err := conn.UpdateUser(input)

		if err != nil {
			return fmt.Errorf("error modifying IAM user (%s): %w", name, err)
		}
	}

	if d.HasChange("permissions_boundary") {
		permissionsBoundary := d.Get("permissions_boundary").(string)
		if permissionsBoundary != "" {
			input := &iam.PutUserPermissionsBoundaryInput{
				PermissionsBoundary: aws.String(permissionsBoundary),
				UserName:            aws.String(d.Id()),
			}
			_, err := conn.PutUserPermissionsBoundary(input)
			if err != nil {
				return fmt.Errorf("error updating IAM user permissions boundary: %w", err)
			}
		} else {
			input := &iam.DeleteUserPermissionsBoundaryInput{
				UserName: aws.String(d.Id()),
			}
			_, err := conn.DeleteUserPermissionsBoundary(input)
			if err != nil {
				return fmt.Errorf("error deleting IAM user permissions boundary: %w", err)
			}
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := userUpdateTags(conn, name, o, n); err != nil {
			return fmt.Errorf("error updating IAM user (%s) tags: %w", name, err)
		}
	}

	return resourceUserRead(d, meta)
}

func resourceUserDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn
	name := d.Id()

	// All access keys, MFA devices and login profile for the user must be removed
	if d.Get("force_destroy").(bool) {
		if err := DeleteUserAccessKeys(conn, name); err != nil {
			return fmt.Errorf("error removing IAM user (%s) access keys: %w", name, err)
		}

		if err := DeleteUserSSHKeys(conn, name); err != nil {
			return fmt.Errorf("error removing IAM user (%s) SSH keys: %w", name, err)
		}

		if err := DeleteUserVirtualMFADevices(conn, name); err != nil {
			return fmt.Errorf("error removing IAM user (%s) Virtual MFA devices: %w", name, err)
		}

		if err := DeactivateUserMFADevices(conn, name); err != nil {
			return fmt.Errorf("error removing IAM user (%s) MFA devices: %w", name, err)
		}

		if err := DeleteUserLoginProfile(conn, name); err != nil {
			return fmt.Errorf("error removing IAM user (%s) login profile: %w", name, err)
		}

		if err := deleteUserSigningCertificates(conn, name); err != nil {
			return fmt.Errorf("error removing IAM user (%s) signing certificate: %w", name, err)
		}

		if err := DeleteServiceSpecificCredentials(conn, name); err != nil {
			return fmt.Errorf("error removing IAM user (%s) Service Specific Credentials: %w", name, err)
		}
	}

	input := &iam.DeleteUserInput{
		UserName: aws.String(name),
	}

	log.Printf("[DEBUG] Deleting IAM user: %s", input)
	_, err := conn.DeleteUser(input)

	if tfawserr.ErrCodeEquals(err, UserNotFoundCode) {
		log.Printf("[WARN] IAM user (%s) not found, removing from state", name)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting IAM user %s: %w", name, err)
	}

	return nil
}

func DeleteUserGroupMemberships(conn *iam.IAM, username string) error {
	var groups []string
	listGroups := &iam.ListGroupsForUserInput{
		UserName: aws.String(username),
	}
	pageOfGroups := func(page *iam.ListGroupsForUserOutput, lastPage bool) (shouldContinue bool) {
		for _, g := range page.Groups {
			groups = append(groups, *g.GroupName)
		}
		return !lastPage
	}
	err := conn.ListGroupsForUserPages(listGroups, pageOfGroups)
	if err != nil {
		return fmt.Errorf("Error removing IAM user %q from all IAM groups: %s", username, err)
	}
	for _, g := range groups {
		// use iam group membership func to remove user from all groups
		log.Printf("[DEBUG] Removing IAM user %s from IAM group %s", username, g)
		// FIXME: Method isn't used, empty project name is added to fix compilation error.
		if err := removeUsersFromGroup(conn, []*string{aws.String(username)}, "", g); err != nil {
			return err
		}
	}

	return nil
}

func DeleteUserSSHKeys(svc *iam.IAM, username string) error {
	var publicKeys []string
	var err error

	listSSHPublicKeys := &iam.ListSSHPublicKeysInput{
		UserName: aws.String(username),
	}
	pageOfListSSHPublicKeys := func(page *iam.ListSSHPublicKeysOutput, lastPage bool) (shouldContinue bool) {
		for _, k := range page.SSHPublicKeys {
			publicKeys = append(publicKeys, *k.SSHPublicKeyId)
		}
		return !lastPage
	}
	err = svc.ListSSHPublicKeysPages(listSSHPublicKeys, pageOfListSSHPublicKeys)
	if err != nil {
		return fmt.Errorf("Error removing public SSH keys of IAM user %s: %w", username, err)
	}
	for _, k := range publicKeys {
		_, err := svc.DeleteSSHPublicKey(&iam.DeleteSSHPublicKeyInput{
			UserName:       aws.String(username),
			SSHPublicKeyId: aws.String(k),
		})
		if err != nil {
			return fmt.Errorf("Error deleting public SSH key %s: %w", k, err)
		}
	}

	return nil
}

func DeleteUserVirtualMFADevices(svc *iam.IAM, username string) error {
	var VirtualMFADevices []string
	var err error

	listVirtualMFADevices := &iam.ListVirtualMFADevicesInput{
		AssignmentStatus: aws.String("Assigned"),
	}
	pageOfVirtualMFADevices := func(page *iam.ListVirtualMFADevicesOutput, lastPage bool) (shouldContinue bool) {
		for _, m := range page.VirtualMFADevices {
			// UserName is `nil` for the root user
			if aws.StringValue(m.User.UserName) == username {
				VirtualMFADevices = append(VirtualMFADevices, *m.SerialNumber)
			}
		}
		return !lastPage
	}
	err = svc.ListVirtualMFADevicesPages(listVirtualMFADevices, pageOfVirtualMFADevices)
	if err != nil {
		return fmt.Errorf("Error removing Virtual MFA devices of IAM user %s: %w", username, err)
	}
	for _, m := range VirtualMFADevices {
		_, err := svc.DeactivateMFADevice(&iam.DeactivateMFADeviceInput{
			UserName:     aws.String(username),
			SerialNumber: aws.String(m),
		})
		if err != nil {
			return fmt.Errorf("Error deactivating Virtual MFA device %s: %w", m, err)
		}
		_, err = svc.DeleteVirtualMFADevice(&iam.DeleteVirtualMFADeviceInput{
			SerialNumber: aws.String(m),
		})
		if err != nil {
			return fmt.Errorf("Error deleting Virtual MFA device %s: %w", m, err)
		}
	}

	return nil
}

func DeactivateUserMFADevices(svc *iam.IAM, username string) error {
	var MFADevices []string
	var err error

	listMFADevices := &iam.ListMFADevicesInput{
		UserName: aws.String(username),
	}
	pageOfMFADevices := func(page *iam.ListMFADevicesOutput, lastPage bool) (shouldContinue bool) {
		for _, m := range page.MFADevices {
			MFADevices = append(MFADevices, *m.SerialNumber)
		}
		return !lastPage
	}
	err = svc.ListMFADevicesPages(listMFADevices, pageOfMFADevices)
	if err != nil {
		return fmt.Errorf("Error removing MFA devices of IAM user %s: %w", username, err)
	}
	for _, m := range MFADevices {
		_, err := svc.DeactivateMFADevice(&iam.DeactivateMFADeviceInput{
			UserName:     aws.String(username),
			SerialNumber: aws.String(m),
		})
		if err != nil {
			return fmt.Errorf("Error deactivating MFA device %s: %w", m, err)
		}
	}

	return nil
}

func DeleteUserLoginProfile(svc *iam.IAM, username string) error {
	var err error
	input := &iam.DeleteLoginProfileInput{
		UserName: aws.String(username),
	}
	err = resource.Retry(propagationTimeout, func() *resource.RetryError {
		_, err = svc.DeleteLoginProfile(input)
		if err != nil {
			if tfawserr.ErrCodeEquals(err, iam.ErrCodeNoSuchEntityException) {
				return nil
			}
			// EntityTemporarilyUnmodifiable: Login Profile for User XXX cannot be modified while login profile is being created.
			if tfawserr.ErrCodeEquals(err, iam.ErrCodeEntityTemporarilyUnmodifiableException) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if tfresource.TimedOut(err) {
		_, err = svc.DeleteLoginProfile(input)
	}
	if err != nil {
		return fmt.Errorf("Error deleting Account Login Profile: %w", err)
	}

	return nil
}

func DeleteUserAccessKeys(svc *iam.IAM, username string) error {
	var accessKeys []string
	var err error
	listAccessKeys := &iam.ListAccessKeysInput{
		UserName: aws.String(username),
	}
	pageOfAccessKeys := func(page *iam.ListAccessKeysOutput, lastPage bool) (shouldContinue bool) {
		for _, k := range page.AccessKeyMetadata {
			accessKeys = append(accessKeys, *k.AccessKeyId)
		}
		return !lastPage
	}
	err = svc.ListAccessKeysPages(listAccessKeys, pageOfAccessKeys)
	if err != nil {
		return fmt.Errorf("Error removing access keys of IAM user %s: %w", username, err)
	}
	for _, k := range accessKeys {
		_, err := svc.DeleteAccessKey(&iam.DeleteAccessKeyInput{
			UserName:    aws.String(username),
			AccessKeyId: aws.String(k),
		})
		if err != nil {
			return fmt.Errorf("Error deleting access key %s: %w", k, err)
		}
	}

	return nil
}

func deleteUserSigningCertificates(svc *iam.IAM, userName string) error {
	var certificateIDList []string

	listInput := &iam.ListSigningCertificatesInput{
		UserName: aws.String(userName),
	}
	err := svc.ListSigningCertificatesPages(listInput,
		func(page *iam.ListSigningCertificatesOutput, lastPage bool) bool {
			for _, c := range page.Certificates {
				certificateIDList = append(certificateIDList, aws.StringValue(c.CertificateId))
			}
			return !lastPage
		})
	if err != nil {
		return fmt.Errorf("Error removing signing certificates of IAM user %s: %w", userName, err)
	}

	for _, c := range certificateIDList {
		_, err := svc.DeleteSigningCertificate(&iam.DeleteSigningCertificateInput{
			CertificateId: aws.String(c),
			UserName:      aws.String(userName),
		})
		if err != nil {
			return fmt.Errorf("Error deleting signing certificate %s: %w", c, err)
		}
	}

	return nil
}

func DeleteServiceSpecificCredentials(svc *iam.IAM, username string) error {
	input := &iam.ListServiceSpecificCredentialsInput{
		UserName: aws.String(username),
	}

	output, err := svc.ListServiceSpecificCredentials(input)
	if err != nil {
		return fmt.Errorf("Error listing Service Specific Credentials of IAM user %s: %w", username, err)
	}
	for _, m := range output.ServiceSpecificCredentials {
		_, err := svc.DeleteServiceSpecificCredential(&iam.DeleteServiceSpecificCredentialInput{
			UserName:                    aws.String(username),
			ServiceSpecificCredentialId: m.ServiceSpecificCredentialId,
		})
		if err != nil {
			return fmt.Errorf("Error deleting Service Specific Credentials %s: %w", m, err)
		}
	}

	return nil
}
