CREATE TABLE IF NOT EXISTS training_requests (
    id BIGSERIAL PRIMARY KEY,
    userId BIGSERIAL NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    topic VARCHAR(500) NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_training_requests_user_id ON training_requests(userId);
CREATE INDEX idx_training_requests_status ON training_requests(status);
CREATE INDEX idx_training_requests_created_at ON training_requests(createdAt DESC);
