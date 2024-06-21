package handlers

const (
	codeNotExistArticle = 100100
	codeNotExistImage   = 100101
)

type articleEntry struct {
	Title string `json:"title"`
	CTime string `json:"ctime"`
}

type article struct {
	Title string `json:"title"`
	CTime string `json:"ctime"`

	Content string `json:"content"`
}

type errResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
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

type getImageReq struct {
	Image string `json:"image"`
}
