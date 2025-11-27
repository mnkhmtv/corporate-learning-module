CREATE TABLE IF NOT EXISTS learning_processes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    requestId UUID NOT NULL REFERENCES training_requests(id) ON DELETE CASCADE,
    userId UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    mentorId UUID NOT NULL REFERENCES mentors(id) ON DELETE RESTRICT,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'completed')),
    startDate TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    endDate TIMESTAMP WITH TIME ZONE,
    
    plan JSONB DEFAULT '[]'::jsonb,
    notes TEXT,
    
    -- Feedback (JSONB: {rating: int, comment: string})
    feedback JSONB,
    
    -- Audit fields
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure one learning process per request
    CONSTRAINT unique_request_learning UNIQUE(requestId)
);

CREATE INDEX idx_learning_processes_userId ON learning_processes(userId);
CREATE INDEX idx_learning_processes_mentorId ON learning_processes(mentorId);
CREATE INDEX idx_learning_processes_status ON learning_processes(status);
CREATE INDEX idx_learning_processes_requestId ON learning_processes(requestId);

-- GIN index for JSONB plan field
CREATE INDEX idx_learning_processes_plan ON learning_processes USING GIN (plan);
