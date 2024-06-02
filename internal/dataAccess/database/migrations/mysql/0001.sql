-- +migrate Up

-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS users (
                                     user_id BIGINT PRIMARY KEY,
                                     username VARCHAR(100),
    password VARCHAR(100)
    );

CREATE TABLE IF NOT EXISTS download_task (
                                             of_user_id BIGINT PRIMARY KEY,
                                             download_type SMALLINT NOT NULL,
                                             url TEXT NOT NULL,
                                             download_status SMALLINT NOT NULL,
                                             metadata TEXT NOT NULL,
                                             FOREIGN KEY (of_user_id) REFERENCES users(user_id)
    );

CREATE TABLE IF NOT EXISTS account_passwords (
    of_account_id BIGINT UNSIGNED PRIMARY KEY,
    hash VARCHAR(128) NOT NULL,
    FOREIGN KEY (of_account_id) REFERENCES account(id)
)

CREATE TABLE IF NOT EXISTS token_public_keys (
    id BIGINT UNSIGNED PRIMARY KEY,
    public_key VARBINARY(4096) NOT NULL
)
-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin
DROP TABLE IF EXISTS download_tasks;
DROP TABLE IF EXISTS token_public_keys;
DROP TABLE IF EXISTS account_passwords;
DROP TABLE IF EXISTS users;
-- +migrate StatementEnd