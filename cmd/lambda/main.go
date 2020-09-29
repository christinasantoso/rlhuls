package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	baseIndexURI = "https://github.com/ronaudinho/rlhuls/tree/master/static"
	baseRawURI   = "https://raw.githubusercontent.com/ronaudinho/rlhuls/tree/master"
)

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.Path {
	case "/":
		var pages []string
		buf := bytes.NewBuffer([]byte{})
		tmpl, err := template.New("index").Parse(`
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{ .Title }}</title>
	</head>
	<body>
		{{ range .Pages }}
			<a href="{{ . }}">{{ . }}</a></br>
		{{end}}
	</body>
</html>
			`)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusPaymentRequired,
				Body:       err.Error(),
			}, err
		}
		err = tmpl.Execute(buf, struct {
			Title string
			Pages []string
		}{
			Title: baseIndexURI,
			Pages: pages,
		})
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusPaymentRequired,
				Body:       err.Error(),
			}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusPaymentRequired,
			Body:       buf.String(),
		}, nil
	default:
		page := req.Path
		uri := fmt.Sprintf("%s/%s", baseRawURI, page)
		res, err := http.Get(uri)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusPaymentRequired,
				Body:       err.Error(),
			}, err
		}
		defer res.Body.Close()
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusPaymentRequired,
				Body:       err.Error(),
			}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusPaymentRequired,
			Body:       string(b),
		}, nil
	}
}

func main() {
	lambda.Start(handler)
}
