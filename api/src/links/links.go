package links

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	"url-shortener/src/database"
	"url-shortener/src/types"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type createRequest struct {
	URL string `json:"url"`
}

func generateCode(n int) (string, error) {
	b := make([]byte, n)
	for i := range b {
		idxBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		b[i] = alphabet[idxBig.Int64()]
	}
	return string(b), nil
}

func CreateLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "this method is not supported!", http.StatusMethodNotAllowed)
		return
	}

	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "url required, but empty", http.StatusBadRequest)
		return
	}

	codeLenStr := os.Getenv("CODE_LENGTH")
	codeLen, err := strconv.Atoi(codeLenStr)
	if err != nil || codeLen <= 0 {
		codeLen = 6
	}

	link := types.Link{
		URL: req.URL,
	}

	for {
		code, err := generateCode(codeLen)
		if err != nil {
			http.Error(w, "failed to generate code", http.StatusInternalServerError)
			return
		}
		link.Code = code

		if err := database.DB.Create(&link).Error; err != nil {
			continue
		}

		break
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(link)
}

func RedirectLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "this method is not supported!", http.StatusMethodNotAllowed)
		return
	}

	code := mux.Vars(r)["code"]
	var link types.Link
	if err := database.DB.Where("code = ?", code).First(&link).Error; err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	database.DB.Model(&link).Update("clicks", link.Clicks+1)
	http.Redirect(w, r, link.URL, http.StatusFound)
}
