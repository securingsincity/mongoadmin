package main

import (
	"encoding/json"
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/maxwellhealth/mgo"
	// "github.com/maxwellhealth/mgo/bson"
	// auth "github.com/nabeken/negroni-auth"
	"io"
	"log"
	"net/http"
	"os"
	// "strconv"
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
	router.HandleFunc("/databases/{db}/collections", collections).Methods("GET")
	router.HandleFunc("/databases/{db}/collections/{col}", insert).Methods("POST")

	router.HandleFunc("/databases/{db}/collections/{col}/indexes", indexes).Methods("GET")
	router.HandleFunc("/databases/{db}/collections/{col}/find", find).Methods("POST")
	router.HandleFunc("/databases/{db}/collections/{col}/total", total).Methods("GET")
	router.HandleFunc("/databases/{db}/collections/{col}/newIndex", addIndex).Methods("POST")
	router.HandleFunc("/databases/{db}/collections/{col}/dropIndex", dropIndex).Methods("POST")
	router.HandleFunc("/databases/{db}/collections/{col}/findById/{mongoId}", findById).Methods("GET")
	router.HandleFunc("/databases/{db}/collections/{col}/update/{mongoId}", updateById).Methods("PUT", "POST")
	router.HandleFunc("/databases/{db}/collections/{col}/delete/{mongoId}", deleteById).Methods("DELETE")
	// http.Handle("/", router)
	n := negroni.Classic()
	// n.Use(negroni.HandlerFunc(auth.Basic(appConfig.AuthUsername, appConfig.AuthPassword)))
	n.Use(negroni.NewStatic(http.Dir("public")))
	n.UseHandler(router)
	n.Run(":" + conf.Port)
	// panic(http.ListenAndServe(":"+conf.Port, http.DefaultServeMux))
}
