package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
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
		r.Get("/{num}", getBlock)
	})
	http.ListenAndServe(":8080", r)
}

func listBlocks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bl := &models.BlockList{
		Blocks: []*models.Block{},
	}
	query := datastore.NewQuery("Blocks")
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
			// todo: sendError()
			break
		}
		fmt.Printf("Block %v, Hash %v\n", block.Number, block.BlockHash)
		bl.Blocks = append(bl.Blocks, &block)
	}
	writeJSON(w, http.StatusOK, bl)
}

func getBlock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bnumS := chi.URLParam(r, "num")
	bnum, err := strconv.Atoi(bnumS)
	if err != nil {
		log.Error().Err(err).Str("bnumS", bnumS).Msg("Error converting bnumS to num")
		// todo: sendError()
		return
	}
	log.Info().Int("bnum", bnum).Msg("looking up block")
	var block models.Block
	query := datastore.NewQuery("Blocks").
		Filter("num =", bnum)
	// Filter("Priority >=", 4).
	// Order("-Priority")

	it := ds.Run(ctx, query)
	for {
		log.Info().Msg("iter")
		_, err := it.Next(&block)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Error().Err(err).Msg("Error fetching block")
			// todo: sendError()
			break
		}
		fmt.Printf("Block %v, Hash %v\n", block.Number, block.BlockHash)
	}
	writeJSON(w, http.StatusOK, block)
}

type Item struct {
	Name string
	Typ  string `json:"type" datastore:"type"`
	Bnum int    `datastore:"bnum"`
}
