# JSON Parser

This is a json parser package. I build it just for learning purposes.

Only exported function is `Parse(b []byte)`. It parses json text and returns a `map[string]interface{}` and an error if any.

- All numbers are treated as `float64` numbers.
- Arrays are fed into a list of interafces `[]interface{}`
- Objects are fed into a `map[string]interface{}`
