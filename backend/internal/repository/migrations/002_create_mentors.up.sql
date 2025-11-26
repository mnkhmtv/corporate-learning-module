CREATE TABLE IF NOT EXISTS mentors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    jobTitle VARCHAR(255) NOT NULL,
    experience TEXT,
    workload INTEGER DEFAULT 0 CHECK (workload >= 0 AND workload <= 5),
    email VARCHAR(255) UNIQUE NOT NULL,
    telegram VARCHAR(100),
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_mentors_workload ON mentors(workload);
CREATE INDEX idx_mentors_email ON mentors(email);
