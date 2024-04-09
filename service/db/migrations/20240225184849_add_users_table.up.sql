CREATE TYPE USER_ACCOUNT as (
    email VARCHAR,
    username VARCHAR,
    hashed_password VARCHAR
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    user_acc USER_ACCOUNT
);





