package cloud

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"log"
	"os"
)

var LOG_GROUP_NAME = "LOG_GROUP_NAME"
var FAILURE = "failure"
var SUCCESS = "success"
var FILTER_NAME = "RootAccountUsage"
var METRIC_NAME = "RootAccountUsageCount"
var FILTER = "{ $.userIdentity.type = \"Root\" && $.userIdentity.invokedBy NOT EXISTS && $.eventType != \"AwsServiceEvent\" }"
var METRIC_VALUE = "1"
var NAMESPACE = "CDAS"

type AWSAlertOps struct {
	alertClient     cloudwatchiface.CloudWatchAPI
	alertLogsClient cloudwatchlogsiface.CloudWatchLogsAPI
}

func (ops *AWSAlertOps) CreateRootLoginAlert() string {
	logGroupName := os.Getenv(LOG_GROUP_NAME)
	r := ops.createLogGroup(logGroupName)
	if r == FAILURE {
		return FAILURE
	}
	r = ops.createMetricFilter(logGroupName)
	if r == FAILURE {
		return FAILURE
	}

	return "success"
}

func (ops *AWSAlertOps) createLogGroup(logGroupName string) string {
	request := &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: &logGroupName,
	}
	response, err := ops.alertLogsClient.CreateLogGroup(request)
	if err != nil {
		log.Printf("error creating log group: %s", err)
		return FAILURE
	}
	return response.GoString()
}

func (ops *AWSAlertOps) createMetricFilter(logGroupName string) string {

	filter := &cloudwatchlogs.MetricTransformation{
		MetricName:      &METRIC_NAME,
		MetricNamespace: &NAMESPACE,
		MetricValue:     &METRIC_VALUE,
		DefaultValue:    nil,
	}
	transformations := []*cloudwatchlogs.MetricTransformation{filter}

	input := &cloudwatchlogs.PutMetricFilterInput{
		LogGroupName:          &logGroupName,
		FilterName:            &FILTER_NAME,
		FilterPattern:         &FILTER,
		MetricTransformations: transformations,
	}
	response, err := ops.alertLogsClient.PutMetricFilter(input)
	if err != nil {
		log.Printf("error creating metric filter: %s", err)
		return FAILURE
	}
	return response.GoString()
}

func NewAwsAlertOps(session *session.Session) *AWSAlertOps {
	alertClient := cloudwatch.New(session)
	alertLogsClient := cloudwatchlogs.New(session)

	return &AWSAlertOps{
		alertClient:     alertClient,
		alertLogsClient: alertLogsClient,
	}
}
