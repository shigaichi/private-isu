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
SELECT p.id, user_id, body, mime, p.created_at
FROM posts p
         STRAIGHT_JOIN users u ON u.id = p.user_id
WHERE u.del_flg = 0
ORDER BY p.created_at DESC
LIMIT 20;
`
	err := db.Select(&results, query)
	if err != nil {
		log.Print(err)
		return
	}

	posts, err := makePosts(results, getCSRFToken(r), false)
	if err != nil {
		log.Print(err)
		return
	}

	//fmap := template.FuncMap{
	//	"imageURL": imageURL,
	//}

	templates["get_index"].Execute(w, struct {
		Posts     []Post
		Me        User
		CSRFToken string
		Flash     string
	}{posts, me, getCSRFToken(r), getFlash(w, r, "notice")})
}
