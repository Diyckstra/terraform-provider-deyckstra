package paas

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	sdkpaas "github.com/aws/aws-sdk-go/service/paas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestResourceLogstashPipelineSchema(t *testing.T) {
	t.Parallel()

	resource := ResourceLogstashPipeline()
	if resource.CreateContext == nil || resource.ReadContext == nil ||
		resource.UpdateContext == nil || resource.DeleteContext == nil {
		t.Fatal("Logstash pipeline resource must implement full CRUD")
	}
	if resource.Importer == nil || resource.Importer.StateContext == nil {
		t.Fatal("Logstash pipeline resource must support import")
	}

	serviceID := resource.Schema["service_id"]
	if !serviceID.Required || !serviceID.ForceNew {
		t.Fatal("service_id must be required and ForceNew")
	}

	pipelineID := resource.Schema["pipeline_id"]
	if !pipelineID.Computed {
		t.Fatal("pipeline_id must be computed")
	}

	name := resource.Schema["name"]
	if !name.Required || !name.ForceNew {
		t.Fatal("name must be required and ForceNew")
	}

	configuration := resource.Schema["configuration"]
	if !configuration.Required || configuration.ForceNew || !configuration.Sensitive {
		t.Fatal("configuration must be required, sensitive, and editable in place")
	}
}

func TestResourceLogstashPipelineValidation(t *testing.T) {
	t.Parallel()

	resource := ResourceLogstashPipeline()

	for _, testCase := range []struct {
		field string
		value string
		valid bool
	}{
		{field: "service_id", value: "fm-cluster-12345678", valid: true},
		{field: "service_id", value: "", valid: false},
		{field: "name", value: "terraform-qa", valid: true},
		{field: "name", value: "", valid: false},
		{field: "name", value: "beats-to-elasticsearch", valid: false},
		{field: "configuration", value: "input { stdin {} }", valid: true},
		{field: "configuration", value: "input {\n  stdin {}\n}", valid: true},
		{field: "configuration", value: "", valid: false},
	} {
		testCase := testCase
		t.Run(testCase.field+"/"+testCase.value, func(t *testing.T) {
			_, errors := resource.Schema[testCase.field].ValidateFunc(testCase.value, testCase.field)
			if testCase.valid && len(errors) != 0 {
				t.Fatalf("expected %q to be accepted for %s, got %v", testCase.value, testCase.field, errors)
			}
			if !testCase.valid && len(errors) == 0 {
				t.Fatalf("expected %q to be rejected for %s", testCase.value, testCase.field)
			}
		})
	}
}

func TestParseLogstashPipelineImportID(t *testing.T) {
	t.Parallel()

	serviceID, pipelineID, err := parseLogstashPipelineImportID(
		"fm-cluster-12345678/logstash-pipeline-87654321",
	)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if serviceID != "fm-cluster-12345678" || pipelineID != "logstash-pipeline-87654321" {
		t.Fatalf("unexpected parsed IDs: %q %q", serviceID, pipelineID)
	}

	for _, id := range []string{
		"",
		"fm-cluster-12345678",
		"/logstash-pipeline-87654321",
		"fm-cluster-12345678/",
		"fm-cluster-12345678/logstash-pipeline-87654321/extra",
	} {
		_, _, err := parseLogstashPipelineImportID(id)
		if err == nil {
			t.Errorf("expected import ID %q to fail", id)
			continue
		}
		if !strings.Contains(err.Error(), "service_id/pipeline_id") {
			t.Errorf("unexpected error for %q: %s", id, err)
		}
	}
}

func TestResourceLogstashPipelineImport(t *testing.T) {
	t.Parallel()

	resource := ResourceLogstashPipeline()
	resourceData := schema.TestResourceDataRaw(t, resource.Schema, nil)
	resourceData.SetId("fm-cluster-12345678/logstash-pipeline-87654321")

	imported, err := resourceLogstashPipelineImport(context.Background(), resourceData, nil)
	if err != nil {
		t.Fatalf("unexpected import error: %s", err)
	}
	if len(imported) != 1 {
		t.Fatalf("unexpected imported resource count: %d", len(imported))
	}
	if got := imported[0].Id(); got != "logstash-pipeline-87654321" {
		t.Fatalf("unexpected state ID: %q", got)
	}
	if got := imported[0].Get("service_id").(string); got != "fm-cluster-12345678" {
		t.Fatalf("unexpected service_id: %q", got)
	}
	if got := imported[0].Get("pipeline_id").(string); got != "logstash-pipeline-87654321" {
		t.Fatalf("unexpected pipeline_id: %q", got)
	}
}

func TestFindLogstashPipelineInListByID(t *testing.T) {
	t.Parallel()

	pipelines := []*sdkpaas.LogstashPipeline{
		nil,
		{
			Id:            aws.String("logstash-pipeline-first"),
			Name:          aws.String("first"),
			Configuration: aws.String("input { stdin {} }"),
		},
		{
			Id:            aws.String("logstash-pipeline-target"),
			Name:          aws.String("target"),
			Configuration: aws.String("input { http { port => 4567 } }"),
		},
	}

	got, err := findLogstashPipelineInListByID(pipelines, "logstash-pipeline-target")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if aws.StringValue(got.Name) != "target" {
		t.Fatalf("unexpected pipeline: %#v", got)
	}
}

func TestFindLogstashPipelineInListByIDNotFound(t *testing.T) {
	t.Parallel()

	for _, pipelines := range [][]*sdkpaas.LogstashPipeline{
		nil,
		{},
		{nil, {Id: aws.String("logstash-pipeline-other")}},
	} {
		_, err := findLogstashPipelineInListByID(pipelines, "logstash-pipeline-target")
		if !tfresource.NotFound(err) {
			t.Fatalf("expected Terraform not-found error, got %v", err)
		}
	}
}

func TestResourceLogstashPipelineReadRejectsReservedPipeline(t *testing.T) {
	t.Parallel()

	httpClient := &http.Client{
		Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(`{
					"logstashPipelines": [
						{
							"id": "logstash-pipeline-system",
							"name": "beats-to-elasticsearch",
							"configuration": "input { beats { port => 5044 } }"
						}
					]
				}`)),
				Request: request,
			}, nil
		}),
	}

	conn := sdkpaas.New(session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("access-key", "secret-key", ""),
		Endpoint:    aws.String("http://paas.test"),
		HTTPClient:  httpClient,
		Region:      aws.String("ru-msk"),
	})))
	resourceData := schema.TestResourceDataRaw(t, ResourceLogstashPipeline().Schema, map[string]interface{}{
		"service_id": "fm-cluster-12345678",
	})
	resourceData.SetId("logstash-pipeline-system")

	diagnostics := resourceLogstashPipelineRead(
		context.Background(),
		resourceData,
		testPaaSClientMeta(conn),
	)
	if !diagnostics.HasError() {
		t.Fatal("expected reserved pipeline read to fail")
	}
	if !strings.Contains(diagnostics[0].Summary, "reserved name") {
		t.Fatalf("unexpected diagnostic: %s", diagnostics[0].Summary)
	}
}

func TestResourceLogstashPipelineDeleteRejectsReservedPipeline(t *testing.T) {
	t.Parallel()

	resourceData := schema.TestResourceDataRaw(t, ResourceLogstashPipeline().Schema, map[string]interface{}{
		"service_id":    "fm-cluster-12345678",
		"name":          reservedLogstashPipelineName,
		"configuration": "input { beats { port => 5044 } }",
	})
	resourceData.SetId("logstash-pipeline-system")

	diagnostics := resourceLogstashPipelineDelete(
		context.Background(),
		resourceData,
		&conns.AWSClient{},
	)
	if !diagnostics.HasError() {
		t.Fatal("expected reserved pipeline delete to fail")
	}
	if !strings.Contains(diagnostics[0].Summary, "cannot be deleted") {
		t.Fatalf("unexpected diagnostic: %s", diagnostics[0].Summary)
	}
}

func TestResourceLogstashPipelineDeleteIgnoresAPI404(t *testing.T) {
	t.Parallel()

	httpClient := &http.Client{
		Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
			if request.Method != http.MethodDelete {
				t.Errorf("unexpected method: got %s want DELETE", request.Method)
			}
			if request.URL.Path != "/services/fm-cluster-12345678/logstash-pipelines/logstash-pipeline-missing" {
				t.Errorf("unexpected path: %s", request.URL.Path)
			}

			return &http.Response{
				StatusCode: http.StatusNotFound,
				Status:     "404 Not Found",
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(`{
					"code": "PipelineNotFound",
					"message": "pipeline not found"
				}`)),
				Request: request,
			}, nil
		}),
	}

	conn := sdkpaas.New(session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("access-key", "secret-key", ""),
		Endpoint:    aws.String("http://paas.test"),
		HTTPClient:  httpClient,
		Region:      aws.String("ru-msk"),
	})))
	resourceData := schema.TestResourceDataRaw(t, ResourceLogstashPipeline().Schema, map[string]interface{}{
		"service_id":    "fm-cluster-12345678",
		"name":          "missing",
		"configuration": "input { stdin {} }",
	})
	resourceData.SetId("logstash-pipeline-missing")

	diagnostics := resourceLogstashPipelineDelete(
		context.Background(),
		resourceData,
		testPaaSClientMeta(conn),
	)
	if diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %#v", diagnostics)
	}
}

func TestFindLogstashPipelineByID(t *testing.T) {
	t.Parallel()

	httpClient := &http.Client{
		Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
			if request.Method != http.MethodGet {
				t.Errorf("unexpected method: got %s want GET", request.Method)
			}
			if request.URL.Path != "/services/fm-cluster-12345678/logstash-pipelines" {
				t.Errorf("unexpected path: %s", request.URL.Path)
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(`{
					"logstashPipelines": [
						{
							"id": "logstash-pipeline-other",
							"name": "other",
							"configuration": "input { stdin {} }"
						},
						{
							"id": "logstash-pipeline-target",
							"name": "target",
							"configuration": "input { http { port => 4567 } }"
						}
					]
				}`)),
				Request: request,
			}, nil
		}),
	}

	conn := sdkpaas.New(session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("access-key", "secret-key", ""),
		Endpoint:    aws.String("http://paas.test"),
		HTTPClient:  httpClient,
		Region:      aws.String("ru-msk"),
	})))

	got, err := FindLogstashPipelineByID(
		context.Background(),
		conn,
		"fm-cluster-12345678",
		"logstash-pipeline-target",
	)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if aws.StringValue(got.Name) != "target" {
		t.Fatalf("unexpected pipeline: %#v", got)
	}
}

func TestFindLogstashPipelineByIDAPI404(t *testing.T) {
	t.Parallel()

	httpClient := &http.Client{
		Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Status:     "404 Not Found",
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(`{
					"code": "PipelineNotFound",
					"message": "pipeline not found"
				}`)),
				Request: request,
			}, nil
		}),
	}

	conn := sdkpaas.New(session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("access-key", "secret-key", ""),
		Endpoint:    aws.String("http://paas.test"),
		HTTPClient:  httpClient,
		Region:      aws.String("ru-msk"),
	})))

	_, err := FindLogstashPipelineByID(
		context.Background(),
		conn,
		"fm-cluster-12345678",
		"logstash-pipeline-target",
	)
	if !tfresource.NotFound(err) {
		t.Fatalf("expected Terraform not-found error, got %v", err)
	}
}

func testPaaSClientMeta(conn *sdkpaas.PaaS) interface{} {
	return &conns.AWSClient{PaaSConn: conn}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestIsLogstashPipelineNotFoundError(t *testing.T) {
	t.Parallel()

	requestFailure := awserr.NewRequestFailure(
		awserr.New("RequestFailed", "not found", nil),
		http.StatusNotFound,
		"request-id",
	)

	for _, err := range []error{
		requestFailure,
		awserr.New("PipelineNotFoundException", "not found", nil),
	} {
		if !isLogstashPipelineNotFoundError(err) {
			t.Errorf("expected %T to be treated as not-found: %v", err, err)
		}
	}

	for _, err := range []error{
		nil,
		awserr.New("ValidationException", "invalid", nil),
	} {
		if isLogstashPipelineNotFoundError(err) {
			t.Errorf("did not expect %T to be treated as not-found: %v", err, err)
		}
	}
}
