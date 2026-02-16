# Skyra Unsupervised Task Estimator

Internal engineering documentation for implementation and operations.

## 1. Overview

The Skyra Unsupervised Task Estimator predicts request size and runtime behavior without hard-coded heuristics. It learns from real executions by clustering request embeddings and maintaining rolling per-cluster statistics.

Primary outcomes:

- predict output token range
- predict latency range
- estimate tool-use probability
- estimate failure/interrupt risk
- inform orchestration lane and sync/async choice

### Why unsupervised clustering instead of heuristics

Heuristic rules (keywords, manual intent buckets) drift quickly in mixed personal workloads. Unsupervised clustering adapts to real usage and captures latent task patterns:

- no manual taxonomy required up front
- updates naturally with telemetry
- supports novelty detection by distance/spread
- reduces policy churn in orchestrator code

## 2. System Role in the Skyra Pipeline

Estimator position in runtime:

```text
Listener -> Front-door -> Estimator -> Orchestrator -> Agents/GPU
```

Detailed path:

```text
+-------------------+      +-------------------+      +----------------------+
| Listener Node     | ---> | Front-door Model  | ---> | Estimator (Control)  |
| Wake/VAD/STT/Gate |      | quick ack + parse |      | cluster-based predict|
+-------------------+      +-------------------+      +----------------------+
                                                            |
                                                            v
                                                    +------------------+
                                                    | Orchestrator     |
                                                    | lane + mode pick |
                                                    +------------------+
                                                            |
                                                            v
                                                    +------------------+
                                                    | Agents / GPU     |
                                                    +------------------+
```

The listener remains responsive because estimator work happens in the control plane.

## 3. Telemetry Model

The estimator relies on two event types.

### 3.1 Request telemetry

Captured before execution:

- `request_id`
- `timestamp`
- `user_text`
- `tokens_in_est`
- `lane_chosen`
- `execution_mode`
- `model_target`

Purpose by field:

- `request_id`: joins request and outcome events.
- `timestamp`: enables recency windows and drift checks.
- `user_text`: embedding source for clustering.
- `tokens_in_est`: early complexity signal.
- `lane_chosen`: policy feedback and retrospective analysis.
- `execution_mode`: compare sync vs async choices by workload.
- `model_target`: segment behavior by model/backend.

### 3.2 Outcome telemetry

Captured after execution completion/interruption:

- `tokens_out_actual`
- `wall_time_ms`
- `tool_calls_count`
- `success` / `failure`
- `cancelled`
- `user_interrupted`

Purpose by field:

- `tokens_out_actual`: output size target for predictions.
- `wall_time_ms`: user-visible latency signal.
- `tool_calls_count`: complexity and orchestration overhead proxy.
- `success/failure`: reliability tracking per cluster.
- `cancelled`: incomplete run signal.
- `user_interrupted`: user satisfaction/latency mismatch indicator.

## 4. Embedding Strategy

### What is embedded

Embed a compact request representation:

- normalized user text
- optional front-door intent hint
- optional project/domain hint (high-confidence only)

Canonical input format:

```text
intent_hint|project_hint|normalized_user_text
```

### Why full context is not embedded

Full conversation/memory context is intentionally excluded from estimator embeddings:

- too expensive for online decisions
- adds noise unrelated to task shape
- destabilizes clusters across sessions
- increases tail latency in control plane

Estimator embeddings should represent request pattern, not final answer context.

### Performance considerations

- keep online embedding path low-latency
- cache embeddings by normalized-text hash
- batch offline embeddings for retraining
- isolate embedding model from listener node

## 5. Clustering Strategy

### 5.1 Initial approach: k-means

Use k-means first for operational simplicity.

How clusters are built:

1. collect embedding corpus for rolling training window
2. fit k-means with configured `k`
3. assign cluster IDs and store centroids
4. attach rolling outcome stats to each cluster

Retraining cadence:

- periodic full retrain (for example every 6-24h)
- continuous stats updates between retrains

### 5.2 Future approach: HDBSCAN

Planned upgrade for adaptive cluster count and outlier handling.

Benefits:

- no fixed `k` requirement
- better for variable-density workloads
- stronger novelty/outlier treatment

### 5.3 Cluster evolution

- version each cluster model
- keep assignment metadata for replay/debug
- roll forward with canary validation
- maintain compatibility windows for old model versions

## 6. Cluster Profile Model

Per-cluster rolling profile fields:

- `sample_count`
- `tokens_out_p50`, `tokens_out_p90`, `tokens_out_p95`
- `wall_time_ms_p50`, `wall_time_ms_p90`, `wall_time_ms_p95`
- `tool_call_rate`
- `failure_rate`
- `interrupt_rate`
- `cancel_rate`
- `cluster_spread`

### Derived buckets

`tokens_out_bucket`:

- `S`: `<300`
- `M`: `300-1199`
- `L`: `1200-3999`
- `XL`: `>=4000`

`latency_bucket`:

- `lt2s`: `<2000 ms`
- `2to10s`: `2000-9999 ms`
- `10to60s`: `10000-59999 ms`
- `gt60s`: `>=60000 ms`

Use percentiles, not averages, for scheduling safety.

## 7. Runtime Estimation Flow

1. new request arrives from front-door/orchestrator
2. normalize and embed request text
3. find nearest cluster centroid
4. load cluster profile stats
5. derive cost buckets from percentile profile
6. compute confidence score
7. return estimate to scheduler/orchestrator
8. scheduler selects execution lane and mode

## 8. Confidence and Novelty Detection

Confidence factors:

- distance to centroid (lower is better)
- cluster sample size (higher is better)
- cluster spread/tightness (tighter is better)

Low-confidence conditions:

- centroid distance above threshold
- low sample count
- high spread + ambiguous assignment

Low-confidence behavior:

- mark estimate as low confidence
- route conservatively (prefer async or stronger lane)
- emit novelty telemetry for retraining

## 9. APIs

### `POST /estimate`

Request:

```json
{
  "request_id": "req_abc123",
  "timestamp": "2026-02-16T22:00:00Z",
  "user_text": "summarize crash logs from last night and suggest root cause",
  "tokens_in_est": 740,
  "intent_hint": "ops.log_analysis",
  "project_hint": "servers"
}
```

Response:

```json
{
  "request_id": "req_abc123",
  "cluster_id": "c_007",
  "cluster_model_version": "km_v1_2026-02-16",
  "confidence": 0.82,
  "novel": false,
  "prediction": {
    "tokens_out_bucket": "L",
    "latency_bucket": "10to60s",
    "tokens_out_p50": 1800,
    "tokens_out_p90": 4200,
    "wall_time_ms_p50": 21000,
    "wall_time_ms_p90": 64000,
    "tool_call_rate": 0.73,
    "failure_rate": 0.08,
    "interrupt_rate": 0.05
  },
  "recommended_execution_mode": "async",
  "recommended_lane": "gpu"
}
```

### `POST /telemetry/request`

Request:

```json
{
  "request_id": "req_abc123",
  "timestamp": "2026-02-16T22:00:00Z",
  "user_text": "summarize crash logs from last night and suggest root cause",
  "tokens_in_est": 740,
  "lane_chosen": "gpu",
  "execution_mode": "async",
  "model_target": "deepseek-r1-32b"
}
```

### `POST /telemetry/outcome`

Request:

```json
{
  "request_id": "req_abc123",
  "timestamp": "2026-02-16T22:00:36Z",
  "tokens_out_actual": 2350,
  "wall_time_ms": 36120,
  "tool_calls_count": 4,
  "success": true,
  "cancelled": false,
  "user_interrupted": false
}
```

## 10. Failure Modes and Safeguards

- No cluster match:
  - return `novel=true`, low confidence
  - apply conservative lane policy
- Low confidence:
  - prefer safe async execution
  - avoid aggressive resource pinning
- Estimator offline:
  - fall back to orchestrator default policy
  - queue telemetry for later ingestion
- Cold start (no telemetry):
  - use neutral defaults
  - mark all predictions low confidence until minimum sample threshold

## 11. Implementation Roadmap

### Phase 1

- telemetry logging endpoints
- offline embedding corpus
- baseline k-means training job

### Phase 2

- rolling cluster stats store
- percentile and rate updates
- retrain cadence automation

### Phase 3

- runtime `/estimate` endpoint
- orchestrator integration for lane decisions

### Phase 4

- HDBSCAN evaluation path
- novelty/outlier handling
- canary model rollout

### Phase 5

- adaptive scheduling policies using estimator feedback
- policy tuning from interrupt/failure outcomes
