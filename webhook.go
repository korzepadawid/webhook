package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"
)

func signBody(secret, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	computed.Write(body)
	return []byte(computed.Sum(nil))
}

func verifySignature(secret []byte, signature string, body []byte) bool {

	const signaturePrefix = "sha1="
	const signatureLength = 45

	if len(signature) != signatureLength || !strings.HasPrefix(signature, signaturePrefix) {
		return false
	}

	actual := make([]byte, 20)
	hex.Decode(actual, []byte(signature[5:]))

	return hmac.Equal(signBody(secret, body), actual)
}

func webhook(w http.ResponseWriter, r *http.Request) {
	signature := r.Header.Get("X-Hub-Signature")
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err)
	}

	if verifySignature([]byte(os.Getenv("WEBHOOK_SECRET")), signature, body) {
		cmd := exec.Command("bash", "update.sh", os.Getenv("REPO_PATH"), os.Getenv("PM2_PID"))
		err := cmd.Run()

		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
	}

	w.WriteHeader(http.StatusForbidden)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/webhook", webhook).Methods("POST")
	err := http.ListenAndServe(":3000", router)

	if err != nil {
		log.Fatal(err)
	}
}
