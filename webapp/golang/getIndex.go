package main

import (
	"log"
	"net/http"
)

func getIndex(w http.ResponseWriter, r *http.Request) {
	me := getSessionUser(r)

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
	log.Println("token: " + token)
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

	templates["get_index"].Execute(w, struct {
		Posts     []Post
		Me        User
		CSRFToken string
		Flash     string
	}{posts, me, token, getFlash(w, r, "notice")})
}
