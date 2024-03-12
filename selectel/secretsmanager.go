package selectel

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/selectel/secretsmanager-go"
)

func getSecretsManagerClient(d *schema.ResourceData, meta interface{}) (*secretsmanager.Client, diag.Diagnostics) {
	config := meta.(*Config)
	// region := d.Get("region").(string)
	selvpcClient, err := config.GetSelVPCClientWithProjectScope(d.Get("project_id").(string))
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("can't get project-scope selvpc client for craas: %w", err))
	}

	cl, err := secretsmanager.New(
		secretsmanager.WithAuthOpts(
			&secretsmanager.AuthOpts{KeystoneToken: selvpcClient.GetXAuthToken()},
		),

		// secretsmanager.WithCustomURLSecrets(endpoint),
		// secretsmanager.WithCustomURLCertificates(endpoint),
	)

	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("can't init secretsmanager client: %w", err))
	}

	return cl, nil
}

func convertToStringSlice(sl []interface{}) []string {
	result := make([]string, len(sl))
	for i := range sl {
		result[i] = sl[i].(string)
	}
	return result
}

func convertToInterfaceSlice(sl []string) []interface{} {
	result := make([]interface{}, len(sl))
	for i := range sl {
		result[i] = sl[i]
	}
	return result
}
