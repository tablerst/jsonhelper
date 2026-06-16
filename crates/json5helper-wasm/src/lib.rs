use wasm_bindgen::prelude::*;

#[wasm_bindgen]
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum WasmSyntax {
    Json,
    Jsonc,
    Json5,
}

impl From<WasmSyntax> for json5helper_core::Syntax {
    fn from(value: WasmSyntax) -> Self {
        match value {
            WasmSyntax::Json => Self::Json,
            WasmSyntax::Jsonc => Self::Jsonc,
            WasmSyntax::Json5 => Self::Json5,
        }
    }
}

#[wasm_bindgen]
pub fn format_json(input: &str, syntax: WasmSyntax, pretty: bool) -> Result<String, String> {
    json5helper_core::format(input, syntax.into(), pretty).map_err(|error| error.to_string())
}

#[wasm_bindgen]
pub fn repr_json(input: &str, pretty: bool) -> Result<String, String> {
    let value = repr_json::parse_repr(input).map_err(|error| error.to_string())?;
    if pretty {
        serde_json::to_string_pretty(&value)
    } else {
        serde_json::to_string(&value)
    }
    .map_err(|error| error.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn formats_json5_as_pretty_json() {
        let output = format_json(
            "{unquoted: 'value', trailing: [1, 2,]}",
            WasmSyntax::Json5,
            true,
        )
        .unwrap();

        assert!(output.contains("\"unquoted\": \"value\""));
        assert!(output.contains("\"trailing\""));
    }

    #[test]
    fn formats_jsonc_as_compact_json() {
        let output = format_json(r#"{"value": 1,}// comment"#, WasmSyntax::Jsonc, false).unwrap();

        assert_eq!(output, r#"{"value":1}"#);
    }

    #[test]
    fn converts_repr_to_pretty_json() {
        let output = repr_json("AgentExecutor(verbose=True)", true).unwrap();

        assert!(output.contains("\"$type\": \"AgentExecutor\""));
        assert!(output.contains("\"verbose\": true"));
    }

    #[test]
    fn returns_error_strings() {
        let error = format_json("{", WasmSyntax::Json5, true).unwrap_err();

        assert!(!error.is_empty());
    }
}
