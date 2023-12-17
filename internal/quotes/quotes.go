package quotes

import (
	"math/rand"
	"time"
)

var Quotes = []string{
	"There is nothing impossible to they who will try.",
	"Success is not final, failure is not fatal: it is the courage to continue that counts.",
	"At the end of the day, whether or not those people are comfortable with how you're living your life doesn't matter. What matters is whether you're comfortable with it.",
	"It is during our darkest moments that we must focus to see the light.",
	"Believe you can and you're halfway there.",
}

type QuoteProvider struct {
	quotes []string
	rand   *rand.Rand
}

func (q *QuoteProvider) Provide() string {
	n := q.rand.Int() % len(q.quotes)
	return q.quotes[n]
}

func New() *QuoteProvider {
	return &QuoteProvider{
		quotes: Quotes,
		rand:   rand.New(rand.NewSource(time.Now().Unix())),
	}
}
