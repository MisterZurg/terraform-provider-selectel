package selectel

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/selectel/secretsmanager-go/service/certs"
)

func resourceSecretsmanagerCertificateV1() *schema.Resource {
	return &schema.Resource{
		Description: "represents a Certificate — entity from SecretsManager service",

		CreateContext: resourceSecretsmanagerCertificateV1Create,
		ReadContext:   resourceSecretsmanagerCertificateV1Read,
		UpdateContext: resourceSecretsmanagerCertificateV1Update,
		DeleteContext: resourceSecretsmanagerCertificateV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSecretsmanagerCertificateV1ImportState,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "name — name of the certificate",
				Type:        schema.TypeString,
				Required:    true,
			},
			"certificates": {
				Description: "certificates — ca_chain list of public certs for certificate",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Description: "certificate — that must begin with -----BEGIN CERTIFICATE----- and end with -----END CERTIFICATE-----",
					Type:        schema.TypeString,
				},
				Required: true,
			},
			"private_key": {
				Description: "private_key — that should start with -----BEGIN PRIVATE KEY----- and end with -----END PRIVATE KEY-----",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"project_id": {
				Description: "project_id — id of a project where secret is used",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"dns_names": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"id": {
				Description: "id — computed id of a certificate",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"issued_by": {
				Description: "issued_by — information that is incorporated into certificate",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"country": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
						"locality": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
						"serial_number": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"street_address": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
					},
				},
				Computed: true,
			},
			"serial": {
				Description: "serial — number written in the certificate that was chosen by the CA which issued the certificate",
				Type:     schema.TypeString,
				Computed: true,
			},
			"validity": {
				Description: "validity — validity of a certificate in terms of notBefore and notAfter timestamps.",
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"basic_constraints": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"not_after": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"not_before": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Computed: true,
			},
			"version": {
				Description: "version — of the certificate",
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceSecretsmanagerCertificateV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cl, diagErr := getSecretsManagerClient(d, meta)
	if diagErr != nil {
		return diagErr
	}

	name := d.Get("name").(string)
	rawCertificates := d.Get("certificates").([]interface{})
	certificates := convertToStringSlice(rawCertificates)

	privateKey := d.Get("private_key").(string)

	cert := certs.CreateCertificateRequest{
		Name: name,
		Pem: certs.Pem{
			Certificates: certificates,
			PrivateKey:   privateKey,
		},
	}

	log.Print(msgCreate(objectCertificate, cert))

	createdCert, errCr := cl.Certificates.Create(ctx, cert)
	if errCr != nil {
		return diag.FromErr(fmt.Errorf("can't create a secret resource: %w", errCr))
	}

	d.SetId(createdCert.ID)

	return resourceSecretsmanagerCertificateV1Read(ctx, d, meta)
}

func resourceSecretsmanagerCertificateV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cl, diagErr := getSecretsManagerClient(d, meta)
	if diagErr != nil {
		return diagErr
	}

	log.Print(msgGet(objectCertificate, d.Id()))
	cert, errGet := cl.Certificates.Get(ctx, d.Get("id").(string))
	if errGet != nil {
		return diag.FromErr(errGettingObject(objectCertificate, d.Id(), errGet))
	}
	d.Set("dns_names", cert.DNSNames)

	issuedByFlatten := resourceSecretsmanagerCertificateV1IssuedByToSet(cert.IssuedBy)
	if err := d.Set("issued_by", issuedByFlatten); err != nil {
		return diag.FromErr(err)
	}

	d.Set("serial", cert.Serial)

	vaidityFlatten := resourceSecretsmanagerCertificateV1ValidityByToSet(cert.Validity)
	if err := d.Set("validity", vaidityFlatten); err != nil {
		return diag.FromErr(err)
	}

	d.Set("version", cert.Version)

	return nil
}

func resourceSecretsmanagerCertificateV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cl, diagErr := getSecretsManagerClient(d, meta)
	if diagErr != nil {
		return diagErr
	}

	if d.HasChange("name") {
		newName := d.Get("name").(string)

		log.Print(msgUpdate(objectCertificate, d.Id(), newName))

		errUpd := cl.Certificates.UpdateName(ctx, d.Id(), newName)
		if errUpd != nil {
			return diag.FromErr(errUpdatingObject(objectCertificate, d.Id(), errUpd))
		}
	}

	if d.HasChange("certificates") && d.HasChange("private_key") {
		rawCertificates := d.Get("certificates").([]interface{})
		certificates := convertToStringSlice(rawCertificates)

		upd := certs.UpdateCertificateVersionRequest{
			Pem: certs.Pem{
				Certificates: certificates,
				PrivateKey:   d.Get("private_key").(string),
			},
		}

		log.Print(msgUpdate(objectCertificate, d.Id(), upd))

		errUpd := cl.Certificates.UpdateVersion(ctx, d.Id(), upd)
		if errUpd != nil {
			return diag.FromErr(errUpdatingObject(objectCertificate, d.Id(), errUpd))
		}
	}

	return resourceSecretsmanagerCertificateV1Read(ctx, d, meta)
}

func resourceSecretsmanagerCertificateV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cl, diagErr := getSecretsManagerClient(d, meta)
	if diagErr != nil {
		return diagErr
	}

	log.Print(msgDelete(objectCertificate, d.Id()))

	errDel := cl.Certificates.Delete(ctx, d.Id())
	if errDel != nil {
		return diag.FromErr(errDeletingObject(objectCertificate, d.Id(), errDel))
	}

	return nil
}

func resourceSecretsmanagerCertificateV1ImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if config.ProjectID == "" {
		return nil, fmt.Errorf("SEL_PROJECT_ID must be set for the resource import")
	}

	d.Set("project_id", config.ProjectID)

	return []*schema.ResourceData{d}, nil
}
