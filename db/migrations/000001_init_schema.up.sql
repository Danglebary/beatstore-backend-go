CREATE TABLE "users" (
    "id" SERIAL PRIMARY KEY,
    "username" VARCHAR UNIQUE NOT NULL,
    "password" VARCHAR NOT NULL,
    "email" VARCHAR UNIQUE NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "beats" (
    "id" SERIAL PRIMARY KEY,
    "creator_id" bigint NOT NULL,
    "title" VARCHAR NOT NULL,
    "genre" VARCHAR NOT NULL,
    "key" VARCHAR NOT NULL,
    "bpm" bigint NOT NULL,
    "tags" VARCHAR NOT NULL,
    "s3_key" VARCHAR NOT NULL,
    "likes_count" bigint NOT NULL DEFAULT 0,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "likes" (
    "id" SERIAL PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "beat_id" bigint NOT NULL
);

ALTER TABLE
    "beats"
ADD
    FOREIGN KEY ("creator_id") REFERENCES "users" ("id");

ALTER TABLE
    "likes"
ADD
    FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE
    "likes"
ADD
    FOREIGN KEY ("beat_id") REFERENCES "beats" ("id");

CREATE INDEX ON "beats" ("creator");

CREATE INDEX ON "beats" ("key");

CREATE INDEX ON "beats" ("genre");

CREATE INDEX ON "beats" ("bpm");

CREATE INDEX ON "beats" ("tags");

CREATE INDEX ON "likes" ("user_id");

CREATE INDEX ON "likes" ("beat_id");

CREATE INDEX ON "likes" ("user_id", "beat_id");