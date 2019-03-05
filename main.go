package kata_go_cloudwatch

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/fpmoles/kata_go_cloudwatch/cloud"
	"log"
	"os"
)

func main() {
	access := os.Getenv("AWS_ACCESS_KEY_ID")
	secret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	if access == "" || secret == "" {
		log.Fatal("error loading credentials")
	}
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(access, secret, ""),
	})
	if err != nil {
		log.Fatal("error creating  session")
	}

	alertOps := cloud.NewAwsAlertOps(sess)
	r := alertOps.CreateRootLoginAlert()

	fmt.Printf("log alert was created with response: %s", r)
}
