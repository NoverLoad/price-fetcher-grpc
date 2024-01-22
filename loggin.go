package main

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

type loggingServiceNew struct {
	next PriceService
}

func NewLoggingService(next PriceService) PriceService {
	return &loggingServiceNew{
		next: next,
	}
}

func (s *loggingServiceNew) FetchPrice(ctx context.Context, ticker string) (price float64, err error) {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"requestID": ctx.Value("requestID"),
			"took":      time.Since(begin),
			"err":       err,
			"price":     price,
		}).Info("fetchPrice")
	}(time.Now())
	return s.next.FetchPrice(ctx, ticker)
}
