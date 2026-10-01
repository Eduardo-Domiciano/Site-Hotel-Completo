CREATE TABLE event_quotes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type TEXT NOT NULL,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT NOT NULL,
    rooms INT NOT NULL DEFAULT 0,
    people INT NOT NULL DEFAULT 0,
    event_at TIMESTAMPTZ NULL,
    coffee BOOLEAN NOT NULL DEFAULT FALSE,
    lunch BOOLEAN NOT NULL DEFAULT FALSE,
    dinner BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE job_applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    area TEXT NOT NULL,
    name TEXT NOT NULL,
    phone TEXT NULL,
    email TEXT NOT NULL,
    resume_path TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX event_quotes_created_at_idx ON event_quotes(created_at DESC);
CREATE INDEX job_applications_created_at_idx ON job_applications(created_at DESC);
