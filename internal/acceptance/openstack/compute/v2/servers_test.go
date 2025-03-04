//go:build acceptance || compute || servers
// +build acceptance compute servers

package v2

import (
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networks "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/extensions/attachinterfaces"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/extensions/availabilityzones"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/extensions/extendedserverattributes"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/extensions/extendedstatus"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/extensions/lockunlock"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/extensions/pauseunpause"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/extensions/serverusage"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/extensions/suspendresume"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/extensions/tags"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestServersCreateDestroy(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	choices, err := clients.AcceptanceTestChoicesFromEnv()
	th.AssertNoErr(t, err)

	server, err := CreateServer(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	allPages, err := servers.List(client, servers.ListOpts{}).AllPages()
	th.AssertNoErr(t, err)

	allServers, err := servers.ExtractServers(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, s := range allServers {
		tools.PrintResource(t, server)

		if s.ID == server.ID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)

	allAddressPages, err := servers.ListAddresses(client, server.ID).AllPages()
	th.AssertNoErr(t, err)

	allAddresses, err := servers.ExtractAddresses(allAddressPages)
	th.AssertNoErr(t, err)

	for network, address := range allAddresses {
		t.Logf("Addresses on %s: %+v", network, address)
	}

	allInterfacePages, err := attachinterfaces.List(client, server.ID).AllPages()
	th.AssertNoErr(t, err)

	allInterfaces, err := attachinterfaces.ExtractInterfaces(allInterfacePages)
	th.AssertNoErr(t, err)

	for _, iface := range allInterfaces {
		t.Logf("Interfaces: %+v", iface)
	}

	allNetworkAddressPages, err := servers.ListAddressesByNetwork(client, server.ID, choices.NetworkName).AllPages()
	th.AssertNoErr(t, err)

	allNetworkAddresses, err := servers.ExtractNetworkAddresses(allNetworkAddressPages)
	th.AssertNoErr(t, err)

	t.Logf("Addresses on %s:", choices.NetworkName)
	for _, address := range allNetworkAddresses {
		t.Logf("%+v", address)
	}
}

func TestServersWithExtensionsCreateDestroy(t *testing.T) {
	clients.RequireLong(t)

	var extendedServer struct {
		servers.Server
		availabilityzones.ServerAvailabilityZoneExt
		extendedstatus.ServerExtendedStatusExt
		serverusage.UsageExt
	}

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	server, err := CreateServer(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	err = servers.Get(client, server.ID).ExtractInto(&extendedServer)
	th.AssertNoErr(t, err)
	tools.PrintResource(t, extendedServer)

	th.AssertEquals(t, extendedServer.AvailabilityZone, "nova")
	th.AssertEquals(t, int(extendedServer.PowerState), extendedstatus.RUNNING)
	th.AssertEquals(t, extendedServer.TaskState, "")
	th.AssertEquals(t, extendedServer.VmState, "active")
	th.AssertEquals(t, extendedServer.LaunchedAt.IsZero(), false)
	th.AssertEquals(t, extendedServer.TerminatedAt.IsZero(), true)
}

func TestServersWithoutImageRef(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	server, err := CreateServerWithoutImageRef(t, client)
	if err != nil {
		if err400, ok := err.(*gophercloud.ErrUnexpectedResponseCode); ok {
			if !strings.Contains(string(err400.Body), "Missing imageRef attribute") {
				defer DeleteServer(t, client, server)
			}
		}
	}
}

func TestServersUpdate(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	server, err := CreateServer(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	alternateName := tools.RandomString("ACPTTEST", 16)
	for alternateName == server.Name {
		alternateName = tools.RandomString("ACPTTEST", 16)
	}

	t.Logf("Attempting to rename the server to %s.", alternateName)

	updateOpts := servers.UpdateOpts{
		Name: alternateName,
	}

	updated, err := servers.Update(client, server.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, updated.ID, server.ID)

	err = tools.WaitFor(func() (bool, error) {
		latest, err := servers.Get(client, updated.ID).Extract()
		if err != nil {
			return false, err
		}

		return latest.Name == alternateName, nil
	})
}

func TestServersMetadata(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	server, err := CreateServer(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	tools.PrintResource(t, server)

	metadata, err := servers.UpdateMetadata(client, server.ID, servers.MetadataOpts{
		"foo":  "bar",
		"this": "that",
	}).Extract()
	th.AssertNoErr(t, err)
	t.Logf("UpdateMetadata result: %+v\n", metadata)

	server, err = servers.Get(client, server.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, server)

	expectedMetadata := map[string]string{
		"abc":  "def",
		"foo":  "bar",
		"this": "that",
	}
	th.AssertDeepEquals(t, expectedMetadata, server.Metadata)

	err = servers.DeleteMetadatum(client, server.ID, "foo").ExtractErr()
	th.AssertNoErr(t, err)

	server, err = servers.Get(client, server.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, server)

	expectedMetadata = map[string]string{
		"abc":  "def",
		"this": "that",
	}
	th.AssertDeepEquals(t, expectedMetadata, server.Metadata)

	metadata, err = servers.CreateMetadatum(client, server.ID, servers.MetadatumOpts{
		"foo": "baz",
	}).Extract()
	th.AssertNoErr(t, err)
	t.Logf("CreateMetadatum result: %+v\n", metadata)

	server, err = servers.Get(client, server.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, server)

	expectedMetadata = map[string]string{
		"abc":  "def",
		"this": "that",
		"foo":  "baz",
	}
	th.AssertDeepEquals(t, expectedMetadata, server.Metadata)

	metadata, err = servers.Metadatum(client, server.ID, "foo").Extract()
	th.AssertNoErr(t, err)
	t.Logf("Metadatum result: %+v\n", metadata)
	th.AssertEquals(t, "baz", metadata["foo"])

	metadata, err = servers.Metadata(client, server.ID).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Metadata result: %+v\n", metadata)

	th.AssertDeepEquals(t, expectedMetadata, metadata)

	metadata, err = servers.ResetMetadata(client, server.ID, servers.MetadataOpts{}).Extract()
	th.AssertNoErr(t, err)
	t.Logf("ResetMetadata result: %+v\n", metadata)
	th.AssertDeepEquals(t, map[string]string{}, metadata)
}

func TestServersActionChangeAdminPassword(t *testing.T) {
	clients.RequireLong(t)
	clients.RequireGuestAgent(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	server, err := CreateServer(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	randomPassword := tools.MakeNewPassword(server.AdminPass)
	res := servers.ChangeAdminPassword(client, server.ID, randomPassword)
	th.AssertNoErr(t, res.Err)

	if err = WaitForComputeStatus(client, server, "PASSWORD"); err != nil {
		t.Fatal(err)
	}

	if err = WaitForComputeStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}
}

func TestServersActionReboot(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	server, err := CreateServer(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	rebootOpts := servers.RebootOpts{
		Type: servers.SoftReboot,
	}

	t.Logf("Attempting reboot of server %s", server.ID)
	res := servers.Reboot(client, server.ID, rebootOpts)
	th.AssertNoErr(t, res.Err)

	if err = WaitForComputeStatus(client, server, "REBOOT"); err != nil {
		t.Fatal(err)
	}

	if err = WaitForComputeStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}
}

func TestServersActionRebuild(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	choices, err := clients.AcceptanceTestChoicesFromEnv()
	th.AssertNoErr(t, err)

	server, err := CreateServer(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	t.Logf("Attempting to rebuild server %s", server.ID)

	rebuildOpts := servers.RebuildOpts{
		Name:      tools.RandomString("ACPTTEST", 16),
		AdminPass: tools.MakeNewPassword(server.AdminPass),
		ImageRef:  choices.ImageID,
	}

	rebuilt, err := servers.Rebuild(client, server.ID, rebuildOpts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, rebuilt.ID, server.ID)

	if err = WaitForComputeStatus(client, rebuilt, "REBUILD"); err != nil {
		t.Fatal(err)
	}

	if err = WaitForComputeStatus(client, rebuilt, "ACTIVE"); err != nil {
		t.Fatal(err)
	}
}

func TestServersActionResizeConfirm(t *testing.T) {
	clients.RequireLong(t)

	choices, err := clients.AcceptanceTestChoicesFromEnv()
	th.AssertNoErr(t, err)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	server, err := CreateServer(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	t.Logf("Attempting to resize server %s", server.ID)
	ResizeServer(t, client, server)

	t.Logf("Attempting to confirm resize for server %s", server.ID)
	if res := servers.ConfirmResize(client, server.ID); res.Err != nil {
		t.Fatal(res.Err)
	}

	if err = WaitForComputeStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}

	server, err = servers.Get(client, server.ID).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, server.Flavor["id"], choices.FlavorIDResize)
}

func TestServersActionResizeRevert(t *testing.T) {
	clients.RequireLong(t)

	choices, err := clients.AcceptanceTestChoicesFromEnv()
	th.AssertNoErr(t, err)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	server, err := CreateServer(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	t.Logf("Attempting to resize server %s", server.ID)
	ResizeServer(t, client, server)

	t.Logf("Attempting to revert resize for server %s", server.ID)
	if res := servers.RevertResize(client, server.ID); res.Err != nil {
		t.Fatal(res.Err)
	}

	if err = WaitForComputeStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}

	server, err = servers.Get(client, server.ID).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, server.Flavor["id"], choices.FlavorID)
}

func TestServersActionPause(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	server, err := CreateServer(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	t.Logf("Attempting to pause server %s", server.ID)
	err = pauseunpause.Pause(client, server.ID).ExtractErr()
	th.AssertNoErr(t, err)

	err = WaitForComputeStatus(client, server, "PAUSED")
	th.AssertNoErr(t, err)

	err = pauseunpause.Unpause(client, server.ID).ExtractErr()
	th.AssertNoErr(t, err)

	err = WaitForComputeStatus(client, server, "ACTIVE")
	th.AssertNoErr(t, err)
}

func TestServersActionSuspend(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	server, err := CreateServer(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	t.Logf("Attempting to suspend server %s", server.ID)
	err = suspendresume.Suspend(client, server.ID).ExtractErr()
	th.AssertNoErr(t, err)

	err = WaitForComputeStatus(client, server, "SUSPENDED")
	th.AssertNoErr(t, err)

	err = suspendresume.Resume(client, server.ID).ExtractErr()
	th.AssertNoErr(t, err)

	err = WaitForComputeStatus(client, server, "ACTIVE")
	th.AssertNoErr(t, err)
}

func TestServersActionLock(t *testing.T) {
	clients.RequireLong(t)
	clients.RequireNonAdmin(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	server, err := CreateServer(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	t.Logf("Attempting to Lock server %s", server.ID)
	err = lockunlock.Lock(client, server.ID).ExtractErr()
	th.AssertNoErr(t, err)

	t.Logf("Attempting to delete locked server %s", server.ID)
	err = servers.Delete(client, server.ID).ExtractErr()
	th.AssertEquals(t, err != nil, true)

	t.Logf("Attempting to unlock server %s", server.ID)
	err = lockunlock.Unlock(client, server.ID).ExtractErr()
	th.AssertNoErr(t, err)

	err = WaitForComputeStatus(client, server, "ACTIVE")
	th.AssertNoErr(t, err)
}

func TestServersConsoleOutput(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	server, err := CreateServer(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	outputOpts := &servers.ShowConsoleOutputOpts{
		Length: 4,
	}
	output, err := servers.ShowConsoleOutput(client, server.ID, outputOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, output)
}

func TestServersTags(t *testing.T) {
	clients.RequireLong(t)

	choices, err := clients.AcceptanceTestChoicesFromEnv()
	th.AssertNoErr(t, err)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)
	client.Microversion = "2.52"

	networkClient, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	networkID, err := networks.IDFromName(networkClient, choices.NetworkName)
	th.AssertNoErr(t, err)

	// Create server with tags.
	server, err := CreateServerWithTags(t, client, networkID)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	// All the following calls should work with "2.26" microversion.
	client.Microversion = "2.26"

	// Check server tags in body.
	serverWithTags, err := servers.Get(client, server.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, []string{"tag1", "tag2"}, *serverWithTags.Tags)

	// Check all tags.
	allTags, err := tags.List(client, server.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, []string{"tag1", "tag2"}, allTags)

	// Check single tag.
	exists, err := tags.Check(client, server.ID, "tag2").Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, true, exists)

	// Add new tag.
	newTags, err := tags.ReplaceAll(client, server.ID, tags.ReplaceAllOpts{Tags: []string{"tag3", "tag4"}}).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, []string{"tag3", "tag4"}, newTags)

	// Add new single tag.
	err = tags.Add(client, server.ID, "tag5").ExtractErr()
	th.AssertNoErr(t, err)

	// Check current tags.
	newAllTags, err := tags.List(client, server.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, []string{"tag3", "tag4", "tag5"}, newAllTags)

	// Remove single tag.
	err = tags.Delete(client, server.ID, "tag4").ExtractErr()
	th.AssertNoErr(t, err)

	// Check that tag doesn't exist anymore.
	exists, err = tags.Check(client, server.ID, "tag4").Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, false, exists)

	// Remove all tags.
	err = tags.DeleteAll(client, server.ID).ExtractErr()
	th.AssertNoErr(t, err)

	// Check that there are no more tags.
	currentTags, err := tags.List(client, server.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 0, len(currentTags))
}

func TestServersWithExtendedAttributesCreateDestroy(t *testing.T) {
	clients.RequireLong(t)
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)
	client.Microversion = "2.3"

	server, err := CreateServer(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	type serverAttributesExt struct {
		servers.Server
		extendedserverattributes.ServerAttributesExt
	}
	var serverWithAttributesExt serverAttributesExt

	err = servers.Get(client, server.ID).ExtractInto(&serverWithAttributesExt)
	th.AssertNoErr(t, err)

	t.Logf("Server With Extended Attributes: %#v", serverWithAttributesExt)

	th.AssertEquals(t, *serverWithAttributesExt.ReservationID != "", true)
	th.AssertEquals(t, *serverWithAttributesExt.LaunchIndex, 0)
	th.AssertEquals(t, *serverWithAttributesExt.RAMDiskID == "", true)
	th.AssertEquals(t, *serverWithAttributesExt.KernelID == "", true)
	th.AssertEquals(t, *serverWithAttributesExt.Hostname != "", true)
	th.AssertEquals(t, *serverWithAttributesExt.RootDeviceName != "", true)
	th.AssertEquals(t, serverWithAttributesExt.Userdata == nil, true)
}

func TestServerNoNetworkCreateDestroy(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	choices, err := clients.AcceptanceTestChoicesFromEnv()
	th.AssertNoErr(t, err)

	client.Microversion = "2.37"

	server, err := CreateServerNoNetwork(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	allPages, err := servers.List(client, servers.ListOpts{}).AllPages()
	th.AssertNoErr(t, err)

	allServers, err := servers.ExtractServers(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, s := range allServers {
		tools.PrintResource(t, server)

		if s.ID == server.ID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)

	allAddressPages, err := servers.ListAddresses(client, server.ID).AllPages()
	th.AssertNoErr(t, err)

	allAddresses, err := servers.ExtractAddresses(allAddressPages)
	th.AssertNoErr(t, err)

	for network, address := range allAddresses {
		t.Logf("Addresses on %s: %+v", network, address)
	}

	allInterfacePages, err := attachinterfaces.List(client, server.ID).AllPages()
	th.AssertNoErr(t, err)

	allInterfaces, err := attachinterfaces.ExtractInterfaces(allInterfacePages)
	th.AssertNoErr(t, err)

	for _, iface := range allInterfaces {
		t.Logf("Interfaces: %+v", iface)
	}

	_, err = servers.ListAddressesByNetwork(client, server.ID, choices.NetworkName).AllPages()
	if err == nil {
		t.Fatalf("Instance must not be a member of specified network")
	}
}
