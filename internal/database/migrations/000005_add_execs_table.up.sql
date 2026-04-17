CREATE TABLE IF NOT EXISTS execs (
    id SERIAL PRIMARY KEY, 
    first_name VARCHAR(255) NOT NULL, 
    last_name VARCHAR(255) NOT NULL, 
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    password_changed_at TIMESTAMPTZ,
    password_reset_token VARCHAR(255),
    password_code_expires TIMESTAMPTZ,
    role VARCHAR(50) NOT NULL,
    inactive_status BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_execs_email ON execs(email);
CREATE INDEX IF NOT EXISTS idx_execs_username ON execs(username);

ALTER SEQUENCE IF EXISTS execs_id_seq RESTART WITH 100;
