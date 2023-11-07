package parser

// func main() {
// 	output, err := Run()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Printf(output)

// }

// func Run() (string, error) {
// 	if len(os.Args) != 1 {
// 		return "", fmt.Errorf("Usage: json-parser [FILE | TEXT]")
// 	}

// 	jsonBytes, err := getJson()
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get json bytes: %w", err)
// 	}

// }

// func getJson() ([]byte, error) {
// 	arg := os.Args[1]
// 	if _, err := os.Stat(arg); err != nil {
// 		// provided arg is not a valid path, treat it as a text
// 		return []byte(arg), nil
// 	}

// 	f, err := os.Open(arg)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open file %s: %w", arg, err)
// 	}

// 	var text []byte
// 	for {
// 		buff := make([]byte, 16*1024)
// 		n, err := f.Read(buff)
// 		if errors.Is(err, io.EOF) {
// 			text = append(text, buff[:n]...)
// 			return text, nil
// 		}
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to read from file %s: %w", arg, err)
// 		}

// 		text = append(text, buff[:n]...)
// 	}

// 	return text, nil
// }

func Decode([]byte) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	for {
		// get next token
		// check token type:
		// 		- opening curly brace -> new object
		// 		- closing curly brace -> end of object
		// 		- double quotes -> begin or end of key or string value
		// 		- opening square bracket -> begin of list
		// 		- closing square bracked -> end of list
		//		- true -> true val
		// 		- false -> false val
		// 		- integer -> integer val
	}
}

func decodeObject([]byte) (map[string]interface{}, error) {
	/*
		{ [whitespace] [key:val [,key:val]...] [whitespace]}
		get opening curly brace
		get whitespace

		l1: check next token:
			- if closing curly brace, close object and return
			- if whitespace, discard
			- if string, this is a key:
				- get colon.
				- get value
				- check next token:
					- if comma, repeat from l1
					- if closing curly brace, close object and return
		check next token:
			- if comma, repeat from l1


	*/
}
