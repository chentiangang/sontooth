package main

import (
	"fmt"
	"net/http"
	"project/myself/xlog"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gin-gonic/gin"
)

func getHostNameByInstanceID(instanceid string) string {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	//awsRegion := "ap-northeast-1"
	awsRegion := "us-east-1"
	svc := ec2.New(sess, &aws.Config{
		Region: aws.String(awsRegion),
		//Endpoint:                      aws.String("ec2.ap-northeast-1.amazonaws.com"),
		Endpoint:                      aws.String("ec2-fips.us-east-1.amazonaws.com"),
		CredentialsChainVerboseErrors: aws.Bool(true),
	})

	var params ec2.DescribeInstancesInput
	params.Filters = []*ec2.Filter{
		{
			Name: aws.String("tag:Name"),
			Values: []*string{
				aws.String(strings.Join([]string{"*", "", "*"}, "")),
			},
		},
	}

	resp, err := svc.DescribeInstances(&params)
	if err != nil {
		xlog.LogError("%s", err)
	}

	for _, i := range resp.Reservations {
		for _, j := range i.Instances {
			if *j.InstanceId == instanceid {
				for _, k := range j.Tags {
					if *k.Key == "Name" {
						return *k.Value
					}
				}
			}
		}
	}
	return ""
}

func GetHostNameByInstanceID(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.Error(fmt.Errorf("instance id is not exist,id: %s", id))
		return
	}
	hostname := getHostNameByInstanceID(id)
	c.String(http.StatusOK, hostname)
}
