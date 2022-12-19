use std::str::FromStr;

#[derive(EnumString, Display, Debug, PartialEq)]
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

pub struct Scanner {
    source: String,
    hasError: bool,
}

impl Scanner {
    pub fn new_from(source: String) -> Scanner {
        Scanner {
            source,
            hasError: false,
        }
    }

    pub fn scan_tokens(&self) -> Vec<Token> {
        let ret = Vec::new();
        ret
    }
}

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
