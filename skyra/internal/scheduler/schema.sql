CREATE TABLE IF NOT EXISTS jobs (
  job_id        TEXT PRIMARY KEY,
  event_id      TEXT NOT NULL,
  project_id    TEXT,

  status        TEXT NOT NULL DEFAULT 'queued',
  -- queued | running | completed | failed

  lane          TEXT,
  -- fast_local | deep_reasoning

  received_at   TEXT NOT NULL,
  started_at    TEXT,
  completed_at  TEXT
);

CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs(status);
CREATE INDEX IF NOT EXISTS idx_jobs_project ON jobs(project_id);
