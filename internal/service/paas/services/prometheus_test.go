package services

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestPrometheusManagerSchema(t *testing.T) {
	t.Parallel()

	resourceSchema := Prometheus.ResourceSchema()
	if resourceSchema.ForceNew {
		t.Fatal("the Prometheus service block must allow editable service parameters")
	}

	serviceResource, ok := resourceSchema.Elem.(*schema.Resource)
	if !ok {
		t.Fatalf("unexpected Prometheus resource schema element type %T", resourceSchema.Elem)
	}

	remoteWriteReceiver, ok := serviceResource.Schema["remote_write_receiver"]
	if !ok {
		t.Fatal("remote_write_receiver is missing from the Prometheus resource schema")
	}
	if !remoteWriteReceiver.Optional || remoteWriteReceiver.Required || remoteWriteReceiver.Computed {
		t.Fatalf("unexpected resource schema flags for remote_write_receiver: %#v", remoteWriteReceiver)
	}
	if remoteWriteReceiver.ForceNew {
		t.Fatal("remote_write_receiver must remain updateable")
	}
	if got, want := remoteWriteReceiver.Default, false; got != want {
		t.Fatalf("unexpected remote_write_receiver default: got %#v, want %#v", got, want)
	}

	dataSourceSchema := Prometheus.DataSourceSchema()
	dataSourceResource, ok := dataSourceSchema.Elem.(*schema.Resource)
	if !ok {
		t.Fatalf("unexpected Prometheus data source schema element type %T", dataSourceSchema.Elem)
	}

	remoteWriteReceiver, ok = dataSourceResource.Schema["remote_write_receiver"]
	if !ok {
		t.Fatal("remote_write_receiver is missing from the Prometheus data source schema")
	}
	if !remoteWriteReceiver.Computed || remoteWriteReceiver.Optional || remoteWriteReceiver.Required {
		t.Fatalf("unexpected data source schema flags for remote_write_receiver: %#v", remoteWriteReceiver)
	}
}

func TestPrometheusManagerExpandServiceParameters(t *testing.T) {
	t.Parallel()

	got := Prometheus.ExpandServiceParameters(map[string]interface{}{
		"remote_write_receiver": true,
	})
	want := ServiceParameters{
		"remote_write_receiver": true,
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected expanded parameters (-want +got):\n%s", diff)
	}
	if _, ok := got["logging"]; ok {
		t.Fatal("Prometheus parameters must not contain the unsupported logging flag")
	}
	if _, ok := got["monitoring"]; ok {
		t.Fatal("Prometheus parameters must not contain the unsupported monitoring flag")
	}
}

func TestPrometheusManagerFlattenServiceParameters(t *testing.T) {
	t.Parallel()

	got := Prometheus.FlattenServiceParametersUsersDatabases(
		ServiceParameters{"remoteWriteReceiver": true},
		nil,
		nil,
	)
	want := map[string]interface{}{
		"remote_write_receiver": true,
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected flattened parameters (-want +got):\n%s", diff)
	}
}
