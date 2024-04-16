CREATE TABLE user_chatrooms (
    user_id BIGINT NOT NULL,
    chatroom_id BIGINT NOT NULL,
    PRIMARY KEY (user_id, chatroom_id),
    FOREIGN KEY (user_id) REFERENCES users,
    FOREIGN KEY (chatroom_id) REFERENCES chatrooms
);


