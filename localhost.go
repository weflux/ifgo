package localhost

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

func MacAddr() (string, error) {
	localIP, err := IP()
	if err != nil {
		return "", err
	}

	ifaces, err := IntranetInterface()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		ips, err := intranetIPs(iface)
		if err != nil {
			return "", err
		}
		for _, ip := range ips {
			if ip == localIP {
				return iface.HardwareAddr.String(), nil
			}
		}
	}

	return "", nil
}

func IP() (string, error) {
	ifaces, err := IntranetInterface()
	if err != nil {
		return "", err
	}
	ips := []string{}
	for _, iface := range ifaces {
		_ips, err := intranetIPs(iface)
		if err != nil {
			continue
		}
		for _, ip := range _ips {
			ips = append(ips, ip)
		}
	}
	if len(ips) == 0 {
		return "", errors.New("no local intranet ip found")
	}

	return ips[0], nil
}

func intranetIPs(iface net.Interface) (ips []string, err error) {
	addresses, e := iface.Addrs()
	if e != nil {
		return ips, e
	}

	for _, addr := range addresses {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		if ip == nil || ip.IsLoopback() {
			continue
		}
		ip = ip.To4()
		if ip == nil {
			// not an ipv4
			continue
		}
		ipStr := ip.String()
		ips = append(ips, ipStr)
	}
	return ips, e
}

func IntranetInterface() ([]net.Interface, error) {
	avails, e := AvailableInterfaces()
	if e != nil {
		return nil, e
	}

	ifaces := []net.Interface{}
	for _, iface := range avails {

		addresses, e := iface.Addrs()
		if e != nil {
			return ifaces, e
		}

		for _, addr := range addresses {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				// not an ipv4
				continue
			}
			ipStr := ip.String()
			if IsIntranet(ipStr) {
				ifaces = append(ifaces, iface)
			}
		}
	}

	return ifaces, nil
}

func AvailableInterfaces() ([]net.Interface, error) {
	rets := []net.Interface{}
	ifaces, e := net.Interfaces()
	if e != nil {
		return rets, e
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			// interface down
			continue
		}

		if iface.Flags&net.FlagLoopback != 0 {
			// loopback interface
			continue
		}

		// ignore warden bridge
		if strings.HasPrefix(iface.Name, "w-") {
			continue
		}

		rets = append(rets, iface)
	}
	return rets, nil
}

// IsIntranet checks and returns whether given ip an intranet ip.
//
// Local: 127.0.0.1
// A: 10.0.0.0--10.255.255.255
// B: 172.16.0.0--172.31.255.255
// C: 192.168.0.0--192.168.255.255
func IsIntranet(ip string) bool {
	if ip == "127.0.0.1" {
		return true
	}
	array := strings.Split(ip, ".")
	if len(array) != 4 {
		return false
	}
	// A
	if array[0] == "10" || (array[0] == "192" && array[1] == "168") {
		return true
	}
	// C
	if array[0] == "192" && array[1] == "168" {
		return true
	}
	// B
	if array[0] == "172" {
		second, err := strconv.ParseInt(array[1], 10, 64)
		if err != nil {
			return false
		}
		if second >= 16 && second <= 31 {
			return true
		}
	}
	return false
}
