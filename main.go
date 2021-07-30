package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

type Response struct {
	Message string `json:"message"`
}

func Handler() error {
	sess := session.Must(session.NewSession())

	aurora := rds.New(sess)
	desc, err := aurora.DescribeDBClusters(&rds.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(os.Getenv("CLUSTER_NAME")),
	})
	if err != nil { // TODO: エラーハンドリング
		return err
	}
	for _, c := range desc.DBClusters {
		status := *c.Status
		log.Println("Status: ", status)
		if status == "available" {
			if _, err := aurora.StopDBCluster(&rds.StopDBClusterInput{
				DBClusterIdentifier: c.DBClusterIdentifier,
			}); err != nil {
				return err
			}
		} else {
			if _, err := aurora.StartDBCluster(&rds.StartDBClusterInput{
				DBClusterIdentifier: c.DBClusterIdentifier,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
