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
	// "strconv"
	// "strings"
)

func insert(w http.ResponseWriter, req *http.Request) {

	sess, col, err := getCollection(req)
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
		return
	}

	defer sess.Close()
	r := bson.M{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&r)
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
		return
	}
	err = col.Insert(&r)
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
		return
	}
	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json")
	marshaled, _ := json.Marshal(r)
	io.WriteString(w, string(marshaled))

}

func updateById(w http.ResponseWriter, req *http.Request) {
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
		decoder := json.NewDecoder(req.Body)
		err = decoder.Decode(&r)
		if err != nil {
			w.WriteHeader(400)
			io.WriteString(w, err.Error())
			return
		}
		delete(r, "_id")
		idHex := bson.ObjectIdHex(id)
		err = col.UpdateId(idHex, &r)
		if err != nil {
			w.WriteHeader(400)
			io.WriteString(w, err.Error())
			return
		}
		w.WriteHeader(201)
		w.Header().Set("Content-Type", "application/json")
		marshaled, _ := json.Marshal(r)
		io.WriteString(w, string(marshaled))
	}

}
