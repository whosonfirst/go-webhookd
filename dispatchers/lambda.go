package dispatchers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/aaronland/go-aws-session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/whosonfirst/go-webhookd/v2"
	"net/url"
)

func init() {

	ctx := context.Background()
	err := RegisterDispatcher(ctx, "lambda", NewLambdaDispatcher)

	if err != nil {
		panic(err)
	}
}

type LambdaDispatcher struct {
	webhookd.WebhookDispatcher
	LambdaFunction  string
	LambdaService   *lambda.Lambda
	invocation_type string
}

func NewLambdaDispatcher(ctx context.Context, uri string) (webhookd.WebhookDispatcher, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	lambda_function := u.Host

	q := u.Query()

	lambda_dsn := q.Get("dsn")

	lambda_sess, err := session.NewSessionWithDSN(lambda_dsn)

	if err != nil {
		return nil, err
	}

	invocation_type := q.Get("invocation_type")

	switch invocation_type {
	case "":
		invocation_type = "RequestResponse"
	case "RequestResponse", "Event", "DryRun":
		// pass
	default:
		return nil, errors.New("Invalid invocation_type parameter")
	}

	lambda_svc := lambda.New(lambda_sess)

	d := LambdaDispatcher{
		LambdaFunction:  lambda_function,
		LambdaService:   lambda_svc,
		invocation_type: invocation_type,
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

	// I don't understand why I need to base64 encode this...
	// (20200526/thisisaaronland)

	enc_body := base64.StdEncoding.EncodeToString(body)

	payload, err := json.Marshal(enc_body)

	if err != nil {
		return &webhookd.WebhookError{Code: 999, Message: err.Error()}
	}

	input := &lambda.InvokeInput{
		FunctionName:   aws.String(d.LambdaFunction),
		Payload:        payload,
		InvocationType: aws.String(d.invocation_type),
	}

	_, err = d.LambdaService.Invoke(input)

	if err != nil {
		return &webhookd.WebhookError{Code: 999, Message: err.Error()}
	}

	return nil
}
