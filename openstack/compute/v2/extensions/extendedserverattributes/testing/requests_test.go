package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/extensions/extendedserverattributes"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestServerWithUsageExt(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/servers/d650a0ce-17c3-497d-961a-43c4af80998a", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		fmt.Fprintf(w, ServerWithAttributesExtResult)
	})

	type serverAttributesExt struct {
		servers.Server
		extendedserverattributes.ServerAttributesExt
	}
	var serverWithAttributesExt serverAttributesExt

	// Extract basic fields.
	err := servers.Get(fake.ServiceClient(), "d650a0ce-17c3-497d-961a-43c4af80998a").ExtractInto(&serverWithAttributesExt)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, serverWithAttributesExt.Host, "compute01")
	th.AssertEquals(t, serverWithAttributesExt.InstanceName, "instance-00000001")
	th.AssertEquals(t, serverWithAttributesExt.HypervisorHostname, "compute01")
	th.AssertEquals(t, *serverWithAttributesExt.Userdata, "Zm9v")
	th.AssertEquals(t, *serverWithAttributesExt.ReservationID, "r-ky9gim1l")
	th.AssertEquals(t, *serverWithAttributesExt.LaunchIndex, 0)
	th.AssertEquals(t, *serverWithAttributesExt.Hostname, "test00")
	th.AssertEquals(t, *serverWithAttributesExt.RootDeviceName, "/dev/sda")
}
