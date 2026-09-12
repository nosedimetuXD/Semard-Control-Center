-- Rollback Migración 000002
DROP TABLE IF EXISTS event_registrations;
DELETE FROM users WHERE email = 'dfrancob1@unicartagena.edu.co';
