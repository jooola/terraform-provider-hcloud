package placementgroup

import (
	"context"
	"log"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/terraform-provider-hcloud/internal/util"
	"github.com/hetznercloud/terraform-provider-hcloud/internal/util/hcloudutil"
)

const ResourceType = "hcloud_placement_group"

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePlacementGroupCreate,
		ReadContext:   resourcePlacementGroupRead,
		UpdateContext: resourcePlacementGroupUpdate,
		DeleteContext: resourcePlacementGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics { // nolint:revive
					if ok, err := hcloud.ValidateResourceLabels(i.(map[string]interface{})); !ok {
						return diag.FromErr(err)
					}
					return nil
				},
			},
			"servers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics { // nolint:revive
					placementGroupType := i.(string)
					switch hcloud.PlacementGroupType(placementGroupType) {
					case hcloud.PlacementGroupTypeSpread:
						return nil
					default:
						return diag.Errorf("%s is not a valid placement group type", placementGroupType)
					}
				},
			},
		},
	}
}

func resourcePlacementGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hcloud.Client)

	opts := hcloud.PlacementGroupCreateOpts{
		Name: d.Get("name").(string),
		Type: hcloud.PlacementGroupType(d.Get("type").(string)),
	}
	if labels, ok := d.GetOk("labels"); ok {
		tmpLabels := make(map[string]string)
		for k, v := range labels.(map[string]interface{}) {
			tmpLabels[k] = v.(string)
		}
		opts.Labels = tmpLabels
	}

	res, _, err := client.PlacementGroup.Create(ctx, opts)
	if err != nil {
		return hcloudutil.ErrorToDiag(err)
	}
	d.SetId(util.FormatID(res.PlacementGroup.ID))

	if res.Action != nil {
		if err := hcloudutil.WaitForAction(ctx, &client.Action, res.Action); err != nil {
			return hcloudutil.ErrorToDiag(err)
		}
	}

	return resourcePlacementGroupRead(ctx, d, m)
}

func resourcePlacementGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hcloud.Client)

	id, err := util.ParseID(d.Id())
	if err != nil {
		log.Printf("[WARN] invalid placement group id (%s), removing from state: %v", d.Id(), err)
		d.SetId("")
		return nil
	}

	placementGroup, _, err := client.PlacementGroup.GetByID(ctx, id)
	if err != nil {
		return hcloudutil.ErrorToDiag(err)
	}
	if placementGroup == nil {
		log.Printf("[WARN] placement group (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	setSchema(d, placementGroup)
	return nil
}

func resourcePlacementGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hcloud.Client)

	id, err := util.ParseID(d.Id())
	if err != nil {
		log.Printf("[WARN] invalid placement group id (%s), removing from state: %v", d.Id(), err)
		d.SetId("")
		return nil
	}

	placementGroup, _, err := client.PlacementGroup.GetByID(ctx, id)
	if err != nil {
		if handleNotFound(err, d) {
			return nil
		}
		return hcloudutil.ErrorToDiag(err)
	}

	d.Partial(true)

	if d.HasChange("name") {
		description := d.Get("name").(string)
		_, _, err := client.PlacementGroup.Update(ctx, placementGroup, hcloud.PlacementGroupUpdateOpts{
			Name: description,
		})
		if err != nil {
			if handleNotFound(err, d) {
				return nil
			}
			return hcloudutil.ErrorToDiag(err)
		}
	}

	if d.HasChange("labels") {
		labels := d.Get("labels")
		tmpLabels := make(map[string]string)
		for k, v := range labels.(map[string]interface{}) {
			tmpLabels[k] = v.(string)
		}
		_, _, err := client.PlacementGroup.Update(ctx, placementGroup, hcloud.PlacementGroupUpdateOpts{
			Labels: tmpLabels,
		})
		if err != nil {
			if handleNotFound(err, d) {
				return nil
			}
			return hcloudutil.ErrorToDiag(err)
		}
	}
	d.Partial(false)

	return resourcePlacementGroupRead(ctx, d, m)
}

func resourcePlacementGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hcloud.Client)

	id, err := util.ParseID(d.Id())
	if err != nil {
		log.Printf("[WARN] invalid placement group id (%s), removing from state: %v", d.Id(), err)
		d.SetId("")
		return nil
	}
	if _, err := client.PlacementGroup.Delete(ctx, &hcloud.PlacementGroup{ID: id}); err != nil {
		return hcloudutil.ErrorToDiag(err)
	}

	return nil
}

func handleNotFound(err error, d *schema.ResourceData) bool {
	if hcloud.IsError(err, hcloud.ErrorCodeNotFound) {
		log.Printf("[WARN] placement group (%s) not found, removing from state", d.Id())
		d.SetId("")
		return true
	}
	return false
}

func setSchema(d *schema.ResourceData, pg *hcloud.PlacementGroup) {
	util.SetSchemaFromAttributes(d, getAttributes(pg))
}

func getAttributes(pg *hcloud.PlacementGroup) map[string]interface{} {
	return map[string]interface{}{
		"id":      pg.ID,
		"name":    pg.Name,
		"labels":  pg.Labels,
		"type":    pg.Type,
		"servers": pg.Servers,
	}
}
