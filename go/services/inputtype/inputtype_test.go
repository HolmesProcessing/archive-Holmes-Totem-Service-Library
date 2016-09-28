package inputtype

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/mail"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInputtype(T *testing.T) {
	T.Log("=== Testing module inputtype ===")
	// prepare test TLD map
	file, err := ioutil.TempFile("", "holmeslibrary_test_tldmap")
	if err != nil {
		T.Error("Unable to create temporary file for testing")
		return
	}
	filepath := file.Name()
	defer os.Remove(filepath)
	file.Write([]byte(`
        # some sample tlds
        COM
        DE
        XN--PUNYCODEY2342
    `))
	file.Close()
	InitializeTLDMap(filepath)

	// seed pseudo random number generator using current nano time
	rand.Seed(time.Now().UTC().UnixNano())

	// run all tests
	success := true &&
		test0_Type2String(T) &&
		test1_IP(T) &&
		test2_IPNet(T) &&
		test3_DomainAndEmail(T) &&
		test4_File(T)

	if !success {
		T.Error("Inputtype test failed")
		return
	}
}

// extended test of IP detection functionality
func test1_IP(T *testing.T) bool {
	T.Log("--- IP Test ---")

	type Test struct {
		ip    net.IP
		valid bool
		err   error
	}

	testcases := []Test{
		{ // random public ip, just between the two private network ranges
			net.IP{byte(11 + rand.Int31n(127-11)), rand255(), rand255(), rand255()},
			true, nil,
		},
		{ // random private ip in 10.0.0.0/8
			net.IP{10, rand255(), rand255(), rand255()},
			false, IPisNotPublicError,
		},
		{ // random private ip in 172.16.0.0/12
			net.IP{172, byte(16 + rand.Int31n(31-16)), rand255(), rand255()},
			false, IPisNotPublicError,
		},
		{ // random private ip in 192.168.0.0/16
			net.IP{192, 168, rand255(), rand255()},
			false, IPisNotPublicError,
		},
		{ // random loopback ip in 127.0.0.0/8
			net.IP{127, rand255(), rand255(), rand255()},
			false, IPisLoopbackError,
		},
		{ // random this-host-this-net ip in 0.0.0.0/8
			net.IP{0, rand255(), rand255(), rand255()},
			false, IPisNotPublicError,
		},
		{ // broadcast
			net.IP{255, 255, 255, 255},
			false, IPisNotPublicError,
		},
		{ // ipv6 loopback
			net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			false, IPisLoopbackError,
		},
		{ // ipv6 unspecified
			net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			false, IPisUnspecifiedError,
		},
		{ // ipv6 multicast
			net.IP{255, rand255(), rand255(), rand255(), rand255(), rand255(),
				rand255(), rand255(), rand255(), rand255(), rand255(), rand255(),
				rand255(), rand255(), rand255(), rand255()},
			false, IPisNotPublicError,
		},
		{ // ipv6 link-local
			net.IP{255, byte(0x80 + rand.Int31n(63)), rand255(), rand255(), rand255(),
				rand255(), rand255(), rand255(), rand255(), rand255(), rand255(),
				rand255(), rand255(), rand255(), rand255(), rand255()},
			false, IPisNotPublicError,
		},
	}

	bMap := make(map[bool]string)
	bMap[true] = "true"
	bMap[false] = "false"
	for _, test := range testcases {
		T.Logf("Detect(\"%s\")\n", test.ip.String())
		t, o, err := Detect(test.ip.String())
		ip := o.(net.IP)
		if len(test.ip) == net.IPv4len {
			ip = ip.To4()
			if ip == nil {
				T.Log("IP address length mismatch (IPv4 expected)")
				return false
			}
		}
		// var bytes []byte = ip
		if !(assertTrue(T, err == nil, "err!=nil") &&
			assertTrue(T, t == IP, "type!=IP") &&
			assertTrue(T, ip.Equal(test.ip), "ip!="+test.ip.String())) {
			return false
		}
		ok, err := ValidateIP(ip)
		if !(assertTrue(T, ok == test.valid, fmt.Sprintf("ValidateIP(ip)!=%v, %v", test.valid, err)) &&
			assertTrue(T, err == test.err, fmt.Sprintf("ValidateIP(ip) err(%v)!=%v", err, test.err))) {
			return false
		}
	}
	return true
}

// simple test of IPNet detection functionality
func test2_IPNet(T *testing.T) bool {
	T.Log("--- IPNet Test ---")

	testcases := []uint{16, 4}

	for _, addrSize := range testcases {
		ip := net.IP(genRandomBytes(addrSize))
		prefixLen := uint(rand.Int31n(int32(addrSize) * 8))
		cidr := fmt.Sprintf("%s/%d", ip.String(), prefixLen)
		netmask := genNetmask(addrSize, prefixLen)
		T.Logf("Detect(\"%s\")\n", cidr)
		t, o, err := Detect(cidr)
		if !(assertTrue(T, t == IPNet, "type!=IPNet") &&
			assertTrue(T, err == nil, "err!=nil") &&
			assertTrue(T, byteArrayEqual(o.(*net.IPNet).Mask, netmask), "Incorrect netmask")) {
			return false
		}
	}
	return true
}

// simple test of Domain detection functionality
func test3_DomainAndEmail(T *testing.T) bool {
	T.Log("--- Domain And Email Test ---")

	type Test struct {
		input           string
		outType         Type
		outError        error
		validationError error
	}

	testcases := []Test{
		{"www.domain.de", Domain, nil, nil},
		{"yet-another-domain.com", Domain, nil, nil},
		{"www.domain.eu", Domain, nil, InvalidTLDError},
		{"xn--punycode.subdomain2.subdomain1.domain.com", Domain, nil, nil},
		{"test.-nodomain.de", Unknown, UnknownTypeError, nil},
		{"somename@somedomain.com", Email, nil, nil},
		{"webmaster123ya##ösäf@somedomain.com", Unknown, UnknownTypeError, nil},
		{"somename@invalidtld.eu", Email, nil, InvalidTLDError},
		{"somename@nodomain-.com", Email, nil, InvalidDomainError},
		{"somename@34.128.94.77", Email, nil, InvalidDomainError},
		{"Max Musterman <max@musterman.com>", Email, nil, nil},
	}

	for _, test := range testcases {
		T.Logf("Detect(\"%s\")\n", test.input)
		t, o, err := Detect(test.input)
		if !assertTrue(T, t == test.outType, fmt.Sprintf("type(%s)!=%s", t, test.outType.String())) {
			return false
		}
		if test.outType == Domain &&
			!assertTrue(T, o != nil && o.(string) == test.input, "input!=output") {
			return false
		}
		if test.outError != err {
			if err != nil && test.outError != nil {
				T.Logf("err(%s)!=%s\n", err.Error(), test.outError.Error())
				return false

			} else if test.outError != nil {
				T.Logf("err(nil)!=%s\n", test.outError.Error())
				return false

			} else {
				T.Logf("err(%s)!=nil\n", err.Error())
				return false
			}
		}
		if t == Domain {
			_, err = ValidateDomain(o.(string))
		} else if t == Email {
			_, err = ValidateEmail(o.(*mail.Address))
		} else {
			err = test.validationError
		}
		if err != test.validationError {
			if err != nil && test.validationError != nil {
				T.Logf("err(%s)!=%s\n", err.Error(), test.validationError.Error())

			} else if test.validationError != nil {
				T.Logf("err(nil)!=%s\n", test.validationError.Error())

			} else {
				T.Logf("err(%s)!=nil\n", err.Error())
			}
			return false
		}
	}
	return true
}

// additional test verifying correct string output when calling String(type)
func test0_Type2String(T *testing.T) bool {
	T.Log("--- Type.String() Test ---")

	if assertTrue(T, IP.String() == "IP", "String(IP)!=\"IP\"") &&
		assertTrue(T, IPNet.String() == "IPNet", "String(IPNet)!=\"IPNet\"") &&
		assertTrue(T, Domain.String() == "Domain", "String(Domain)!=\"Domain\"") &&
		assertTrue(T, Email.String() == "Email", "String(Email)!=\"Email\"") {
		return true
	}
	return false
}

// Simple test if file detection works.
// TODO: the file detection test is assuming Linux, implement windows variant too
func test4_File(T *testing.T) bool {
	T.Log("--- File Test ---")

	type Test struct {
		filepath string
		outType  Type
		outErr   error
		valErr   error
	}

	// Create one test file and retrieve its path.
	etc := "/tmp"
	prefix := "holmeslibrary_file_detection_test"
	file, err := ioutil.TempFile(etc, prefix)
	if err != nil {
		T.Error("Unable to create temporary file for testing")
		return false
	}
	path := file.Name()
	defer os.Remove(path)
	file.Write([]byte(`
        Some blah blah.
    `))
	file.Close()

	// Define all test cases.
	// TODO: additional test-cases
	testcases := []Test{
		{filepath.Base(path), File, nil, nil},
		{"../etc/shadow", File, nil, FileAccessDeniedError},
		{"../invalidtopleveldirectory/invalidfile", Unknown, UnknownTypeError, nil},
	}

	for _, test := range testcases {
		T.Logf("Detect(\"%s\")\n", test.filepath)
		t, o, err := Detect(test.filepath)
		if !(assertTrue(T, t == test.outType, fmt.Sprintf("type(%v)!=%v", t, test.outType)) &&
			assertTrue(T, err == test.outErr, fmt.Sprintf("err(%v)!=%v", err, test.outErr))) {
			return false
		}
		if t == File {
			_, err := ValidateFile(o.(string))
			if err != test.valErr {
				fmt.Println(o.(string), ",", err, ",", test.valErr)
			}
		}
	}

	// TODO need a better testcase
	_, err = ValidateFile("/invalidtopleveldirectory/invalidfile")
	if !assertTrue(T, err == FileNotFoundError, fmt.Sprintf("err(%v)!=%v", err, FileNotFoundError)) {
		return false
	}

	return true
}

/* Helper functions for testing */
// Simple helper function that expects b to be true, otherwise it will trigger
// an error on the testing.T object using the provided failmsg.
func assertTrue(T *testing.T, b bool, failmsg string) bool {
	if !b {
		T.Error(failmsg)
	}
	return b
}

// Simple comparison function returning true if two byte arrays are identical.
func byteArrayEqual(arr1, arr2 []byte) bool {
	if len(arr1) != len(arr2) {
		return false
	}
	for i := 0; i < len(arr1); i++ {
		if arr1[i] != arr2[i] {
			return false
		}
	}
	return true
}

// helper function to create a netmask for IPNet test
func genNetmask(maskLen, prefixLen uint) []byte {
	mask := make([]byte, maskLen)
	ffs := uint(prefixLen / 8)
	i := uint(0)
	for ; i < ffs && i < maskLen; i++ {
		mask[i] = byte(0xff)
	}
	if i < maskLen {
		if prefixLen%8 != 0 {
			mask[i] = 0xff << (8 - prefixLen%8)
			i++
		}
		for ; i < maskLen; i++ {
			mask[i] = 0
		}
	}
	return mask
}

// helper function to create a slice of random bytes of length n for IPNet test
func genRandomBytes(n uint) []byte {
	bytes := make([]byte, n)
	for i := uint(0); i < n; i++ {
		bytes[i] = rand255()
	}
	return bytes
}

// helper function to get a random byte value
func rand255() byte {
	return byte(rand.Int31n(255))
}
