package repository

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/alexis-wizeline/ondemand-go-bootcamp/domain/model"
	repo "github.com/alexis-wizeline/ondemand-go-bootcamp/usecase/repository"
)

type pokemonRepository struct {
}

func NewPokemonRepository() repo.PokemonRepository {
	return pokemonRepository{}
}

func (p pokemonRepository) GetPokemons() ([]*model.Pokemon, error) {
	pokemons, err := getAllPokemons()

	if err != nil {
		return pokemons, err
	}

	return pokemons, nil
}

func (p pokemonRepository) GetPokemonById(id uint64) (*model.Pokemon, error) {
	pokemons, err := getAllPokemons()

	if err != nil {
		return nil, err
	}

	for _, pokemon := range pokemons {
		if id == pokemon.ID {
			return pokemon, nil
		}
	}

	return nil, errors.New("Pokemon Not Found")
}

func (p pokemonRepository) StorePokemons(pokemons []*model.Pokemon) error {
	existingIds, err := getPokemonIdMap()
	if err != nil {
		return err
	}

	writer, err := getCsvWriter()
	defer writer.Flush()
	if err != nil {
		return err
	}
	for _, pokemon := range pokemons {
		if ok := existingIds[pokemon.ID]; !ok {
			id := fmt.Sprintf("%v", pokemon.ID)
			if err = writer.Write([]string{id, pokemon.Name, pokemon.Type}); err != nil {
				return err
			}
		}
	}

	return nil
}

func getAllPokemons() ([]*model.Pokemon, error) {
	rows, err := openAndGetCSVData()

	if err != nil {
		return nil, err
	}

	pokemons, err := transformRowsToPokemons(rows)

	if err != nil {
		return nil, err
	}

	return pokemons, nil

}

func transformRowsToPokemons(rows [][]string) ([]*model.Pokemon, error) {
	var pokemons []*model.Pokemon
	for _, row := range rows {
		pokemon := new(model.Pokemon)
		stringId := row[0]
		id, err := strconv.ParseUint(stringId, 10, 32)

		if err != nil {
			return nil, errors.New("invalid csv format")
		}

		pokemon.ID = id
		pokemon.Name = row[1]
		pokemon.Type = row[2]

		pokemons = append(pokemons, pokemon)

	}

	return pokemons, nil
}

func openAndGetCSVData() ([][]string, error) {
	f, err := os.OpenFile("./data/Pokemon.csv", os.O_RDONLY, 0755)

	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(f)
	_, _ = reader.Read()

	rows, err := reader.ReadAll()

	if err != nil {
		return nil, err
	}

	return rows, nil
}

func getCsvWriter() (*csv.Writer, error) {
	f, err := os.OpenFile("./data/Pokemon.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return csv.NewWriter(f), nil
}

func getPokemonIdMap() (map[uint64]bool, error) {
	pokemons, err := getAllPokemons()
	result := make(map[uint64]bool)
	if err != nil {
		return nil, err
	}

	for _, pokemon := range pokemons {
		result[pokemon.ID] = true
	}

	return result, nil
}
