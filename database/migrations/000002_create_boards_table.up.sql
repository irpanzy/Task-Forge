CREATE TABLE IF NOT EXISTS boards
    (
        internal_id     BIGSERIAL PRIMARY KEY                         ,
        public_id       UUID DEFAULT uuid_generate_v4() NOT NULL      ,
        owner_id        BIGINT NOT NULL                               ,
        owner_public_id UUID NOT NULL                                 ,
        title           VARCHAR(255) NOT NULL                         ,
        description     TEXT                                          ,
        due_date        TIMESTAMPTZ                                   ,
        created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
        updated_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
        deleted_at      TIMESTAMPTZ                                   ,
        CONSTRAINT uq_boards_public_id UNIQUE (public_id)             ,
        CONSTRAINT fk_boards_owner_id FOREIGN KEY (owner_id) REFERENCES users(internal_id) ON
        DELETE
            CASCADE
            ,
            CONSTRAINT fk_boards_owner_public_id FOREIGN KEY (owner_public_id) REFERENCES users(public_id) ON
        DELETE
            CASCADE );
CREATE INDEX IF
NOT EXISTS idx_boards_owner_id ON boards
    (
        owner_id
    )
;
CREATE INDEX IF
NOT EXISTS idx_boards_owner_public_id ON boards
    (
        owner_public_id
    )
;
CREATE INDEX IF
NOT EXISTS idx_boards_deleted_at ON boards
    (
        deleted_at
    );