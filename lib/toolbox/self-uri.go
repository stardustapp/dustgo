package toolbox

import (
	"errors"
	"log"
	"net"
)

// Build cluster-wide absolute URI prefix to the current process
func SelfURI(suffix string) (string, error) {
	selfIp := GetLocalIP()
	if selfIp == "" {
		log.Println("SelfURI didn't locate any host IPs")
		return "", errors.New("SelfURI didn't locate any host IPs")
	}

	return "skylink+ws://" + selfIp + suffix, nil
}

// https://stackoverflow.com/a/31551220
// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }
    for _, address := range addrs {
        // check the address type and if it is not a loopback the display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
}
