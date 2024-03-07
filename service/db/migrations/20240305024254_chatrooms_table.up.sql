CREATE TABLE "chatrooms" (
    "roomid" bigserial PRIMARY KEY,
    "roomname" varchar NOT NULL,
    "creator" varchar NOT NULL,
    "admin" varchar NOT NULL
    "location" varchar NOT NULL
)