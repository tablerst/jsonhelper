use serde_json::Value;
use thiserror::Error;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Syntax {
    Json,
    Jsonc,
    Json5,
}

impl Syntax {
    pub fn from_hint(value: &str) -> Option<Self> {
        match value.to_ascii_lowercase().as_str() {
            "json" => Some(Self::Json),
            "jsonc" => Some(Self::Jsonc),
            "json5" => Some(Self::Json5),
            _ => None,
        }
    }
}

#[derive(Debug, Error)]
pub enum Error {
    #[error("failed to parse JSON: {0}")]
    Json(#[from] serde_json::Error),
    #[error("failed to parse JSON5: {0}")]
    Json5(#[from] json5::Error),
    #[error("failed to parse JSONC: {0}")]
    Jsonc(#[from] jsonc_parser::errors::ParseError),
}

pub type Result<T> = std::result::Result<T, Error>;

pub fn parse(input: &str, syntax: Syntax) -> Result<Value> {
    match syntax {
        Syntax::Json => Ok(serde_json::from_str(input)?),
        Syntax::Jsonc => Ok(jsonc_parser::parse_to_serde_value(
            input,
            &Default::default(),
        )?),
        Syntax::Json5 => Ok(json5::from_str(input)?),
    }
}

pub fn format(input: &str, syntax: Syntax, pretty: bool) -> Result<String> {
    let value = parse(input, syntax)?;
    if pretty {
        Ok(serde_json::to_string_pretty(&value)?)
    } else {
        Ok(serde_json::to_string(&value)?)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn parses_json() {
        let value = parse(r#"{"name":"json","ok":true}"#, Syntax::Json).unwrap();
        assert_eq!(value["name"], "json");
        assert_eq!(value["ok"], true);
    }

    #[test]
    fn parses_jsonc_comments() {
        let value = parse(
            r#"
            {
              // line comment
              "name": "jsonc",
              "items": [1, 2,],
            }
            "#,
            Syntax::Jsonc,
        )
        .unwrap();

        assert_eq!(value["name"], "jsonc");
        assert_eq!(value["items"][1], 2);
    }

    #[test]
    fn parses_json5_features() {
        let value = parse(
            r#"
            {
              unquoted: 'json5',
              hexadecimal: 0xdecaf,
              trailingComma: [1, 2,],
            }
            "#,
            Syntax::Json5,
        )
        .unwrap();

        assert_eq!(value["unquoted"], "json5");
        assert_eq!(value["hexadecimal"], 912559);
        assert_eq!(value["trailingComma"][0], 1);
    }

    #[test]
    fn formats_pretty_json() {
        let output = format(r#"{unquoted:"json5"}"#, Syntax::Json5, true).unwrap();
        assert!(output.contains("\n  \"unquoted\": \"json5\"\n"));
    }
}
