package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

func getPostsID(w http.ResponseWriter, r *http.Request) {
	pidStr := chi.URLParam(r, "id")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	results := []Post{}
	query := `
SELECT p.id         AS post_id,
       user_id      AS post_user_id,
       body,
       mime,
       p.created_at AS post_created_at,
       u.id         AS user_id,
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

	me := getSessionUser(r)

	//fmap := template.FuncMap{
	//	"imageURL": imageURL,
	//}

	templates["get_post_id"].Execute(w, struct {
		Post Post
		Me   User
	}{p, me})
}
