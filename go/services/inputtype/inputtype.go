package inputtype

// Imports required to determine the type of the object contained in a *Object.
import (
	"net"
	"net/mail"
)

/*
 *
 * General exported functions.
 * Call InitializeTLDMap before using Detect on a domain or email string.
 * Otherwise that call will panic.
 *
 */

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
// *Any attempt to detect a domain without a proper TLD map set will
// result in a panic!* (see inputtype.InitializeTLDMap(path))
//
// TODO: provide methods to check if an input is a specific type
//
func Detect(input string) (Type, interface{}, error) {
	return detectInputType(input)
}

// Initialize the internal TLD map.
// Accepts the path to a file containing a newline separated list of top level
// domains. Lines beginning with # are ignored.
func InitializeTLDMap(path string) error {
	return initTldMap(path)
}

/*
 * Validation functions.
 * Please note that this collection contains no functions that replicate any
 * functionality already available within the net package (net.IP offers quite
 * a range of functionality to detect different types of IP)
 *
 */

// Check if an IP is public or not.
// Public IP is any IP that is neither of:
//
// - private
// - loopback
// - broadcast / multicast
// - link-local
//
// Please note that any IP created by net.IPv4(...) is a IPv6 mapped IPv4 address
// which will not be detected and handled properly. It is your own responsibility
// to ensure that your address has the correct length (4/16 bytes)
// (Can use `ip := net.IP{byte,byte,byte,byte}` or `ip := net.IP([]byte)` for example)
//
func ValidateIP(ip net.IP) (bool, error) {
	return validateIP(ip)
}

// Checks whether or not the given string is valid domain name.
// Additionally checks whether or not the given string contains a valid top
// level domain.
//
// *Any attempt to identify a domain without a proper TLD map set will
// result in a panic!* (see inputtype.InitializeTLDMap(path))
//
// Returns (true, nil) if the domain is valid. Otherwise it returns either
// (false, InvalidDomainError) or (false, InvalidTLDError)
func ValidateDomain(domain string) (bool, error) {
	return validateDomain(domain)
}

// Checks whether or not the given emails address contains a valid domain name.
// Additionally checks whether or not the given domain contains a valid top
// level domain.
//
// *Any attempt to identify a eamil without a proper TLD map set will
// result in a panic!* (see inputtype.InitializeTLDMap(path))
//
// Returns (true, nil) if the domain is valid. Otherwise it returns either
// (false, InvalidDomainError) or (false, InvalidTLDError)
func ValidateEmail(email *mail.Address) (bool, error) {
	return validateEmail(email)
}

// Checks if the given filepath exists and can be opened for reading.
// Returns (true, nil) on success ad either (false, FileNotFoundError) or
// (false, FileAccessDeniedError) on failure.
func ValidateFile(filepath string) (bool, error) {
	return validateFile(filepath)
}
