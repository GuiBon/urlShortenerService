package domain

// URLStatistic represents an URL statistic with a shortened counter and an accessed counter
type URLStatistic struct {
	URL              string
	ShortenedCounter int
	AccessedCounter  int
}
