# Elastic Distribution for OpenTelemetry Collector

This is an Elastic supported distribution of the [OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector).

## Running the Elastic Distribution for OpenTelemetry Collector

To run the Elastic Distribution for OpenTelemetry Collector you can use Elastic-Agent binary downloaded for your OS and architecture.
Running command

```bash
./elastic-agent otel --config otel.yml
```

from unpacked Elastic Agent package will run Elastic-Agent as an OpenTelemetry Collector. The `--config` flag needs to point to [OpenTelemetry Collector Configuration file](https://opentelemetry.io/docs/collector/configuration/). OTel mode is available only using `otel` subcommand. Elastic Agent will not do any autodetection of configuration file passed when used without `otel` subcommand and will try to run normally.

To validate OTel configuration run `otel validate` subcommand:

```bash
./elastic-agent otel validate --config otel.yml
```

Use the components command to get the list of components included in the binary:

```bash
./elastic-agent otel components
```

[feature gates](https://github.com/open-telemetry/opentelemetry-collector/blob/main/featuregate/README.md#controlling-gates) are supported using `--feature-gates` flag.

## Components

This section provides a summary of components included in the Elastic Distribution for OpenTelemetry Collector.

### Receivers

| Component | Version |
|---|---|
| [filelogreceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/receiver/filelogreceiver/v0.120.1/receiver/filelogreceiver/README.md) | v0.120.1 |
| [hostmetricsreceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/receiver/hostmetricsreceiver/v0.120.1/receiver/hostmetricsreceiver/README.md) | v0.120.1 |
| [httpcheckreceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/receiver/httpcheckreceiver/v0.120.1/receiver/httpcheckreceiver/README.md) | v0.120.1 |
| [jaegerreceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/receiver/jaegerreceiver/v0.120.1/receiver/jaegerreceiver/README.md) | v0.120.1 |
| [jmxreceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/receiver/jmxreceiver/v0.120.1/receiver/jmxreceiver/README.md) | v0.120.1 |
| [k8sclusterreceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/receiver/k8sclusterreceiver/v0.120.1/receiver/k8sclusterreceiver/README.md) | v0.120.1 |
| [k8sobjectsreceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/receiver/k8sobjectsreceiver/v0.120.1/receiver/k8sobjectsreceiver/README.md) | v0.120.1 |
| [kafkareceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/receiver/kafkareceiver/v0.120.1/receiver/kafkareceiver/README.md) | v0.120.1 |
| [kubeletstatsreceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/receiver/kubeletstatsreceiver/v0.120.1/receiver/kubeletstatsreceiver/README.md) | v0.120.1 |
| [nginxreceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/receiver/nginxreceiver/v0.120.1/receiver/nginxreceiver/README.md) | v0.120.1 |
| [nopreceiver](https://github.com/open-telemetry/opentelemetry-collector/blob/receiver/nopreceiver/v0.119.0/receiver/nopreceiver/README.md) | v0.119.0 |
| [otlpreceiver](https://github.com/open-telemetry/opentelemetry-collector/blob/receiver/otlpreceiver/v0.120.0/receiver/otlpreceiver/README.md) | v0.120.0 |
| [prometheusreceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/receiver/prometheusreceiver/v0.120.1/receiver/prometheusreceiver/README.md) | v0.120.1 |
| [prometheusreceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/receiver/prometheusreceiver/=&gt; github.com/elastic/opentelemetry-collector-contrib/receiver/prometheusreceiver v0.0.0-20250317163643-19cd4e80024f/receiver/prometheusreceiver/README.md) | =&gt; github.com/elastic/opentelemetry-collector-contrib/receiver/prometheusreceiver v0.0.0-20250317163643-19cd4e80024f |
| [receivercreator](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/receiver/receivercreator/v0.120.1/receiver/receivercreator/README.md) | v0.120.1 |
| [redisreceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/receiver/redisreceiver/v0.120.1/receiver/redisreceiver/README.md) | v0.120.1 |
| [zipkinreceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/receiver/zipkinreceiver/v0.120.1/receiver/zipkinreceiver/README.md) | v0.120.1 |

### Exporters

| Component | Version |
|---|---|
| [debugexporter](https://github.com/open-telemetry/opentelemetry-collector/blob/exporter/debugexporter/v0.120.0/exporter/debugexporter/README.md) | v0.120.0 |
| [elasticsearchexporter](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/exporter/elasticsearchexporter/v0.120.1/exporter/elasticsearchexporter/README.md) | v0.120.1 |
| [fileexporter](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/exporter/fileexporter/v0.120.1/exporter/fileexporter/README.md) | v0.120.1 |
| [kafkaexporter](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/exporter/kafkaexporter/v0.120.1/exporter/kafkaexporter/README.md) | v0.120.1 |
| [loadbalancingexporter](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/exporter/loadbalancingexporter/v0.120.1/exporter/loadbalancingexporter/README.md) | v0.120.1 |
| [otlpexporter](https://github.com/open-telemetry/opentelemetry-collector/blob/exporter/otlpexporter/v0.120.0/exporter/otlpexporter/README.md) | v0.120.0 |
| [otlphttpexporter](https://github.com/open-telemetry/opentelemetry-collector/blob/exporter/otlphttpexporter/v0.120.0/exporter/otlphttpexporter/README.md) | v0.120.0 |

### Processors

| Component | Version |
|---|---|
| [attributesprocessor](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/processor/attributesprocessor/v0.120.1/processor/attributesprocessor/README.md) | v0.120.1 |
| [batchprocessor](https://github.com/open-telemetry/opentelemetry-collector/blob/processor/batchprocessor/v0.120.0/processor/batchprocessor/README.md) | v0.120.0 |
| [elasticinframetricsprocessor](https://github.com/elastic/opentelemetry-collector-components/blob/processor/elasticinframetricsprocessor/v0.13.0/processor/elasticinframetricsprocessor/README.md) | v0.13.0 |
| [elastictraceprocessor](https://github.com/elastic/opentelemetry-collector-components/blob/processor/elastictraceprocessor/v0.4.1/processor/elastictraceprocessor/README.md) | v0.4.1 |
| [filterprocessor](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/processor/filterprocessor/v0.120.1/processor/filterprocessor/README.md) | v0.120.1 |
| [geoipprocessor](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/processor/geoipprocessor/v0.120.1/processor/geoipprocessor/README.md) | v0.120.1 |
| [k8sattributesprocessor](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/processor/k8sattributesprocessor/v0.120.1/processor/k8sattributesprocessor/README.md) | v0.120.1 |
| [memorylimiterprocessor](https://github.com/open-telemetry/opentelemetry-collector/blob/processor/memorylimiterprocessor/v0.119.0/processor/memorylimiterprocessor/README.md) | v0.119.0 |
| [resourcedetectionprocessor](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/processor/resourcedetectionprocessor/v0.120.1/processor/resourcedetectionprocessor/README.md) | v0.120.1 |
| [resourceprocessor](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/processor/resourceprocessor/v0.120.1/processor/resourceprocessor/README.md) | v0.120.1 |
| [transformprocessor](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/processor/transformprocessor/v0.120.1/processor/transformprocessor/README.md) | v0.120.1 |

### Extensions

| Component | Version |
|---|---|
| [filestorage](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/extension/storage/filestorage/v0.120.1/extension/storage/filestorage/README.md) | v0.120.1 |
| [healthcheckextension](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/extension/healthcheckextension/v0.120.1/extension/healthcheckextension/README.md) | v0.120.1 |
| [k8sobserver](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/extension/observer/k8sobserver/v0.120.1/extension/observer/k8sobserver/README.md) | v0.120.1 |
| [memorylimiterextension](https://github.com/open-telemetry/opentelemetry-collector/blob/extension/memorylimiterextension/v0.120.0/extension/memorylimiterextension/README.md) | v0.120.0 |
| [pprofextension](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/extension/pprofextension/v0.120.1/extension/pprofextension/README.md) | v0.120.1 |

### Connectors

| Component | Version |
|---|---|
| [elasticapmconnector](https://github.com/elastic/opentelemetry-collector-components/blob/connector/elasticapmconnector/v0.2.0/connector/elasticapmconnector/README.md) | v0.2.0 |
| [routingconnector](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/connector/routingconnector/v0.120.1/connector/routingconnector/README.md) | v0.120.1 |
| [spanmetricsconnector](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/connector/spanmetricsconnector/v0.120.1/connector/spanmetricsconnector/README.md) | v0.120.1 |
## Persistence in OpenTelemetry Collector

By default, the OpenTelemetry Collector is stateless, which means it doesn't store offsets on disk while reading files. As a result, if you restart the collector, it won't retain the last read offset, potentially leading to data duplication or loss. However, we have configured persistence in the settings provided with the Elastic Agent package.

To enable persistence for the `filelogreceiver`, we add the `file_storage` extension and activate it for `filelog`.
Execute `export STATE_PATH=/path/to/store/otel/offsets` and use the following configuration to enable persistence:

```yaml
receivers:
  filelog/platformlogs:
    include: [ /var/log/system.log ]
    start_at: beginning
    storage: file_storage/filelogreceiver
extensions:
  file_storage/filelogreceiver:
    directory: ${env:STATE_PATH}
    create_directory: true
exporters:
  ...
processors:
  ...
service:
  extensions: [file_storage]
  pipelines:
    logs/platformlogs:
      receivers: [filelog/platformlogs]
      processors: [...]
      exporters: [...]
```

> [!WARNING]
Removing the storage key from the filelog section will disable persistence, which will lead to data duplication or loss when the collector restarts.

> [!IMPORTANT]
If you remove the `create_directory: true` option, you'll need to manually create a directory to store the data. You can ignore this option if the directory already exists.

### Persistence in standalone Docker mode

By default, when running Elastic Distribution for OpenTelemetry Collector in Docker, checkpoints are stored in `/usr/share/elastic-agent/otel_registry` by default. To ensure data persists across container restarts, you can use the following command:

```bash
docker run --rm -ti --entrypoint="elastic-agent" --mount type=bind,source=/path/on/host,target=/usr/share/elastic-agent/otel_registry  docker.elastic.co/elastic-agent/elastic-agent:9.0.0-SNAPSHOT otel
```

### Known issues:
-  You face following `failed to build extensions: failed to create extension "file_storage/filelogreceiver": mkdir ...: permission denied` error while running the otel mode
	- Cause: This issue is likely because the user running the executable lacks sufficient permissions to create the directory.
	- Resolution: You can either create the directory manually or specify a path with necessary permissions.
