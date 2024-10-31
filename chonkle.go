package main

import (
	"chonkle/cli"
	"chonkle/utils"
	"fmt"
	"math"
	"os"
	"strconv"

	prompt "github.com/c-bata/go-prompt"
	env "github.com/caarlos0/env/v11"
	"github.com/manifoldco/promptui"
	"github.com/olekukonko/tablewriter"
)

const RandomGuess = "random"

type ChonkleOptions struct {
	RandomGuessCount   int  `env:"RANDOM_GUESS_COUNT"`
	MaxGuessCount      int  `env:"MAX_GUESS_COUNT"`
	GenMin             int  `env:"GENERATION_MIN"`
	GenMax             int  `env:"GENERATION_MAX"`
	HideHeightCol      bool `env:"HIDE_HEIGHT_COLUMN"`
	HideWeightCol      bool `env:"HIDE_WEIGHT_COLUMN"`
	HideGenCol         bool `env:"HIDE_GENERATION_COLUMN"`
	HideRemainingTypes bool `env:"HIDE_REMAINING_TYPES"`
	HideSecondTypeCol  bool `env:"HIDE_SECOND_TYPE_COLUMN"`
}

// Main Menu Options
const (
	SingleMode    = "Single Mode"
	PhlarcMode    = "Phlarc Special"
	UnlimitedMode = "∞ unlimited mode"
)

func MainMenu(pokemon []Pokemon) {
	var opts ChonkleOptions
	utils.CheckError(env.Parse(&opts))

	var err error
	pokemon, err = GetPokemonFromGenerationRange(opts.GenMin, opts.GenMax, pokemon)
	utils.CheckError(err)

	prompt := promptui.Select{
		Label: "Main Menu",
		Items: []string{SingleMode, PhlarcMode, UnlimitedMode, cli.Exit},
	}
	_, result, err := prompt.Run()
	utils.CheckError(err)

	gameCount := 0
	targetCount := 0

	switch result {
	case SingleMode:
		targetCount = 1
	case PhlarcMode:
		targetCount = 3 // should auto random the first guess
	case UnlimitedMode:
		targetCount = math.MaxInt
	default:
		targetCount = -1
	}

	for gameCount < targetCount {
		Chonkle(pokemon, opts)
		gameCount++
		cli.PressKeyToContinue()
	}
}

func Chonkle(pokemon []Pokemon, opts ChonkleOptions) {
	sugg := []prompt.Suggest{{Text: cli.Exit}, {Text: RandomGuess}}
	for i := range pokemon {
		sugg = append(sugg, prompt.Suggest{Text: pokemon[i].Name})
	}

	completer := func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterFuzzy(sugg, d.Text, true)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Gen", "Type 1", "Type 2", "Height", "Weight", ""})
	table.SetBorder(false)

	target := GetRandomPokemon(pokemon)
	var guesses []Pokemon
	guessCount := 0
	res := "lost"

	for i := 0; i < int(opts.RandomGuessCount); i++ {
		guess := GetRandomPokemon(pokemon)
		table.Append(BuildGuessResultRow(target, guess))
		guesses = append(guesses, guess)
		guessCount++
	}

	fmt.Print(cli.Clear)
	table.Render()

	for guessCount < opts.MaxGuessCount {
		var guess Pokemon
		for guess.Generation == 0 {
			PrintTypingHintTable(target, guesses)
			fmt.Println("Who's that pokemon?")
			t := prompt.Input("> ", completer)
			guess = ValidateUserGuess(t, pokemon)
		}

		guesses = append(guesses, guess)
		fmt.Print(cli.Clear)

		table.Append(BuildGuessResultRow(target, guess))
		table.Render()

		if guess == target {
			res = "won"
			break
		}

		guessCount++
	}

	fmt.Printf("You %s! The secret Pokémon was %s!\n", res, target.Name)
}

func ValidateUserGuess(guessStr string, pokemon []Pokemon) Pokemon {
	if guessStr == cli.Exit {
		os.Exit(0)
	}

	if guessStr == RandomGuess {
		return GetRandomPokemon(pokemon)
	}

	var guess Pokemon
	for i := range pokemon {
		if pokemon[i].Name == guessStr {
			guess = pokemon[i]
			break
		}
	}

	return guess
}

func PrintTypingHintTable(target Pokemon, guesses []Pokemon) {
	data := [][]string{
		{"Normal", "Fire", "Water", "Grass", "Electric", "Ice"},
		{"Fighting", "Poison", "Ground", "Flying", "Psychic", "Bug"},
		{"Rock", "Ghost", "Dark", "Dragon", "Steel", "Fairy", SingleType},
	}

	correctTypes := make(map[string]struct{})
	incorrectTypes := make(map[string]struct{})

	for _, g := range guesses {
		if g.Type1 == target.Type1 || g.Type1 == target.Type2 {
			correctTypes[g.Type1] = struct{}{}
		} else {
			incorrectTypes[g.Type1] = struct{}{}
		}

		if g.Type2 == target.Type2 || g.Type2 == target.Type1 {
			correctTypes[utils.CoalesceString(g.Type2, SingleType)] = struct{}{}
		} else {
			incorrectTypes[utils.CoalesceString(g.Type2, SingleType)] = struct{}{}
		}
	}

	correctType := tablewriter.FgGreenColor
	availableType := tablewriter.Bold
	incorrectType := tablewriter.FgRedColor

	table := tablewriter.NewWriter(os.Stdout)
	for _, row := range data {
		var rowColors []tablewriter.Colors
		for _, v := range row {
			if _, exists := correctTypes[v]; exists {
				rowColors = append(rowColors, []int{correctType})
			} else if _, exists := incorrectTypes[v]; exists {
				rowColors = append(rowColors, []int{incorrectType})
			} else {
				rowColors = append(rowColors, []int{availableType})
			}
		}
		table.Rich(row, rowColors)
	}

	table.Render()
}

func BuildGuessResultRow(target Pokemon, guess Pokemon) []string {
	genIndication := GetNumberComparisonResult(target.Generation, guess.Generation)
	heightIndication := GetNumberComparisonResult(target.Height, guess.Height)
	weightIndication := GetNumberComparisonResult(target.Weight, guess.Weight)

	var type1Indication string
	var type2Indication string

	switch {
	case guess.Type1 == target.Type1:
		type1Indication = cli.Correct
	case guess.Type1 == target.Type2:
		type1Indication = cli.ToRight
	default:
		type1Indication = cli.Incorrect
	}

	switch {
	case guess.Type2 == target.Type2:
		type2Indication = cli.Correct
	case guess.Type2 == target.Type1:
		type2Indication = cli.ToLeft
	default:
		type2Indication = cli.Incorrect
	}

	tableFmt := "%s %s"
	res := []string{
		fmt.Sprintf(tableFmt, genIndication, strconv.Itoa(guess.Generation)),
		fmt.Sprintf(tableFmt, type1Indication, guess.Type1),
		fmt.Sprintf(tableFmt, type2Indication, utils.CoalesceString(guess.Type2, NoneType)),
		fmt.Sprintf(tableFmt, heightIndication, strconv.FormatFloat(guess.Height, 'f', 2, 64)),
		fmt.Sprintf(tableFmt, weightIndication, strconv.FormatFloat(guess.Weight, 'f', 2, 64)),
		guess.Name,
	}

	return res
}

func GetNumberComparisonResult[T int | float64](l T, r T) string {
	switch {
	case l < r:
		return cli.Below
	case l > r:
		return cli.Above
	default:
		return cli.Correct
	}
}
