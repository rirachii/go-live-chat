CREATE TYPE USER_ACCOUNT as (
    email VARCHAR,
    hashed_password VARCHAR
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR,
    user_account USER_ACCOUNT
);





