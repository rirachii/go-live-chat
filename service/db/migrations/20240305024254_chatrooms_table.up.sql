CREATE TYPE Message AS (
    sender_id BIGINT,
    username VARCHAR,
    msg TEXT
);

CREATE TABLE chatrooms (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    owner_id BIGINT NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES users (id),
    admins UserAccount[],
    location JSON,
    logs Message[]
);