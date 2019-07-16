DROP INDEX IF EXISTS "Profiles_username_uindex", "Profiles_phone_uindex", "Profiles_id_uindex", "Profiles_email_uindex";
DROP TABLE IF EXISTS "Profiles";

DROP EXTENSION IF EXISTS "uuid-ossp";