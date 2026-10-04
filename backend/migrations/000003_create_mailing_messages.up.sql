CREATE TABLE mailing_messages (
    id BIGSERIAL PRIMARY KEY,
    mailing_id BIGINT NOT NULL REFERENCES mailings (id) ON DELETE CASCADE,
    contact_id BIGINT NOT NULL REFERENCES contacts (id),
    status TEXT NOT NULL DEFAULT 'pending' CHECK (
        status IN (
            'pending',
            'sending',
            'sent',
            'failed'
        )
    ),
    attempts INT NOT NULL DEFAULT 0,
    next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    locked_at TIMESTAMPTZ,
    sent_at TIMESTAMPTZ,
    UNIQUE (mailing_id, contact_id)
);

CREATE INDEX idx_mailing_messages_queue ON mailing_messages (next_attempt_at)
WHERE
    status = 'pending';

CREATE INDEX idx_mailing_messages_stuck ON mailing_messages (locked_at)
WHERE
    status = 'sending';

CREATE INDEX idx_mailing_messages_mailing_status ON mailing_messages (mailing_id, status);