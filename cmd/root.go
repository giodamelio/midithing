package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"log"
	"os"

	"github.com/giodamelio/midithing/midi"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "midithing",
	Short: "Automate computer actions with a MIDI controller",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func selectInput(m *midi.Midi) {
	// Get the current MIDI inputs
	inputs := m.GetInputs()

	if len(inputs) < 1 {
		log.Fatalf("No MIDI inputs connected")
	} else if len(inputs) == 1 {
		log.Printf("Using the only MIDI input (%s)\n", inputs[0].String())
		m.SetInput(inputs[0])
	} else {
		// Get a list of the inputs names as strings
		inputsNames := make([]string, len(inputs))
		for i, in := range inputs {
			inputsNames[i] = in.String()
		}

		// Ask the user which one they would like to log
		question := &survey.Select{
			Message: "What input would you like to log?",
			Options: inputsNames,
		}
		var answer string
		err := survey.AskOne(question, &answer)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		// Get the actual input based on the name
		m.SetInputByName(answer)
	}
}

func init() {

}
