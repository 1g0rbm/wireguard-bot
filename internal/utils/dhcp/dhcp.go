package dhcp

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

type DHCP struct {
	gatewayIP net.IP
	subnet    *net.IPNet

	mu          sync.Mutex
	assignedIPs map[string]bool
	reservedIPs map[string]bool
}

func NewDHCP(mask, gatewayIP string, assignedIPs map[string]bool) (*DHCP, error) {
	parsedGatewayIP := net.ParseIP(gatewayIP)
	if parsedGatewayIP == nil {
		return nil, errors.New("invalid gateway ip")
	}

	_, subnet, err := net.ParseCIDR(mask)
	if err != nil {
		return nil, fmt.Errorf("dhcp.new %w", err)
	}

	return &DHCP{
		subnet:      subnet,
		gatewayIP:   parsedGatewayIP,
		assignedIPs: assignedIPs,
		reservedIPs: make(map[string]bool),
	}, nil
}

func (d *DHCP) Reserve() (net.IP, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	ip := d.subnet.IP.Mask(d.subnet.Mask)
	for d.subnet.Contains(ip) {
		if d.subnet.IP.Equal(ip) {
			incrementIP(ip)
			continue
		}

		if d.gatewayIP.Equal(ip) {
			incrementIP(ip)
			continue
		}

		ipStr := ip.String()

		if !d.reservedIPs[ipStr] && !d.assignedIPs[ipStr] {
			d.reservedIPs[ipStr] = true
			return ip, nil
		}
		incrementIP(ip)
	}

	return nil, fmt.Errorf("there is no available ip address in subnet %s", d.subnet.Mask)
}

func (d *DHCP) Assign(ip net.IP) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	ipStr := ip.String()

	if !d.subnet.Contains(ip) {
		return fmt.Errorf("ip %s is not in the subnet", ipStr)
	}

	if !d.reservedIPs[ipStr] {
		return fmt.Errorf("ip %s is not reserved", ipStr)
	}

	delete(d.reservedIPs, ipStr)
	d.assignedIPs[ipStr] = true

	return nil
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
