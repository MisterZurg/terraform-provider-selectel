package selectel

import (
	"fmt"

	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSecretsmanagerSecretV1ImportBasic(t *testing.T) {
	projectID := os.Getenv("SEL_PROJECT_ID")

	resourceName := "selectel_secretsmanager_secret_v1.secret_tf_acc_test_1"

	secretKey := acctest.RandomWithPrefix("tf-acc")
	secretValue := acctest.RandomWithPrefix("tf-acc")
	secretDescription := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckVPCV2ProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecretsmanagerSecretV1WithoutProjectBasic(secretKey, secretDescription, secretValue, projectID),
				Check:  testAccCheckSelectelImportEnv(resourceName),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
			},
		},
	})
}

func testAccSecretsmanagerSecretV1WithoutProjectBasic(secretKey, secretDescription, secretValue, projectID string) string {
	return fmt.Sprintf(`
		resource "selectel_secretsmanager_secret_v1" "secret_tf_acc_test_1" {
				key = "%s"
				description = "%s"
				value = "%s"
				project_id = "%s"
		}`,
		secretKey,
		secretDescription,
		secretValue,
		projectID,
	)
}
