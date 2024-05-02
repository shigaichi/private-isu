package main

import (
	"log"
	"net/http"
	"net/url"
	"time"
)

func getPosts(w http.ResponseWriter, r *http.Request) {
	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}
	maxCreatedAt := m.Get("max_created_at")
	if maxCreatedAt == "" {
		return
	}

	t, err := time.Parse(ISO8601Format, maxCreatedAt)
	if err != nil {
		log.Print(err)
		return
	}

	results := []Post{}
	query := `
	SELECT p.id         AS post_id,
		   user_id      AS post_user_id,
		   body,
		   mime,
		   p.created_at AS post_created_at,
# 		   u.id         AS user_id,
		   u.account_name,
		   u.passhash,
		   u.authority,
		   u.del_flg,
		   u.created_at AS user_created_at
	FROM posts AS p STRAIGHT_JOIN users u ON u.id = p.user_id
	WHERE p.created_at <= ? AND u.del_flg = 0
	ORDER BY p.created_at DESC
	LIMIT 20
`
	err = db.Select(&results, query, t.Format(ISO8601Format))
	if err != nil {
		log.Print(err)
		return
	}

	posts, err := makePosts(results, getCSRFToken(r), false)
	if err != nil {
		log.Print(err)
		return
	}

	if len(posts) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//fmap := template.FuncMap{
	//	"imageURL": imageURL,
	//}

	templates["get_post"].Execute(w, posts)
}
