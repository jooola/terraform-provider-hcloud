package volume_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/terraform-provider-hcloud/internal/server"
	"github.com/hetznercloud/terraform-provider-hcloud/internal/teste2e"
	"github.com/hetznercloud/terraform-provider-hcloud/internal/testsupport"
	"github.com/hetznercloud/terraform-provider-hcloud/internal/testtemplate"
	"github.com/hetznercloud/terraform-provider-hcloud/internal/volume"
)

func TestAccVolumeResource(t *testing.T) {
	var vol hcloud.Volume

	res := VolumeRData()
	resRenamed := &volume.RData{
		Name:         res.Name + "-renamed",
		LocationName: teste2e.TestLocationName,
		Size:         10,
		Labels: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}
	resRenamed.SetRName(res.RName())

	tmplMan := testtemplate.Manager{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 teste2e.PreCheck(t),
		ProtoV6ProviderFactories: teste2e.ProtoV6ProviderFactories(),
		CheckDestroy:             testsupport.CheckResourcesDestroyed(volume.ResourceType, volume.ByID(t, &vol)),
		Steps: []resource.TestStep{
			{
				// Create a new Volume using the required values
				// only.
				Config: tmplMan.Render(t, "testdata/r/hcloud_volume", res),
				Check: resource.ComposeTestCheckFunc(
					testsupport.CheckResourceExists(res.TFID(), volume.ByID(t, &vol)),
					resource.TestCheckResourceAttr(res.TFID(), "name",
						fmt.Sprintf("basic-volume--%d", tmplMan.RandInt)),
					resource.TestCheckResourceAttr(res.TFID(), "size", "10"),
					resource.TestCheckResourceAttr(res.TFID(), "location", res.LocationName),
				),
			},
			{
				// Try to import the newly created volume
				ResourceName:      res.TFID(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Update the Volume created in the previous step by
				// setting all optional fields and renaming the volume.
				Config: tmplMan.Render(t,
					"testdata/r/hcloud_volume", resRenamed,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resRenamed.TFID(), "name",
						fmt.Sprintf("basic-volume-renamed--%d", tmplMan.RandInt)),
					resource.TestCheckResourceAttr(resRenamed.TFID(), "size", "10"),
					resource.TestCheckResourceAttr(resRenamed.TFID(), "location", resRenamed.LocationName),
					resource.TestCheckResourceAttr(resRenamed.TFID(), "labels.key1", "value1"),
					resource.TestCheckResourceAttr(resRenamed.TFID(), "labels.key2", "value2"),
				),
			},
		},
	})
}

func TestAccVolumeResource_Resize(t *testing.T) {
	var vol hcloud.Volume

	res := VolumeRData()
	res.Name = "resized-volume"
	resResized := &volume.RData{
		Name:         res.Name,
		LocationName: teste2e.TestLocationName,
		Size:         25,
		Labels: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}
	resResized.SetRName(res.RName())

	tmplMan := testtemplate.Manager{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 teste2e.PreCheck(t),
		ProtoV6ProviderFactories: teste2e.ProtoV6ProviderFactories(),
		CheckDestroy:             testsupport.CheckResourcesDestroyed(volume.ResourceType, volume.ByID(t, &vol)),
		Steps: []resource.TestStep{
			{
				// Create a new Volume using the required values
				// only.
				Config: tmplMan.Render(t, "testdata/r/hcloud_volume", res),
				Check: resource.ComposeTestCheckFunc(
					testsupport.CheckResourceExists(res.TFID(), volume.ByID(t, &vol)),
					resource.TestCheckResourceAttr(res.TFID(), "name",
						fmt.Sprintf("resized-volume--%d", tmplMan.RandInt)),
					resource.TestCheckResourceAttr(res.TFID(), "size", "10"),
					resource.TestCheckResourceAttr(res.TFID(), "location", res.LocationName),
				),
			},
			{
				// Update the Volume created in the previous step by
				// changing the size.
				Config: tmplMan.Render(t,
					"testdata/r/hcloud_volume", resResized,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resResized.TFID(), "name",
						fmt.Sprintf("resized-volume--%d", tmplMan.RandInt)),
					resource.TestCheckResourceAttr(resResized.TFID(), "size", "25"),
					resource.TestCheckResourceAttr(resResized.TFID(), "location", resResized.LocationName),
					resource.TestCheckResourceAttr(resResized.TFID(), "labels.key1", "value1"),
					resource.TestCheckResourceAttr(resResized.TFID(), "labels.key2", "value2"),
				),
			},
		},
	})
}

func TestAccVolumeResource_WithServer(t *testing.T) {
	var vol hcloud.Volume
	tmplMan := testtemplate.Manager{}
	resServer1 := &server.RData{
		Name:         "some-server",
		Type:         teste2e.TestServerType,
		Image:        teste2e.TestImage,
		LocationName: teste2e.TestLocationName,
	}
	resServer1.SetRName("some-server")

	resServer2 := &server.RData{
		Name:         "another-server",
		Type:         teste2e.TestServerType,
		Image:        teste2e.TestImage,
		LocationName: teste2e.TestLocationName,
	}
	resServer2.SetRName("another-server")

	res := VolumeRData()
	res.Name = "volume-with-server"
	res.LocationName = ""
	res.ServerID = resServer1.TFID() + ".id"

	resAnotherServer := VolumeRData()
	resAnotherServer.Name = "volume-with-server"
	resAnotherServer.LocationName = ""
	resAnotherServer.ServerID = resServer2.TFID() + ".id"
	resAnotherServer.SetRName(res.RName())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 teste2e.PreCheck(t),
		ProtoV6ProviderFactories: teste2e.ProtoV6ProviderFactories(),
		CheckDestroy:             testsupport.CheckResourcesDestroyed(volume.ResourceType, volume.ByID(t, &vol)),
		Steps: []resource.TestStep{
			{
				// Create a new Volume using the required values
				// only.
				Config: tmplMan.Render(t,
					"testdata/r/hcloud_server", resServer1,
					"testdata/r/hcloud_server", resServer2,
					"testdata/r/hcloud_volume", res),
				Check: resource.ComposeTestCheckFunc(
					testsupport.CheckResourceExists(res.TFID(), volume.ByID(t, &vol)),
					resource.TestCheckResourceAttr(res.TFID(), "name",
						fmt.Sprintf("volume-with-server--%d", tmplMan.RandInt)),
					resource.TestCheckResourceAttr(res.TFID(), "size", "10"),
					resource.TestCheckResourceAttr(res.TFID(), "location", resServer1.LocationName),
				),
			},
			{
				// Try to import the newly created volume
				ResourceName:      res.TFID(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Update the Volume created in the previous step by
				// changing the attached server.
				Config: tmplMan.Render(t,
					"testdata/r/hcloud_server", resServer1,
					"testdata/r/hcloud_server", resServer2,
					"testdata/r/hcloud_volume", resAnotherServer),
				Check: resource.ComposeTestCheckFunc(
					testsupport.CheckResourceExists(res.TFID(), volume.ByID(t, &vol)),
					resource.TestCheckResourceAttr(res.TFID(), "name",
						fmt.Sprintf("volume-with-server--%d", tmplMan.RandInt)),
					resource.TestCheckResourceAttr(res.TFID(), "size", "10"),
					resource.TestCheckResourceAttr(res.TFID(), "location", resServer2.LocationName),
				),
			},
		},
	})
}

func TestAccVolumeResource_WithServerMultipleVolumes(t *testing.T) {
	var vol, vol2 hcloud.Volume
	tmplMan := testtemplate.Manager{}
	resServer1 := &server.RData{
		Name:         "some-server",
		Type:         teste2e.TestServerType,
		Image:        teste2e.TestImage,
		LocationName: teste2e.TestLocationName,
	}
	resServer1.SetRName("some-server")

	res := VolumeRData()
	res.Name = "volume-with-server"
	res.LocationName = ""
	res.ServerID = resServer1.TFID() + ".id"
	res.SetRName("first-volume")

	resAnotherVolume := &volume.RData{
		Name:         "volume-with-server-2",
		LocationName: "",
		Size:         10,
		ServerID:     resServer1.TFID() + ".id",
	}
	resAnotherVolume.SetRName("another-volume")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 teste2e.PreCheck(t),
		ProtoV6ProviderFactories: teste2e.ProtoV6ProviderFactories(),
		CheckDestroy:             testsupport.CheckResourcesDestroyed(volume.ResourceType, volume.ByID(t, &vol)),
		Steps: []resource.TestStep{
			{
				// Create a new Volume using the required values
				// only.
				Config: tmplMan.Render(t,
					"testdata/r/hcloud_server", resServer1,
					"testdata/r/hcloud_volume", res,
					"testdata/r/hcloud_volume", resAnotherVolume),
				Check: resource.ComposeTestCheckFunc(
					testsupport.CheckResourceExists(res.TFID(), volume.ByID(t, &vol)),
					testsupport.CheckResourceExists(resAnotherVolume.TFID(), volume.ByID(t, &vol2)),
					resource.TestCheckResourceAttr(res.TFID(), "name",
						fmt.Sprintf("volume-with-server--%d", tmplMan.RandInt)),
					resource.TestCheckResourceAttr(res.TFID(), "size", "10"),
					resource.TestCheckResourceAttr(res.TFID(), "location", resServer1.LocationName),
				),
			},
		},
	})
}

func TestAccVolumeResource_Protection(t *testing.T) {
	var (
		vol hcloud.Volume

		res = &volume.RData{
			Name:             "basic-volume",
			LocationName:     teste2e.TestLocationName,
			Size:             10,
			DeleteProtection: true,
		}

		updateProtection = func(d *volume.RData, protection bool) *volume.RData {
			d.DeleteProtection = protection
			return d
		}
	)

	tmplMan := testtemplate.Manager{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 teste2e.PreCheck(t),
		ProtoV6ProviderFactories: teste2e.ProtoV6ProviderFactories(),
		CheckDestroy:             testsupport.CheckResourcesDestroyed(volume.ResourceType, volume.ByID(t, &vol)),
		Steps: []resource.TestStep{
			{
				// Create a new Volume using the required values
				// only.
				Config: tmplMan.Render(t,
					"testdata/r/hcloud_volume", res,
				),
				Check: resource.ComposeTestCheckFunc(
					testsupport.CheckResourceExists(res.TFID(), volume.ByID(t, &vol)),
					resource.TestCheckResourceAttr(res.TFID(), "name",
						fmt.Sprintf("basic-volume--%d", tmplMan.RandInt)),
					resource.TestCheckResourceAttr(res.TFID(), "size", fmt.Sprintf("%d", res.Size)),
					resource.TestCheckResourceAttr(res.TFID(), "location", res.LocationName),
					resource.TestCheckResourceAttr(res.TFID(), "delete_protection", fmt.Sprintf("%t", res.DeleteProtection)),
				),
			},
			{
				// Update the Volume created in the previous step by
				// setting all optional fields and renaming the volume.
				Config: tmplMan.Render(t,
					"testdata/r/hcloud_volume", updateProtection(res, false),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(res.TFID(), "delete_protection", fmt.Sprintf("%t", res.DeleteProtection)),
				),
			},
		},
	})
}
