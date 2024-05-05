package main

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

var indexCounts int

func getIndex(w http.ResponseWriter, r *http.Request) {
	me := getSessionUser(r)

	const cacheKey = "index-page"
	if me.ID == 0 {
		// ログインしていない場合はキャッシュを試みる
		item, err := memcacheClient.Get(cacheKey)
		indexCounts++
		// 2024/05/06 02:47:29 getIndex.go:21: cache hit in getindex count: 7269(10271)
		log.Printf("cache hit in getindex count: %d\n", indexCounts)
		if err == nil {
			w.Write(item.Value)
			return
		}
	}

	// 投稿一覧画面ではAccountNameしか使われていない
	results := []Post{}

	query := `
SELECT p.id         AS post_id,
       user_id      AS post_user_id,
       body,
       mime,
       p.created_at AS post_created_at,
       comment_count,
#        u.id         AS user_id,
       u.account_name,
       u.passhash,
       u.authority,
       u.del_flg,
       u.created_at AS user_created_at
FROM posts p
         STRAIGHT_JOIN users u ON u.id = p.user_id
WHERE u.del_flg = 0
ORDER BY p.created_at DESC
LIMIT 20;
`
	token := getCSRFToken(r)
	err := db.Select(&results, query)
	if err != nil {
		log.Print(err)
		return
	}

	posts, err := makePosts(results, token, false)
	if err != nil {
		log.Print(err)
		return
	}

	// キャッシュになければテンプレートをレンダリング
	var tpl bytes.Buffer
	if err := templates["get_index"].Execute(&tpl, struct {
		Posts     []Post
		Me        User
		CSRFToken string
		Flash     string
	}{
		Posts:     posts,
		Me:        me,
		CSRFToken: token,
		Flash:     getFlash(w, r, "notice"),
	}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	htmlContent := tpl.Bytes()
	if me.ID == 0 {
		// HTMLコンテンツをキャッシュに保存
		item := &memcache.Item{
			Key:        cacheKey,
			Value:      htmlContent,
			Expiration: int32(30 * time.Minute / time.Second),
		}
		memcacheClient.Set(item)
	}

	w.Write(htmlContent)

	//templates["get_index"].Execute(w, struct {
	//	Posts     []Post
	//	Me        User
	//	CSRFToken string
	//	Flash     string
	//}{posts, me, token, getFlash(w, r, "notice")})
}
