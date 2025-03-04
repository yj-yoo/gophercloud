// Package v3 contains common functions for creating block storage based
// resources for use in acceptance tests. See the `*_test.go` files for
// example usages.
package v3

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/qos"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/snapshots"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumetypes"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// CreateSnapshot will create a snapshot of the specified volume.
// Snapshot will be assigned a random name and description.
func CreateSnapshot(t *testing.T, client *gophercloud.ServiceClient, volume *volumes.Volume) (*snapshots.Snapshot, error) {
	snapshotName := tools.RandomString("ACPTTEST", 16)
	snapshotDescription := tools.RandomString("ACPTTEST", 16)
	t.Logf("Attempting to create snapshot: %s", snapshotName)

	createOpts := snapshots.CreateOpts{
		VolumeID:    volume.ID,
		Name:        snapshotName,
		Description: snapshotDescription,
	}

	snapshot, err := snapshots.Create(client, createOpts).Extract()
	if err != nil {
		return snapshot, err
	}

	err = snapshots.WaitForStatus(client, snapshot.ID, "available", 60)
	if err != nil {
		return snapshot, err
	}

	snapshot, err = snapshots.Get(client, snapshot.ID).Extract()
	if err != nil {
		return snapshot, err
	}

	tools.PrintResource(t, snapshot)
	th.AssertEquals(t, snapshot.Name, snapshotName)
	th.AssertEquals(t, snapshot.VolumeID, volume.ID)

	t.Logf("Successfully created snapshot: %s", snapshot.ID)

	return snapshot, nil
}

// CreateVolume will create a volume with a random name and size of 1GB. An
// error will be returned if the volume was unable to be created.
func CreateVolume(t *testing.T, client *gophercloud.ServiceClient) (*volumes.Volume, error) {
	volumeName := tools.RandomString("ACPTTEST", 16)
	volumeDescription := tools.RandomString("ACPTTEST-DESC", 16)
	t.Logf("Attempting to create volume: %s", volumeName)

	createOpts := volumes.CreateOpts{
		Size:        1,
		Name:        volumeName,
		Description: volumeDescription,
	}

	volume, err := volumes.Create(client, createOpts).Extract()
	if err != nil {
		return volume, err
	}

	err = volumes.WaitForStatus(client, volume.ID, "available", 60)
	if err != nil {
		return volume, err
	}

	volume, err = volumes.Get(client, volume.ID).Extract()
	if err != nil {
		return volume, err
	}

	tools.PrintResource(t, volume)
	th.AssertEquals(t, volume.Name, volumeName)
	th.AssertEquals(t, volume.Description, volumeDescription)
	th.AssertEquals(t, volume.Size, 1)

	t.Logf("Successfully created volume: %s", volume.ID)

	return volume, nil
}

// CreateVolumeWithType will create a volume of the given volume type
// with a random name and size of 1GB. An error will be returned if
// the volume was unable to be created.
func CreateVolumeWithType(t *testing.T, client *gophercloud.ServiceClient, vt *volumetypes.VolumeType) (*volumes.Volume, error) {
	volumeName := tools.RandomString("ACPTTEST", 16)
	volumeDescription := tools.RandomString("ACPTTEST-DESC", 16)
	t.Logf("Attempting to create volume: %s", volumeName)

	createOpts := volumes.CreateOpts{
		Size:        1,
		Name:        volumeName,
		Description: volumeDescription,
		VolumeType:  vt.Name,
	}

	volume, err := volumes.Create(client, createOpts).Extract()
	if err != nil {
		return volume, err
	}

	err = volumes.WaitForStatus(client, volume.ID, "available", 60)
	if err != nil {
		return volume, err
	}

	volume, err = volumes.Get(client, volume.ID).Extract()
	if err != nil {
		return volume, err
	}

	tools.PrintResource(t, volume)
	th.AssertEquals(t, volume.Name, volumeName)
	th.AssertEquals(t, volume.Description, volumeDescription)
	th.AssertEquals(t, volume.Size, 1)
	th.AssertEquals(t, volume.VolumeType, vt.Name)

	t.Logf("Successfully created volume: %s", volume.ID)

	return volume, nil
}

// CreateVolumeType will create a volume type with a random name. An
// error will be returned if the volume was unable to be created.
func CreateVolumeType(t *testing.T, client *gophercloud.ServiceClient) (*volumetypes.VolumeType, error) {
	name := tools.RandomString("ACPTTEST", 16)
	description := "create_from_gophercloud"
	t.Logf("Attempting to create volume type: %s", name)

	createOpts := volumetypes.CreateOpts{
		Name:        name,
		ExtraSpecs:  map[string]string{"volume_backend_name": "fake_backend_name"},
		Description: description,
	}

	vt, err := volumetypes.Create(client, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	tools.PrintResource(t, vt)
	th.AssertEquals(t, vt.IsPublic, true)
	th.AssertEquals(t, vt.Name, name)
	th.AssertEquals(t, vt.Description, description)
	// TODO: For some reason returned extra_specs are empty even in API reference: https://developer.openstack.org/api-ref/block-storage/v3/?expanded=create-a-volume-type-detail#volume-types-types
	// "extra_specs": {}
	// th.AssertEquals(t, vt.ExtraSpecs, createOpts.ExtraSpecs)

	t.Logf("Successfully created volume type: %s", vt.ID)

	return vt, nil
}

// CreateVolumeTypeNoExtraSpecs will create a volume type with a random name and
// no extra specs. This is required to bypass cinder-scheduler filters and be able
// to create a volume with this volumeType. An error will be returned if the volume
// type was unable to be created.
func CreateVolumeTypeNoExtraSpecs(t *testing.T, client *gophercloud.ServiceClient) (*volumetypes.VolumeType, error) {
	name := tools.RandomString("ACPTTEST", 16)
	description := "create_from_gophercloud"
	t.Logf("Attempting to create volume type: %s", name)

	createOpts := volumetypes.CreateOpts{
		Name:        name,
		ExtraSpecs:  map[string]string{},
		Description: description,
	}

	vt, err := volumetypes.Create(client, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	tools.PrintResource(t, vt)
	th.AssertEquals(t, vt.IsPublic, true)
	th.AssertEquals(t, vt.Name, name)
	th.AssertEquals(t, vt.Description, description)

	t.Logf("Successfully created volume type: %s", vt.ID)

	return vt, nil
}

// CreateVolumeTypeMultiAttach will create a volume type with a random name and
// extra specs for multi-attach. An error will be returned if the volume type was
// unable to be created.
func CreateVolumeTypeMultiAttach(t *testing.T, client *gophercloud.ServiceClient) (*volumetypes.VolumeType, error) {
	name := tools.RandomString("ACPTTEST", 16)
	description := "create_from_gophercloud"
	t.Logf("Attempting to create volume type: %s", name)

	createOpts := volumetypes.CreateOpts{
		Name:        name,
		ExtraSpecs:  map[string]string{"multiattach": "<is> True"},
		Description: description,
	}

	vt, err := volumetypes.Create(client, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	tools.PrintResource(t, vt)
	th.AssertEquals(t, vt.IsPublic, true)
	th.AssertEquals(t, vt.Name, name)
	th.AssertEquals(t, vt.Description, description)
	th.AssertEquals(t, vt.ExtraSpecs["multiattach"], "<is> True")

	t.Logf("Successfully created volume type: %s", vt.ID)

	return vt, nil
}

// CreatePrivateVolumeType will create a private volume type with a random
// name and no extra specs. An error will be returned if the volume type was
// unable to be created.
func CreatePrivateVolumeType(t *testing.T, client *gophercloud.ServiceClient) (*volumetypes.VolumeType, error) {
	name := tools.RandomString("ACPTTEST", 16)
	description := "create_from_gophercloud"
	isPublic := false
	t.Logf("Attempting to create volume type: %s", name)

	createOpts := volumetypes.CreateOpts{
		Name:        name,
		ExtraSpecs:  map[string]string{},
		Description: description,
		IsPublic:    &isPublic,
	}

	vt, err := volumetypes.Create(client, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	tools.PrintResource(t, vt)
	th.AssertEquals(t, vt.IsPublic, false)
	th.AssertEquals(t, vt.Name, name)
	th.AssertEquals(t, vt.Description, description)

	t.Logf("Successfully created volume type: %s", vt.ID)

	return vt, nil
}

// DeleteSnapshot will delete a snapshot. A fatal error will occur if the
// snapshot failed to be deleted.
func DeleteSnapshot(t *testing.T, client *gophercloud.ServiceClient, snapshot *snapshots.Snapshot) {
	err := snapshots.Delete(client, snapshot.ID).ExtractErr()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			t.Logf("Snapshot %s is already deleted", snapshot.ID)
			return
		}
		t.Fatalf("Unable to delete snapshot %s: %+v", snapshot.ID, err)
	}

	// Volumes can't be deleted until their snapshots have been,
	// so block until the snapshoth as been deleted.
	err = tools.WaitFor(func() (bool, error) {
		_, err := snapshots.Get(client, snapshot.ID).Extract()
		if err != nil {
			return true, nil
		}

		return false, nil
	})
	if err != nil {
		t.Fatalf("Error waiting for snapshot to delete: %v", err)
	}

	t.Logf("Deleted snapshot: %s", snapshot.ID)
}

// DeleteVolume will delete a volume. A fatal error will occur if the volume
// failed to be deleted. This works best when used as a deferred function.
func DeleteVolume(t *testing.T, client *gophercloud.ServiceClient, volume *volumes.Volume) {
	t.Logf("Attempting to delete volume: %s", volume.ID)

	err := volumes.Delete(client, volume.ID, volumes.DeleteOpts{}).ExtractErr()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			t.Logf("Volume %s is already deleted", volume.ID)
			return
		}
		t.Fatalf("Unable to delete volume %s: %v", volume.ID, err)
	}

	// VolumeTypes can't be deleted until their volumes have been,
	// so block until the volume is deleted.
	err = tools.WaitFor(func() (bool, error) {
		_, err := volumes.Get(client, volume.ID).Extract()
		if err != nil {
			return true, nil
		}

		return false, nil
	})
	if err != nil {
		t.Fatalf("Error waiting for volume to delete: %v", err)
	}

	t.Logf("Successfully deleted volume: %s", volume.ID)
}

// DeleteVolumeType will delete a volume type. A fatal error will occur if the
// volume type failed to be deleted. This works best when used as a deferred
// function.
func DeleteVolumeType(t *testing.T, client *gophercloud.ServiceClient, vt *volumetypes.VolumeType) {
	t.Logf("Attempting to delete volume type: %s", vt.ID)

	err := volumetypes.Delete(client, vt.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete volume type %s: %v", vt.ID, err)
	}

	t.Logf("Successfully deleted volume type: %s", vt.ID)
}

// CreateQoS will create a QoS with one spec and a random name. An
// error will be returned if the volume was unable to be created.
func CreateQoS(t *testing.T, client *gophercloud.ServiceClient) (*qos.QoS, error) {
	name := tools.RandomString("ACPTTEST", 16)
	t.Logf("Attempting to create QoS: %s", name)

	createOpts := qos.CreateOpts{
		Name:     name,
		Consumer: qos.ConsumerFront,
		Specs: map[string]string{
			"read_iops_sec": "20000",
		},
	}

	qs, err := qos.Create(client, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	tools.PrintResource(t, qs)
	th.AssertEquals(t, qs.Consumer, "front-end")
	th.AssertEquals(t, qs.Name, name)
	th.AssertDeepEquals(t, qs.Specs, createOpts.Specs)

	t.Logf("Successfully created QoS: %s", qs.ID)

	return qs, nil
}

// DeleteQoS will delete a QoS. A fatal error will occur if the QoS
// failed to be deleted. This works best when used as a deferred function.
func DeleteQoS(t *testing.T, client *gophercloud.ServiceClient, qs *qos.QoS) {
	t.Logf("Attempting to delete QoS: %s", qs.ID)

	deleteOpts := qos.DeleteOpts{
		Force: true,
	}

	err := qos.Delete(client, qs.ID, deleteOpts).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete QoS %s: %v", qs.ID, err)
	}

	t.Logf("Successfully deleted QoS: %s", qs.ID)
}
