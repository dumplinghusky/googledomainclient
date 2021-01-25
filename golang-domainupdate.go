

package main

import (
    "fmt"
    "os/exec"
    "runtime"
)
func rfc1918private(ip net.IP) bool {
	for _, cidr := range []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"} {
		_, subnet, err := net.ParseCIDR(cidr)
		if err != nil {
			panic("failed to parse hardcoded rfc1918 cidr: " + err.Error())
		}
		if subnet.Contains(ip) {
			return true
		}
	}
	return false
}

func rfc4193private(ip net.IP) bool {
	_, subnet, err := net.ParseCIDR("fd00::/8")
	if err != nil {
		panic("failed to parse hardcoded rfc4193 cidr: " + err.Error())
	}
	return subnet.Contains(ip)
}

func isLoopback(ip net.IP) bool {
	for _, cidr := range []string{"127.0.0.0/8", "::1/128"} {
		_, subnet, err := net.ParseCIDR(cidr)
		if err != nil {
			panic("failed to parse hardcoded loopback cidr: " + err.Error())
		}
		if subnet.Contains(ip) {
			return true
		}
	}
	return false
}

func mightBePublic(ip net.IP) bool {
	return !rfc1918private(ip) && !rfc4193private(ip) && !isLoopback(ip)
}

type name_and_ip struct {
	string
	net.IP
}

func heuristic(ni name_and_ip) (ret int) {
	a := strings.ToLower(ni.string)
	ip := ni.IP
	if isLoopback(ip) {
		ret += 1000
	}
	if rfc1918private(ip) || rfc4193private(ip) {
		ret += 500
	}
	if strings.Contains(a, "dyn") {
		ret += 100
	}
	if strings.Contains(a, "dhcp") {
		ret += 99
	}
	for i := 0; i < len(ip); i++ {
		if strings.Contains(a, strconv.Itoa(int(ip[i]))) {
			ret += 5
		}
	}
	return ret
}

type nameAndIPByStabilityHeuristic []name_and_ip
func (nis nameAndIPByStabilityHeuristic) Len() int { return len(nis) }
func (nis nameAndIPByStabilityHeuristic) Swap(i, j int) { nis[i], nis[j] = nis[j], nis[i] }
func (nis nameAndIPByStabilityHeuristic) Less(i, j int) bool { return heuristic(nis[i]) < heuristic(nis[j]) }

func publicAddresses() ([]name_and_ip, error) {
	var ret []name_and_ip

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			return nil, err
		}
		// ignore unresolvable addresses
		names, err := net.LookupAddr(ip.String())
		if err != nil {
			continue
		}
		for _, name := range names {
			ret = append(ret, name_and_ip{name, ip})
		}
	}

	sort.Sort(nameAndIPByStabilityHeuristic(ret))
	return ret, nil
}

func main() {
	if names, err := publicAddresses(); err == nil {
		for _, ni := range names {
			println(ni.string, ni.IP.String(), heuristic(ni))
		}
	} else {
		panic(err)
	}

}
func execute() {

    // here we perform the pwd command.
    // we can store the output of this in our out variable
    // and catch any errors in err
    out, err := exec.Command("ls").Output()

    // if there is an error with our execution
    // handle it here
    if err != nil {
		fmt.Printf("%s", err)
		fmt.Printf("Try Again")
    }
    // as the out variable defined above is of type []byte we need to convert
    // this to a string or else we will see garbage printed out in our console
    // this is how we convert it to a string
    fmt.Println("Command Successfully Executed")
    output := string(out[:])
    fmt.Println(output)

    // let's try the pwd command herer
    out, err = exec.Command("pwd").Output()
    if err != nil {
        fmt.Printf("%s", err)
    }
    fmt.Println("Command Successfully Executed")
	
    output = string(out[:])
    fmt.Println(output)
}

func main() {
    if runtime.GOOS == "windows" {
        fmt.Println("Can't Execute this on a windows machine")
    } else {
        execute()
    }
}

