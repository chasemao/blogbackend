package handlers

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	artcileDir = "../blogarticle"
	imageDir   = "../blogarticle/static"

	ctimeLength     = 10
	titleStartIndex = 11
)

type ArticleLogic interface {
	ListArticles(c *gin.Context)
	GetArticle(c *gin.Context)
	GetImage(c *gin.Context)
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
		if a.getTitleFromFileName(entry.Name()) == req.Title {
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
	for i := len(fileEntries) - 1; i >= 0; i-- {
		fileEntry := fileEntries[i]
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
			Title: a.getTitleFromFileName(fileEntry.Name()),
			CTime: a.getCTimeFromFileName(fileEntry.Name()),
		})
	}
	return res
}

func (a *articleLogicImpl) getTitleFromFileName(name string) string {
	extLen := len(filepath.Ext(name))
	if len(name) < titleStartIndex {
		return name
	}
	// remove time prefix and extension suffix
	return name[titleStartIndex : len(name)-extLen]
}

func (a *articleLogicImpl) getCTimeFromFileName(name string) string {
	if len(name) < ctimeLength {
		return name
	}
	return name[:ctimeLength]
}

func (a *articleLogicImpl) getArticleFromDisk(file fs.DirEntry) (*article, error) {
	contentBytes, err := os.ReadFile(artcileDir + "/" + file.Name())
	if err != nil {
		return nil, err
	}
	return &article{
		Title:   a.getTitleFromFileName(file.Name()),
		CTime:   a.getCTimeFromFileName(file.Name()),
		Content: string(contentBytes),
	}, nil
}

func (a *articleLogicImpl) GetImage(c *gin.Context) {
	var req getImageReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	imagePath := filepath.Join(imageDir, req.Image) // Assuming imageDir is defined elsewhere

	// Read the image file
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		a.buildErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	// Determine content type based on file extension
	contentType := getImageContentType(imagePath)

	// Set the appropriate Content-Type header
	c.Header("Content-Type", contentType)

	// Send the image data as response
	c.Data(http.StatusOK, contentType, imageData)
}

// Function to get content type based on file extension
func getImageContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	// Add more cases as needed for other image formats
	default:
		return "application/octet-stream" // Default to binary data if content type is unknown
	}
}
