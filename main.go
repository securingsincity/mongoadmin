package main

import (
	"encoding/json"
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/maxwellhealth/mgo"
	"github.com/maxwellhealth/mgo/bson"
	auth "github.com/nabeken/negroni-auth"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type DatabaseConfig struct {
	Label            string `json:"label"`
	Database         string `json:"-"`
	ConnectionString string `json:"-"`
}

type Config struct {
	DB           []*DatabaseConfig
	Port         string
	AuthUsername string
	AuthPassword string
}

var appConfig *Config

func parseFile(file string) (*Config, error) {

	conf := &Config{}

	_, err := toml.DecodeFile(file, conf)

	return conf, err
}

func getDatabase(req *http.Request) (*mgo.Session, *mgo.Database, error) {
	vars := mux.Vars(req)

	if db, ok := vars["db"]; ok {
		for _, d := range appConfig.DB {
			if d.Label == db {
				sess, err := mgo.Dial(d.ConnectionString)
				if err != nil {
					return &mgo.Session{}, &mgo.Database{}, err
				}

				return sess, sess.DB(d.Database), nil
			}
		}
		return &mgo.Session{}, &mgo.Database{}, errors.New("Could not find DB by name " + db)
	}

	return &mgo.Session{}, &mgo.Database{}, errors.New("Missing db parameter")
}

func getCollection(req *http.Request) (*mgo.Session, *mgo.Collection, error) {
	vars := mux.Vars(req)

	sess, db, err := getDatabase(req)
	if err != nil {
		return sess, &mgo.Collection{}, err
	}

	if colname, ok := vars["col"]; ok {
		return sess, db.C(colname), nil
	}

	sess.Close()

	return sess, &mgo.Collection{}, errors.New("Missing collection parameter")
}

func main() {
	log.Println("booting")
	args := os.Args
	var conf *Config
	var err error
	// Need at least one (command is included)...
	if len(args) == 1 {
		log.Panic("Please specify toml config file")
	}
	if len(args) >= 2 && strings.HasSuffix(args[1], ".toml") {
		log.Println("has toml arg")
		conf, err = parseFile(args[1])
		if err != nil {
			log.Panic(err.Error())
		}

	}
	// authenticator := auth.NewBasicAuthenticator("MongoAdmin", basicAuth)

	appConfig = conf

	router := mux.NewRouter()

	router.HandleFunc("/databases", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		marshaled, _ := json.Marshal(conf.DB)

		io.WriteString(w, string(marshaled))
	}).Methods("GET")

	router.HandleFunc("/databases/{db}/collections", func(w http.ResponseWriter, req *http.Request) {

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
	}).Methods("GET")

	router.HandleFunc("/databases/{db}/collections/{col}/indexes", func(w http.ResponseWriter, req *http.Request) {

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

	}).Methods("GET")
	router.HandleFunc("/databases/{db}/collections/{col}/find", func(w http.ResponseWriter, req *http.Request) {
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

	}).Methods("GET")
	router.HandleFunc("/databases/{db}/collections/{col}", func(w http.ResponseWriter, req *http.Request) {

		// sess, col, err := getCollection(req)
		// if err != nil {
		// 	w.WriteHeader(400)
		// 	io.WriteString(w, err.Error())
		// 	return
		// }

		// defer sess.Close()
		// v := make(map[string]string)
		r := bson.M{}
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&r)
		if err != nil {
			w.WriteHeader(400)
			io.WriteString(w, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		marshaled, _ := json.Marshal(r)
		io.WriteString(w, string(marshaled))

	}).Methods("POST")
	router.HandleFunc("/databases/{db}/collections/{col}/findById/{mongoId}", func(w http.ResponseWriter, req *http.Request) {
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

	}).Methods("GET")
	router.HandleFunc("/databases/{db}/collections/{col}/total", func(w http.ResponseWriter, req *http.Request) {

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

	}).Methods("GET")

	router.HandleFunc("/databases/{db}/collections/{col}/newIndex", func(w http.ResponseWriter, req *http.Request) {

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

	}).Methods("POST")

	router.HandleFunc("/databases/{db}/collections/{col}/dropIndex", func(w http.ResponseWriter, req *http.Request) {

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

	}).Methods("POST")

	// http.Handle("/", router)
	n := negroni.Classic()
	n.Use(negroni.NewStatic(http.Dir("public")))
	n.Use(negroni.HandlerFunc(auth.Basic(appConfig.AuthUsername, appConfig.AuthPassword)))
	n.UseHandler(router)
	n.Run(":" + conf.Port)
	// panic(http.ListenAndServe(":"+conf.Port, http.DefaultServeMux))
}
