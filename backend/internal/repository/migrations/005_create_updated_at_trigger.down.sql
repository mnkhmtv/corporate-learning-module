DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_mentors_updated_at ON mentors;
DROP TRIGGER IF EXISTS update_training_requests_updated_at ON training_requests;
DROP TRIGGER IF EXISTS update_learning_processes_updated_at ON learning_processes;
DROP FUNCTION IF EXISTS update_updated_at_column();
