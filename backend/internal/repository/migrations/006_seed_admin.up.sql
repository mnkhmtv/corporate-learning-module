INSERT INTO users (id, name, email, password_hash, role, createdAt, updatedAt)
VALUES (
    gen_random_uuid(),
    'System Admin',
    'admin@example.com',
    '$2a$10$0flahX7a6Viyc8i8YU/t6.obkr.3ZUFO3RRGqIL5p19GqSIKhiA8i',
    'admin',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
ON CONFLICT (email) DO NOTHING;
