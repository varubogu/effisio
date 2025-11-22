-- Drop audit_logs table
BEGIN;

DROP TABLE IF EXISTS audit_logs CASCADE;

COMMIT;
