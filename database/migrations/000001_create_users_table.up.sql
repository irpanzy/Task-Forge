CREATE EXTENSION
IF NOT EXISTS "uuid-ossp";
    CREATE TABLE IF NOT EXISTS users
        (
            internal_id BIGSERIAL PRIMARY KEY                         ,
            public_id   UUID DEFAULT uuid_generate_v4() NOT NULL      ,
            name        VARCHAR(255) NOT NULL                         ,
            email       VARCHAR(255) NOT NULL                         ,
            password    VARCHAR(255) NOT NULL                         ,
            role        VARCHAR(50) DEFAULT 'user' NOT NULL           ,
            created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
            updated_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
            deleted_at  TIMESTAMPTZ                                   ,
            CONSTRAINT uq_users_public_id UNIQUE (public_id)          ,
            CONSTRAINT uq_users_email UNIQUE (email)                  ,
            CONSTRAINT chk_users_role CHECK (role IN ('admin', 'user'))
        )
    ;
    CREATE INDEX IF
    NOT EXISTS idx_users_deleted_at ON users
        (
            deleted_at
        );