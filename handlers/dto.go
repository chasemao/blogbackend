package handlers

const (
	codeNotExistArticle = 100100
)

type articleEntry struct {
	Title string `json:"title"`
}

type article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type listArticleResp struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data []*articleEntry `json:"data"`
}

type getArticleReq struct {
	Title string `json:"title"`
}

type getArticleResp struct {
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
	Data *article `json:"data"`
}

type errResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
