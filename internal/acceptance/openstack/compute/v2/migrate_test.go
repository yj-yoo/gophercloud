//go:build acceptance || compute || servers
// +build acceptance compute servers

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/extensions/migrate"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestMigrate(t *testing.T) {
	clients.RequireLong(t)
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	server, err := CreateServer(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	t.Logf("Attempting to migrate server %s", server.ID)

	err = migrate.Migrate(client, server.ID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestLiveMigrate(t *testing.T) {
	clients.RequireLong(t)
	clients.RequireAdmin(t)
	clients.RequireLiveMigration(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	server, err := CreateServer(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	t.Logf("Attempting to migrate server %s", server.ID)

	blockMigration := false
	diskOverCommit := false

	liveMigrateOpts := migrate.LiveMigrateOpts{
		BlockMigration: &blockMigration,
		DiskOverCommit: &diskOverCommit,
	}

	err = migrate.LiveMigrate(client, server.ID, liveMigrateOpts).ExtractErr()
	th.AssertNoErr(t, err)
}
