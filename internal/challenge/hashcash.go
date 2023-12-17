package challenge

import (
	"crypto/sha1"
	"errors"
	"fmt"
)

var (
	ErrTooManyIterations = errors.New("too many iterations")
)

// hashcash is a cryptographic hash-based proof-of-work algorithm
// that requires a selectable amount of work to compute,
// but the proof can be verified efficiently.
// https://en.wikipedia.org/wiki/Hashcash
type hashcash struct {
	// Resource data string being transmitted, e.g., an IP address or email address.
	Resource string

	// String of random characters, encoded in base-64 format.
	Rand string

	// The time that the message was sent, in the format YYMMDD[hhmm[ss]]
	Date uint64

	// Binary counter, encoded in base-64 format.
	Counter uint64

	// format version, 1 (which supersedes version 0).
	Ver uint8

	// Number of "partial pre-image" (zero) bits in the hashed code.
	Bits uint8
}

func (hc *hashcash) Header() string {
	return fmt.Sprintf(
		"%d|%d|%d|%s|%s|%d",
		hc.Ver, hc.Bits, hc.Date, hc.Resource, hc.Rand, hc.Counter,
	)
}

func (hc *hashcash) Hash() string {
	hasher := sha1.New()
	hasher.Write([]byte(hc.Header()))
	s := hasher.Sum(nil)
	return fmt.Sprintf("%x", s)
}

func validateZeroBits(hash string, bits uint8) bool {
	if int(bits) > len(hash) {
		return false
	}

	// check that hash actually contains zeroes
	for _, b := range hash[:bits] {
		if b != '0' {
			return false
		}
	}

	return true
}

func (hc *hashcash) Bruteforce(iterations uint64) error {
	for hc.Counter <= iterations {
		hash := hc.Hash()
		if validateZeroBits(hash, hc.Bits) {
			return nil
		}

		hc.Counter++
	}

	return fmt.Errorf("%w: could not solve %s with %d max iterations", ErrTooManyIterations, hc.Header(), iterations)
}
