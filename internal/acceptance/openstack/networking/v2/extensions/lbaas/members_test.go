//go:build acceptance || networking || lbaas || member
// +build acceptance networking lbaas member

package lbaas

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/lbaas/members"
)

func TestMembersList(t *testing.T) {
	t.Skip("Neutron LBaaS was replaced by Octavia and the API will be removed in a future release")
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	allPages, err := members.List(client, members.ListOpts{}).AllPages()
	if err != nil {
		t.Fatalf("Unable to list members: %v", err)
	}

	allMembers, err := members.ExtractMembers(allPages)
	if err != nil {
		t.Fatalf("Unable to extract members: %v", err)
	}

	for _, member := range allMembers {
		tools.PrintResource(t, member)
	}
}

func TestMembersCRUD(t *testing.T) {
	t.Skip("Neutron LBaaS was replaced by Octavia and the API will be removed in a future release")
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	network, err := networking.CreateNetwork(t, client)
	if err != nil {
		t.Fatalf("Unable to create network: %v", err)
	}
	defer networking.DeleteNetwork(t, client, network.ID)

	subnet, err := networking.CreateSubnet(t, client, network.ID)
	if err != nil {
		t.Fatalf("Unable to create subnet: %v", err)
	}
	defer networking.DeleteSubnet(t, client, subnet.ID)

	pool, err := CreatePool(t, client, subnet.ID)
	if err != nil {
		t.Fatalf("Unable to create pool: %v", err)
	}
	defer DeletePool(t, client, pool.ID)

	member, err := CreateMember(t, client, pool.ID)
	if err != nil {
		t.Fatalf("Unable to create member: %v", err)
	}
	defer DeleteMember(t, client, member.ID)

	tools.PrintResource(t, member)

	updateOpts := members.UpdateOpts{
		AdminStateUp: gophercloud.Enabled,
	}

	_, err = members.Update(client, member.ID, updateOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to update member: %v", err)
	}

	newMember, err := members.Get(client, member.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get member: %v", err)
	}

	tools.PrintResource(t, newMember)
}
