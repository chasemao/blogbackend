package handlers

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	artcileDir = "../article"
)

type ArticleLogic interface {
	ListArticles(c *gin.Context)
	GetArticle(c *gin.Context)
}

func NewArticleLogic() ArticleLogic {
	return &articleLogicImpl{}
}

type articleLogicImpl struct{}

// ListArticles list article, source is file system
func (a *articleLogicImpl) ListArticles(c *gin.Context) {
	fileEntries, err := a.readFileEntries()
	if err != nil {
		a.buildErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	articleEntries := a.convertToArticleEntries(fileEntries)

	c.JSON(http.StatusOK, &listArticleResp{
		Data: articleEntries,
	})
}

// GetArticle get a article detail
func (a *articleLogicImpl) GetArticle(c *gin.Context) {
	req := &getArticleReq{}
	if err := c.BindJSON(req); err != nil {
		a.buildErrorResp(c, http.StatusBadRequest, err)
		return
	}

	fileEntries, err := a.readFileEntries()
	if err != nil {
		a.buildErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	var file fs.DirEntry
	for _, entry := range fileEntries {
		if entry.Name() == req.Title {
			file = entry
		}
	}
	if file == nil {
		c.JSON(http.StatusOK, &getArticleResp{
			Code: codeNotExistArticle,
			Msg:  fmt.Sprintf("req article name=[%s] not exist", req.Title),
		})
		return
	}

	article, err := a.getArticleFromDisk(file)
	if err != nil {
		a.buildErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &getArticleResp{
		Data: article,
	})
}

func (a *articleLogicImpl) readFileEntries() ([]fs.DirEntry, error) {
	fileEntries, err := os.ReadDir(artcileDir)
	if err != nil {
		return nil, err
	}
	var res []fs.DirEntry
	for _, fileEntry := range fileEntries {
		if fileEntry.IsDir() {
			continue
		}
		res = append(res, fileEntry)
	}
	return res, nil
}

func (a *articleLogicImpl) buildErrorResp(c *gin.Context, code int, err error) {
	log.Println("req=", c.Request, " err=", err)
	c.JSON(code, &errResp{
		Code: code,
		Msg:  err.Error(),
	})
}

func (a *articleLogicImpl) convertToArticleEntries(fileEntries []fs.DirEntry) []*articleEntry {
	var res []*articleEntry
	for _, fileEntry := range fileEntries {
		res = append(res, &articleEntry{
			Title: fileEntry.Name(),
		})
	}
	return res
}

func (a *articleLogicImpl) getArticleFromDisk(file fs.DirEntry) (*article, error) {
	contentBytes, err := os.ReadFile(artcileDir + "/" + file.Name())
	if err != nil {
		return nil, err
	}
	return &article{
		Title:   file.Name(),
		Content: string(contentBytes),
	}, nil
}
