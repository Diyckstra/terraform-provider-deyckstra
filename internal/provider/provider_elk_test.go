package provider

import "testing"

func TestProviderRegistersLogstashPipeline(t *testing.T) {
	t.Parallel()

	provider := Provider()
	resource, ok := provider.ResourcesMap["aws_paas_logstash_pipeline"]
	if !ok || resource == nil {
		t.Fatal("expected aws_paas_logstash_pipeline to be registered")
	}

	if err := provider.InternalValidate(); err != nil {
		t.Fatalf("provider schema validation failed: %s", err)
	}
}
