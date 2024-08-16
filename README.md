---

# Custom Grep Implementation

This project is a custom implementation of the Unix `grep` tool, built from scratch using a scanner, parser, and interpreter to find matching substrings based on regular expressions. The goal is to replicate and extend the functionality of `grep` while offering a deeper understanding of pattern matching.

## Features

-   Custom-built regex engine with support for a variety of regular expression patterns.
-   Lightweight and fast, designed to work seamlessly with the command line.
-   Handles complex pattern matching, including backreferences and match groups, without relying on external libraries.

## Supported Regex Syntax

The following regex features are supported:

-   **`.`**: Matches any single character except newline.
-   **`*`**: Matches zero or more occurrences of the preceding character.
-   **`+`**: Matches one or more occurrences of the preceding character.
-   **`?`**: Matches zero or one occurrence of the preceding character.
-   **`|`**: Logical OR between patterns.
-   **`()`:** Groups patterns to capture matches or apply operators.
-   **`\\n`**: Matches the same text as most recently matched by the nth capturing group (backreference).
-   **`{n}`**: Matches exactly `n` occurrences of the preceding character.
-   **`{n,}`**: Matches `n` or more occurrences of the preceding character.
-   **`{n,m}`**: Matches between `n` and `m` occurrences of the preceding character.
-   **`[]`**: Matches any one of the characters enclosed.
-   **`[^]`**: Matches any character not enclosed.
-   **`-`**: Represents a range of characters within `[]`.
-   **`^`**: Asserts position at the start of the line.
-   **`$`**: Asserts position at the end of the line.
-   **`\`**: Escapes special characters.
-   **`\w`**: Matches any word character (alphanumeric + underscore).
-   **`\d`**: Matches any digit.

## Build

To build the `grep` executable:

```shell
$ go build -o grep
```

This will compile the source code and create an executable named `grep`.

## Usage

To use the custom `grep`, pipe in the string you want to search and provide the pattern:

```shell
$ echo -n <line> | ./grep -E <pattern>
```

### Command Options

-   **`-E`**: Interpret the pattern as an extended regular expression.

## Examples

Here's how to use this custom `grep` implementation:

### Basic Example

```shell
$ echo -n "Hello, World\!" | ./grep -E "\w+, \w+\!"
```

Output:

```
found: Hello, World!
```

### Advanced Example

```shell
$ echo -n "aaabbbcccaaabbbccc" | ./grep -E "(a{2,3})(b{2,3})(c{2,3})\1\2\3"
```

Output:

```
found: aaabbbcccaaabbbccc
```

## Limitations

-   This implementation is a work in progress and may not support all edge cases found in traditional `grep`.
-   Performance optimizations are ongoing.

## Contributing

Contributions are welcome! If you have ideas for new features or find a bug, feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---
