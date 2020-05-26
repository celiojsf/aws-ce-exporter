package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics(sess *session.Session) {
	go func() {
		for {
			ceclient := costexplorer.New(sess)

			dateNow := time.Now().Format("2006-01-02")
			yesterday := time.Now().Add(time.Hour * -24).Format("2006-01-02")
			gran := costexplorer.GranularityDaily
			blendedcost := costexplorer.MetricBlendedCost
			groupdef := costexplorer.GroupDefinition{}
			groupdef.SetKey("SERVICE")
			groupdef.SetType(costexplorer.GroupDefinitionTypeDimension)

			req, resp := ceclient.GetCostAndUsageRequest(&costexplorer.GetCostAndUsageInput{
				TimePeriod: &costexplorer.DateInterval{
					Start: &yesterday,
					End:   &dateNow,
				},
				Granularity: &gran,
				Metrics:     []*string{&blendedcost},
				GroupBy:     []*costexplorer.GroupDefinition{&groupdef},
			})

			err := req.Send()
			if err == nil {

				for _, item := range resp.ResultsByTime[0].Groups {
					fmt.Println(item)
					value, _ := strconv.ParseFloat(*item.Metrics["BlendedCost"].Amount, 64)
					serviceCost.With(prometheus.Labels{"service": *item.Keys[0]}).Set(value)
				}

				time.Sleep(24 * time.Hour)

			} else {
				fmt.Println("Error while requesting values from AWS, trying again in 5 minutes")
				fmt.Println(err)
				time.Sleep(5 * time.Minute)
			}

		}
	}()
}

var (
	serviceCost = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aws_service_daily_cost",
		},
		[]string{
			"service",
		},
	)
)

func main() {

	sess, err := session.NewSession()

	if err != nil {
		fmt.Println("Error ocurred while getting a new session")
		fmt.Println(err)
		os.Exit(1)
	}

	recordMetrics(sess)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)

}
