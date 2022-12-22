use std::collections::HashMap;
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

pub struct Scanner<'a> {
    source: String,
    tokens: Vec<Token>,
    hash_map: HashMap<&'a str, TokenType>,
    start: usize,
    current: usize,
    line: usize,
}

impl Scanner<'_> {
    pub fn new_from(source: String) -> Scanner<'static> {
        let mut hm = HashMap::new();
        hm.insert("and", TokenType::And);
        hm.insert("class", TokenType::Class);
        hm.insert("else", TokenType::Else);
        hm.insert("false", TokenType::False);
        hm.insert("for", TokenType::For);
        hm.insert("fun", TokenType::Fun);
        hm.insert("if", TokenType::If);
        hm.insert("nil", TokenType::Nil);
        hm.insert("or", TokenType::Or);
        hm.insert("print", TokenType::Print);
        hm.insert("return", TokenType::Return);
        hm.insert("super", TokenType::Super);
        hm.insert("this", TokenType::This);
        hm.insert("true", TokenType::True);
        hm.insert("var", TokenType::Var);
        hm.insert("while", TokenType::While);
        Scanner {
            source,
            tokens: Vec::new(),
            hash_map: hm,
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
            String::new(),
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
            match self
                .hash_map
                .get(self.source.index(self.start..self.current))
            {
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
        self.add_to_token_raw(
            TokenType::Number,
            String::from_str(self.source.index(self.start..self.current)).unwrap(),
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

        // Trim the surrounding quotes.
        self.add_to_token_raw(
            TokenType::Strings,
            String::from_str(self.source.index(self.start + 1..self.current - 1)).unwrap(),
        );
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
        self.add_to_token_raw(kind, String::new())
    }

    fn add_to_token_raw(&mut self, kind: TokenType, literal: String) {
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
    literal: String,
}

impl Token {
    pub fn new(token_type: TokenType, lexeme: String, line: i32, literal: String) -> Token {
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
        ret.push_str(self.literal.as_str());
        ret
    }
}
