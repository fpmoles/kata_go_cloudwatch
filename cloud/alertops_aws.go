package cloud

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"os"
)

var LOG_GROUP_NAME = "LOG_GROUP_NAME"
var FAILURE = "failure"
var SUCCESS = "success"

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

	return "success"
}

func (ops *AWSAlertOps) createLogGroup(logGroupName string) string {
	request := &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: &logGroupName,
	}
	response, err := ops.alertLogsClient.CreateLogGroup(request)
	if err != nil {
		return FAILURE
	}
	return response.GoString()
}

func (ops *AWSAlertOps) createMetricFilter(logGroupName string) string {

}

func NewAwsAlertOps(session *session.Session) *AWSAlertOps {
	alertClient := cloudwatch.New(session)
	alertLogsClient := cloudwatchlogs.New(session)

	return &AWSAlertOps{
		alertClient:     alertClient,
		alertLogsClient: alertLogsClient,
	}
}
