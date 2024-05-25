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