package lib

import (
	"log"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// TODO: return parse error
func ParseHexInt(x string) int64 {
	x = strings.Replace(x, "0x", "", -1)
	value, err := strconv.ParseInt(x, 16, 64)
	if err != nil {
		log.Fatalf("failed to parse hex number '%v': %v", x, err)
	}
	return value
}

// IsValidAddress validate hex address
func IsValidAddress(address string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
}

// IsValidAddressSlice returns true if each
// element of slice is a valid Ethereum address
func IsValidAddressSlice(addresses []string) bool {
	for _, address := range addresses {
		if !IsValidAddress(address) {
			return false
		}
	}
	return true
}

func ParseTokenAmountSlice(amounts []string) ([]*big.Int, bool) {
	as := make([]*big.Int, len(amounts))
	for i, amount := range amounts {
		a, ok := ParseBig256(amount)
		if !ok {
			return nil, false
		}
		as[i] = a
	}
	return as, true
}

// GetIntEnv retreives int64 environment variable named by the key
func GetIntEnv(key string, def int64) int64 {
	envValue := os.Getenv(key)
	if envValue == "" {
		return def
	}
	value, err := strconv.ParseInt(envValue, 10, 64)
	if err != nil {
		log.Fatalf("failed to parse %v %v: %v", key, envValue, err)
	}
	return value
}
