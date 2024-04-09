CREATE TABLE user_display_names(
    user_id BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users
    display_name VARCHAR(20)
);