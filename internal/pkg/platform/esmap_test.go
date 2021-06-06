package platform

import (
	"go.uber.org/zap"
	"testing"
)

var esClustermock = `
{
    "_nodes": {
        "total": 1,
        "successful": 1,
        "failed": 0
    },
    "cluster_name": "elasticsearch",
    "nodes": {
        "b-gmel9SReK8Ku4C8gUkzA": {
            "name": "b-gmel9",
            "transport_address": "127.0.0.1:9300",
            "host": "127.0.0.1",
            "ip": "127.0.0.1",
            "version": "6.8.11",
            "build_flavor": "default",
            "build_type": "deb",
            "build_hash": "00bf386",
            "roles": [
                "master",
                "data",
                "ingest"
            ],
            "attributes": {
                "ml.machine_memory": "33512615936",
                "xpack.installed": "true",
                "ml.max_open_jobs": "20",
                "ml.enabled": "true"
            },
            "http": {
                "bound_address": [
                    "{REPLACE_URL}"
                ],
                "publish_address": "{REPLACE_URL}",
                "max_content_length_in_bytes": 104857600
            }
        }
    }
}`

func TestNewES(t *testing.T) {
	log, _ := zap.NewProduction()
	server := mockServer(200, esClustermock)
	defer server.Close()
	e := NewES(server.URL, "elastic", "", log)
	if e == nil {
		t.Errorf("Unable to initialize Elasticsearch")
	}
}
