DELETE
FROM cron.job
WHERE jobid IN (SELECT jobid
                FROM cron.job
                WHERE command LIKE '%Sessions%');

DROP EXTENSION IF EXISTS "pg_cron";

DROP TABLE IF EXISTS "Sessions", "Logins", "Users";

DROP EXTENSION IF EXISTS "uuid-ossp";


