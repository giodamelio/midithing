package cmd

import (
	"fmt"
	"sync"

	"gitlab.com/gomidi/midi"
	. "gitlab.com/gomidi/midi/midimessage/channel" // (Channel Messages)
	"gitlab.com/gomidi/midi/midimessage/sysex"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/rtmididrv"

	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Log midi messages",
	Long:  `Currently only logs ProgramChange, ControlChange and SysExt messages.`,
	Run: func(cmd *cobra.Command, args []string) {
		drv, err := rtmididrv.New()
		if err != nil {
			panic(err.Error())
		}

		// make sure to close the driver at the end
		defer func(drv *rtmididrv.Driver) {
			err := drv.Close()
			if err != nil {
				panic(err.Error())
			}
		}(drv)

		ins, err := drv.Ins()
		if err != nil {
			panic(err.Error())
		}

		var in midi.In
		for _, i := range ins {
			if i.String() == "WORLDE easy CTRL 0" {
				in = i
				break
			}
		}
		fmt.Printf("opening MIDI Port %v\n", in)
		err = in.Open()
		if err != nil {
			panic(err.Error())
		}
		// Make sure to close the midi input
		defer func(in midi.In) {
			err := in.Close()
			if err != nil {
				panic(err.Error())
			}
		}(in)

		var wg sync.WaitGroup

		wg.Add(1)
		rd := reader.New(
			reader.NoLogger(),
			// print every message
			reader.Each(func(pos *reader.Position, msg midi.Message) {
				// inspect
				//fmt.Printf("%s: % 02X\n", msg.String(), msg.Raw())

				// Cast to the various types we are interested in
				switch message := msg.(type) {
				case ControlChange:
					// If it is a specific button, set the WaitGroup, so we can exit
					if message.Controller() == 67 && message.Value() == 0 {
						wg.Done()
					}

					fmt.Printf("ControlChange: %v, controller: %v, value: %v\n", message.String(), message.Controller(), message.Value())
				case ProgramChange:
					fmt.Printf("ProgramChange: %v, program: %v\n", message.String(), message.Program())
				case sysex.SysEx:
					fmt.Printf("Sysex: %v, raw: % 02X, data: % 02X\n", message.String(), message.Raw(), message.Data())
				default:
					fmt.Printf("Unknown type: %v\n", message.String())
				}
			}),
		)

		// listen for MIDI
		err = rd.ListenTo(in)
		if err != nil {
			panic(err.Error())
		}

		wg.Wait()

		err = in.StopListening()
		if err != nil {
			panic(err.Error())
		}

		fmt.Printf("closing MIDI Port %v\n", in)
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
