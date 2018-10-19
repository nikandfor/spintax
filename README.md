# spintax
Spintax is a text generation library at go

Spintax is an template with alternative parts that can be evaluated into one of many texts.
Depth is not limited

## Example

```
Spintax                 -> results
"string"                -> "string"
"{opt_a|opt_b}"         -> "opt_a" or "opt_b"
"this is {|so} awesome" -> "this is so awesome" or
                           "this is  awesome" (two spaces here)
"this is{| so} awesome" -> "this is so awesome" or
                           "this is awesome" (one space here)
"{{a|b}|{c|d}}"
the same as "{a|b|c|d}" -> "a" or "b" or "c" or "d"
"a{b{c|d}e|f{g|h}i}j"   -> "abcej" or "abdej" or "afgij" or "afhij"
```

## Specification

Using Extended Backus-Naur Form Spintax is an Expression where
```
Expression  ::= { String | Alternative }
Alternative ::= "{" [ AltList ] "}"
AltList     ::= Expression { "|" Expression }
String is an empty or non-empty string
```
