CREATE TYPE UserAccount as (
    email VARCHAR,
    username VARCHAR,
    hashed_password VARCHAR
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    user_acc UserAccount
);









