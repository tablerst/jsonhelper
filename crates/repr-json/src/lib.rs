use serde_json::{Map, Number, Value, json};
use thiserror::Error;

#[derive(Debug, Error, PartialEq, Eq)]
pub enum Error {
    #[error("unexpected end of input")]
    UnexpectedEof,
    #[error("unexpected token at byte {position}: {message}")]
    UnexpectedToken { position: usize, message: String },
}

pub type Result<T> = std::result::Result<T, Error>;

pub fn parse_repr(input: &str) -> Result<Value> {
    let mut parser = Parser::new(input);
    let value = parser.parse_pipe()?;
    parser.skip_ws();
    if !parser.is_eof() {
        return Err(parser.unexpected("expected end of input"));
    }
    Ok(value)
}

struct Parser<'a> {
    input: &'a str,
    pos: usize,
}

impl<'a> Parser<'a> {
    fn new(input: &'a str) -> Self {
        Self { input, pos: 0 }
    }

    fn parse_pipe(&mut self) -> Result<Value> {
        let first = self.parse_value()?;
        let mut items = vec![first];
        loop {
            self.skip_ws();
            if !self.consume('|') {
                break;
            }
            items.push(self.parse_value()?);
        }

        if items.len() == 1 {
            Ok(items.remove(0))
        } else {
            Ok(json!({
                "$op": "pipe",
                "items": items,
            }))
        }
    }

    fn parse_value(&mut self) -> Result<Value> {
        self.skip_ws();
        match self.peek_char() {
            Some('\'') | Some('"') => self.parse_string().map(Value::String),
            Some('{') => self.parse_map(),
            Some('[') => self.parse_sequence('[', ']'),
            Some('(') => self.parse_sequence('(', ')'),
            Some('-') | Some('0'..='9') => self.parse_number_or_raw(),
            Some(ch) if is_ident_start(ch) => self.parse_identifier_value(),
            Some('<') => Ok(Value::String(self.parse_angle_repr())),
            Some(_) => Err(self.unexpected("expected value")),
            None => Err(Error::UnexpectedEof),
        }
    }

    fn parse_identifier_value(&mut self) -> Result<Value> {
        let name = self.parse_identifier()?;
        match name.as_str() {
            "True" | "true" => return Ok(Value::Bool(true)),
            "False" | "false" => return Ok(Value::Bool(false)),
            "None" | "null" => return Ok(Value::Null),
            _ => {}
        }

        self.skip_ws();
        if self.consume('(') {
            self.parse_call(name)
        } else {
            Ok(Value::String(name))
        }
    }

    fn parse_call(&mut self, name: String) -> Result<Value> {
        let mut object = Map::new();
        let mut args = Vec::new();
        object.insert("$type".to_string(), Value::String(name));

        loop {
            self.skip_ws();
            if self.consume(')') {
                break;
            }

            if self.starts_lambda() {
                args.push(json!({
                    "$repr": self.parse_raw_until_delimiter(&[',', ')'])
                }));
            } else if let Some(key) = self.try_parse_key_assignment()? {
                let value = self.parse_pipe()?;
                object.insert(key, value);
            } else {
                args.push(self.parse_pipe()?);
            }

            self.skip_ws();
            if self.consume(',') {
                continue;
            }
            if self.consume(')') {
                break;
            }
            return Err(self.unexpected("expected ',' or ')'"));
        }

        if !args.is_empty() {
            object.insert("$args".to_string(), Value::Array(args));
        }
        Ok(Value::Object(object))
    }

    fn parse_map(&mut self) -> Result<Value> {
        self.expect('{')?;
        let mut object = Map::new();
        loop {
            self.skip_ws();
            if self.consume('}') {
                break;
            }

            let key = if matches!(self.peek_char(), Some('\'') | Some('"')) {
                self.parse_string()?
            } else {
                self.parse_identifier()?
            };
            self.skip_ws();
            self.expect(':')?;
            let value = self.parse_pipe()?;
            object.insert(key, value);

            self.skip_ws();
            if self.consume(',') {
                continue;
            }
            if self.consume('}') {
                break;
            }
            return Err(self.unexpected("expected ',' or '}'"));
        }
        Ok(Value::Object(object))
    }

    fn parse_sequence(&mut self, open: char, close: char) -> Result<Value> {
        self.expect(open)?;
        let mut items = Vec::new();
        loop {
            self.skip_ws();
            if self.consume(close) {
                break;
            }
            items.push(self.parse_pipe()?);
            self.skip_ws();
            if self.consume(',') {
                continue;
            }
            if self.consume(close) {
                break;
            }
            return Err(self.unexpected(format!("expected ',' or '{close}'")));
        }
        Ok(Value::Array(items))
    }

    fn try_parse_key_assignment(&mut self) -> Result<Option<String>> {
        let save = self.pos;
        self.skip_ws();
        let key = if matches!(self.peek_char(), Some('\'') | Some('"')) {
            self.parse_string()?
        } else if matches!(self.peek_char(), Some(ch) if is_ident_start(ch)) {
            self.parse_identifier()?
        } else {
            self.pos = save;
            return Ok(None);
        };

        self.skip_ws();
        if self.consume('=') {
            Ok(Some(key))
        } else {
            self.pos = save;
            Ok(None)
        }
    }

    fn parse_number_or_raw(&mut self) -> Result<Value> {
        let start = self.pos;
        if self.peek_char() == Some('-') {
            self.pos += 1;
        }
        while matches!(self.peek_char(), Some('0'..='9')) {
            self.pos += 1;
        }
        if self.peek_char() == Some('.') {
            self.pos += 1;
            while matches!(self.peek_char(), Some('0'..='9')) {
                self.pos += 1;
            }
        }
        if matches!(self.peek_char(), Some('e' | 'E')) {
            self.pos += 1;
            if matches!(self.peek_char(), Some('+' | '-')) {
                self.pos += 1;
            }
            while matches!(self.peek_char(), Some('0'..='9')) {
                self.pos += 1;
            }
        }

        let raw = &self.input[start..self.pos];
        if raw.contains(['.', 'e', 'E']) {
            raw.parse::<f64>()
                .ok()
                .and_then(Number::from_f64)
                .map(Value::Number)
                .ok_or_else(|| self.unexpected("invalid number"))
        } else {
            raw.parse::<i64>()
                .map(|value| Value::Number(value.into()))
                .map_err(|_| self.unexpected("invalid number"))
        }
    }

    fn parse_string(&mut self) -> Result<String> {
        let quote = self.peek_char().ok_or(Error::UnexpectedEof)?;
        self.expect(quote)?;
        let mut output = String::new();
        while let Some(ch) = self.next_char() {
            if ch == quote {
                return Ok(output);
            }
            if ch == '\\' {
                let escaped = self.next_char().ok_or(Error::UnexpectedEof)?;
                output.push(match escaped {
                    'n' => '\n',
                    'r' => '\r',
                    't' => '\t',
                    '\\' => '\\',
                    '\'' => '\'',
                    '"' => '"',
                    other => other,
                });
            } else {
                output.push(ch);
            }
        }
        Err(Error::UnexpectedEof)
    }

    fn parse_identifier(&mut self) -> Result<String> {
        self.skip_ws();
        let start = self.pos;
        match self.peek_char() {
            Some(ch) if is_ident_start(ch) => self.pos += ch.len_utf8(),
            _ => return Err(self.unexpected("expected identifier")),
        }
        while matches!(self.peek_char(), Some(ch) if is_ident_continue(ch)) {
            self.pos += self.peek_char().unwrap().len_utf8();
        }
        Ok(self.input[start..self.pos].to_string())
    }

    fn parse_angle_repr(&mut self) -> String {
        let start = self.pos;
        while let Some(ch) = self.next_char() {
            if ch == '>' {
                break;
            }
        }
        self.input[start..self.pos].to_string()
    }

    fn parse_raw_until_delimiter(&mut self, delimiters: &[char]) -> String {
        let start = self.pos;
        let mut depth = 0usize;
        while let Some(ch) = self.peek_char() {
            if depth == 0 && delimiters.contains(&ch) {
                break;
            }
            match ch {
                '(' | '[' | '{' => depth += 1,
                ')' | ']' | '}' => depth = depth.saturating_sub(1),
                '\'' | '"' => {
                    let _ = self.parse_string();
                    continue;
                }
                _ => {}
            }
            self.pos += ch.len_utf8();
        }
        self.input[start..self.pos].trim().to_string()
    }

    fn starts_lambda(&mut self) -> bool {
        self.skip_ws();
        self.input[self.pos..].starts_with("lambda ")
    }

    fn skip_ws(&mut self) {
        while matches!(self.peek_char(), Some(ch) if ch.is_whitespace()) {
            self.pos += self.peek_char().unwrap().len_utf8();
        }
    }

    fn expect(&mut self, expected: char) -> Result<()> {
        if self.consume(expected) {
            Ok(())
        } else {
            Err(self.unexpected(format!("expected '{expected}'")))
        }
    }

    fn consume(&mut self, expected: char) -> bool {
        if self.peek_char() == Some(expected) {
            self.pos += expected.len_utf8();
            true
        } else {
            false
        }
    }

    fn next_char(&mut self) -> Option<char> {
        let ch = self.peek_char()?;
        self.pos += ch.len_utf8();
        Some(ch)
    }

    fn peek_char(&self) -> Option<char> {
        self.input[self.pos..].chars().next()
    }

    fn is_eof(&self) -> bool {
        self.pos >= self.input.len()
    }

    fn unexpected(&self, message: impl Into<String>) -> Error {
        Error::UnexpectedToken {
            position: self.pos,
            message: message.into(),
        }
    }
}

fn is_ident_start(ch: char) -> bool {
    ch == '_' || ch.is_ascii_alphabetic()
}

fn is_ident_continue(ch: char) -> bool {
    ch == '_' || ch == '.' || ch.is_ascii_alphanumeric()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn parses_constructor_keywords() {
        let value =
            parse_repr("AgentExecutor(verbose=True, agent=RunnableMultiActionAgent())").unwrap();

        assert_eq!(value["$type"], "AgentExecutor");
        assert_eq!(value["verbose"], true);
        assert_eq!(value["agent"]["$type"], "RunnableMultiActionAgent");
    }

    #[test]
    fn parses_pipe_operator() {
        let value = parse_repr("RunnableAssign(mapper={agent_scratchpad: RunnableLambda(lambda x: x)}) | ChatPromptTemplate(input_variables=['input'])").unwrap();

        assert_eq!(value["$op"], "pipe");
        assert_eq!(value["items"][0]["$type"], "RunnableAssign");
        assert_eq!(value["items"][1]["$type"], "ChatPromptTemplate");
    }

    #[test]
    fn parses_lists_maps_and_scalars() {
        let value = parse_repr("{name: 'repr', flags: [True, False, None], count: 3}").unwrap();

        assert_eq!(value["name"], "repr");
        assert_eq!(value["flags"][0], true);
        assert_eq!(value["flags"][2], Value::Null);
        assert_eq!(value["count"], 3);
    }

    #[test]
    fn parses_positional_args_tuple_numbers_and_angle_repr() {
        let value =
            parse_repr("Thing('line\\ntext', -3, 1.25e2, (<function _get_type at 0xabc>,))")
                .unwrap();

        assert_eq!(value["$type"], "Thing");
        assert_eq!(value["$args"][0], "line\ntext");
        assert_eq!(value["$args"][1], -3);
        assert_eq!(value["$args"][2], 125.0);
        assert_eq!(value["$args"][3][0], "<function _get_type at 0xabc>");
    }

    #[test]
    fn parses_string_assignment_keys_and_identifiers() {
        let value = parse_repr("FieldInfo('annotation'=None, mode=test.mode)").unwrap();

        assert_eq!(value["$type"], "FieldInfo");
        assert_eq!(value["annotation"], Value::Null);
        assert_eq!(value["mode"], "test.mode");
    }

    #[test]
    fn reports_trailing_input() {
        let err = parse_repr("AgentExecutor() trailing").unwrap_err();

        assert_eq!(
            err,
            Error::UnexpectedToken {
                position: 16,
                message: "expected end of input".to_string()
            }
        );
    }

    #[test]
    fn reports_missing_identifier_in_map() {
        let err = parse_repr("{: 1}").unwrap_err();

        assert!(matches!(
            err,
            Error::UnexpectedToken {
                message,
                ..
            } if message == "expected identifier"
        ));
    }

    #[test]
    fn reports_unterminated_string() {
        let err = parse_repr("'unterminated").unwrap_err();

        assert_eq!(err, Error::UnexpectedEof);
    }

    #[test]
    fn reports_bad_number() {
        let err = parse_repr("-").unwrap_err();

        assert!(matches!(
            err,
            Error::UnexpectedToken {
                message,
                ..
            } if message == "invalid number"
        ));
    }
}
