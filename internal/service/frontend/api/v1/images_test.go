package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractImagePaths(t *testing.T) {
	logs := `saved: /data1/project/output/image1.jpg
saved: /data1/project/output/result.png
saved: /data1/project/output/image1.jpg
ignore: /tmp/report.txt`

	images := extractImagePaths(logs, "/data1/project", "/images")
	require.Equal(t, []string{
		"/images/output/image1.jpg",
		"/images/output/result.png",
	}, images)
}

func TestHandleGetImages(t *testing.T) {
	api := &API{}
	req := httptest.NewRequest(
		http.MethodGet,
		"/images?sourcePrefix=/data1/project&servePrefix=/images&log=%2Fdata1%2Fproject%2Foutput%2Fa.jpeg&log=%2Fdata1%2Fproject%2Foutput%2Fb.gif",
		nil,
	)
	rec := httptest.NewRecorder()

	api.handleGetImages(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp imagesResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, []string{
		"/images/output/a.jpeg",
		"/images/output/b.gif",
	}, resp.Images)
}
