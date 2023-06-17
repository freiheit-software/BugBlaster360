CREATE TABLE check_runs (
  id SERIAL PRIMARY KEY,
  check_run_id TEXT,
  commit_hash CHAR(40),
  created_at TIMESTAMPTZ DEFAULT NOW()
);