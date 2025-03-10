package floatingip_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hetznercloud/terraform-provider-hcloud/internal/floatingip"
	"github.com/hetznercloud/terraform-provider-hcloud/internal/teste2e"
	"github.com/hetznercloud/terraform-provider-hcloud/internal/testsupport"
	"github.com/hetznercloud/terraform-provider-hcloud/internal/testtemplate"
)

func TestAccFloatingIPDataSource(t *testing.T) {
	tmplMan := testtemplate.Manager{}

	res := &floatingip.RData{
		Name: "floatingip-ds-test",
		Type: "ipv4",
		Labels: map[string]string{
			"key": strconv.Itoa(acctest.RandInt()),
		},
		HomeLocationName: teste2e.TestLocationName,
	}
	res.SetRName("floatingip-ds-test")
	floatingipByName := &floatingip.DData{
		FloatingIPName: res.TFID() + ".name",
	}
	floatingipByName.SetRName("floatingip_by_name")
	floatingipByID := &floatingip.DData{
		FloatingIPID: res.TFID() + ".id",
	}
	floatingipByID.SetRName("floatingip_by_id")
	floatingipBySel := &floatingip.DData{
		LabelSelector: fmt.Sprintf("key=${%s.labels[\"key\"]}", res.TFID()),
	}
	floatingipBySel.SetRName("floatingip_by_sel")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 teste2e.PreCheck(t),
		ProtoV6ProviderFactories: teste2e.ProtoV6ProviderFactories(),
		CheckDestroy:             testsupport.CheckResourcesDestroyed(floatingip.ResourceType, floatingip.ByID(t, nil)),
		Steps: []resource.TestStep{
			{
				Config: tmplMan.Render(t,
					"testdata/r/hcloud_floating_ip", res,
				),
			},
			{
				Config: tmplMan.Render(t,
					"testdata/r/hcloud_floating_ip", res,
					"testdata/d/hcloud_floating_ip", floatingipByName,
					"testdata/d/hcloud_floating_ip", floatingipByID,
					"testdata/d/hcloud_floating_ip", floatingipBySel,
				),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(floatingipByName.TFID(),
						"name", fmt.Sprintf("%s--%d", res.Name, tmplMan.RandInt)),
					resource.TestCheckResourceAttr(floatingipByName.TFID(), "home_location", res.HomeLocationName),

					resource.TestCheckResourceAttr(floatingipByID.TFID(),
						"name", fmt.Sprintf("%s--%d", res.Name, tmplMan.RandInt)),
					resource.TestCheckResourceAttr(floatingipByID.TFID(), "home_location", res.HomeLocationName),

					resource.TestCheckResourceAttr(floatingipBySel.TFID(),
						"name", fmt.Sprintf("%s--%d", res.Name, tmplMan.RandInt)),
					resource.TestCheckResourceAttr(floatingipBySel.TFID(), "home_location", res.HomeLocationName),
				),
			},
		},
	})
}

func TestAccFloatingIPDataSourceList(t *testing.T) {
	res := &floatingip.RData{
		Name: "floatingip-ds-test",
		Type: "ipv4",
		Labels: map[string]string{
			"key": strconv.Itoa(acctest.RandInt()),
		},
		HomeLocationName: teste2e.TestLocationName,
	}
	res.SetRName("floatingip-ds-test")

	floatingipBySel := &floatingip.DDataList{
		LabelSelector: fmt.Sprintf("key=${%s.labels[\"key\"]}", res.TFID()),
	}
	floatingipBySel.SetRName("floatingips_by_sel")

	allFloatingIPsSel := &floatingip.DDataList{}
	allFloatingIPsSel.SetRName("all_floatingips_sel")

	tmplMan := testtemplate.Manager{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 teste2e.PreCheck(t),
		ProtoV6ProviderFactories: teste2e.ProtoV6ProviderFactories(),
		CheckDestroy:             testsupport.CheckResourcesDestroyed(floatingip.ResourceType, floatingip.ByID(t, nil)),
		Steps: []resource.TestStep{
			{
				Config: tmplMan.Render(t,
					"testdata/r/hcloud_floating_ip", res,
				),
			},
			{
				Config: tmplMan.Render(t,
					"testdata/r/hcloud_floating_ip", res,
					"testdata/d/hcloud_floating_ips", floatingipBySel,
					"testdata/d/hcloud_floating_ips", allFloatingIPsSel,
				),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs(floatingipBySel.TFID(), "floating_ips.*",
						map[string]string{
							"name":          fmt.Sprintf("%s--%d", res.Name, tmplMan.RandInt),
							"home_location": res.HomeLocationName,
						},
					),

					resource.TestCheckTypeSetElemNestedAttrs(allFloatingIPsSel.TFID(), "floating_ips.*",
						map[string]string{
							"name":          fmt.Sprintf("%s--%d", res.Name, tmplMan.RandInt),
							"home_location": res.HomeLocationName,
						},
					),
				),
			},
		},
	})
}
