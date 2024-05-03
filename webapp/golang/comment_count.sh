#!/bin/bash

# データベース接続設定
DB_NAME="isuconp"
DB_USER="isuconp"
DB_PASS="isuconp"
DB_HOST="127.0.0.1"

# p.comment_count が NULL の post の id を取得
POST_IDS=$(mysql -u"$DB_USER" -p"$DB_PASS" -h "$DB_HOST" -D"$DB_NAME" -se "SELECT id FROM posts WHERE comment_count IS NULL;")

# 各 post_id に対してコメント数をカウントし、posts テーブルを更新
for POST_ID in $POST_IDS; do
    # コメント数をカウント
    COUNT=$(mysql -u"$DB_USER" -p"$DB_PASS" -h "$DB_HOST" -D"$DB_NAME" -se "SELECT COUNT(*) FROM comments WHERE post_id = $POST_ID;")

    # posts テーブルを更新
    mysql -u"$DB_USER" -p"$DB_PASS" -h "$DB_HOST" -D"$DB_NAME" -se "UPDATE posts SET comment_count = $COUNT WHERE id = $POST_ID;"
#  echo "$COUNT"
done

echo "Update complete."
