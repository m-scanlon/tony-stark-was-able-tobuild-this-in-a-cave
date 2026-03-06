-- Job Registry: passive lifecycle tracker for all jobs in the system.
-- Does not make placement or routing decisions — those are owned by the Estimator.
-- Source of truth for job state from creation through completion or failure.

CREATE TABLE IF NOT EXISTS jobs (
  job_id        TEXT PRIMARY KEY,
  event_id      TEXT NOT NULL,
  agent_id      TEXT,

  status        TEXT NOT NULL DEFAULT 'created',
  -- created | routed | planning | executing | completed | failed

  lane          TEXT,
  -- fast_local | deep_reasoning

  shard_id      TEXT,
  -- which shard the job was routed to (set by Estimator)

  received_at   TEXT NOT NULL,
  routed_at     TEXT,
  -- when the Estimator made the placement decision
  started_at    TEXT,
  completed_at  TEXT
);

CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs(status);
CREATE INDEX IF NOT EXISTS idx_jobs_agent ON jobs(agent_id);
