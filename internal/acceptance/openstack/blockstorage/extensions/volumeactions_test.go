//go:build acceptance || blockstorage
// +build acceptance blockstorage

package extensions

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	blockstorage "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/blockstorage/v2"
	blockstorageV3 "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/blockstorage/v3"
	compute "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/compute/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/volumes"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestVolumeActionsUploadImageDestroy(t *testing.T) {
	blockClient, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	computeClient, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	volume, err := blockstorage.CreateVolume(t, blockClient)
	th.AssertNoErr(t, err)
	defer blockstorage.DeleteVolume(t, blockClient, volume)

	volumeImage, err := CreateUploadImage(t, blockClient, volume)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, volumeImage)

	err = DeleteUploadedImage(t, computeClient, volumeImage.ImageID)
	th.AssertNoErr(t, err)
}

func TestVolumeActionsAttachCreateDestroy(t *testing.T) {
	blockClient, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	computeClient, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	server, err := compute.CreateServer(t, computeClient)
	th.AssertNoErr(t, err)
	defer compute.DeleteServer(t, computeClient, server)

	volume, err := blockstorage.CreateVolume(t, blockClient)
	th.AssertNoErr(t, err)
	defer blockstorage.DeleteVolume(t, blockClient, volume)

	err = CreateVolumeAttach(t, blockClient, volume, server)
	th.AssertNoErr(t, err)

	newVolume, err := volumes.Get(blockClient, volume.ID).Extract()
	th.AssertNoErr(t, err)

	DeleteVolumeAttach(t, blockClient, newVolume)
}

func TestVolumeActionsReserveUnreserve(t *testing.T) {
	client, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	volume, err := blockstorage.CreateVolume(t, client)
	th.AssertNoErr(t, err)
	defer blockstorage.DeleteVolume(t, client, volume)

	err = CreateVolumeReserve(t, client, volume)
	th.AssertNoErr(t, err)
	defer DeleteVolumeReserve(t, client, volume)
}

func TestVolumeActionsExtendSize(t *testing.T) {
	blockClient, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	volume, err := blockstorage.CreateVolume(t, blockClient)
	th.AssertNoErr(t, err)
	defer blockstorage.DeleteVolume(t, blockClient, volume)

	tools.PrintResource(t, volume)

	err = ExtendVolumeSize(t, blockClient, volume)
	th.AssertNoErr(t, err)

	newVolume, err := volumes.Get(blockClient, volume.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newVolume)
}

func TestVolumeActionsImageMetadata(t *testing.T) {
	blockClient, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	volume, err := blockstorage.CreateVolume(t, blockClient)
	th.AssertNoErr(t, err)
	defer blockstorage.DeleteVolume(t, blockClient, volume)

	err = SetImageMetadata(t, blockClient, volume)
	th.AssertNoErr(t, err)
}

func TestVolumeActionsSetBootable(t *testing.T) {
	blockClient, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	volume, err := blockstorage.CreateVolume(t, blockClient)
	th.AssertNoErr(t, err)
	defer blockstorage.DeleteVolume(t, blockClient, volume)

	err = SetBootable(t, blockClient, volume)
	th.AssertNoErr(t, err)
}

func TestVolumeActionsChangeType(t *testing.T) {
	//	clients.RequireAdmin(t)

	client, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	volumeType1, err := blockstorageV3.CreateVolumeTypeNoExtraSpecs(t, client)
	th.AssertNoErr(t, err)
	defer blockstorageV3.DeleteVolumeType(t, client, volumeType1)

	volumeType2, err := blockstorageV3.CreateVolumeTypeNoExtraSpecs(t, client)
	th.AssertNoErr(t, err)
	defer blockstorageV3.DeleteVolumeType(t, client, volumeType2)

	volume, err := blockstorageV3.CreateVolumeWithType(t, client, volumeType1)
	th.AssertNoErr(t, err)
	defer blockstorageV3.DeleteVolume(t, client, volume)

	tools.PrintResource(t, volume)

	err = ChangeVolumeType(t, client, volume, volumeType2)
	th.AssertNoErr(t, err)

	newVolume, err := volumes.Get(client, volume.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, newVolume.VolumeType, volumeType2.Name)

	tools.PrintResource(t, newVolume)
}

func TestVolumeActionsResetStatus(t *testing.T) {
	client, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	volume, err := blockstorageV3.CreateVolume(t, client)
	th.AssertNoErr(t, err)
	defer blockstorageV3.DeleteVolume(t, client, volume)

	tools.PrintResource(t, volume)

	err = ResetVolumeStatus(t, client, volume, "error")
	th.AssertNoErr(t, err)

	err = ResetVolumeStatus(t, client, volume, "available")
	th.AssertNoErr(t, err)
}

func TestVolumeActionsReImage(t *testing.T) {
	clients.SkipReleasesBelow(t, "stable/yoga")

	choices, err := clients.AcceptanceTestChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	blockClient, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)
	blockClient.Microversion = "3.68"

	volume, err := blockstorage.CreateVolume(t, blockClient)
	th.AssertNoErr(t, err)
	defer blockstorage.DeleteVolume(t, blockClient, volume)

	err = ReImage(t, blockClient, volume, choices.ImageID)
	th.AssertNoErr(t, err)
}

// Note(jtopjian): I plan to work on this at some point, but it requires
// setting up a server with iscsi utils.
/*
func TestVolumeConns(t *testing.T) {
    client, err := newClient()
    th.AssertNoErr(t, err)

    t.Logf("Creating volume")
    cv, err := volumes.Create(client, &volumes.CreateOpts{
        Size: 1,
        Name: "blockv2-volume",
    }).Extract()
    th.AssertNoErr(t, err)

    defer func() {
        err = volumes.WaitForStatus(client, cv.ID, "available", 60)
        th.AssertNoErr(t, err)

        t.Logf("Deleting volume")
        err = volumes.Delete(client, cv.ID, volumes.DeleteOpts{}).ExtractErr()
        th.AssertNoErr(t, err)
    }()

    err = volumes.WaitForStatus(client, cv.ID, "available", 60)
    th.AssertNoErr(t, err)

    connOpts := &volumeactions.ConnectorOpts{
        IP:        "127.0.0.1",
        Host:      "stack",
        Initiator: "iqn.1994-05.com.redhat:17cf566367d2",
        Multipath: false,
        Platform:  "x86_64",
        OSType:    "linux2",
    }

    t.Logf("Initializing connection")
    _, err = volumeactions.InitializeConnection(client, cv.ID, connOpts).Extract()
    th.AssertNoErr(t, err)

    t.Logf("Terminating connection")
    err = volumeactions.TerminateConnection(client, cv.ID, connOpts).ExtractErr()
    th.AssertNoErr(t, err)
}
*/
