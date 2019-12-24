
# go-core 

|                               Key Path                                |                                    ENV                                     | Default |         Type          |
|-----------------------------------------------------------------------|----------------------------------------------------------------------------|---------|-----------------------|
| logger.Debug                                                          | LOGGER_DEBUG                                                               |         | bool                  |
| logger.Verbose                                                        | LOGGER_VERBOSE                                                             |         | bool                  |
| logger.Level                                                          | LOGGER_LEVEL                                                               |         | logger.Level          |
| logger.DebugTags                                                      | LOGGER_DEBUG_TAGS                                                          |         | []string              |
| logger.MapTagsSplitSep                                                | LOGGER_MAP_TAGS_SPLIT_SEP                                                  | :       | string                |
| logger.DisableRedirectStdLog                                          | LOGGER_DISABLE_REDIRECT_STD_LOG                                            |         | bool                  |
| logger.RedirectLevel                                                  | LOGGER_REDIRECT_LEVEL                                                      | 6       | logger.Level          |
| metric.Enabled                                                        | METRIC_ENABLED                                                             |         | bool                  |
| metric.StatsD.Addr                                                    | METRIC_STATS_D_ADDR                                                        |         | string                |
| metric.StatsD.Prefix                                                  | METRIC_STATS_D_PREFIX                                                      |         | string                |
| metric.StatsD.FlushInterval                                           | METRIC_STATS_D_FLUSH_INTERVAL                                              |         | time.Duration         |
| metric.StatsD.FlushBytes                                              | METRIC_STATS_D_FLUSH_BYTES                                                 |         | int                   |
| metric.StatsD.Options.SampleRate                                      | METRIC_STATS_D_OPTIONS_SAMPLE_RATE                                         |         | float32               |
| metric.StatsD.Options.HistogramBucketNamePrecision                    | METRIC_STATS_D_OPTIONS_HISTOGRAM_BUCKET_NAME_PRECISION                     |         | uint                  |
| metric.Scope.Tags                                                     | METRIC_SCOPE_TAGS                                                          |         | map[string]string     |
| metric.Scope.Prefix                                                   | METRIC_SCOPE_PREFIX                                                        |         | string                |
| metric.Scope.Separator                                                | METRIC_SCOPE_SEPARATOR                                                     |         | string                |
| metric.Scope.SanitizeOptions.NameCharacters.Ranges                    | METRIC_SCOPE_SANITIZE_OPTIONS_NAME_CHARACTERS_RANGES                       |         | []tally.SanitizeRange |
| metric.Scope.SanitizeOptions.NameCharacters.Characters                | METRIC_SCOPE_SANITIZE_OPTIONS_NAME_CHARACTERS_CHARACTERS                   |         | []int32               |
| metric.Scope.SanitizeOptions.KeyCharacters.Ranges                     | METRIC_SCOPE_SANITIZE_OPTIONS_KEY_CHARACTERS_RANGES                        |         | []tally.SanitizeRange |
| metric.Scope.SanitizeOptions.KeyCharacters.Characters                 | METRIC_SCOPE_SANITIZE_OPTIONS_KEY_CHARACTERS_CHARACTERS                    |         | []int32               |
| metric.Scope.SanitizeOptions.ValueCharacters.Ranges                   | METRIC_SCOPE_SANITIZE_OPTIONS_VALUE_CHARACTERS_RANGES                      |         | []tally.SanitizeRange |
| metric.Scope.SanitizeOptions.ValueCharacters.Characters               | METRIC_SCOPE_SANITIZE_OPTIONS_VALUE_CHARACTERS_CHARACTERS                  |         | []int32               |
| metric.Scope.SanitizeOptions.ReplacementCharacter                     | METRIC_SCOPE_SANITIZE_OPTIONS_REPLACEMENT_CHARACTER                        |         | int32                 |
| metric.Interval                                                       | METRIC_INTERVAL                                                            |         | time.Duration         |
| tracing.Enabled                                                       | TRACING_ENABLED                                                            |         | bool                  |
| tracing.Jaeger.ServiceName                                            | TRACING_JAEGER_SERVICE_NAME                                                |         | string                |
| tracing.Jaeger.Disabled                                               | TRACING_JAEGER_DISABLED                                                    |         | bool                  |
| tracing.Jaeger.RPCMetrics                                             | TRACING_JAEGER_RPC_METRICS                                                 |         | bool                  |
| tracing.Jaeger.Tags                                                   | TRACING_JAEGER_TAGS                                                        |         | []opentracing.Tag     |
| tracing.Jaeger.Sampler.Type                                           | TRACING_JAEGER_SAMPLER_TYPE                                                |         | string                |
| tracing.Jaeger.Sampler.Param                                          | TRACING_JAEGER_SAMPLER_PARAM                                               |         | float64               |
| tracing.Jaeger.Sampler.SamplingServerURL                              | TRACING_JAEGER_SAMPLER_SAMPLING_SERVER_URL                                 |         | string                |
| tracing.Jaeger.Sampler.MaxOperations                                  | TRACING_JAEGER_SAMPLER_MAX_OPERATIONS                                      |         | int                   |
| tracing.Jaeger.Sampler.SamplingRefreshInterval                        | TRACING_JAEGER_SAMPLER_SAMPLING_REFRESH_INTERVAL                           |         | time.Duration         |
| tracing.Jaeger.Reporter.QueueSize                                     | TRACING_JAEGER_REPORTER_QUEUE_SIZE                                         |         | int                   |
| tracing.Jaeger.Reporter.BufferFlushInterval                           | TRACING_JAEGER_REPORTER_BUFFER_FLUSH_INTERVAL                              |         | time.Duration         |
| tracing.Jaeger.Reporter.LogSpans                                      | TRACING_JAEGER_REPORTER_LOG_SPANS                                          |         | bool                  |
| tracing.Jaeger.Reporter.LocalAgentHostPort                            | TRACING_JAEGER_REPORTER_LOCAL_AGENT_HOST_PORT                              |         | string                |
| tracing.Jaeger.Reporter.CollectorEndpoint                             | TRACING_JAEGER_REPORTER_COLLECTOR_ENDPOINT                                 |         | string                |
| tracing.Jaeger.Reporter.User                                          | TRACING_JAEGER_REPORTER_USER                                               |         | string                |
| tracing.Jaeger.Reporter.Password                                      | TRACING_JAEGER_REPORTER_PASSWORD                                           |         | string                |
| tracing.Jaeger.Headers.JaegerDebugHeader                              | TRACING_JAEGER_HEADERS_JAEGER_DEBUG_HEADER                                 |         | string                |
| tracing.Jaeger.Headers.JaegerBaggageHeader                            | TRACING_JAEGER_HEADERS_JAEGER_BAGGAGE_HEADER                               |         | string                |
| tracing.Jaeger.Headers.TraceContextHeaderName                         | TRACING_JAEGER_HEADERS_TRACE_CONTEXT_HEADER_NAME                           |         | string                |
| tracing.Jaeger.Headers.TraceBaggageHeaderPrefix                       | TRACING_JAEGER_HEADERS_TRACE_BAGGAGE_HEADER_PREFIX                         |         | string                |
| tracing.Jaeger.BaggageRestrictions.DenyBaggageOnInitializationFailure | TRACING_JAEGER_BAGGAGE_RESTRICTIONS_DENY_BAGGAGE_ON_INITIALIZATION_FAILURE |         | bool                  |
| tracing.Jaeger.BaggageRestrictions.HostPort                           | TRACING_JAEGER_BAGGAGE_RESTRICTIONS_HOST_PORT                              |         | string                |
| tracing.Jaeger.BaggageRestrictions.RefreshInterval                    | TRACING_JAEGER_BAGGAGE_RESTRICTIONS_REFRESH_INTERVAL                       |         | time.Duration         |
| tracing.Jaeger.Throttler.HostPort                                     | TRACING_JAEGER_THROTTLER_HOST_PORT                                         |         | string                |
| tracing.Jaeger.Throttler.RefreshInterval                              | TRACING_JAEGER_THROTTLER_REFRESH_INTERVAL                                  |         | time.Duration         |
| tracing.Jaeger.Throttler.SynchronousInitialization                    | TRACING_JAEGER_THROTTLER_SYNCHRONOUS_INITIALIZATION                        |         | bool                  |
