package main

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

var prices = map[string]float64{
	"ETH": 999.99,
	"BTC": 20000.0,
	"GG":  1000000.0,
}

type PriceService interface {
	FetchPrice(context.Context, string) (float64, error)
}

type priceService struct{}

func (s *priceService) FetchPrice(_ context.Context, ticker string) (float64, error) {
	price, ok := prices[ticker]
	if !ok {
		return 0.0, fmt.Errorf("price for ticker (%s) is not available", ticker)
	}
	return price, nil
}

type loggingService struct {
	next priceService
}

func (s loggingService) FetchPrice(ctx context.Context, ticker string) (price float64, err error) {
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

var priceMocks = map[string]float64{
	"BTC": 20_000.0,
	"ETH": 200.0,
	"DGH": 100_1000.0,
}

func MockPriceFetcher(ctx context.Context, ticker string) (float64, error) {
	time.Sleep(100 * time.Millisecond)
	price, ok := priceMocks[ticker]
	if !ok {
		return price, fmt.Errorf("the given ticker (%s) is not supported", ticker)
	}
	return price, nil

}
