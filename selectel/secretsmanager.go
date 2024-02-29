package selectel

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/selectel/secretsmanager-go"
	"github.com/selectel/secretsmanager-go/service/certs"
)

func getSecretsManagerClient(d *schema.ResourceData, meta interface{}) (*secretsmanager.Client, diag.Diagnostics) {
	config := meta.(*Config)
	selvpcClient, err := config.GetSelVPCClientWithProjectScope(d.Get("project_id").(string))
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("can't get project-scope selvpc client for craas: %w", err))
	}

	cl, err := secretsmanager.New(
		secretsmanager.WithAuthOpts(
			&secretsmanager.AuthOpts{KeystoneToken: selvpcClient.GetXAuthToken()},
		),
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

func issuedBySchema() *schema.Resource {
	return resourceSecretsmanagerCertificateV1().Schema["issued_by"].Elem.(*schema.Resource)
}

func issuedByHashSetFunc() schema.SchemaSetFunc {
	return schema.HashResource(issuedBySchema())
}

// resourceSecretsmanagerCertificateV1IssuedByToSet — helper for setting attribute with nested structure.
func resourceSecretsmanagerCertificateV1IssuedByToSet(ib certs.IssuedBy) *schema.Set {
	issuedBySet := &schema.Set{
		F: issuedByHashSetFunc(),
	}

	issuedBySet.Add(map[string]interface{}{
		"country":        convertToInterfaceSlice(ib.Country),
		"locality":       convertToInterfaceSlice(ib.Locality),
		"serial_number":  ib.SerialNumber,
		"street_address": convertToInterfaceSlice(ib.StreetAddress),
	})

	return issuedBySet
}

func validitySchema() *schema.Resource {
	return resourceSecretsmanagerCertificateV1().Schema["validity"].Elem.(*schema.Resource)
}

func validityHashSetFunc() schema.SchemaSetFunc {
	return schema.HashResource(validitySchema())
}

// resourceSecretsmanagerCertificateV1ValidityByToSet — helper for setting attribute with nested structure.
func resourceSecretsmanagerCertificateV1ValidityByToSet(validity certs.Validity) *schema.Set {
	validitySet := &schema.Set{
		F: validityHashSetFunc(),
	}

	validitySet.Add(map[string]interface{}{
		"basic_constraints": validity.BasicConstraints,
		"not_after":         validity.NotAfter,
		"not_before":        validity.NotBefore,
	})

	return validitySet
}
