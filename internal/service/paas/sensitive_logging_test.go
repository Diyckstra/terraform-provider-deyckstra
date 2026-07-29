package paas

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	sdkpaas "github.com/aws/aws-sdk-go/service/paas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceServiceCreateDoesNotLogELKPassword(t *testing.T) {
	t.Parallel()

	const password = "elk-password-secret"

	logger := &recordingSDKLogger{}
	conn := testPaaSClientWithDebugBodyLogging(logger, roundTripFunc(func(request *http.Request) (*http.Response, error) {
		switch {
		case request.Method == http.MethodPost && request.URL.Path == "/services":
			return jsonResponse(request, `{
				"service": {
					"id": "fm-cluster-12345678",
					"serviceType": "elk",
					"serviceClass": "logging"
				}
			}`), nil
		case request.Method == http.MethodGet && request.URL.Path == "/services/fm-cluster-12345678":
			return jsonResponse(request, `{
				"service": {
					"id": "fm-cluster-12345678",
					"name": "tf-elk-test",
					"serviceType": "elk",
					"serviceClass": "logging",
					"status": "READY",
					"instanceType": "c5.large",
					"rootVolumeType": "gp2",
					"rootVolumeSize": 32,
					"dataVolumeType": "gp2",
					"dataVolumeSize": 32,
					"parameters": {
						"password": "`+password+`",
						"version": "8.17"
					}
				}
			}`), nil
		case request.Method == http.MethodPut &&
			strings.HasPrefix(request.URL.Path, "/services/fm-cluster-12345678/"):
			return jsonResponse(request, `{}`), nil
		default:
			t.Fatalf("unexpected PaaS request: %s %s", request.Method, request.URL.Path)
			return nil, nil
		}
	}))

	resourceData := schema.TestResourceDataRaw(t, ResourceService().Schema, map[string]interface{}{
		"name":          "tf-elk-test",
		"instance_type": "c5.large",
		"root_volume": []interface{}{
			map[string]interface{}{
				"type": "gp2",
				"size": 32,
			},
		},
		"data_volume": []interface{}{
			map[string]interface{}{
				"type": "gp2",
				"size": 32,
			},
		},
		"nodes": []interface{}{
			map[string]interface{}{
				"main": []interface{}{
					map[string]interface{}{
						"role": "node",
					},
				},
			},
		},
		"subnet_ids": []interface{}{"subnet-12345678"},
		"elk": []interface{}{
			map[string]interface{}{
				"class":    "logging",
				"password": password,
				"version":  "8.17",
			},
		},
	})

	diagnostics := resourceServiceCreate(
		context.Background(),
		resourceData,
		testPaaSClientMeta(conn),
	)
	if diagnostics.HasError() {
		t.Fatalf("creating ELK service: %#v", diagnostics)
	}
	if got := resourceData.Id(); got != "fm-cluster-12345678" {
		t.Fatalf("unexpected service ID: %q", got)
	}
	if strings.Contains(logger.String(), password) {
		t.Fatal("ELK password was written to AWS SDK debug logs")
	}
}

func TestFindServiceByIDDoesNotLogELKPassword(t *testing.T) {
	t.Parallel()

	const password = "elk-password-secret"

	logger := &recordingSDKLogger{}
	conn := testPaaSClientWithDebugBodyLogging(logger, roundTripFunc(func(request *http.Request) (*http.Response, error) {
		return jsonResponse(request, `{
			"service": {
				"id": "fm-cluster-12345678",
				"serviceType": "elk",
				"serviceClass": "logging",
				"parameters": {
					"password": "`+password+`",
					"version": "8.17"
				}
			}
		}`), nil
	}))

	service, err := FindServiceByID(conn, "fm-cluster-12345678")
	if err != nil {
		t.Fatalf("reading ELK service: %s", err)
	}
	if got := aws.StringValue(service.ServiceType); got != "elk" {
		t.Fatalf("unexpected service type: %q", got)
	}
	if strings.Contains(logger.String(), password) {
		t.Fatal("ELK password was written to AWS SDK debug logs")
	}
}

func TestResourceLogstashPipelineCreateDoesNotLogConfiguration(t *testing.T) {
	t.Parallel()

	const configuration = "input { stdin {} } # pipeline-secret"

	logger := &recordingSDKLogger{}
	conn := testPaaSClientWithDebugBodyLogging(logger, roundTripFunc(func(request *http.Request) (*http.Response, error) {
		switch {
		case request.Method == http.MethodPost &&
			request.URL.Path == "/services/fm-cluster-12345678/logstash-pipelines":
			return jsonResponse(request, `{
				"logstashPipeline": {
					"id": "logstash-pipeline-87654321",
					"name": "terraform-qa",
					"configuration": "`+configuration+`"
				}
			}`), nil
		case request.Method == http.MethodGet &&
			request.URL.Path == "/services/fm-cluster-12345678":
			return jsonResponse(request, `{
				"service": {
					"id": "fm-cluster-12345678",
					"serviceType": "elk",
					"serviceClass": "logging",
					"status": "READY",
					"parameters": {
						"version": "8.17"
					}
				}
			}`), nil
		case request.Method == http.MethodGet &&
			request.URL.Path == "/services/fm-cluster-12345678/logstash-pipelines":
			return jsonResponse(request, `{
				"logstashPipelines": [
					{
						"id": "logstash-pipeline-87654321",
						"name": "terraform-qa",
						"configuration": "`+configuration+`"
					}
				]
			}`), nil
		default:
			t.Fatalf("unexpected PaaS request: %s %s", request.Method, request.URL.Path)
			return nil, nil
		}
	}))

	resourceData := schema.TestResourceDataRaw(t, ResourceLogstashPipeline().Schema, map[string]interface{}{
		"service_id":    "fm-cluster-12345678",
		"name":          "terraform-qa",
		"configuration": configuration,
	})

	diagnostics := resourceLogstashPipelineCreate(
		context.Background(),
		resourceData,
		testPaaSClientMeta(conn),
	)
	if diagnostics.HasError() {
		t.Fatalf("creating Logstash pipeline: %#v", diagnostics)
	}
	if got := resourceData.Id(); got != "logstash-pipeline-87654321" {
		t.Fatalf("unexpected pipeline ID: %q", got)
	}
	if strings.Contains(logger.String(), "pipeline-secret") {
		t.Fatal("Logstash configuration was written to AWS SDK debug logs")
	}
}

func TestFindLogstashPipelineByIDDoesNotLogConfiguration(t *testing.T) {
	t.Parallel()

	const configuration = "input { http { password => pipeline-secret } }"

	logger := &recordingSDKLogger{}
	conn := testPaaSClientWithDebugBodyLogging(logger, roundTripFunc(func(request *http.Request) (*http.Response, error) {
		return jsonResponse(request, `{
			"logstashPipelines": [
				{
					"id": "logstash-pipeline-87654321",
					"name": "terraform-qa",
					"configuration": "`+configuration+`"
				}
			]
		}`), nil
	}))

	pipeline, err := FindLogstashPipelineByID(
		context.Background(),
		conn,
		"fm-cluster-12345678",
		"logstash-pipeline-87654321",
	)
	if err != nil {
		t.Fatalf("reading Logstash pipeline: %s", err)
	}
	if got := aws.StringValue(pipeline.Name); got != "terraform-qa" {
		t.Fatalf("unexpected pipeline name: %q", got)
	}
	if strings.Contains(logger.String(), configuration) {
		t.Fatal("Logstash configuration was written to AWS SDK debug logs")
	}
}

func TestResourceLogstashPipelineUpdateDoesNotLogConfiguration(t *testing.T) {
	t.Parallel()

	const configuration = "input { http { password => pipeline-update-secret } }"

	logger := &recordingSDKLogger{}
	modifyCalled := false
	conn := testPaaSClientWithDebugBodyLogging(logger, roundTripFunc(func(request *http.Request) (*http.Response, error) {
		switch {
		case request.Method == http.MethodPut &&
			request.URL.Path == "/services/fm-cluster-12345678/logstash-pipelines/logstash-pipeline-87654321":
			modifyCalled = true
			return jsonResponse(request, `{}`), nil
		case request.Method == http.MethodGet &&
			request.URL.Path == "/services/fm-cluster-12345678":
			return jsonResponse(request, `{
				"service": {
					"id": "fm-cluster-12345678",
					"serviceType": "elk",
					"serviceClass": "logging",
					"status": "READY",
					"parameters": {
						"version": "8.17"
					}
				}
			}`), nil
		case request.Method == http.MethodGet &&
			request.URL.Path == "/services/fm-cluster-12345678/logstash-pipelines":
			return jsonResponse(request, `{
				"logstashPipelines": [
					{
						"id": "logstash-pipeline-87654321",
						"name": "terraform-qa",
						"configuration": "`+configuration+`"
					}
				]
			}`), nil
		default:
			t.Fatalf("unexpected PaaS request: %s %s", request.Method, request.URL.Path)
			return nil, nil
		}
	}))

	resourceData := schema.TestResourceDataRaw(t, ResourceLogstashPipeline().Schema, map[string]interface{}{
		"service_id":    "fm-cluster-12345678",
		"name":          "terraform-qa",
		"configuration": configuration,
	})
	resourceData.SetId("logstash-pipeline-87654321")

	diagnostics := resourceLogstashPipelineUpdate(
		context.Background(),
		resourceData,
		testPaaSClientMeta(conn),
	)
	if diagnostics.HasError() {
		t.Fatalf("updating Logstash pipeline: %#v", diagnostics)
	}
	if !modifyCalled {
		t.Fatal("Logstash pipeline update request was not sent")
	}
	if strings.Contains(logger.String(), "pipeline-update-secret") {
		t.Fatal("updated Logstash configuration was written to AWS SDK debug logs")
	}
}

func testPaaSClientWithDebugBodyLogging(logger aws.Logger, transport http.RoundTripper) *sdkpaas.PaaS {
	return sdkpaas.New(session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("access-key", "secret-key", ""),
		Endpoint:    aws.String("http://paas.test"),
		HTTPClient:  &http.Client{Transport: transport},
		Logger:      logger,
		LogLevel:    aws.LogLevel(aws.LogDebugWithHTTPBody),
		Region:      aws.String("ru-msk"),
	})))
}

func jsonResponse(request *http.Request, body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body:    io.NopCloser(strings.NewReader(body)),
		Request: request,
	}
}

type recordingSDKLogger struct {
	mu      sync.Mutex
	builder strings.Builder
}

func (l *recordingSDKLogger) Log(values ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	fmt.Fprint(&l.builder, values...)
}

func (l *recordingSDKLogger) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.builder.String()
}
