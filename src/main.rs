use clap::{AppSettings, Parser, Subcommand};

mod cmd;
mod common;

#[derive(Parser, Debug)]
#[clap(author, version, about, long_about = None)]
#[clap(global_setting(AppSettings::PropagateVersion))]
#[clap(global_setting(AppSettings::UseLongFormatForHelpSubcommand))]
struct Cli {
    #[clap(subcommand)]
    command: Commands,
}

#[derive(Subcommand, Debug)]
enum Commands {
    Log {},
}

fn main() {
    let cli = Cli::parse();

    match &cli.command {
        Commands::Log {} => match cmd::log::run() {
            Ok(_) => (),
            Err(err) => println!("Error: {}", err),
        },
    }
}
