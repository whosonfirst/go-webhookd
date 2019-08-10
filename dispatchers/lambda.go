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

func NewLambdaDispatcher(lambda_dsn string, lambda_function string) (*LambdaDispatcher, error) {

	lambda_sess, err := session.NewSessionWithDSN(lambda_dsn)

	if err != nil {
		return nil, err
	}

	lambda_svc := lambda.New(lambda_sess)

	d := LambdaDispatcher{
		LambdaFunction: lambda_function,
		LambdaService:  lambda_svc,
	}

	return &d, nil
}

func (d *LambdaDispatcher) Dispatch(body []byte) *webhookd.WebhookError {

	payload, err := json.Marshal(string(body))

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
