CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "Profiles" (
    id         UUID                        DEFAULT uuid_generate_v4() NOT NULL,
    username   VARCHAR(100)                                           NOT NULL,
    email      VARCHAR(100)                                           NOT NULL,
    phone      VARCHAR(20),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now()              NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now()              NOT NULL,
    CONSTRAINT "Profiles_pk"
        PRIMARY KEY (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS "Profiles_email_uindex"
    ON "Profiles" (email);
CREATE UNIQUE INDEX IF NOT EXISTS "Profiles_id_uindex"
    ON "Profiles" (id);
CREATE UNIQUE INDEX IF NOT EXISTS "Profiles_phone_uindex"
    ON "Profiles" (phone);
CREATE UNIQUE INDEX IF NOT EXISTS "Profiles_username_uindex"
    ON "Profiles" (username);



