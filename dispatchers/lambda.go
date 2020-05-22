package dispatchers

import (
	"context"
	"encoding/json"
	"github.com/aaronland/go-aws-session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/whosonfirst/go-webhookd/v2"
	"net/url"
)

type LambdaDispatcher struct {
	webhookd.WebhookDispatcher
	LambdaFunction string
	LambdaService  *lambda.Lambda
}

func NewLambdaDispatcher(ctx context.Context, uri string) (webhookd.WebhookDispatcher, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	lambda_dsn := q.Get("dsn")
	lambda_function := q.Get("function")

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

func (d *LambdaDispatcher) Dispatch(ctx context.Context, body []byte) *webhookd.WebhookError {

	select {
	case <-ctx.Done():
		return nil
	default:
		// pass
	}

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
