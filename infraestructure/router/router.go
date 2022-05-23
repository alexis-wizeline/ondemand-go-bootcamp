package router

import (
	"github.com/alexis-wizeline/ondemand-go-bootcamp/interface/gateway"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/alexis-wizeline/ondemand-go-bootcamp/interface/controller"
	"github.com/alexis-wizeline/ondemand-go-bootcamp/interface/repository"
)

func NewRouter(api *echo.Echo) *echo.Echo {
	pc := controller.NewPokemonController(repository.NewPokemonRepository(), gateway.NewPokemonGateway())
	api.GET("/", test)

	pokemons := api.Group("pokemons")
	pokemons.GET("*", pc.GetPokemons)
	pokemons.GET("/:id", pc.GetPokemonById)
	pokemons.GET("/external", pc.CallGateway)

	return api
}

func test(c echo.Context) error {
	return c.String(http.StatusOK, "this is a test")
}
