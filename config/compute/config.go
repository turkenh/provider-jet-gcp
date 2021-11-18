package compute

import (
	"fmt"

	"github.com/crossplane-contrib/terrajet/pkg/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/crossplane-contrib/provider-tf-gcp/config/common"
)

// Configure configures individual resources by adding custom
// ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("google_compute_managed_ssl_certificate", func(r *config.Resource) {
		r.Kind = "ManagedSSLCertificate"
		r.UseAsync = true
	})

	p.AddResourceConfigurator("google_compute_subnetwork", func(r *config.Resource) {
		r.References["network"] = config.Reference{
			Type: "Network",
		}
		r.UseAsync = true
	})

	p.AddResourceConfigurator("google_compute_address", func(r *config.Resource) {
		r.References["network"] = config.Reference{
			Type: "Network",
		}
		r.References["subnetwork"] = config.Reference{
			Type: "Subnetwork",
		}
	})

	p.AddResourceConfigurator("google_compute_firewall", func(r *config.Resource) {
		r.References["network"] = config.Reference{
			Type:      "Network",
			Extractor: common.PathSelfLinkExtractor,
		}
	})

	p.AddResourceConfigurator("google_compute_router", func(r *config.Resource) {
		r.References["network"] = config.Reference{
			Type:      "Network",
			Extractor: common.PathSelfLinkExtractor,
		}
	})

	p.AddResourceConfigurator("google_compute_router_nat", func(r *config.Resource) {
		r.References["router"] = config.Reference{
			Type: "Router",
		}
		r.References["subnetwork.name"] = config.Reference{
			Type:      "Subnetwork",
			Extractor: common.PathSelfLinkExtractor,
		}
		r.UseAsync = true
	})

	p.AddResourceConfigurator("google_compute_instance", func(r *config.Resource) {
		// Note(turkenh): We have to modify schema of
		// "boot_disk.initialize_params", since "labels" key here is a map where
		// elements configured as nil, defaulting to map[string]string:
		// https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_instance#nested_initialize_params
		r.TerraformResource.
			Schema["boot_disk"].Elem.(*schema.Resource).
			Schema["initialize_params"].Elem.(*schema.Resource).
			Schema["labels"].Elem = schema.TypeString

		r.ExternalName = config.NameAsIdentifier
		r.ExternalName.GetNameFn = common.GetNameFromFullyQualifiedID
		r.ExternalName.GetIDFn = func(name string, parameters map[string]interface{}, providerConfig map[string]interface{}) string {
			return fmt.Sprintf("projects/%s/zones/%s/instances/%s", common.KeyProject, parameters["zone"].(string), name)
		}
		r.UseAsync = true
	})

	p.AddResourceConfigurator("google_compute_instance_template", func(r *config.Resource) {
		// Note(turkenh): We have to modify schema of
		// "boot_disk.initialize_params.labels", since it is a map where
		// elements configured as nil, defaulting to map[string]string:
		r.TerraformResource.Schema["metadata"].Elem = schema.TypeString

		r.UseAsync = true
	})

	p.AddResourceConfigurator("google_compute_instance_from_template", func(r *config.Resource) {
		// Note(turkenh): We have to modify schema of
		// "boot_disk.initialize_params.labels", since it is a map where
		// elements configured as nil, defaulting to map[string]string:
		// https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_instance#nested_initialize_params
		r.TerraformResource.
			Schema["boot_disk"].Elem.(*schema.Resource).
			Schema["initialize_params"].Elem.(*schema.Resource).
			Schema["labels"].Elem = schema.TypeString

		r.TerraformResource.Schema["metadata"].Elem = schema.TypeString
	})
}
