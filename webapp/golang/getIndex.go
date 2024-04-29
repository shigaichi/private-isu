package main

import (
	"html/template"
	"log"
	"net/http"
)

func getIndex(w http.ResponseWriter, r *http.Request) {
	me := getSessionUser(r)

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

	fmap := template.FuncMap{
		"imageURL": imageURL,
	}

	template.Must(template.New("layout.html").Funcs(fmap).ParseFiles(
		getTemplPath("layout.html"),
		getTemplPath("index.html"),
		getTemplPath("posts.html"),
		getTemplPath("post.html"),
	)).Execute(w, struct {
		Posts     []Post
		Me        User
		CSRFToken string
		Flash     string
	}{posts, me, getCSRFToken(r), getFlash(w, r, "notice")})
}
