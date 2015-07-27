package mongoadmin

import (
	"encoding/json"
	// "errors"
	"github.com/maxwellhealth/mgo"
	// "github.com/maxwellhealth/mgo/bson"
	// auth "github.com/nabeken/negroni-auth"
	"io"
	// "log"
	"net/http"
	// "os"
	// "strconv"
	"strings"
)

func indexes(w http.ResponseWriter, req *http.Request) {

	sess, col, err := getCollection(req)
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
		return
	}

	defer sess.Close()

	idxs, err := col.Indexes()
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	marshaled, _ := json.Marshal(idxs)
	io.WriteString(w, string(marshaled))

}
func dropIndex(w http.ResponseWriter, req *http.Request) {

	sess, col, err := getCollection(req)
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
		return
	}

	defer sess.Close()

	keys := req.URL.Query().Get("keys")

	if len(keys) == 0 {
		w.WriteHeader(400)
		io.WriteString(w, "Missing keys param")
		return
	}

	ks := strings.Split(keys, ",")

	err = col.DropIndex(ks...)

	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
		return
	}
	io.WriteString(w, "OK")

}

func addIndex(w http.ResponseWriter, req *http.Request) {

	sess, col, err := getCollection(req)
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
		return
	}

	defer sess.Close()

	keys := req.URL.Query().Get("keys")
	sparse := req.URL.Query().Get("sparse")
	unique := req.URL.Query().Get("unique")

	if len(keys) == 0 {
		w.WriteHeader(400)
		io.WriteString(w, "Missing keys param")
		return
	}
	idx := mgo.Index{
		Key:        strings.Split(keys, ","),
		Background: true,
	}

	if sparse == "true" {
		idx.Sparse = true
	}

	if unique == "true" {
		idx.Unique = true
	}

	err = col.EnsureIndex(idx)
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	marshaled, _ := json.Marshal(idx)
	io.WriteString(w, string(marshaled))

}
