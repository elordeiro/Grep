This is an implementation of grep using a scanner, a parser and an interpreter to find matching substrings.

## Usage

```shell
$ ./grep <line> <pattern>
```

## Example

```shell
$ ./grep "Hello, World!" "\w+, \w+!"
found: "Hello, World!"
```
