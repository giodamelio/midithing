package cmd

import (
	"bufio"
	"github.com/AlecAivazis/survey/v2"
	"github.com/giodamelio/midithing/midi"
	"github.com/spf13/cobra"
	"log"
	"os"
	"sync"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Log midi messages",
	Long:  `Currently only logs ProgramChange, ControlChange and SysExt messages.`,
	Run: func(cmd *cobra.Command, args []string) {
		m := midi.New()
		defer m.Close()

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

		log.Println("Press enter to stop logging...")

		messagesChan := make(chan *midi.Message)
		var wg sync.WaitGroup

		// Print all the messages
		go func() {
			for message := range messagesChan {
				log.Printf("%+v\n", message)
			}
		}()

		// Watch for the enter key on stdin
		go func() {
			reader := bufio.NewReader(os.Stdin)
			_, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			wg.Done()
		}()

		m.Listen(messagesChan, &wg)
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
