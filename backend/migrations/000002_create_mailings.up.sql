CREATE TABLE mailings (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    subject TEXT NOT NULL,
    body_template TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'draft' CHECK (
        status IN (
            'draft',
            'running',
            'paused',
            'done'
        )
    ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    started_at TIMESTAMPTZ
);