package gridscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/gridscale/gsclient-go"
)

func resourceGridscaleStorage() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleStorageCreate,
		Read:   resourceGridscaleStorageRead,
		Delete: resourceGridscaleStorageDelete,
		Update: resourceGridscaleStorageUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.",
				Required:    true,
			},
			"capacity": {
				Type:         schema.TypeInt,
				Description:  "The capacity of a storage in GB.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"location_uuid": {
				Type:        schema.TypeString,
				Description: "Helps to identify which datacenter an object belongs to.",
				Optional:    true,
				ForceNew:    true,
				Default:     "45ed677b-3702-4b36-be2a-a2eab9827950",
			},
			"storage_type": {
				Type:        schema.TypeString,
				Description: "(one of storage, storage_high, storage_insane)",
				Optional:    true,
				ForceNew:    true,
				Default:     "storage",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					valid := false
					for _, stype := range StorageTypes {
						if v.(string) == stype {
							valid = true
							break
						}
					}
					if !valid {
						errors = append(errors, fmt.Errorf("%v is not a valid storage type. Valid types are: %v", v.(string), strings.Join(StorageTypes, ",")))
					}
					return
				},
			},
			"license_product_no": {
				Type:        schema.TypeInt,
				Description: "If a template has been used that requires a license key (e.g. Windows Servers) this shows the product_no of the license (see the /prices endpoint for more details).",
				Computed:    true,
			},
			"last_used_template": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "status indicates the status of the object.",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "Defines the date and time the object was initially created.",
				Computed:    true,
			},
			"change_time": {
				Type:        schema.TypeString,
				Description: "Defines the date and time of the last object change.",
				Computed:    true,
			},
			"location_name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.",
				Computed:    true,
			},
			"location_country": {
				Type:        schema.TypeString,
				Description: "Formatted by the 2 digit country code (ISO 3166-2) of the host country.",
				Computed:    true,
			},
			"location_iata": {
				Type:        schema.TypeString,
				Description: "Uses IATA airport code, which works as a location identifier.",
				Computed:    true,
			},
			"current_price": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"template": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sshkeys": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"password": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"password_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"hostname": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"template_uuid": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(time.Minute * 3),
		},
	}
}

func resourceGridscaleStorageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	storage, err := client.GetStorage(d.Id())
	if err != nil {
		if requestError, ok := err.(*gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}

	d.Set("change_time", storage.Properties.ChangeTime)
	d.Set("location_iata", storage.Properties.LocationIata)
	d.Set("status", storage.Properties.Status)
	d.Set("license_product_no", storage.Properties.LicenseProductNo)
	d.Set("location_country", storage.Properties.LocationCountry)
	d.Set("usage_in_minutes", storage.Properties.UsageInMinutes)
	d.Set("last_used_template", storage.Properties.LastUsedTemplate)
	d.Set("current_price", storage.Properties.CurrentPrice)
	d.Set("capacity", storage.Properties.Capacity)
	d.Set("location_uuid", storage.Properties.LocationUuid)
	d.Set("storage_type", storage.Properties.StorageType)
	d.Set("parent_uuid", storage.Properties.ParentUuid)
	d.Set("name", storage.Properties.Name)
	d.Set("location_name", storage.Properties.LocationName)
	d.Set("create_time", storage.Properties.CreateTime)

	if err = d.Set("labels", storage.Properties.Labels); err != nil {
		return fmt.Errorf("Error setting labels: %v", err)
	}
	return nil
}

func resourceGridscaleStorageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	requestBody := gsclient.StorageUpdateRequest{
		Name:   d.Get("name").(string),
		Labels: d.Get("labels").(*schema.Set).List(),
	}

	err := client.UpdateStorage(d.Id(), requestBody)
	if err != nil {
		return err
	}

	return resourceGridscaleStorageRead(d, meta)
}

func resourceGridscaleStorageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	requestBody := gsclient.StorageCreateRequest{
		Name:         d.Get("name").(string),
		Capacity:     d.Get("capacity").(int),
		LocationUuid: d.Get("location_uuid").(string),
		StorageType:  d.Get("storage_type").(string),
		Labels:       d.Get("labels").(*schema.Set).List(),
	}

	//since only one template can be used, we can just look at index 0
	if _, ok := d.GetOk("template"); ok {
		template := gsclient.StorageTemplate{
			Password:     d.Get("template.0.password").(string),
			PasswordType: d.Get("template.0.password_type").(string),
			Hostname:     d.Get("template.0.hostname").(string),
			TemplateUuid: d.Get("template.0.template_uuid").(string),
		}

		if attr, ok := d.GetOk("template.0.sshkeys"); ok {
			for _, value := range attr.([]interface{}) {
				template.Sshkeys = append(template.Sshkeys, value.(string))
			}
		}
		requestBody.Template = &template
	}

	response, err := client.CreateStorage(requestBody)
	if err != nil {
		return err
	}

	d.SetId(response.ObjectUuid)

	log.Printf("The id for storage %s has been set to %v", requestBody.Name, response.ObjectUuid)

	return resourceGridscaleStorageRead(d, meta)
}

func resourceGridscaleStorageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		return resource.RetryableError(client.DeleteStorage(d.Id()))
	})
}
