use std::error::Error;
use std::io::{stdin, stdout, Write};

use midir::{MidiInput, MidiInputPort, MidiInputPorts};

pub fn choose_input<'a>(
    midi_in: &MidiInput,
    in_ports: &'a MidiInputPorts,
) -> Result<&'a MidiInputPort, Box<dyn Error>> {
    // Get an input port (read from console if multiple are available)
    match in_ports.len() {
        0 => return Err("no input port found".into()),
        1 => {
            println!(
                "Choosing the only available input port: {}",
                midi_in.port_name(&in_ports[0]).unwrap()
            );
            Ok(&in_ports[0])
        }
        _ => {
            println!("\nAvailable input ports:");
            for (i, p) in in_ports.iter().enumerate() {
                println!("{}: {}", i, midi_in.port_name(p).unwrap());
            }
            print!("Please select input port: ");
            stdout().flush()?;
            let mut input = String::new();
            stdin().read_line(&mut input)?;
            Ok(in_ports
                .get(input.trim().parse::<usize>()?)
                .ok_or("invalid input port selected")?)
        }
    }
}
