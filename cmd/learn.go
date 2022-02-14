package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/giodamelio/midithing/midi"
	"github.com/spf13/cobra"
	"log"
	"path"
	"sync"
	"time"
)

// learnCmd represents the learn command
var learnCmd = &cobra.Command{
	Use:        "learn [file name]",
	Short:      "Allows you to easily create a mapping file of MIDI inputs to names",
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"file name"},
	Run: func(cmd *cobra.Command, args []string) {
		outputFilename := args[0]

		// If the output file doesn't have the toml extension, add one
		if path.Ext(outputFilename) != ".toml" {
			outputFilename = outputFilename + ".toml"
		}

		// Open a midi and select the input
		m := midi.New()
		defer m.Close()
		selectInput(m)

		log.Printf("Do something on %s\n", m.Name())

		// Collect inputs for a second
		wg := oneSecondWaitGroup()
		messages := m.CollectMessagesUntil(wg)

		// Analyze inputs and try to guess the prefix, postfix and valueIndex from a set of messages
		prefix, postfix, valueIndex, err := guessPattern(messages)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		guessedInputType, err := guessInputType(messages, valueIndex)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		// Ask if the type is correct
		inputTypes := []string{string(midi.InputUnknown), midi.InputButton, midi.InputSlider, midi.InputSelector}
		inputTypeQuestion := &survey.Select{
			Message: "What type of input is this?",
			Options: inputTypes,
			Default: string(guessedInputType),
		}
		var inputTypeAnswer string
		err = survey.AskOne(inputTypeQuestion, &inputTypeAnswer)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		// Ask for a name
		nameQuestion := &survey.Input{
			Message: "What is the name of this input?",
		}
		var nameAnswer string
		err = survey.AskOne(nameQuestion, &nameAnswer)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		input := &midi.Input{
			Name:         nameAnswer,
			Type:         midi.InputTypeParse(inputTypeAnswer),
			PrefixBytes:  prefix,
			PostFixBytes: postfix,
		}

		log.Println(input)
	},
}

func oneSecondWaitGroup() *sync.WaitGroup {
	var wg sync.WaitGroup

	go func() {
		time.Sleep(1 * time.Second)
		wg.Done()
	}()

	return &wg
}

// Find the common prefix and postfix bytes in a slice of messages
func guessPattern(messages []midi.Message) ([]byte, []byte, int, error) {
	if len(messages) < 2 {
		return nil, nil, -1, fmt.Errorf("more then %d messages needed to guess the pattern", len(messages))
	}

	// Get the matching prefix and postfix
	firstMessage, secondMessage := messages[0], messages[1]
	var prefix []byte
	prefixDone := false
	var postfix []byte
	postfixDone := false
	var valueIndex int
	for i, b := range firstMessage.Raw {
		if secondMessage.Raw[i] == b && !prefixDone {
			prefix = append(prefix, b)
		} else {
			prefixDone = true
			valueIndex = i
		}

		// Get the inverse of the index
		i = len(firstMessage.Raw) - i - 1
		if secondMessage.Raw[i] == firstMessage.Raw[i] && !postfixDone {
			postfix = append(postfix, firstMessage.Raw[i])
		} else {
			postfixDone = true
		}
	}

	// Ensure if all the messages match the prefix
	for _, m := range messages {
		for i, b := range prefix {
			if b != m.Raw[i] {
				return nil, nil, -1, fmt.Errorf("message %v does not match prefix %v", m.Raw, prefix)
			}
		}

		for i, b := range postfix {
			if b != m.Raw[len(m.Raw)-i-1] {
				return nil, nil, -1, fmt.Errorf("message %v does not match postfix %v", m.Raw, prefix)
			}
		}
	}

	return prefix, postfix, valueIndex, nil
}

// Try to guess what input type it is based on a set of messages
func guessInputType(messages []midi.Message, valueIndex int) (midi.InputType, error) {
	// If all the values are 0 and 127 it is probably a button
	all0Or127 := true
	for _, m := range messages {
		if m.Raw[valueIndex] != 0 && m.Raw[valueIndex] != 127 {
			all0Or127 = false
		}
	}
	if all0Or127 {
		return midi.InputButton, nil
	}

	// Otherwise, it is probably a slider
	return midi.InputSlider, nil
}

func init() {
	rootCmd.AddCommand(learnCmd)
}
