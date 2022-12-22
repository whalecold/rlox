use std::io::Write;
use std::{fs, io, process::exit};

mod scanner;

pub struct Lox {
    has_error: bool,
}

impl Lox {
    pub fn new() -> Lox {
        Lox { has_error: false }
    }

    pub fn error(&mut self, line: &str, message: &str) {
        self.report(line, "", message)
    }

    fn report(&mut self, line: &str, wh: &str, message: &str) {
        println!("[line:{:?}] Error {:?}:{:?}", line, wh, message);
        self.has_error = true;
    }

    fn run_file(&self, file: &str) {
        println!("rlox [script] {:?}", file);
        let context = fs::read_to_string(file).unwrap();

        println!("script file: {:?}", context);
        self.run(context)
    }

    fn run_prompt(&self) -> () {
        let stdin = io::stdin();
        loop {
            if self.has_error {
                exit(65);
            }
            print!(">");
            let _ = io::stdout().flush();
            let mut str_buf = String::new();

            let _ = match stdin.read_line(&mut str_buf) {
                Ok(size) => size,
                Err(e) => {
                    println!("error get stdin {:?}", e);
                    break;
                }
            };
            self.run(str_buf);
        }
    }

    fn run(&self, source: String) {
        let mut scan = scanner::Scanner::new_from(source);
        let tokens = scan.scan_tokens();
        for token in tokens {
            print!("token {:?}", token.to_string())
        }
    }

    pub fn start_up(&mut self, args: Vec<String>) {
        if args.len() == 2 {
            self.run_file(args.get(1).as_ref().unwrap());
        } else {
            self.run_prompt();
        }
    }
}
