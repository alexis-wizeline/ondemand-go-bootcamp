package repository

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"
)

var repository = NewPokemonRepository()

func init() {
	_, filename, _, _ := runtime.Caller(0)
	// The ".." may change depending on you folder structure
	dir := path.Join(path.Dir(filename), "../..")
	fmt.Println(dir)
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestPokemonRepository_GetPokemons(t *testing.T) {
	exp, _ := openAndGetCSVData()
	if res, err := repository.GetPokemons(); err != nil || len(res) != len(exp) {
		t.Fatalf(`not matching %v, %v, want "", error, %e`, len(res), len(exp), err)
	}
}

var tests = []struct {
	id   uint64
	name string
}{
	{
		id:   1,
		name: "Bulbasaur",
	}, {
		id:   75,
		name: "Graveler",
	},
}

func TestPokemonRepository_GetPokemonById(t *testing.T) {
	for _, test := range tests {
		if p, err := repository.GetPokemonById(test.id); err != nil || p.Name != test.name {
			t.Fatalf(`not matching %s, %s, want "", error, %e`, p.Name, test.name, err)
		}
	}
}
