use std::{env, process::exit};
mod lox;

#[macro_use]
extern crate strum;

fn main() {
    let args: Vec<String> = env::args().collect();
    println!("{:?}", args);
    if args.len() > 2 {
        println!("Usage: rlox [script]");
        exit(64);
    }

    let mut l = lox::Lox::new();
    l.start_up(args);
}
