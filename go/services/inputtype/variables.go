package inputtype

// Imports for logging
import (
	"log"
	"os"
)

// Imports required to determine the type of an input
import (
	"errors"
	"io/ioutil"
	"net"
	"strings"
)

// Public Type value, basically an int with a String method attached
type Type int

// Returns a string corresponding to the value that the Type variable holds.
func (t Type) String() string {
	switch int(t) {
	case -1:
		return "Unknown"
	case 0:
		return "Domain"
	case 1:
		return "Email"
	case 2:
		return "IP"
	case 3:
		return "IPNet"
	case 4:
		return "File"
	default:
		return "UNDEFINED TYPE CODE"
	}
}

// Public Type constants. Always only use the constant names, not their respective
// values, as these are subject to changes.
const (
	Domain Type = iota
	Email
	IP
	IPNet
	File
	Unknown Type = -1
)

// Public error type, an int with an Error (to string) method attached
type Error int

// Returns a string value corresponding to the contained error value
func (e Error) Error() string {
	switch e {
	case EmptyInputError:
		return "Empty Input String"
	case InvalidTLDError:
		return "Invalid TLD"
	case InvalidDomainError:
		return "Invalid Domain"
	case UnknownTypeError:
		return "Unknown Type"

	case FileNotFoundError:
		return "File Not Found"
	case FileAccessDeniedError:
		return "File Access Denied"

	case IPisLoopbackError:
		return "Loopback IP"
	case IPisMulticastError:
		return "Multicast IP"
	case IPisUnspecifiedError:
		return "Unspecified IP"
	case IPisInterfaceLocalMulticastError:
		return "Interface Local Multicast IP"
	case IPisLinkLocalMulticastError:
		return "Link Local Multicast IP"
	case IPisLinkLocalUnicastError:
		return "Link Local Unicast IP"
	case IPisNotPublicError:
		return "Non-Public IP"

	default:
		return "UNDEFINED ERROR CODE"
	}
}

// Public Error constants. Always only use the constant names, not their respective
// values, as these are subject to changes.
const (
	EmptyInputError Error = iota
	UnknownTypeError

	InvalidTLDError
	InvalidDomainError

	FileNotFoundError
	FileAccessDeniedError

	IPisLoopbackError
	IPisMulticastError
	IPisUnspecifiedError
	IPisInterfaceLocalMulticastError
	IPisLinkLocalMulticastError
	IPisLinkLocalUnicastError
	IPisNotPublicError
)

// Private constant for error output.
const (
	errBaseStr = "InputType Error: "
)

// Private global variables.
var (
	// global error logger instance
	errorLogger *log.Logger

	// address ranges to be filtered out
	ipv4to6       = []*net.IPNet{}
	ipv4Nonpublic = []*net.IPNet{}
	ipv4Private   = []*net.IPNet{}

	ipv6Nonpublic = []*net.IPNet{}
	ipv6Private   = []*net.IPNet{}

	// global tld map, fetched from iana
	tldMap            = make(map[string]bool)
	tldMapInitialized = false
)

// Initialize all private global variables
func init() {
	// initialize the logger instance
	errorLogger = log.New(os.Stderr, "", log.Ltime|log.Lshortfile)

	// init address ranges to be filtered out
	// for ipv4 see: https://tools.ietf.org/html/rfc6890
	// for ipv6 see: https://tools.ietf.org/html/rfc4291
	// further see:
	// http://www.iana.org/assignments/iana-ipv4-special-registry/iana-ipv4-special-registry.xhtml
	// http://www.iana.org/assignments/iana-ipv6-special-registry/iana-ipv6-special-registry.xhtml
	ipv4Private = parseCIDRList([]string{
		// private use networks
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	})
	ipv4Nonpublic = append(parseCIDRList([]string{
		// this-host-this-net,
		// loopback,
		// link-local,
		// shared address space (for cgn devices for providers, see iana rfc6598)
		// IETF protocol assignments (iana rfc6890)
		// TEST-NET-1 (iana rfc5737)
		// benchmarking
		// TEST-NET-2 (iana rfc5737)
		// TEST-NET-3 (iana rfc5737)
		// Reserved (iana rfc1112, section 4)
		// limited broadcast
		"0.0.0.0/8",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"100.64.0.0/10",
		"192.0.0.0/24",
		"192.0.2.0/24",
		"198.18.0.0/15",
		"198.51.100.0/24",
		"203.0.113.0/24",
		"240.0.0.0/4",
		"255.255.255.255/32",
		// local network control block (iana rfc5771)
		"224.0.0.0/24",
		// possible other canditates listed in rfc5771 for multicast and
		// detailed allocations see, some of them seem to be routable though:
		// http://www.iana.org/assignments/multicast-addresses/multicast-addresses.xhtml#multicast-addresses-2
	}), ipv4Private...)

	ipv6Private = parseCIDRList([]string{
		// private ipv6 networks (ietf rfc1918)
		// ORCHIDv2 non-routable (ietf rfc7343)
		"fc00::/7",
		"2001:20::/28",
	})
	ipv6Nonpublic = append(parseCIDRList([]string{
		// Unspecified Address
		// Loopback Address
		// IPv4-mapped Address
		// Discard-Only Address Block
		// TEREDO
		// Benchmarking
		// Documentation
		// Unique-Local
		// Link-Scoped Unicast
		"::/128",
		"::1/128",
		"::ffff:0:0/96",
		"100::/64",
		"2001::/32",
		"2001:2::/48",
		"2001:db8::/32",
		"fc00::/7",
		"fe80::/10",

		// link-local (ietf rfc4193)
		// multicast
		// see:
		// http://www.iana.org/assignments/ipv6-multicast-addresses/ipv6-multicast-addresses.xhtml
		"FE80::/10",
		"FF00::/8",
	}), ipv6Private...)
}

func initTldMap(path string) error {
	// read iana tld registry file
	// wget -O TLDList.txt "http://data.iana.org/TLD/tlds-alpha-by-domain.txt"
	tlds, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.New("Error reading TLD file: " + err.Error())
	} else {
		for _, tld := range strings.Split(string(tlds), "\n") {
			tld = strings.TrimSpace(tld)
			if tld != "" && tld[0] != '#' {
				tldMap[tld] = true
			}
		}
		tldMapInitialized = true
	}
	if len(tldMap) == 0 {
		return errors.New("Error reading TLD file: File did not contain any entries. Map not initialized.")
	}
	return nil
}

func parseCIDRList(cidrs []string) []*net.IPNet {
	var l []*net.IPNet
	for _, cidr := range cidrs {
		if _, ipnet, err := net.ParseCIDR(cidr); err != nil {
			panic(err)
		} else {
			l = append(l, ipnet)
		}
	}
	return l
}
