package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gochain-io/explorer/api/models"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"

	"cloud.google.com/go/datastore"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"golang.org/x/oauth2/google"
)

var ds *datastore.Client

func main() {
	// Use oauth2.NoContext if there isn't a good context to pass in.
	ctx := context.Background()

	creds, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/datastore")
	if err != nil {
		log.Fatal().Err(err).Msg("oops")
	}
	ds, err = datastore.NewClient(ctx, creds.ProjectID)
	if err != nil {
		log.Fatal().Err(err).Msg("oops")
	}

	query := datastore.NewQuery("items")
	// Filter("type =", "foo")
	// Filter("Priority >=", 4).
	// Order("-Priority")

	it := ds.Run(ctx, query)
	for {
		var task Item
		_, err := it.Next(&task)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal().Err(err).Msg("Error fetching next task")
		}
		fmt.Printf("Task %v, Priority %v\n", task.Name, task.Typ)
	}

	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	cors2 := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors2.Handler)
	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Route("/blocks", func(r chi.Router) {
		r.Get("/", listBlocks)
		// r.Get("/{num}", getBlock)
	})
	http.ListenAndServe(":8080", r)
}

func listBlocks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bl := &models.BlockList{
		Blocks: []*models.Block{},
	}
	query := datastore.NewQuery("blocks")
	// Filter("type =", "foo")
	// Filter("Priority >=", 4).
	// Order("-Priority")

	it := ds.Run(ctx, query)
	for {
		var block models.Block
		_, err := it.Next(&block)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Error().Err(err).Msg("Error fetching next block")
			break
		}
		fmt.Printf("Block %v, Hash %v\n", block.Number, block.BlockHash)
		bl.Blocks = append(bl.Blocks, &block)
	}
	// w.Write([]byte(fmt.Sprintf("title:%s", article.Title)))
	writeJSON(w, http.StatusOK, bl)
}

// func getArticle(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// articleID := chi.URLParam(r, "articleID")
// 	article, ok := ctx.Value("article").(*Article)
// 	if !ok {
// 	  http.Error(w, http.StatusText(422), 422)
// 	  return
// 	}
// 	w.Write([]byte(fmt.Sprintf("title:%s", article.Title)))
//   }

type Item struct {
	Name string
	Typ  string `json:"type" datastore:"type"`
	Bnum int    `datastore:"bnum"`
}
