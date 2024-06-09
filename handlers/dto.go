package handlers

type articleEntry struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type errResp struct {
	ErrMsg string `json:"err_msg"`
}
