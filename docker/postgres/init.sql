CREATE TABLE IF NOT EXISTS tasks (
  id TEXT PRIMARY KEY,
  message TEXT NOT NULL,
  binary_code TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS task_outbox_events (
  task_id TEXT PRIMARY KEY,
  status TEXT NOT NULL,
  message TEXT NOT NULL,
  binary_code TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_task_outbox_events_status_created_at
  ON task_outbox_events (status, created_at);
