package repository

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"

	"github.com/alexis-wizeline/ondemand-go-bootcamp/domain/model"
	repo "github.com/alexis-wizeline/ondemand-go-bootcamp/usecase/repository"
)

type pokemonRepository struct {
}

func NewPokemonRepository() repo.PokemonRepository {
	return pokemonRepository{}
}

var waitGroup = new(sync.WaitGroup)
var mu = new(sync.Mutex)

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

func (p pokemonRepository) GetPokemonsConcurrently(t string, items, itemsPerWorker int) ([]*model.Pokemon, error) {
	var pokemons []*model.Pokemon
	f, err := openAndGetFile()
	if err != nil {
		return nil, err
	}

	workers := items / itemsPerWorker
	pokemonChan := make(chan *model.Pokemon)
	reader := csv.NewReader(f)
	reader.Read()

	for i := 0; i < workers; i++ {
		waitGroup.Add(1)
		go worker(pokemonChan, waitGroup, mu, reader, t, itemsPerWorker)
	}

	go func(waitingGroup *sync.WaitGroup, pokemonChan chan *model.Pokemon) {
		defer close(pokemonChan)
		waitingGroup.Wait()
	}(waitGroup, pokemonChan)

	for pokemon := range pokemonChan {
		pokemons = append(pokemons, pokemon)
	}

	return pokemons, nil
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
		pokemon, err := rowToPokemon(row)
		if err != nil {
			return nil, err
		}

		pokemons = append(pokemons, pokemon)
	}
	return pokemons, nil
}

func worker(pokemonChan chan<- *model.Pokemon, wg *sync.WaitGroup, mu *sync.Mutex, reader *csv.Reader, t string, pokemonsPerWorker int) {
	var processedPokemon int
	defer fmt.Println("worker closed....")
	defer wg.Done()
	for {
		if processedPokemon == pokemonsPerWorker {
			return
		}
		mu.Lock()
		row, err := reader.Read()
		mu.Unlock()
		if err == io.EOF {
			return
		}

		pokemon, err := rowToPokemon(row)
		if err != nil {
			return
		}

		if shouldPokemonBeAdded(t, pokemon.ID) {
			pokemonChan <- pokemon
			processedPokemon++
		}
	}

}

func rowToPokemon(row []string) (*model.Pokemon, error) {
	pokemon := new(model.Pokemon)
	stringId := row[0]
	id, err := strconv.ParseUint(stringId, 10, 32)

	if err != nil {
		return nil, errors.New("invalid csv format")
	}

	pokemon.ID = id
	pokemon.Name = row[1]
	pokemon.Type = row[2]

	return pokemon, nil
}

func openAndGetCSVData() ([][]string, error) {
	f, err := openAndGetFile()
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

func openAndGetFile() (*os.File, error) {
	return os.OpenFile("./data/Pokemon.csv", os.O_RDONLY, 0755)
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

func shouldPokemonBeAdded(t string, id uint64) bool {
	switch t {
	case "odd":
		return id%2 == 0
	case "even":
		return id%2 != 0
	default:
		return true
	}
}
