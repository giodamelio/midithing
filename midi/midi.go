package midi

import (
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/rtmididrv"
	"log"
	"sync"
)

type Midi struct {
	driver midi.Driver
	in     midi.In
}

func New() *Midi {
	myMidi := &Midi{}

	drv, err := rtmididrv.New()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	myMidi.driver = drv

	return myMidi
}

func (m Midi) Close() {
	if m.in != nil {
		err := m.in.Close()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}
	err := m.driver.Close()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

// GetInputs retrieves a list of the currently connected inputs
func (m Midi) GetInputs() []midi.In {
	ins, err := m.driver.Ins()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return ins
}

func (m *Midi) SetInputByName(name string) {
	for _, input := range m.GetInputs() {
		if input.String() == name {
			m.in = input
			break
		}
	}
}

func (m *Midi) Listen(messageChan chan *Message, closeGroup *sync.WaitGroup) {
	closeGroup.Add(1)

	rd := reader.New(
		reader.NoLogger(),
		// Send every message to the channel
		reader.Each(func(pos *reader.Position, msg midi.Message) {
			message := &Message{
				description: msg.String(),
				raw:         msg.Raw(),
			}
			messageChan <- message
		}),
	)

	if m.in == nil {
		log.Fatalln("Error: no input set")
	}

	// Open the port
	err := m.in.Open()
	if err != nil {
		log.Fatalln("Error: no input set")
	}
	defer func(in midi.In) {
		err := in.Close()
		if err != nil {
			log.Fatalln("Error: no input set")
		}
	}(m.in)

	// Listen for MIDI message
	err = rd.ListenTo(m.in)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Wait until the waitgroup tells us we are done
	closeGroup.Wait()

	err = m.in.StopListening()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

type Message struct {
	description string
	raw         []byte
}
