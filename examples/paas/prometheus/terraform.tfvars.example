# `ssh_public_key` and a unique `paas_service_name` are intentionally supplied
# through TF_VAR_* environment variables in README.md.

region            = "ru-msk"
availability_zone = "ru-msk-vol52"
instance_type     = "c5.large"
# The reserved .test domain cannot receive alert payloads. Replace this
# value with an HTTPS endpoint controlled by you to test actual delivery.
webhook_url = "https://prometheus-webhook.test/alerts"
