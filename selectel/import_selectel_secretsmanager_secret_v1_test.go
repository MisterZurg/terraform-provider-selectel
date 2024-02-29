package selectel

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSecretsmanagerSecretV1ImportBasic(t *testing.T) {
	resourceName := "selectel_secretsmanager_secret_v1.secret_tf_acc_test_1"

	projectName := acctest.RandomWithPrefix("tf-acc")
	secretKey := acctest.RandomWithPrefix("tf-acc")
	secretValue := acctest.RandomWithPrefix("tf-acc")
	secretDescription := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckVPCV2ProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecretsmanagerSecretV1BasicConfig(projectName, secretKey, secretValue, secretDescription),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
