CREATE TABLE IF NOT EXISTS learning_processes (
    id BIGSERIAL PRIMARY KEY,
    requestId BIGSERIAL NOT NULL REFERENCES training_requests(id) ON DELETE CASCADE,
    userId BIGSERIAL NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    mentorId BIGSERIAL NOT NULL REFERENCES mentors(id) ON DELETE RESTRICT,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'completed')),
    plan JSONB DEFAULT '[]'::jsonb,
    notes TEXT,
    
    -- Feedback fields (filled when status = 'completed')
    feedbackRating INTEGER CHECK (feedbackRating >= 1 AND feedbackRating <= 5),
    feedbackComment TEXT,
    
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completedAt TIMESTAMP WITH TIME ZONE,
    
    -- Ensure one learning process per request
    CONSTRAINT unique_request_learning UNIQUE(request_id)
);

CREATE INDEX idx_learning_processes_user_id ON learning_processes(userId);
CREATE INDEX idx_learning_processes_mentor_id ON learning_processes(mentorId);
CREATE INDEX idx_learning_processes_status ON learning_processes(status);
CREATE INDEX idx_learning_processes_request_id ON learning_processes(requestId);

-- GIN index for JSONB plan field (for efficient querying)
CREATE INDEX idx_learning_processes_plan ON learning_processes USING GIN (plan);
