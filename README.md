# spintax
Spintax is a text generation library on go. And also command line tool is included.

Spintax is an template with alternative parts that can be evaluated into one of many texts.
Depth is not limited

## Command
There is ready to use command line tool that can generate text from spintax template.
Usage:
```
$ spintax -h
usage: spintax [-c] [-f] [-a] <spintax> // generates one random text out of given spintax template
  -a	print all the possible texts instead of one random
  -c	print count of distinct possible texts instead of generating one
  -f	interpret argument as a file name with spintax not as spintax itself
  -h	print help and exit
```

### Install
You need to have go installed, or [install it now](https://golang.org/doc/install).

`go get -u github.com/nikandfor/spintax/...`

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
command line tool usage examples
```
$ spintax "a{b{c|d}e|f{g|h}i}j"
abdej
$ spintax "a{b{c|d}e|f{g|h}i}j"
afhij
$ echo "a{b{c|d}e|f{g|h}i}j" >spintax.txt
$ spintax -a -f spintax.txt
abcej
abdej
afgij
afhij
```

## Specification

Using Extended Backus-Naur Form Spintax is an Expression where
```
Expression  ::= { String | Alternative }
Alternative ::= "{" [ AltList ] "}"
AltList     ::= Expression { "|" Expression }
String is an empty or non-empty string
```
