package stone_swallow

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"golang.org/x/net/context"

	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

type runtimeEnv struct {
	NumCPU       int
	GOMAXPROCS   int
	NumGoroutine int
}

type requestParam struct {
	Host       string
	Method     string
	UrlHost    string
	Fragment   string
	Path       string
	Scheme     string
	Opaque     string
	RawQuery   string
	RemoteAddr string
	RequestURI string
	UserAgent  string
}

type testCookie struct {
	Domain string
}

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	switch r.URL.Path {
	default:
		http.Error(w, "not found.", http.StatusNotFound)
	case "/entity":
		listEntity(w, r, c)
	case "/kind":
		allKind(w, r, c)
	case "/sample":
		putHoge(w, r, c)
	case "/env":
		listEnvironment(w, r, c)
	case "/log":
		postLog(w, r, c)
	case "/param":
		getParam(w, r, c)
	case "/importBilling":
		importBilling(w, r, c)
	case "/queryBilling":
		listBilling4chart(w, r, c)
	case "/testcookie":
		handleTestCookie(w, r, c)
	case "/static":
		writeStaticFile(w, r, c)
	case "/":
		writeStaticFile(w, r, c)
	}
}

func listEnvironment(w http.ResponseWriter, r *http.Request, c context.Context) {
	fmt.Println(runtime.NumCPU())
	fmt.Println(runtime.GOMAXPROCS(0))
	fmt.Println(runtime.NumGoroutine())

	re := &runtimeEnv{
		NumCPU:       runtime.NumCPU(),
		GOMAXPROCS:   runtime.GOMAXPROCS((0)),
		NumGoroutine: runtime.NumGoroutine(),
	}

	json, err := json.Marshal(re)
	if err != nil {
		log.Errorf(c, "handler error: %#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(json)
	if err != nil {
		log.Errorf(c, "write response error: %#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func allKind(w http.ResponseWriter, r *http.Request, c context.Context) {
	t := datastore.NewQuery("__kind__").KeysOnly().Run(c)
	kinds := make([]string, 0)
	for {
		key, err := t.Next(nil)
		if err == datastore.Done {
			break // No further entities match the query.
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		kinds = append(kinds, key.StringID())
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(kinds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type entity struct {
	Key      *datastore.Key
	KeyValue string
	List     datastore.PropertyList
}

func listEntity(w http.ResponseWriter, r *http.Request, c context.Context) {
	kind := r.FormValue("kind")
	log.Infof(c, "kind=%s", kind)

	order := r.FormValue("order")
	log.Infof(c, "order=%s", order)

	limit := r.FormValue("limit")
	log.Infof(c, "limit=%s", limit)

	q := datastore.NewQuery(kind)
	if order != "" {
		q = q.Order(order)
	}
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			log.Errorf(c, "handler error: %#v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if l != 0 {
			log.Infof(c, "set limit : %d", l)
			q = q.Limit(l)
		}
	}

	var dst []datastore.PropertyList
	keys, err := q.GetAll(c, &dst)
	if err != nil {
		log.Errorf(c, "handler error: %#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]*entity, len(keys))
	for idx, key := range keys {
		resp[idx] = &entity{
			key, fmt.Sprintf("%v", key), dst[idx],
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getParam(w http.ResponseWriter, r *http.Request, c context.Context) {
	p := &requestParam{
		r.Host,
		r.Method,
		r.URL.Host,
		r.URL.Fragment,
		r.URL.Path,
		r.URL.Scheme,
		r.URL.Opaque,
		r.URL.RawQuery,
		r.RemoteAddr,
		r.RequestURI,
		r.UserAgent(),
	}

	json, err := json.Marshal(p)
	if err != nil {
		log.Errorf(c, "handler error: %#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(json)
	if err != nil {
		log.Errorf(c, "write response error: %#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func writeStaticFile(w http.ResponseWriter, r *http.Request, c context.Context) {
	log.Infof(c, "%q\n", strings.Split(r.URL.Host, "."))
	w.Header().Set("Cache-Control:public", "max-age=120")
	sd := strings.Split(r.URL.Host, ".")[0]
	if sd != "fuga" && sd != "hoge" {
		c, err := r.Cookie("testdomain")
		if err == nil {
			sd = c.Value
		} else {
			http.Redirect(w, r, "/testcookie", 302)
			return
		}
	}
	http.ServeFile(w, r, sd)
}

func handleTestCookie(w http.ResponseWriter, r *http.Request, c context.Context) {
	switch r.Method {
	default:
		http.Error(w, "not support method.", http.StatusMethodNotAllowed)
	case "POST":
		postTestCookie(w, r, c)
	case "GET":
		getTestCookie(w, r, c)
	}
}

func postTestCookie(w http.ResponseWriter, r *http.Request, c context.Context) {
	defer r.Body.Close()
	var tc testCookie

	err := json.NewDecoder(r.Body).Decode(&tc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cookie := http.Cookie{Name: "testdomain", Value: tc.Domain}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
}

func getTestCookie(w http.ResponseWriter, r *http.Request, c context.Context) {
	tc, err := r.Cookie("testdomain")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json, err := json.Marshal(tc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)
}
