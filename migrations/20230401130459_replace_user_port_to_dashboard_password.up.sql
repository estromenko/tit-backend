ALTER TABLE users DROP COLUMN port;
ALTER TABLE users ADD COLUMN dashboard_password VARCHAR(64) DEFAULT '';
