CREATE TABLE comments (
                         "id" bigserial PRIMARY KEY,
                         "movie_title" varchar NOT NULL,
                         "movie_id" integer NOT NULL,
                         "author" varchar NOT NULL,
                         "content" varchar NOT NULL,
                         "created_at" timestamptz NOT NULL DEFAULT (now())
);