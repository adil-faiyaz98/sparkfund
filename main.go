package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// Response is of type APIGatewayProxyResponse since we're leveraging API Gateway
type Response events.APIGatewayProxyResponse

// CloudWatchClient represents the CloudWatch client
type CloudWatchClient struct {
	client *cloudwatch.CloudWatch
}

// NewCloudWatchClient creates a new CloudWatch client
func NewCloudWatchClient() (*CloudWatchClient, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Update with your AWS region
	})
	if err != nil {
		return nil, err
	}

	return &CloudWatchClient{
		client: cloudwatch.New(sess),
	}, nil
}

// PutMetric sends a custom metric to CloudWatch
func (c *CloudWatchClient) PutMetric(namespace, metricName string, value float64, unit string) error {
	_, err := c.client.PutMetricData(&cloudwatch.PutMetricDataInput{
		Namespace: aws.String(namespace),
		MetricData: []*cloudwatch.MetricDatum{
			{
				MetricName: aws.String(metricName),
				Value:      aws.Float64(value),
				Unit:       aws.String(unit),
				Timestamp:  aws.Time(time.Now()),
			},
		},
	})
	return err
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	startTime := time.Now()
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	// Initialize CloudWatch client
	cwClient, err := NewCloudWatchClient()
	if err != nil {
		log.Printf("Error creating CloudWatch client: %v", err)
		// Continue execution even if CloudWatch client creation fails
	}

	// Your application logic goes here
	responseBody := map[string]interface{}{
		"message": "Money-Pulse Lambda function executed successfully!",
		"request": request.Path,
	}

	body, err := json.Marshal(responseBody)
	if err != nil {
		return Response{StatusCode: 500}, err
	}

	// Record execution time and invocation count as CloudWatch metrics
	if cwClient != nil {
		execTime := time.Since(startTime).Milliseconds()
		if err := cwClient.PutMetric("MoneyPulseMetrics", "LambdaExecutionTime", float64(execTime), "Milliseconds"); err != nil {
			log.Printf("Error publishing execution time metric: %v", err)
		}

		if err := cwClient.PutMetric("MoneyPulseMetrics", "APIInvocationCount", 1.0, "Count"); err != nil {
			log.Printf("Error publishing invocation count metric: %v", err)
		}
	}

	return Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
