package dispatchers

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/whosonfirst/go-webhookd"
	"github.com/whosonfirst/go-whosonfirst-aws/session"
)

type LambdaDispatcher struct {
	webhookd.WebhookDispatcher
	LambdaFunction string
	LambdaService  *lambda.Lambda
}

func NewLambdaDispatcher(session_dsn string, function string) (*LambdaDispatcher, error) {

	sess, err := session.NewSessionWithDSN(session_dsn)

	if err != nil {
		return nil, err
	}

	svc := lambda.New(sess)

	d := LambdaDispatcher{
		LambdaFunction: function,
		LambdaService:  svc,
	}

	return &d, nil
}

func (d *LambdaDispatcher) Dispatch(body []byte) *webhookd.WebhookError {

	payload, err := json.Marshal(body)

	if err != nil {
		return &webhookd.WebhookError{Code: 999, Message: err.Error()}
	}

	input := &lambda.InvokeInput{
		FunctionName: aws.String(d.LambdaFunction),
		Payload:      payload,
	}

	_, err = d.LambdaService.Invoke(input)

	if err != nil {
		return &webhookd.WebhookError{Code: 999, Message: err.Error()}
	}

	return nil
}
