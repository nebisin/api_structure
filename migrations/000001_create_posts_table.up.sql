CREATE TABLE IF NOT EXISTS posts (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    body text NOT NULL,
    tags text[],
    version integer NOT NULL DEFAULT 1
);