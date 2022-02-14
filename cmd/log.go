package cmd

import (
	"bufio"
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

		// Select the input
		selectInput(m)

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
