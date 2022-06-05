package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/alexis-wizeline/ondemand-go-bootcamp/interface/gateway"
	"github.com/alexis-wizeline/ondemand-go-bootcamp/usecase/repository"
)

type pokemonController struct {
	pokemonRepository repository.PokemonRepository
	pokemonGateway    gateway.PokemonGateway
}

type PokemonController interface {
	GetPokemons(c echo.Context) error
	GetPokemonsConcurrently(c echo.Context) error
	GetPokemonById(c echo.Context) error
	CallGateway(c echo.Context) error
}

func NewPokemonController(pr repository.PokemonRepository, pg gateway.PokemonGateway) PokemonController {
	return pokemonController{pokemonRepository: pr, pokemonGateway: pg}
}

func (p pokemonController) GetPokemonById(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id")
	}

	pokemon, err := p.pokemonRepository.GetPokemonById(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, pokemon)
}

func (p pokemonController) GetPokemons(c echo.Context) error {
	pokemons, err := p.pokemonRepository.GetPokemons()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, pokemons)
}

func (p pokemonController) GetPokemonsConcurrently(c echo.Context) error {
	t := c.QueryParam("type")
	items, err := strconv.Atoi(c.QueryParam("items"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	itemsPerWorker, err := strconv.Atoi(c.QueryParam("items_per_worker"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	pokemons, err := p.pokemonRepository.GetPokemonsConcurrently(t, items, itemsPerWorker)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, pokemons)
}

func (p pokemonController) CallGateway(c echo.Context) error {
	pokemons, err := p.pokemonGateway.GetAndAddPokemons(c.QueryParams())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if err = p.pokemonRepository.StorePokemons(pokemons); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, pokemons)
}
