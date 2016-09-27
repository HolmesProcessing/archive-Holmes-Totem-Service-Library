package inputtype

import (
	"net"
	"net/mail"
	"os"
	"strings"
)

// Helper function to determine if a given domains tld is registered with iana.
// Will panic if tldMap is not initialized.
func inTldMap(domain string) bool {
	if len(tldMap) == 0 {
		panic("tldMap not (or not properly) initialized - use inputtype.InitializeTLDMap(path)")
	}
	end := len(domain) - 1
	// strip a possible trailing dot
	if domain[end] == '.' {
		end--
	}
	// account for initial decrease
	end++
	// start one position below end
	// go until either all characters read or a dot is found
	start := end - 1
	for ; start >= 0 && domain[start] != '.'; start-- {
	}
	// increase position by one to account for the additional decrease
	start++
	// TODO: empty strings invalid? (possibly only a dot supplied for example)
	if start == end {
		return false
	}
	// otherwise check in the TLD map
	tld := domain[start:end]
	_, ok := tldMap[strings.ToUpper(tld)]
	return ok
}

// Part of golang standard package net:
// https://golang.org/src/net/dnsclient.go?m=text
func isDomainName(s string) bool {
	// See RFC 1035, RFC 3696.
	if len(s) == 0 {
		return false
	}
	if len(s) > 255 {
		return false
	}

	last := byte('.')
	ok := false // Ok once we've seen a letter.
	partlen := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		default:
			return false
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_':
			ok = true
			partlen++
		case '0' <= c && c <= '9':
			// fine
			// TODO: check code here
			// here's probably missing a `if not ok: return false` to avoid
			// labels starting with numbers
			// see rfc https://tools.ietf.org/html/rfc1035 page 8
			partlen++
		case c == '-':
			// Byte before dash cannot be dot.
			if last == '.' {
				return false
			}
			partlen++
		case c == '.':
			// Byte before dot cannot be dot, dash.
			if last == '.' || last == '-' {
				return false
			}
			if partlen > 63 || partlen == 0 {
				return false
			}
			partlen = 0
			// here's probably missing a `ok = false` to avoid labels starting
			// with numbers
			// see rfc https://tools.ietf.org/html/rfc1035 page 8
		}
		last = c
	}
	if last == '-' || partlen > 63 {
		return false
	}

	return ok
}

// Helper function, checks if the given IP is within one of the reserved IP
// address ranges. See ipv4NetFiltered and ipv6NetFiltered variable
// initialization in init().
func isInAddressRange(ip net.IP, ipv4ranges, ipv6ranges []*net.IPNet) bool {
	if len(ip) == net.IPv4len {
		for _, net := range ipv4ranges {
			if net.Contains(ip) {
				return true
			}
		}
	} else if len(ip) == net.IPv6len {
		for _, net := range ipv6ranges {
			if net.Contains(ip) {
				return true
			}
		}
	}
	return false
}

// Helper function to check whether a domain is valid or not.
func validateDomain(domain string) (bool, error) {
	if !isDomainName(domain) {
		return false, InvalidDomainError
	}
	if !inTldMap(domain) {
		return false, InvalidTLDError
	}
	return true, nil
}

// Helper function to check whether an email address is valid or not.
func validateEmail(email *mail.Address) (bool, error) {
	parts := strings.SplitN(email.Address, "@", 2)
	domain := parts[len(parts)-1]
	parts = strings.Split(domain, ".")
	if len(parts) < 2 || !isDomainName(domain) {
		return false, InvalidDomainError
	}
	if !inTldMap(domain) {
		return false, InvalidTLDError
	}
	return true, nil
}

// Helper function to check if a given file is valid (can be opened for reading)
func validateFile(filepath string) (bool, error) {
	if _, err := os.Stat(filepath); err == nil {
		if file, err := os.Open(filepath); err == nil {
			file.Close()
			return true, nil
		}
		return false, FileAccessDeniedError
	}
	return false, FileNotFoundError
}

// Helper function to check if a given IP is not public (non-routable).
func validateIP(ip net.IP) (bool, error) {
	if ip.IsLoopback() {
		return false, IPisLoopbackError
	}
	if ip.IsUnspecified() {
		return false, IPisUnspecifiedError
	}
	if isInAddressRange(ip, ipv4Nonpublic, ipv6Nonpublic) {
		return false, IPisNotPublicError
	}
	return true, nil
}
