package main

import (
	"flag"
)

func main() {
	// client := client.New("http://localhost:3000")
	// price, err := client.FetchPrice(context.Background(), "ET")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%+v\n", price)
	// return
	var (
		jsonAddr = flag.String("listenaddr", ":3000", "listen address of the json transport")
		grpcAddr = flag.String("listenaddr", ":4000", "listen address of the grpc transport")
	)
	flag.Parse()
	svc := NewLoggingService(NewMetricService(&priceFetcher{}))
	//svc := loggingService{priceService{}}
	// price, err := svc.FetchPrice(context.Background(), "ETH")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	go makeGRPCServerAndRun(*grpcAddr, svc)
	server := NewJSONAPIServer(*jsonAddr, svc)
	server.RUN()

	//fmt.Println(price)
}
