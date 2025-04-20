-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    description TEXT,
    owner_id UUID NOT NULL,
    notify_before INTERVAL
);

-- Индексы
CREATE INDEX idx_events_owner ON events(owner_id);
CREATE INDEX idx_events_time_range ON events(start_time, end_time);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
