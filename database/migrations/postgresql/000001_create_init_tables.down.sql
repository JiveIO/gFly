-- Delete tables
-- Make sure user have permission `SUPERUSER`.
-- Example user `user`. SQL> ALTER USER "user" WITH SUPERUSER;

SET session_replication_role = 'replica';
DROP TABLE IF EXISTS address CASCADE;
DROP TABLE IF EXISTS user_roles CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TYPE IF EXISTS user_status;
DROP TYPE IF EXISTS address_type;
SET session_replication_role = 'origin';