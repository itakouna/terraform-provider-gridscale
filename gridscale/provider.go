package gridscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GRIDSCALE_UUID", nil),
				Description: "User-UUID for the gridscale API.",
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GRIDSCALE_TOKEN", nil),
				Description: "API-token for the gridscale API.",
			},
			"api_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GRIDSCALE_URL", "https://api.gridscale.io"),
				Description: "the url for the gridscale API.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"gridscale_storage":      dataSourceGridscaleStorage(),
			"gridscale_network":      dataSourceGridscaleNetwork(),
			"gridscale_ipv4":         dataSourceGridscaleIp(),
			"gridscale_ipv6":         dataSourceGridscaleIp(),
			"gridscale_sshkey":       dataSourceGridscaleSshkey(),
			"gridscale_template":     dataSourceGridscaleTemplate(),
			"gridscale_loadbalancer": dataSourceGridscaleLoadBalancer(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"gridscale_server":       resourceGridscaleServer(),
			"gridscale_storage":      resourceGridscaleStorage(),
			"gridscale_network":      resourceGridscaleNetwork(),
			"gridscale_ipv4":         resourceGridscaleIpv4(),
			"gridscale_ipv6":         resourceGridscaleIpv6(),
			"gridscale_sshkey":       resourceGridscaleSshkey(),
			"gridscale_loadbalancer": resourceGridscaleLoadBalancer(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		UserUUID: d.Get("uuid").(string),
		APIToken: d.Get("token").(string),
		APIUrl:   d.Get("api_url").(string),
	}

	return config.Client()
}
