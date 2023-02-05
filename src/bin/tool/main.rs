use std::fs;
use std::fs::File;
use std::io::prelude::*;
use std::{env, process::exit, str::FromStr};

fn main() {
    let args: Vec<String> = env::args().collect();
    println!("{:?}", args);
    if args.len() != 2 {
        println!("Usage: generate_ast <output directory>");
        exit(64);
    }
    define_ast(
        args.get(1).as_ref().unwrap(),
        "expr",
        vec![
            "Binary   : Expr left, Token operator, Expr right",
            "Grouping : Expr expression",
            "Literal  : Object value",
            "Unary    : Token operator, Expr right",
        ],
    )
}

fn define_ast(file: &str, base_name: &str, types: Vec<&str>) {
    let mut path: String = String::from_str(file).unwrap();
    path.push_str("/");
    path.push_str(base_name);
    path.push_str(".rs");

    let mut file = fs::File::create(path.as_str()).unwrap();

    let mut out = String::from_str("pub trait ").unwrap();
    out.push_str(base_name);
    out.push_str(" {}\n");

    file.write_all(out.as_bytes()).unwrap();

    for t in types {
        let mut split = t.split(":");
        // TODD validate the result
        let class_name = split.next().unwrap().trim();
        let fields = split.next().unwrap().trim();
        define_type(&mut file, base_name, class_name, fields);
    }
}

fn define_type(file: &mut File, base_name: &str, class_name: &str, filed_list: &str) {
    let mut context = String::from_str("pub struct").unwrap();

    context.push_str(class_name);

    file.write_all(context.as_bytes()).unwrap();
}
