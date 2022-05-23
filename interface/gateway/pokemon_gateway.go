package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/alexis-wizeline/ondemand-go-bootcamp/domain/model"
)

type PokemonApi struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokeApiResponse struct {
	Results []*PokemonApi `json:"results"`
}
type PokemonDetail struct {
	ID    uint64 `json:"id"`
	Types []TypesResponse
}

type TypesResponse struct {
	Slot int        `json:"slot"`
	Type TypeDetail `json:"type"`
}

type TypeDetail struct {
	Name string `json:"name"`
}

type pokemonGateway struct {
}

type PokemonGateway interface {
	GetAndAddPokemons(params url.Values) ([]*model.Pokemon, error)
}

func NewPokemonGateway() PokemonGateway {
	return pokemonGateway{}
}

func (pg pokemonGateway) GetAndAddPokemons(params url.Values) ([]*model.Pokemon, error) {
	apiUrl := "https://pokeapi.co/api/v2/pokemon?"
	if param := fmt.Sprintf("limit=%v&offset=%v", params.Get("limit"), params.Get("offset")); param != "&" {
		apiUrl += param
	}

	res, err := http.Get(apiUrl)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	pokeApiResult := new(PokeApiResponse)
	if err = json.Unmarshal(body, pokeApiResult); err != nil {
		return nil, err
	}

	pokemons, err := getPokemonDetails(pokeApiResult.Results)
	if err != nil {
		return nil, err
	}

	return pokemons, nil
}

func getPokemonDetails(pokeApiResponse []*PokemonApi) ([]*model.Pokemon, error) {
	var pokemons []*model.Pokemon

	for _, pokemon := range pokeApiResponse {
		res, err := http.Get(pokemon.Url)
		if err != nil {
			return nil, err
		}
		body, err := io.ReadAll(res.Body)
		_ = res.Body.Close()
		pokemonDetail := new(PokemonDetail)
		if err = json.Unmarshal(body, pokemonDetail); err != nil {
			return nil, err
		}

		pokemons = append(pokemons, &model.Pokemon{
			ID:   pokemonDetail.ID,
			Name: pokemon.Name,
			Type: pokemonDetail.Types[0].Type.Name,
		})
	}

	return pokemons, nil
}
