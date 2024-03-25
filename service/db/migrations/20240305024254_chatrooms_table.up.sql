CREATE TYPE CHAT_MESSAGE AS (
    sender_id BIGINT,
    username VARCHAR,
    msg TEXT
);

CREATE TABLE chatrooms (
    id BIGSERIAL PRIMARY KEY,
    room_name VARCHAR NOT NULL,
    owner_id BIGINT NOT NULL,
        FOREIGN KEY (owner_id) REFERENCES users (id),
    is_public BOOLEAN,
    is_active BOOLEAN,
    logs CHAT_MESSAGE[]
    -- location JSON
);
