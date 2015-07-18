package main

import (
	"encoding/json"
	// "errors"
	// "github.com/maxwellhealth/mgo"
	"github.com/gorilla/mux"
	"github.com/maxwellhealth/mgo/bson"
	// auth "github.com/nabeken/negroni-auth"
	"io"
	// "log"
	"net/http"
	// "os"
	"strconv"
	// "strings"
)

func collections(w http.ResponseWriter, req *http.Request) {

	sess, db, err := getDatabase(req)
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
		return
	}

	defer sess.Close()

	cols, err := db.CollectionNames()

	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	marshaled, _ := json.Marshal(cols)
	io.WriteString(w, string(marshaled))
}

func find(w http.ResponseWriter, req *http.Request) {
	limit := 50
	skip := 0

	sess, col, err := getCollection(req)
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
		return
	}

	defer sess.Close()
	// v := make(map[string]string)

	limitQuery := req.URL.Query().Get("limit")
	if limitQuery != "" {
		limit, _ = strconv.Atoi(limitQuery)
	}
	skipQuery := req.URL.Query().Get("skip")
	if skipQuery != "" {
		skip, _ = strconv.Atoi(skipQuery)
	}

	r := []bson.M{}
	err = col.Find(bson.M{}).Skip(skip).Limit(limit).All(&r)
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	marshaled, _ := json.Marshal(r)
	io.WriteString(w, string(marshaled))

}

func findById(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	if id, ok := vars["mongoId"]; ok {
		sess, col, err := getCollection(req)
		if err != nil {
			w.WriteHeader(400)
			io.WriteString(w, err.Error())
			return
		}

		defer sess.Close()

		r := bson.M{}
		idHex := bson.ObjectIdHex(id)
		err = col.Find(bson.D{{"_id", idHex}}).One(&r)
		if err != nil {
			w.WriteHeader(400)
			io.WriteString(w, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		marshaled, _ := json.Marshal(r)
		io.WriteString(w, string(marshaled))
	}
	// v := make(map[string]string)

}

func total(w http.ResponseWriter, req *http.Request) {

	sess, col, err := getCollection(req)
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
		return
	}

	defer sess.Close()
	v := make(map[string]string)
	count, err := col.Find(v).Count()
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
		return
	}
	io.WriteString(w, strconv.Itoa(count))

}
