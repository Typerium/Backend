CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "Users" (
    id         UUID                        DEFAULT uuid_generate_v4() PRIMARY KEY,
    password   VARCHAR(150),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS "Logins" (
    id         UUID                     DEFAULT uuid_generate_v4() NOT NULL,
    login      VARCHAR(80)                                         NOT NULL,
    user_id    UUID                                                NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()              NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()              NOT NULL,
    CONSTRAINT "Logins_pk"
        PRIMARY KEY (id),
    CONSTRAINT "Logins_Users_id_fk"
        FOREIGN KEY (user_id) REFERENCES "Users"
            ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS "Logins_id_uindex"
    ON "Logins" (id);
CREATE UNIQUE INDEX IF NOT EXISTS "Logins_login_uindex"
    ON "Logins" (login);
CREATE UNIQUE INDEX IF NOT EXISTS "Logins_login_user_id_uindex"
    ON "Logins" (login, user_id);

CREATE TABLE IF NOT EXISTS "Sessions" (
    id            UUID                        DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id       UUID                                      NOT NULL,
    time_exp      TIMESTAMP WITHOUT TIME ZONE,
    key_signature BYTEA,
    created_at    TIMESTAMP WITHOUT TIME ZONE DEFAULT now() NOT NULL,
    updated_at    TIMESTAMP WITHOUT TIME ZONE DEFAULT now() NOT NULL,
    CONSTRAINT "Sessions_Users_id_fk"
        FOREIGN KEY (user_id) REFERENCES "Users"
            ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE EXTENSION IF NOT EXISTS "pg_cron";

SELECT cron.schedule('@hourly', $$DELETE FROM "Sessions" WHERE time_exp < now()$$);
