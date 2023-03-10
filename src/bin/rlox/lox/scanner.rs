use phf::phf_map;
use std::marker::Copy;
use std::{ops::Index, str::FromStr};

#[derive(Clone, Copy, EnumString, Display, Debug, PartialEq)]
pub enum TokenType {
    // Single-character tokens.
    LeftParen,
    RightParen,
    LeftBrace,
    RightBrace,
    Comma,
    Dot,
    Minus,
    Plus,
    Semicolon,
    Slash,
    Star,

    // One or two character tokens.
    Bang,
    BangEqual,
    Equal,
    EqualEqual,
    Greater,
    GreaterEqual,
    Less,
    LessEqual,

    // Literals.
    Identifier,
    Strings,
    Number,

    // Keywords.
    And,
    Class,
    Else,
    False,
    Fun,
    For,
    If,
    Nil,
    Or,
    Print,
    Return,
    Super,
    This,
    True,
    Var,
    While,

    Eof,
}

static KEYWORDS: phf::Map<&str, TokenType> = phf_map! {
    "and" => TokenType::And,
    "class"=> TokenType::Class,
    "else" => TokenType::Else,
    "false" => TokenType::False,
    "for" => TokenType::For,
    "fun" => TokenType::Fun,
    "if" => TokenType::If,
    "nil" => TokenType::Nil,
    "or" => TokenType::Or,
    "print" => TokenType::Print,
    "return" => TokenType::Return,
    "super" => TokenType::Super,
    "this" => TokenType::This,
    "true" => TokenType::True,
    "var" => TokenType::Var,
    "while" => TokenType::While
};

#[derive(Clone, EnumString, Display, Debug, PartialEq)]
pub enum Literal {
    String(String),
    Number(u64),
    Boolean(bool),
    Nil,
}

pub struct Scanner {
    source: String,
    tokens: Vec<Token>,
    start: usize,
    current: usize,
    line: usize,
}

impl Scanner {
    pub fn new_from(source: String) -> Scanner {
        Scanner {
            source,
            tokens: Vec::new(),
            start: 0,
            current: 0,
            line: 1,
        }
    }

    fn is_digital(c: char) -> bool {
        c >= '0' && c <= '9'
    }

    fn is_alpha(c: char) -> bool {
        (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
    }

    fn is_alpha_numeric(c: char) -> bool {
        Scanner::is_alpha(c) || Scanner::is_alpha(c)
    }

    pub fn scan_tokens(&mut self) -> &Vec<Token> {
        while !self.is_at_end() {
            self.start = self.current;
            self.scan_token();
        }
        self.tokens.push(Token::new(
            TokenType::Eof,
            String::new(),
            self.line as i32,
            Literal::Nil,
        ));
        &self.tokens
    }

    fn scan_token(&mut self) {
        let ch = self.advance();
        match ch {
            '(' => self.add_to_token(TokenType::LeftParen),
            ')' => self.add_to_token(TokenType::RightParen),
            '{' => self.add_to_token(TokenType::LeftBrace),
            '}' => self.add_to_token(TokenType::RightBrace),
            ',' => self.add_to_token(TokenType::Comma),
            '.' => self.add_to_token(TokenType::Dot),
            '-' => self.add_to_token(TokenType::Minus),
            '+' => self.add_to_token(TokenType::Plus),
            ';' => self.add_to_token(TokenType::Semicolon),
            '*' => self.add_to_token(TokenType::Star),
            '!' => {
                let token = match self.match_with_char('=') {
                    true => TokenType::BangEqual,
                    _ => TokenType::Bang,
                };
                self.add_to_token(token);
            }
            '=' => {
                let token = match self.match_with_char('=') {
                    true => TokenType::EqualEqual,
                    _ => TokenType::Equal,
                };
                self.add_to_token(token);
            }
            '<' => {
                let token = match self.match_with_char('=') {
                    true => TokenType::LessEqual,
                    _ => TokenType::Less,
                };
                self.add_to_token(token);
            }
            '>' => {
                let token = match self.match_with_char('=') {
                    true => TokenType::GreaterEqual,
                    _ => TokenType::Greater,
                };
                self.add_to_token(token);
            }
            '/' => {
                if self.match_with_char('/') {
                    while self.peek() != '\n' && !self.is_at_end() {
                        self.advance();
                    }
                } else {
                    self.add_to_token(TokenType::Slash);
                }
            }
            // do nothing
            ' ' | '\r' | '\t' => (),
            '\n' => self.line += 1,
            '"' => self.string(),
            _ => {
                if Scanner::is_digital(ch) {
                    self.number();
                } else if Scanner::is_alpha(ch) {
                    self.identifier();
                } else {
                    // TODO handle error
                }
            }
        }
    }

    fn identifier(&mut self) {
        while Scanner::is_alpha_numeric(self.peek()) {
            self.advance();
        }

        self.add_to_token(
            match KEYWORDS.get(self.source.index(self.start..self.current)) {
                None => TokenType::Identifier,
                Some(val) => *val,
            },
        )
    }

    fn number(&mut self) {
        while Scanner::is_digital(self.peek()) {
            self.advance();
        }
        // Look for a fractional part.
        if self.peek() == '.' && Scanner::is_digital(self.peek_next()) {
            self.advance();
            while Scanner::is_digital(self.peek()) {
                self.advance();
            }
        }
        let val = self.source.index(self.start..self.current);
        self.add_to_token_raw(
            TokenType::Number,
            Literal::Number(val.parse::<u64>().unwrap()),
        )
    }

    fn string(&mut self) {
        while self.peek() != '"' && !self.is_at_end() {
            if self.peek() == '\n' {
                self.line += 1
            }
            self.advance();
        }
        if self.is_at_end() {
            // TODO add log
            return;
        }

        // The closing ".
        self.advance();

        let val = String::from_str(self.source.index(self.start + 1..self.current - 1)).unwrap();
        // Trim the surrounding quotes.
        self.add_to_token_raw(TokenType::Strings, Literal::String(val));
    }

    fn match_with_char(&mut self, expected: char) -> bool {
        if self.is_at_end() || !self.source.chars().nth(self.current).unwrap().eq(&expected) {
            false
        } else {
            self.current += 1;
            true
        }
    }
    fn is_at_end(&self) -> bool {
        return self.current >= self.source.len();
    }

    fn advance(&mut self) -> char {
        let ret = self.source.chars().nth(self.current).unwrap();
        self.current += 1;
        ret
    }

    fn add_to_token(&mut self, kind: TokenType) {
        self.add_to_token_raw(kind, Literal::Nil)
    }

    fn add_to_token_raw(&mut self, kind: TokenType, literal: Literal) {
        let text = self.source.index(self.start..self.current);
        self.tokens.push(Token::new(
            kind,
            String::from_str(text).unwrap(),
            self.line as i32,
            literal,
        ));
    }

    fn peek_next(&self) -> char {
        if self.current + 1 >= self.source.len() {
            '\0'
        } else {
            self.source.chars().nth(self.current + 1).unwrap()
        }
    }
    fn peek(&self) -> char {
        if self.is_at_end() {
            '\0'
        } else {
            self.source.chars().nth(self.current).unwrap()
        }
    }
}

#[derive(Debug)]
pub struct Token {
    token_type: TokenType,
    lexeme: String,
    line: i32,
    literal: Literal,
}

impl Token {
    pub fn new(token_type: TokenType, lexeme: String, line: i32, literal: Literal) -> Token {
        Token {
            token_type,
            lexeme,
            line,
            literal,
        }
    }
    pub fn to_string(&self) -> String {
        let mut ret = self.token_type.to_string();
        ret.push_str("/");
        ret.push_str(self.lexeme.as_str());
        ret.push_str("/");
        match self.literal {
            Literal::String(ref val) => ret.push_str(val.as_str()),
            Literal::Nil => ret.push_str("nil"),
            Literal::Number(num) => ret.push_str(num.to_string().as_str()),
            Literal::Boolean(b) => {
                if b == true {
                    ret.push_str("true")
                } else {
                    ret.push_str("false")
                }
            }
        };
        ret
    }
}
