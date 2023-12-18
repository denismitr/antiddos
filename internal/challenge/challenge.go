package challenge

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidHeader             = errors.New("invalid header")
	ErrChallengeDurationExceeded = errors.New("challenge duration exceeded")
)

const (
	HeaderDelimiter = "|"
)

// validator serves to verify rand values in hashcash
type validator interface {
	Validate(key string) bool
	Remember(key string)
}

type Challenge struct {
	zeroes        uint8
	maxDuration   uint64
	maxIterations uint64
	r             *rand.Rand
	now           func() time.Time
	randomizer    func() int
	validator     validator
}

func createDefaultRandomizer() func() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return func() int {
		return r.Intn(100_000)
	}
}

func New(
	store validator,
	zeroes uint8,
	maxDuration uint64,
) *Challenge {
	return &Challenge{
		validator:     store,
		zeroes:        zeroes,
		maxDuration:   maxDuration,
		maxIterations: math.MaxUint64,
		now:           time.Now,
		randomizer:    createDefaultRandomizer(),
	}
}

func (c *Challenge) SetRandomizer(r func() int) {
	c.randomizer = r
}

func (c *Challenge) SetNow(now func() time.Time) {
	c.now = now
}

func (c *Challenge) SetMaxIterations(maxIterations uint64) {
	c.maxIterations = maxIterations
}

func (c *Challenge) Create(resource string) (string, error) {
	random := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", c.randomizer())))

	c.validator.Remember(random)

	hc := hashcash{
		Ver:      1,
		Bits:     c.zeroes,
		Date:     uint64(c.now().Unix()),
		Resource: resource,
		Rand:     random,
		Counter:  0,
	}

	return hc.Header(), nil
}

func (c *Challenge) Solve(header string) (string, error) {
	hc, err := c.headerToHashcash(header)
	if err != nil {
		return "", err
	}

	if err := c.validate(hc); err != nil {
		return "", err
	}

	iterations := hc.Counter
	if iterations == 0 {
		iterations = c.maxIterations
	}

	if err := hc.Bruteforce(iterations); err != nil {
		return "", fmt.Errorf("%w failed to solve hashcash: %v", ErrTooManyIterations, err)
	}

	return hc.Header(), nil
}

func (c *Challenge) validate(hc *hashcash) error {
	if !c.validator.Validate(hc.Rand) {
		return fmt.Errorf("header seems to be milicious")
	}

	if c.zeroes != hc.Bits {
		return fmt.Errorf("amount of zeroes does not match the config")
	}

	if uint64(c.now().Unix())-hc.Date > c.maxDuration {
		return ErrChallengeDurationExceeded
	}

	return nil
}

func (c *Challenge) headerToHashcash(header string) (*hashcash, error) {
	segments := strings.Split(header, HeaderDelimiter)
	if len(segments) != 6 {
		return nil, fmt.Errorf("%w: expected 6 segments in header but got %s", ErrInvalidHeader, header)
	}

	ver, err := strconv.Atoi(segments[0])
	if err != nil {
		return nil, fmt.Errorf("%w: version is invalid: %v", ErrInvalidHeader, err)
	}

	bits, err := strconv.Atoi(segments[1])
	if err != nil {
		return nil, fmt.Errorf("%w: bits are invalid: %v", ErrInvalidHeader, err)
	}

	date, err := strconv.ParseUint(segments[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: date is invalid: %v", ErrInvalidHeader, err)
	}

	counter, err := strconv.Atoi(segments[5])
	if err != nil {
		return nil, fmt.Errorf("%w: date is invalid: %v", ErrInvalidHeader, err)
	}

	return &hashcash{
		Ver:      uint8(ver),  // todo: validate int size
		Bits:     uint8(bits), // todo: validate int size
		Date:     date,
		Resource: segments[3],
		Rand:     segments[4],
		Counter:  uint64(counter),
	}, nil
}
