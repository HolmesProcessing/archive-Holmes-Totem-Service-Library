package inputtype

// Imports required to determine the type of the object contained in a *Object.
import (
	"net"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
)

// Determine the type of the object and return an integer corresponding to the
// global constants Domain, Email, or IP.
//
// Supports IPv4, IPv6 in punctuated string representation (returns net.IP)
// Supports CIDR IP address range (returns *net.IPNet)
// Supports Email and "Name <Email>" (returns *mail.Address)
// Supports Domain (returns the input string)
// Supports detection of files (by checking if it is available in /tmp/)
//
// On success returns:
//
//      Type(int): Domain, Email, IP, IPNet, File
//      interface{}: The parsed version of the input (*net.IP / *net.IPNet / *net.Email / *net.Domain)
//      Error: nil
//
// On failure returns:
//
//      Type(int): Unknown (-1)
//      interface{}: nil
//      Error: EmptyInputError / UnknownTypeError
//
func detectInputType(obj string) (Type, interface{}, error) {
	var detectErr error = nil
	// Avoid problems in functions below.
	if obj == "" {
		detectErr = EmptyInputError

	} else if ip := net.ParseIP(obj); ip != nil {
		// TODO: support for IP's in decimal form
		return IP, ip, nil

	} else if _, ipnet, err := net.ParseCIDR(obj); err == nil {
		return IPNet, ipnet, nil

	} else if strings.Contains(obj, ".") && isDomainName(obj) {
		return Domain, obj, nil

	} else if email, err := mail.ParseAddress(obj); err == nil {
		return Email, email, nil

	} else {
		path := string(os.PathSeparator) + filepath.Join("tmp", filepath.Clean(obj))
		if _, err := os.Stat(path); err == nil {
			return File, path, nil

		} else {
			detectErr = UnknownTypeError

		}
	}
	return Unknown, nil, detectErr
}
