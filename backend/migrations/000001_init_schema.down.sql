-- Revert Initial Schema Migration
DROP TABLE IF EXISTS print3d_requests CASCADE;
DROP TABLE IF EXISTS loan_requests CASCADE;
DROP TABLE IF EXISTS inventory_items CASCADE;
DROP TABLE IF EXISTS resource_requests CASCADE;
DROP TABLE IF EXISTS project_updates CASCADE;
DROP TABLE IF EXISTS project_members CASCADE;
DROP TABLE IF EXISTS projects CASCADE;
DROP TABLE IF EXISTS events CASCADE;
DROP TABLE IF EXISTS registration_requests CASCADE;
DROP TABLE IF EXISTS users CASCADE;

DROP TYPE IF EXISTS event_visibility;
DROP TYPE IF EXISTS print3d_status;
DROP TYPE IF EXISTS loan_status;
DROP TYPE IF EXISTS resource_status;
DROP TYPE IF EXISTS resource_type;
DROP TYPE IF EXISTS update_status;
DROP TYPE IF EXISTS project_status;
DROP TYPE IF EXISTS registration_status;
DROP TYPE IF EXISTS user_role;
