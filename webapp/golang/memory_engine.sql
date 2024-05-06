CREATE TABLE isuconp.comments_backup AS SELECT * FROM isuconp.comments;
CREATE TABLE isuconp.posts_backup AS SELECT * FROM isuconp.posts;
CREATE TABLE isuconp.users_backup AS SELECT * FROM isuconp.users;

DROP TABLE IF EXISTS isuconp.comments;
CREATE TABLE isuconp.comments
(
    id         int auto_increment
        primary key,
    post_id    int not null,
    user_id    int not null,
    comment    varchar(1024) not null, -- TEXTから変更
    created_at timestamp default CURRENT_TIMESTAMP not null
)
    ENGINE=MEMORY
    CHARSET=utf8mb4;

CREATE INDEX comments_post_id_created_at_index ON isuconp.comments (post_id ASC, created_at) USING BTREE;
CREATE INDEX comments_user_id_index ON isuconp.comments (user_id);

DROP TABLE IF EXISTS isuconp.posts;
CREATE TABLE isuconp.posts
(
    id            int auto_increment
        primary key,
    user_id       int not null,
    mime          varchar(64) not null,
    body          varchar(1024) not null, -- TEXTから変更
    created_at    timestamp default CURRENT_TIMESTAMP not null,
    comment_count int default 0 null
)
    ENGINE=MEMORY
    CHARSET=utf8mb4;

CREATE INDEX posts_created_at_index ON isuconp.posts (created_at) USING BTREE;
CREATE INDEX posts_user_id_index ON isuconp.posts (user_id);

DROP TABLE IF EXISTS isuconp.users;
CREATE TABLE isuconp.users
(
    id           int auto_increment
        primary key,
    account_name varchar(64) not null,
    passhash     varchar(128) not null,
    authority    tinyint(1) default 0 not null,
    del_flg      tinyint(1) default 0 not null,
    created_at   timestamp default CURRENT_TIMESTAMP not null,
    constraint account_name
        unique (account_name)
)
    ENGINE=MEMORY
    CHARSET=utf8mb4;

INSERT INTO isuconp.comments (id, post_id, user_id, comment, created_at)
SELECT id, post_id, user_id, comment, created_at FROM isuconp.comments_backup;

INSERT INTO isuconp.posts (id, user_id, mime, body, created_at, comment_count)
SELECT id, user_id, mime, body, created_at, comment_count FROM isuconp.posts_backup;

INSERT INTO isuconp.users (id,account_name,passhash,authority,del_flg,created_at)
SELECT id,account_name,passhash,authority,del_flg,created_at FROM isuconp.users_backup;
