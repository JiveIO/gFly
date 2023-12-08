package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

// MsgForTagFunc Customize error fields to build messages
type MsgForTagFunc func(fe validator.FieldError) string

// MsgForTag is a type of MsgForTagFunc function.
func MsgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	// ------------------------------ Comparisons ------------------------------
	case "eq":
		return fmt.Sprintf("equals %s", fe.Param())
	case "eq_ignore_case":
		return fmt.Sprintf("equal ignoring case %s", fe.Param())
	case "gt":
		return fmt.Sprintf("greater than %s", fe.Param())
	case "gte":
		return fmt.Sprintf("greater than or equal %s", fe.Param())
	case "lt":
		return fmt.Sprintf("less than %s", fe.Param())
	case "lte":
		return fmt.Sprintf("less than or equal %s", fe.Param())
	case "ne":
		return fmt.Sprintf("not equal %s", fe.Param())
	case "ne_ignore_case":
		return fmt.Sprintf("not equal ignoring case %s", fe.Param())
	// ------------------------------ Strings ------------------------------
	case "alpha":
		return "invalid alpha"
	case "alphanum":
		return "invalid alphanumeric"
	case "alphanumunicode":
		return "invalid alphanumeric unicode"
	case "alphaunicode":
		return "invalid alpha unicode"
	case "ascii":
		return "invalid ASCII"
	case "boolean":
		return "invalid boolean type"
	case "contains":
		return fmt.Sprintf("contains %s", fe.Param())
	case "containsany":
		return "contains any"
	case "containsrune":
		return "contains rune"
	case "endsnotwith":
		return fmt.Sprintf("ends not with %s", fe.Param())
	case "endswith":
		return fmt.Sprintf("ends with %s", fe.Param())
	case "excludes":
		return fmt.Sprintf("excludes %s", fe.Param())
	case "excludesall":
		return fmt.Sprintf("excludes all %s", fe.Param())
	case "excludesrune":
		return fmt.Sprintf("excludes rune %s", fe.Param())
	case "lowercase":
		return "lowercase only"
	case "multibyte":
		return "invalid multi-byte"
	case "number":
		return "invalid number"
	case "numeric":
		return "invalid numeric"
	case "printascii":
		return "invalid a printable ASCII"
	case "startsnotwith":
		return fmt.Sprintf("starts not with %s", fe.Param())
	case "startswith":
		return fmt.Sprintf("starts with %s", fe.Param())
	case "uppercase":
		return "uppercase only"
	// ------------------------------ Format ------------------------------
	case "base64":
		return "invalid Base64"
	case "base64url":
		return "invalid Base64 URL"
	case "base64rawurl":
		return "invalid Base64 raw URL"
	case "bic":
		return "invalid Business Identifier Code (ISO 9362)"
	case "bcp47_language_tag":
		return "invalid Language tag (BCP 47)"
	case "btc_addr":
		return "invalid Bitcoin address"
	case "btc_addr_bech32":
		return "invalid Bitcoin Bech32 Address (segwit)"
	case "credit_card":
		return "invalid Credit Card Number"
	case "mongodb":
		return "invalid MongoDB ObjectID"
	case "cron":
		return "invalid cron"
	case "spicedb":
		return "invalid SpiceDb ObjectID/Permission/Type"
	case "datetime":
		return "invalid Datetime"
	case "e164":
		return "invalid e164 formatted phone number"
	case "email":
		return "invalid email"
	case "eth_addr":
		return "invalid Ethereum address"
	case "hexadecimal":
		return "invalid Hexadecimal string"
	case "hexcolor":
		return "invalid Hexcolor string"
	case "hsl":
		return "invalid HSL string"
	case "hsla":
		return "invalid HSLA string"
	case "html":
		return "invalid HTTP tags"
	case "html_encoded":
		return "invalid HTML Encoded"
	case "isbn":
		return "invalid International Standard Book Number"
	case "isbn10":
		return "invalid International Standard Book Number 10"
	case "isbn13":
		return "invalid International Standard Book Number 13"
	case "issn":
		return "invalid International Standard Serial Number"
	case "iso3166_1_alpha2":
		return "invalid Two-letter country code (ISO 3166-1 alpha-2)"
	case "iso3166_1_alpha3":
		return "invalid Three-letter country code (ISO 3166-1 alpha-3)"
	case "iso3166_1_alpha_numeric":
		return "invalid Numeric country code (ISO 3166-1 numeric)"
	case "iso3166_2":
		return "invalid Country subdivision code (ISO 3166-2)"
	case "iso4217":
		return "invalid Currency code (ISO 4217)"
	case "json":
		return "invalid JSON"
	case "jwt":
		return "invalid JSON Web Token (JWT)"
	case "latitude":
		return "invalid Latitude"
	case "longitude":
		return "invalid Longitude"
	case "luhn_checksum":
		return "invalid Luhn Algorithm Checksum (for strings and (u)int)"
	case "postcode_iso3166_alpha2":
		return "invalid Postcode"
	case "postcode_iso3166_alpha2_field":
		return "invalid Postcode"
	case "rgb":
		return "invalid RGB string"
	case "rgba":
		return "invalid RGBA string"
	case "ssn":
		return "invalid Social Security Number SSN"
	case "timezone":
		return "invalid Timezone"
	case "uuid":
		return "invalid Universally Unique Identifier UUID"
	case "uuid3":
		return "invalid Universally Unique Identifier UUID v3"
	case "uuid3_rfc4122":
		return "invalid Universally Unique Identifier UUID v3 RFC4122"
	case "uuid4":
		return "invalid Universally Unique Identifier UUID v4"
	case "uuid4_rfc4122":
		return "invalid Universally Unique Identifier UUID v4 RFC4122"
	case "uuid5":
		return "invalid Universally Unique Identifier UUID v5"
	case "uuid5_rfc4122":
		return "invalid Universally Unique Identifier UUID v5 RFC4122"
	case "uuid_rfc4122":
		return "invalid Universally Unique Identifier UUID RFC4122"
	case "md4":
		return "invalid MD4 hash"
	case "md5":
		return "invalid MD5 hash"
	case "sha256":
		return "invalid SHA256 hash"
	case "sha384":
		return "invalid SHA384 hash"
	case "sha512":
		return "invalid SHA512 hash"
	case "ripemd128":
		return "invalid RIPEMD-128 hash"
	case "ripemd160":
		return "invalid RIPEMD-160 hash"
	case "tiger128":
		return "invalid TIGER128 hash"
	case "tiger160":
		return "invalid TIGER160 hash"
	case "tiger192":
		return "invalid TIGER192 hash"
	case "semver":
		return "invalid Semantic Versioning 2.0.0"
	case "ulid":
		return "invalid Universally Unique Lexicographically Sortable Identifier ULID"
	case "cve":
		return "Common Vulnerabilities and Exposures Identifier (CVE id)"

	// --------------- Fields ---------------
	case "eqcsfield":
		return fmt.Sprintf("field equals another field (relative) %s", fe.Param())
	case "eqfield":
		return fmt.Sprintf("field equals another field %s", fe.Param())
	case "fieldcontains":
		return fmt.Sprintf("field contains %s", fe.Param())
	case "fieldexcludes":
		return fmt.Sprintf("field excludes content %s", fe.Param())
	case "gtcsfield":
		return fmt.Sprintf("field greater than another relative field %s", fe.Param())
	case "gtecsfield":
		return fmt.Sprintf("field greater than or equal to another relative field %s", fe.Param())
	case "gtefield":
		return fmt.Sprintf("field greater than or equal to another field %s", fe.Param())
	case "gtfield":
		return fmt.Sprintf("field greater than another field %s", fe.Param())
	case "ltcsfield":
		return fmt.Sprintf("field less than another relative field %s", fe.Param())
	case "ltecsfield":
		return fmt.Sprintf("field less than or equal to another relative field %s", fe.Param())
	case "ltefield":
		return fmt.Sprintf("field less than another relative field %s", fe.Param())
	case "ltfield":
		return fmt.Sprintf("field less than another field %s", fe.Param())
	case "necsfield":
		return fmt.Sprintf("field does not equal another field (relative) %s", fe.Param())
	case "nefield":
		return fmt.Sprintf("field does not equal another field %s", fe.Param())

	// --------------- Network ---------------
	case "cidr":
		return "invalid Classless Inter-Domain Routing CIDR"
	case "cidrv4":
		return "invalid Classless Inter-Domain Routing CIDRv4"
	case "cidrv6":
		return "invalid Classless Inter-Domain Routing CIDRv6"
	case "datauri":
		return "invalid Data URL"
	case "fqdn":
		return "invalid Full Qualified Domain Name (FQDN)"
	case "hostname":
		return "invalid Hostname RFC 952"
	case "hostname_port":
		return "invalid HostPort"
	case "hostname_rfc1123":
		return "invalid Hostname RFC 1123"
	case "ip":
		return "invalid Internet Protocol Address IP"
	case "ip4_addr":
		return "invalid Internet Protocol Address IPv4"
	case "ip6_addr":
		return "invalid Internet Protocol Address IPv6"
	case "ip_addr":
		return "invalid Internet Protocol Address IP"
	case "ipv4":
		return "invalid Internet Protocol Address IPv4"
	case "ipv6":
		return "invalid Internet Protocol Address IPv6"
	case "mac":
		return "invalid Media Access Control Address MAC"
	case "tcp4_addr":
		return "invalid Transmission Control Protocol Address TCPv4"
	case "tcp6_addr":
		return "Transmission Control Protocol Address TCPv6"
	case "tcp_addr":
		return "invalid Transmission Control Protocol Address TCP"
	case "udp4_addr":
		return "invalid User Datagram Protocol Address UDPv4"
	case "udp6_addr":
		return "invalid User Datagram Protocol Address UDPv6"
	case "udp_addr":
		return "invalid User Datagram Protocol Address UDP"
	case "unix_addr":
		return "invalid Unix domain socket end point Address"
	case "uri":
		return "invalid URI string"
	case "url":
		return "invalid URL string"
	case "http_url":
		return "invalid HTTP URL string"
	case "url_encoded":
		return "invalid URL encoded"
	case "urn_rfc2141":
		return "invalid Urn RFC 2141 string"

	// --------------- Other ---------------
	case "dir":
		return "existing directory"
	case "dirpath":
		return "invalid directory path"
	case "file":
		return "existing file"
	case "filepath":
		return "invalid file path"
	case "image":
		return "invalid image"
	case "isdefault":
		return "is default"
	case "len":
		return fmt.Sprintf("length %s", fe.Param())
	case "max":
		return fmt.Sprintf("maximum %s", fe.Param())
	case "min":
		return fmt.Sprintf("minimum %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("one of %s", fe.Param())
	case "required":
		return "is required"
	case "required_if":
		return "required if"
	case "required_unless":
		return "required unless"
	case "required_with":
		return "required with"
	case "required_with_all":
		return "required with all"
	case "required_without":
		return "required without"
	case "required_without_all":
		return "required without all"
	case "excluded_if":
		return "excluded if"
	case "excluded_unless":
		return "excluded unless"
	case "excluded_with":
		return "excluded with"
	case "excluded_with_all":
		return "excluded with all"
	case "excluded_without":
		return "excluded without"
	case "excluded_without_all":
		return "excluded without all"
	case "unique":
		return "not unique"
		// --------------- Other ---------------
	case "iscolor":
		return "color format hexcolor|rgb|rgba|hsl|hsla"
	case "country_code":
		return "country format iso3166_1_alpha2|iso3166_1_alpha3|iso3166_1_alpha_numeric"
	}

	return fe.Error() // default error
}
