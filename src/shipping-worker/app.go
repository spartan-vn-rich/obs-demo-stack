package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	// ... Init Tracer ...
    tracer := otel.Tracer("shipping-worker")

	// AWS/Localstack Config
	cfg, _ := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: os.Getenv("SQS_ENDPOINT")}, nil
			})),
	)
	client := sqs.NewFromConfig(cfg)
	queueURL := os.Getenv("SQS_QUEUE_URL")

	for {
		// Start a root span for the polling cycle
		ctx, span := tracer.Start(context.Background(), "poll-sqs")
		
		output, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            &queueURL,
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     5,
		})

		if err == nil && len(output.Messages) > 0 {
			for _, msg := range output.Messages {
				processMessage(ctx, tracer, msg.Body)
				client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
					QueueUrl:      &queueURL,
					ReceiptHandle: msg.ReceiptHandle,
				})
			}
		}
		span.End()
		time.Sleep(1 * time.Second)
	}
}

func processMessage(ctx context.Context, tracer trace.Tracer, body *string) {
	_, span := tracer.Start(ctx, "process_message")
	defer span.End()
	
    // Simulate work
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("{\"level\":\"info\",\"msg\":\"Consumed message\",\"body\":%q}\n", *body)
}