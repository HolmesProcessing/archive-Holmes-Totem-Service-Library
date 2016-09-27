# Import for error output
# from tornado.web import HTTPError

# Imports for parsing
import ipaddress
import re
import os

# Import for error and type enums
import enum


"""
Exported functions. For public use
"""

def Detect(_input):
    """
    Determine the type of the object and return an integer corresponding to the
    global constants Domain, Email, or IP.

    Any string input is assumed to be utf-8 encoded and transformed to a bytes
    object by calling _input.encode("utf-8") if necessary (domain/email
    detection).

    Supports IPv4, IPv6 in punctuated string representation or decimal number
        (<32 bit is considered IPv4)
        -> Return value: ipaddress.IPv4Address or ipaddress.IPv6Address
    Supports CIDR IP address range (returns: ipaddress.IPv4Network or
        ipaddress.IPv6Network)
    Supports Email and "Name <Email>" (returns Email)
    Supports Domain (returns the input string)
    Supports detection of files (by checking if it is available in /tmp/)

    On success returns:

         Type(int): Domain, Email, IP, IPNet, File
         Parsed Object: Email, ipaddress.ip_address, ipaddress.ip_network, os.Fi
         Error: None

    On failure returns:

         Type(int): Unknown / None (None if EmptyInputError)
         Parsed Object: None
         Error: EmptyInputError / UnknownTypeError

    *Any attempt to detect a domain without a proper TLD map set will
    result in a panic!* (see inputtype.InitializeTLDMap(path))

    TODO: provide methods to check if an input is a specific type
    """
    return __detectType(_input)

def InitializeTLDMap(path):
    """
    Initialize the internal TLD map.
    Accepts the path to a file containing a newline separated list of top level
    domains. Lines beginning with # are ignored.
    Raises an error if the file cannot be opened, returns false if the file did
    not contain any TLDs (lines starting with # are ignored), returns True if
    the the length of the resulting TLD map is greater than one.

    Returns True on success and False on failure.
    """
    return __initTldMap(path)

def ValidateIP(ip):
    """
    Check if an IP is public or not.
    Public IP is any IP that is neither of:

    - private
    - loopback
    - broadcast / multicast
    - link-local

    Returns True on success and False on failure.

    Raises a ValueError exception if the input is not of type
    ipaddress.IPv4Address or ipaddress.IPv6Address.
    """
    return __validateIP(ip)

def ValidateDomain(domain):
    """
    Checks whether or not the given string is valid domain name.
    Additionally checks whether or not the given string contains a valid top
    level domain.

    *Any attempt to identify a domain without a proper TLD map set will
    result in an exception being raised!* (see inputtype.InitializeTLDMap(path))

    Returns (True, None) if the domain is valid. Otherwise it returns either
    (False, Errors.InvalidDomainError) or (False, Errors.InvalidTLDError)

    Raises a ValueError exception if the input is not of type str or bytes.
    """
    return __validateDomain(domain)

def ValidateEmail(email):
    """
    Checks whether or not the given emails address contains a valid domain name.
    Additionally checks whether or not the given domain contains a valid top
    level domain.

    *Any attempt to identify a domain without a proper TLD map set will
    result in an exception being raised!* (see inputtype.InitializeTLDMap(path))

    Returns (True, None) if the email is valid. Otherwise it returns either
    (False, Errors.InvalidEmailError), (False, Errors.InvalidDomainError),
    or (False, Errors.InvalidTLDError)

    TODO: output differs a bit from Go version, should be unified long-term

    raises a ValueError exception if the input is not of type str, bytes or
    inputtypes.Email.
    """
    return __validateEmail(email)

def ValidateFile(file):
    """
    Checks if the given filepath exists and can be opened for reading.

    Returns (True, None) on success and either of
    (False, Errors.FileNotFoundError) or (False, Errors.FileAccessDeniedError)
    """
    return __validateFile(file)


"""
Custom classes for input types
"""

class Email(object):
    def __init__(self, fullname, user, domain, ip=None):
        self.fullname = fullname  # "Peter Parker" in "Peter Parker <peter@parker.com>"
        self.user = user          # "peter"        in "Peter Parker <peter@parker.com"
        self.domain = domain      # "parker.com"   in "Peter Parker <peter@parker.com"
        self.ip = ip              # "53.24.88.19"  in "tester@[53.24.88.19]"
    def validate(self):
        return ValidateEmail(self)


"""
********************************************************************************
Private functions. Not for public use.
"""

"""
Detection functions.
"""

def is_ascii(s):
    return isinstance(s, str) and len(s) == len(s.encode())

def __detectType(_input):
    if not _input:
        return None, None, Errors.EmptyInputError
    if isinstance(_input, bytes):
        _input = _input.decode()

    if not isinstance(_input, (str, int)):
        raise ValueError("Invalid parameter type supplied to __detectType: {}, must be str or int".format(type(_input)))

    ip = __detectIP(_input)
    if ip:
        return Types.IP, ip, None

    ipnet = __detectIPNet(_input)
    if ipnet:
        return Types.IPNet, ipnet, None

    if isinstance(_input, str):
        domain = __detectDomain(_input)
        if domain:
            return Types.Domain, domain, None

        email = __detectEmail(_input)
        if email:
            return Types.Email, email, None

        file = __detectFile(_input)
        if file:
            return Types.File, file, None

    return Types.Unknown, None, Errors.UnknownTypeError


def __detectIP(_input):
    if not _input:
        return None
    try:
        ip = ipaddress.ip_address(_input)
    except ValueError:
        return None
    return ip

def __detectIPNet(_input):
    if not _input:
        return False
    try:
        ipnet = ipaddress.ip_network(_input, strict=False)
    except ValueError:
        return False
    return ipnet

def __detectDomain(_input):
    if _input and is_ascii(_input) and __isDomainName(_input):
        return _input
    return None

# based on validators library which in turn is based on djangos email validator
# http://validators.readthedocs.io/en/latest/_modules/validators/email.html#email
# improved by a domain name parsing function ported from Go
def __detectEmail(_input):
    if _input:
        fullname = ""
        if _input[-1] == '>':
            pos = _input.rfind('<')
            if pos >= 0:
                fullname = _input[:pos].strip()
                _input = (_input[pos+1:-1]).strip()
        if '@' in _input:
            user, domain = _input.rsplit('@', 1)
            if is_ascii(_input):
                if __isDomainName(domain):
                    return Email(fullname, user, domain)
                if domain[0]=="[" and domain[-1] == "]":
                    ip = __detectIP(domain[1:-1])
                    if ip:
                        return Email(fullname, user, None, ip)
    return None

def __cleanPath(path):
    return os.path.normpath(path) #.lstrip("."+os.path.sep)
    # the commented change --------^
    # would be more restrictive than go version
    # (might be better? would jail it to /tmp) TODO: discuss
def __detectFile(_input):
    path = os.path.sep + os.path.join("tmp", __cleanPath(_input))
    if os.path.isfile(path):
        return path
    return None


"""
Initialize variables and errors for validation functionality.
"""

Types = enum.Enum('Types', [
    "Domain",
    "Email",
    "IP",
    "IPNet",
    "File",
    "Unknown",
    ], module=__name__)

Errors = enum.Enum('Errors', [
    "EmptyInputError",
    "UnknownTypeError",

    "NonAsciiCharacters",
    "InvalidEmailError",
    "InvalidTLDError",
    "InvalidDomainError",

    "FileNotFoundError",
    "FileAccessDeniedError",

    "IPisLoopbackError",
    "IPisMulticastError",
    "IPisUnspecifiedError",
    "IPisInterfaceLocalMulticastError",
    "IPisLinkLocalMulticastError",
    "IPisLinkLocalUnicastError",
    "IPisNotPublicError",
    ], module=__name__)

def __parseCIDRList(cidrlist):
    ipnetlist = []
    for cidr in cidrlist:
        ipnetlist.append(ipaddress.ip_network(cidr))
    return ipnetlist

__ipv4Private = __parseCIDRList([
    # private use networks
    "10.0.0.0/8",
    "172.16.0.0/12",
    "192.168.0.0/16",
])
__ipv4Nonpublic = __parseCIDRList([
    # this-host-this-net,
    # loopback,
    # link-local,
    # shared address space (for cgn devices for providers, see iana rfc6598)
    # IETF protocol assignments (iana rfc6890)
    # TEST-NET-1 (iana rfc5737)
    # benchmarking
    # TEST-NET-2 (iana rfc5737)
    # TEST-NET-3 (iana rfc5737)
    # Reserved (iana rfc1112, section 4)
    # limited broadcast
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
    # local network control block (iana rfc5771)
    "224.0.0.0/24",
    # possible other canditates listed in rfc5771 for multicast and
    # detailed allocations see, some of them seem to be routable though:
    # http://www.iana.org/assignments/multicast-addresses/multicast-addresses.xhtml#multicast-addresses-2
]) + __ipv4Private

__ipv6Private = __parseCIDRList([
    # private ipv6 networks (ietf rfc1918)
    # ORCHIDv2 non-routable (ietf rfc7343)
    "fc00::/7",
    "2001:20::/28",
])
__ipv6Nonpublic = __parseCIDRList([
    # Unspecified Address
    # Loopback Address
    # IPv4-mapped Address
    # Discard-Only Address Block
    # TEREDO
    # Benchmarking
    # Documentation
    # Unique-Local
    # Link-Scoped Unicast
    "::/128",
    "::1/128",
    "::ffff:0:0/96",
    "100::/64",
    "2001::/32",
    "2001:2::/48",
    "2001:db8::/32",
    "fc00::/7",
    "fe80::/10",

    # link-local (ietf rfc4193)
    # multicast
    # see:
    # http://www.iana.org/assignments/ipv6-multicast-addresses/ipv6-multicast-addresses.xhtml
    "FE80::/10",
    "FF00::/8",
]) + __ipv6Private


__tldMap = {}
__tldMapInitialized = False
def __initTldMapHelper(data):
    tlds = [line.strip() for line in data.split("\n")]
    for tld in tlds:
        if not is_ascii(tld):
            raise ValueError("TLD list contains non-ascii TLD: {}".format(tld))
        if tld and tld[0] != "#":
            __tldMap[tld.upper()] = True
    global __tldMapInitialized
    __tldMapInitialized = True

def __initTldMap(path):
    with open(path, "r") as file:
        data = file.read()
        __initTldMapHelper(data.decode())
    if len(__tldMap) == 0:
        return False
    return True


"""
Validation functionality.
"""

def __isDomainName(s):
    # Ported from golang standard package net:
    # https://golang.org/src/net/dnsclient.go?m=text
    # See RFC 1035, RFC 3696.
    # Expectes a ASCII (byte) string (ensure this with a call to is_ascii(s))
    if len(s) == 0 or len(s) > 255:
        return False

    last = '.'
    ok = False # Ok once we've seen a letter.
    partlen = 0
    i = 0
    containsDot = False
    while i < len(s):
        c = s[i]
        i += 1

        if ('a' <= c and c <= 'z') or ('A' <= c and c <= 'Z') or c == '_':
            ok = True
            partlen += 1

        elif '0' <= c and c <= '9':
            # TODO: check code here
            # here's probably missing a `if not ok: return False` to avoid
            # labels starting with numbers
            # see rfc https://tools.ietf.org/html/rfc1035 page 8
            partlen += 1

        elif c == '-':
            if last == '.':
                return False
            partlen += 1

        elif c == '.':
            if last == '.' or last == '-':
                return False
            if partlen > 63 or partlen == 0:
                return False
            partlen = 0
            containsDot = True  # modification to indicate that we have more than just one label
            # here's probably missing a `ok = False` to avoid labels starting
            # with numbers
            # see rfc https://tools.ietf.org/html/rfc1035 page 8

        else:
            return False

        last = c

    if last == '-' or partlen > 63:
        return False
    return ok and containsDot


def __inIPNet(ip, ipv4nets, ipv6nets):
    if isinstance(ip, ipaddress.IPv4Address):
        ipnets = ipv4nets
    if isinstance(ip, ipaddress.IPv6Address):
        ipnets = ipv6nets

    for ipnet in ipnets:
        if ip in ipnet:
            return True
    return False


def __inTldMap(domain):
    if not __tldMapInitialized:
        raise UnboundLocalError("tldMap not (or not properly) initialized - use inputtype.InitializeTLDMap(path)")
    pos = domain.rfind('.')
    if pos >= 0:
        domain = domain[pos+1:]
        if domain:
            return (domain.upper() in __tldMap)
    return False


def __validateIP(ip):
    if not (isinstance(ip, ipaddress.IPv4Address) or isinstance(ip, ipaddress.IPv6Address)):
        raise ValueError("Invalid parameter supplied to __validateIP(ip), must be ipaddress.ip_address")
    if ip.is_loopback:
        return False, Errors.IPisLoopbackError
    if ip.is_unspecified:
        return False, Errors.IPisUnspecifiedError
    if __inIPNet(ip, __ipv4Nonpublic, __ipv6Nonpublic):
        return False, Errors.IPisNotPublicError
    return True, None

def __validateDomain(domain):
    if not is_ascii(domain):
        return False, Errors.NonAsciiCharacters
    if not __isDomainName(domain):
        return False, Errors.InvalidDomainError
    if not __inTldMap(domain):
        return False, Errors.InvalidTLDError
    return True, None

# based on validators library which in turn is based on djangos email validator
# http://validators.readthedocs.io/en/latest/_modules/validators/email.html#email
# improved by a domain name parsing function ported from Go
__email_user_regex = re.compile(
    # dot-atom
    r"(^[-!#$%&'*+/=?^_`{}|~0-9A-Z]+"
    r"(\.[-!#$%&'*+/=?^_`{}|~0-9A-Z]+)*$"
    # quoted-string
    r'|^"([\001-\010\013\014\016-\037!#-\[\]-\177]|'
    r"""\\[\001-\011\013\014\016-\177])*"$)""",
    re.IGNORECASE
)
def __validateEmail(email):
    if not isinstance(email, Email):
        raise ValueError("Invalid parameter supplied to __validateEmail(email), must be inputtypes.Email")
    if not __email_user_regex.match(email.user):
        return False, Errors.InvalidEmailError
    ok, err = __validateDomain(email.domain)
    if ok:
        return True, None
    if not email.ip:
        return False, err
    return __validateIP(email.ip)

def __validateFile(file):
    try:
        with open(file, "r") as fh:
            fh.read(0x1)
    except (FileNotFoundError, IsADirectoryError):
        return False, Errors.FileNotFoundError
    except PermissionError:
        return False, Errors.FileAccessDeniedError
    return True, None
