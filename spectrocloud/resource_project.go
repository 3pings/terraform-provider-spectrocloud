package spectrocloud

import (
	"context"
	"time"

	"github.com/spectrocloud/hapi/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/spectrocloud/terraform-provider-spectrocloud/pkg/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		SchemaVersion: 2,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.V1alpha1Client)
	var diags diag.Diagnostics

	uid, err := c.CreateProject(toProject(d))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(uid)

	return diags
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.V1alpha1Client)
	var diags diag.Diagnostics

	project, err := c.GetProject(d.Id())
	if err != nil {
		return diag.FromErr(err)
	} else if project == nil {
		// Deleted - Terraform will recreate it
		d.SetId("")
		return diags
	}

	if err := d.Set("name", project.Metadata.Name); err != nil {
		return diag.FromErr(err)
	}

	if v, found := project.Metadata.Annotations["description"]; found {
		if err := d.Set("description", v); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("tags", flattenTags(project.Metadata.Labels)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.V1alpha1Client)
	var diags diag.Diagnostics

	err := c.UpdateProject(d.Id(), toProject(d))
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.V1alpha1Client)
	var diags diag.Diagnostics

	err := c.DeleteProject(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func toProject(d *schema.ResourceData) *models.V1alpha1ProjectEntity {
	annotations := make(map[string]string)
	if len(d.Get("description").(string)) > 0 {
		annotations["description"] = d.Get("description").(string)
	}
	return &models.V1alpha1ProjectEntity{
		Metadata: &models.V1ObjectMeta{
			Name:        d.Get("name").(string),
			UID:         d.Id(),
			Labels:      toTags(d),
			Annotations: annotations,
		},
	}
}