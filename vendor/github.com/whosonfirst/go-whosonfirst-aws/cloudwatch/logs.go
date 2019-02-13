package cloudwatch

import (
	"github.com/aws/aws-sdk-go/aws"
	aws_cloudwatchlogs "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/whosonfirst/go-whosonfirst-aws/session"
	_ "log"
)

func GetLogEvents(dsn string, cw_group string, cw_stream string) ([]*aws_cloudwatchlogs.OutputLogEvent, error) {

	cw_sess, err := session.NewSessionWithDSN(dsn)

	if err != nil {
		return nil, err
	}

	cw_svc := aws_cloudwatchlogs.New(cw_sess)

	return GetLogEventsWithService(cw_svc, cw_group, cw_stream)
}

func GetLogEventsWithService(cw_svc *aws_cloudwatchlogs.CloudWatchLogs, cw_group string, cw_stream string) ([]*aws_cloudwatchlogs.OutputLogEvent, error) {

	// something something something something that emits to channels something something something
	// (20190213/thisisaaronland)

	events := make([]*aws_cloudwatchlogs.OutputLogEvent, 0)

	var cursor string

	for {

		cw_req := &aws_cloudwatchlogs.GetLogEventsInput{
			LogGroupName:  aws.String(cw_group),
			LogStreamName: aws.String(cw_stream),
			StartFromHead: aws.Bool(true),
		}

		if cursor != "" {
			cw_req.NextToken = aws.String(cursor)
		}

		cw_rsp, err := cw_svc.GetLogEvents(cw_req)

		if err != nil {
			return nil, err
		}

		for _, e := range cw_rsp.Events {
			events = append(events, e)
		}

		// sigh... (20190213/thisisaaronland)

		if *cw_rsp.NextForwardToken != "" && *cw_rsp.NextForwardToken != cursor {
			cursor = *cw_rsp.NextForwardToken
		} else {
			break
		}

	}

	return events, nil
}
