CREATE TABLE IF NOT EXISTS tasks
(
    id          UUID      NOT NULL,
    title       TEXT      NOT NULL,
    description TEXT      NOT NULL,
    status      TEXT      NOT NULL,
    due_date    TIMESTAMP NOT NULL,
    created_at  TIMESTAMP NOT NULL,

    CONSTRAINT PK_TASKS PRIMARY KEY (id)
);