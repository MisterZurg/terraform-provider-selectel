package selectel

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/selectel/secretsmanager-go/service/secrets"
)

func resourceSecretsmanagerSecretV1() *schema.Resource {
	return &schema.Resource{
		Description: "represents a Secret — entity from SecretsManager service",

		CreateContext: resourceSecretsmanagerSecretV1Create,
		ReadContext:   resourceSecretsmanagerSecretV1Read,
		UpdateContext: resourceSecretsmanagerSecretV1Update,
		DeleteContext: resourceSecretsmanagerSecretV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
			// StateContext: resourceSecretsmanagerSecretV1ImportState,
		},

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  false,
				Sensitive: true,
				// DefaultFunc: func() (interface{}, error) {
				// 	return "SENSITIVE_TERRAFORM_IMPORT", nil
				// },
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSecretsmanagerSecretV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cl, diagErr := getSecretsManagerClient(d, meta)
	if diagErr != nil {
		return diagErr
	}

	key := d.Get("key").(string)
	desc := d.Get("description").(string)
	value := d.Get("value").(string)

	secret := secrets.UserSecret{
		Key:         key,
		Description: desc,
		Value:       value,
	}

	log.Print(msgCreate(objectSecret, secret))

	errCr := cl.Secrets.Create(ctx, secret)
	if errCr != nil {
		return diag.FromErr(fmt.Errorf("can't create a secret resource: %w", errCr))
	}

	projectID := d.Get("project_id").(string)
	d.SetId(resourceSecretV1BuildID(projectID, key))

	return resourceSecretsmanagerSecretV1Read(ctx, d, meta)
}

func resourceSecretsmanagerSecretV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	if config.ProjectID == "" {
		return diag.FromErr(errors.New("SEL_PROJECT_ID must be set for the resource import"))
	}
	d.Set("project_id", config.ProjectID)
	// projectID := d.Get("project_id").(string)

	cl, diagErr := getSecretsManagerClient(d, meta)
	if diagErr != nil {
		return diagErr
	}

	log.Print(msgGet(objectSecret, d.Id()))

	var key string
	fullID := strings.Split(d.Id(), "/")
	if len(fullID) > 1 {
		key = fullID[1]
	} else {
		key = fullID[0]
	}

	secret, errGet := cl.Secrets.Get(ctx, key)
	if errGet != nil {
		return diag.FromErr(errGettingObject(objectSecret, d.Id(), errGet))
	}

	d.Set("name", secret.Name)
	d.Set("key", secret.Name)
	d.Set("description", secret.Description)
	if _, ok := d.GetOk("value"); !ok {
		d.Set("value", "UNKNOWN")
	}

	return nil
}

func resourceSecretsmanagerSecretV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cl, diagErr := getSecretsManagerClient(d, meta)
	if diagErr != nil {
		return diagErr
	}

	key := d.Get("key").(string)

	log.Print(msgDelete(objectSecret, d.Id()))

	errDel := cl.Secrets.Delete(ctx, key)
	if errDel != nil {
		return diag.FromErr(errDeletingObject(objectSecret, d.Id(), errDel))
	}

	return nil
}

func resourceSecretsmanagerSecretV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cl, diagErr := getSecretsManagerClient(d, meta)
	if diagErr != nil {
		return diagErr
	}

	key := d.Get("key").(string)
	desc := d.Get("description").(string)

	secret := secrets.UserSecret{
		Key:         key,
		Description: desc,
	}

	log.Print(msgUpdate(objectSecret, d.Id(), secret))

	errUpd := cl.Secrets.Update(ctx, secret)
	if errUpd != nil {
		return diag.FromErr(errUpdatingObject(objectSecret, d.Id(), errUpd))
	}

	return resourceSecretsmanagerSecretV1Read(ctx, d, meta)
}

func resourceSecretV1BuildID(projectID, key string) string {
	return fmt.Sprintf("%s/%s", projectID, key)
}

func resourceSecretsmanagerSecretV1ImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if config.ProjectID == "" {
		return nil, errors.New("SEL_PROJECT_ID must be set for the resource import")
	}

	d.Set("project_id", config.ProjectID)
	// d.Set("value", "SENSITIVE_TERRAFORM_IMPORT")

	cl, diagErr := getSecretsManagerClient(d, meta)
	if diagErr != nil {
		return nil, fmt.Errorf("can't getSecretsManagerClient: %v", diagErr)
	}

	secretName := d.Id()

	secret, errGet := cl.Secrets.Get(ctx, secretName)
	if errGet != nil {
		return nil, errGettingObject(objectSecret, d.Id(), errGet)
	}

	d.Set("key", secret.Name)
	d.Set("description", secret.Description)
	return []*schema.ResourceData{d}, nil
}
