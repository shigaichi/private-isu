package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/go-chi/chi/v5"
)

var postsIdCounts int

func getPostsID(w http.ResponseWriter, r *http.Request) {
	pidStr := chi.URLParam(r, "id")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	me := getSessionUser(r)

	var cacheKey string
	if me.ID == 0 {
		cacheKey = fmt.Sprintf("post-%d", pid)
		cachedContent, found := memcacheClient.Get(cacheKey)
		// 2000/7502以上HITしているので意味はありそう
		postsIdCounts++
		log.Printf("cache hit in get postsId count: %d\n", postsIdCounts)
		if found == nil {
			w.Write(cachedContent.Value)
			return
		}
	}

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
WHERE p.id = ? AND u.del_flg = 0 
ORDER BY p.created_at DESC
LIMIT 20;
`

	err = db.Select(&results, query, pid)
	if err != nil {
		log.Print(err)
		return
	}

	posts, err := makePosts(results, getCSRFToken(r), true)
	if err != nil {
		log.Print(err)
		return
	}

	if len(posts) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	p := posts[0]

	//fmap := template.FuncMap{
	//	"imageURL": imageURL,
	//}

	// キャッシュになければテンプレートをレンダリング
	var tpl bytes.Buffer
	if err := templates["get_post_id"].Execute(&tpl, struct {
		Post Post
		Me   User
	}{p, me}); err != nil {
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
		err := memcacheClient.Set(item)
		if err != nil {
			// エラー処理（必要に応じて）
			log.Printf("Failed to set cache: %v", err)
		}
	}

	w.Write(htmlContent)

	//templates["get_post_id"].Execute(w, struct {
	//	Post Post
	//	Me   User
	//}{p, me})
}
