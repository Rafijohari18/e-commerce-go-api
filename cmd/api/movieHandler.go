package main

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) getOneMovie(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		//app.logger.Print(errors.New("Invalid id parameter"))
		app.errorJSON(rw, err)
		return
	}

	// app.logger.Print("the id is:", id)

	// movie := models.Movie{
	// 	ID:          id,
	// 	Title:       "Some movie title",
	// 	Description: "some description",
	// 	Year:        2022,
	// 	ReleaseDate: time.Date(1990, 01, 01, 01, 0, 0, 0, time.Local),
	// 	Runtime:     112,
	// 	Rating:      5,
	// 	MPAARating:  "PG-13",
	// 	CreatedAt:   time.Now(),
	// 	UpdatedAt:   time.Now(),
	// }

	movie, err := app.models.DB.Get(id)
	err = app.writeJSON(rw, http.StatusOK, movie, "movie")
	if err != nil {
		app.errorJSON(rw, err)
		return
	}
}

func (app *application) getAllMovies(rw http.ResponseWriter, r *http.Request) {
	movies, err := app.models.DB.All()
	if err != nil {
		app.errorJSON(rw, err)
		return
	}
	err = app.writeJSON(rw, http.StatusOK, movies, "movies")
	if err != nil {
		app.errorJSON(rw, err)
		return
	}
}

func (app *application) getAllGenres(rw http.ResponseWriter, r *http.Request) {
	genres, err := app.models.DB.GetGenreAll()
	if err != nil {
		app.errorJSON(rw, err)
		return
	}
	err = app.writeJSON(rw, http.StatusOK, genres, "genres")
	if err != nil {
		app.errorJSON(rw, err)
		return
	}
}
