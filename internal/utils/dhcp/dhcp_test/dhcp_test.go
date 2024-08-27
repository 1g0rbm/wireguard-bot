package dhcp_test

import (
	"net"
	"testing"
	"wireguard-bot/internal/utils/dhcp"

	"github.com/stretchr/testify/require"
)

func TestDHCP(t *testing.T) {
	mask := "10.0.0.0/24"
	gateway := "10.0.0.1"
	assignedIPs := make(map[string]bool)

	_, subnet, _ := net.ParseCIDR(mask)

	t.Run("reserve and assign one ip should be success", func(t *testing.T) {
		dhcpInstance, err := dhcp.NewDHCP(mask, gateway, assignedIPs)
		require.NoError(t, err)

		reservedIP, err := dhcpInstance.Reserve()
		require.NoError(t, err)
		require.True(t, subnet.Contains(reservedIP))

		err = dhcpInstance.Assign(reservedIP)
		require.NoError(t, err)
	})

	t.Run("reserve and assign more than one ip should be success", func(t *testing.T) {
		dhcpInstance, err := dhcp.NewDHCP(mask, gateway, assignedIPs)
		require.NoError(t, err)

		firstReservedIP, err := dhcpInstance.Reserve()
		require.NoError(t, err)
		require.True(t, subnet.Contains(firstReservedIP))

		secondReservedIP, err := dhcpInstance.Reserve()
		require.NoError(t, err)
		require.True(t, subnet.Contains(secondReservedIP))

		require.False(t, firstReservedIP.Equal(secondReservedIP))

		err = dhcpInstance.Assign(firstReservedIP)
		require.NoError(t, err)

		err = dhcpInstance.Assign(secondReservedIP)
		require.NoError(t, err)
	})

	t.Run("reserve and assign not reserved ip should be failed", func(t *testing.T) {
		dhcpInstance, err := dhcp.NewDHCP(mask, gateway, assignedIPs)
		require.NoError(t, err)

		reservedIP, err := dhcpInstance.Reserve()
		require.NoError(t, err)
		require.True(t, subnet.Contains(reservedIP))

		err = dhcpInstance.Assign(reservedIP)
		require.NoError(t, err)

		err = dhcpInstance.Assign(reservedIP)
		require.Error(t, err)
	})
}
