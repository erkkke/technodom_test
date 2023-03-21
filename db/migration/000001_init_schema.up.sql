CREATE TABLE "redirects"
(
    "id"           BIGSERIAL PRIMARY KEY,
    "active_link"  VARCHAR NOT NULL,
    "history_link" VARCHAR NOT NULL
);