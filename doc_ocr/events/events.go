package events

import (
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	uuid "github.com/satori/go.uuid"
)

var (
	stream = flag.String("stream", "TESTE-stream", "FIAP")
)

func PutStream(streamNameRDN string, data string) *string {
	s := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("teste", "teste", ""),
		Region:      aws.String(endpoints.UsWest2RegionID),
		Endpoint:    aws.String("http://localhost:4568"),
	}))
	kc := kinesis.New(s)
	streamName := aws.String(streamNameRDN)
	putOutput, err := kc.PutRecord(&kinesis.PutRecordInput{
		Data:         []byte(data),
		StreamName:   streamName,
		PartitionKey: aws.String("ocr"),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", putOutput)
	return putOutput.ShardId
}
func DeleteStream(streamNameRDN string) {
	s := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("teste", "teste", ""),
		Region:      aws.String(endpoints.UsWest2RegionID),
		Endpoint:    aws.String("http://localhost:4568"),
	}))
	kc := kinesis.New(s)
	streamName := aws.String(streamNameRDN)
	deleteOutput, err := kc.DeleteStream(&kinesis.DeleteStreamInput{
		StreamName: streamName,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", deleteOutput)
}

func GetRecords(streamNameRDN string, ShardId *string) {

	s := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("teste", "teste", ""),
		Region:      aws.String(endpoints.UsWest2RegionID),
		Endpoint:    aws.String("http://localhost:4568"),
	}))
	kc := kinesis.New(s)
	streamName := aws.String(streamNameRDN)

	// retrieve iterator
	iteratorOutput, err := kc.GetShardIterator(&kinesis.GetShardIteratorInput{
		// Shard Id is provided when making put record(s) request.
		ShardId:           ShardId,
		ShardIteratorType: aws.String("TRIM_HORIZON"),
		// ShardIteratorType: aws.String("AT_SEQUENCE_NUMBER"),
		// ShardIteratorType: aws.String("LATEST"),
		StreamName: streamName,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", iteratorOutput)

	// get records use shard iterator for making request
	records, err := kc.GetRecords(&kinesis.GetRecordsInput{
		ShardIterator: iteratorOutput.ShardIterator,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", records)

	// and, you can iteratively make GetRecords request using records.NextShardIterator
	recordsSecond, err := kc.GetRecords(&kinesis.GetRecordsInput{
		ShardIterator: records.NextShardIterator,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", recordsSecond)

}

func CreateStream() *string {
	streamNameRDN := "TESTE"
	s := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("teste", "teste", ""),
		Region:      aws.String(endpoints.UsWest2RegionID),
		Endpoint:    aws.String("http://localhost:4568"),
	}))
	kc := kinesis.New(s)

	streamName := aws.String(streamNameRDN)
	out, err := kc.CreateStream(&kinesis.CreateStreamInput{
		ShardCount: aws.Int64(1),
		StreamName: streamName,
	})
	if err != nil {
		DeleteStream(streamNameRDN)
	}
	fmt.Printf("%v\n", out)
	return streamName
}
func GenerateNewUUID() string {
	u1 := uuid.Must(uuid.NewV4())
	return u1.String()
}
