# TypeR
> A superset language implemented in Go that "types" the R language

![TypeR logo](./logo/typer.png)

## Introduction

According to the language site itself, [R](https://www.r-project.org/) is:

_"R is a language and environment for statistical computing and graphics."_

Because it is a scripting language, R seeks to allow flexibility in the development and prototyping of ideas. Through its weak typing system this can be a problem for applications in production -- it is worth noting that this is not a "problem" of language since it was meant to behave the way it is, but this same behavior may end becoming a difficulty to maintain a great code base in the language.

**TypeR** tries to be for R what [TypeScript](https://www.typescriptlang.org/) is for [JavaScript](https://www.javascript.com/), implementing a strong typing system that allows inference and statically typed -- the idea is to go beyond just being types and also to limit the language only to the functional paradigm, cleaning up a little of the multi paradigm of R.

At the end of the day the idea is to write a "functional and typed R code" which will then spit out a normal code in R after all the checks are done, avoiding possible errors when the code is running in production.

## How should it supposed to be?

A small example of how language it's supposed to be one day.

**WARNING**: these are just ideas, they may be possible or not.

### Functions

```
add := Numeric a => a -> a -> a
add <- function(x, y) x + y
```

The function header is similar to Haskell's -- read more at [functional Approach](#Functional-approach) --; therefore, it also does not need to be declared, the language engine is able to infer the types of the parameters by the operations performed with them in the function body.

```
multiply <- function(x, y) x * y
# multiply := Numeric a => a -> a -> a
```

#### Throw notation
If a function throws an error or a warning, the caller has to know that without having to look at the implementation of the function itself.

```
isURLValid("www.google.com")
# isURLValid := String -> Bool; throws error, warning
```

### Constants

```
result <- add(1, 2)
# result is 3

result <- 4
# This will throw a compiler error, because result is a constant
```

If you want a "constant" to have its value changed, you must use the **let** keyword.

```
let result <- add(1, 2)
# result is 3

result <- 4
# result is 4
```

### Point free notation

```
square <- function(x) x ^ 2

addTwo <- function(x) x + 2

result <- addTwo . square 2
# result is 6
```

### Prefix operations
```
result <- (+) 1 2
# result is 3
```

### Polymorphism
> More on that later on

## Changes from R

- Removing loops as **for**, **while**, **do ... while**;
- ...

## Why

The following topics try to clarify the choice of some design decisions.

### Go

As R alone is not a very performative language, Go was chosen to meet such need.

_"But why Go and not another language?"_

The answer is simple, Go is:

- Easy to read and to write tests
- Has a large community and documentation
- Its concurrent design helps when writing a compiler
- ...


### Functional approach

The choice of just following the functional paradigm is simply a personal decision, since the main use of R for the project author is for mathematical scenarios. Having a background in [Haskell](https://www.haskell.org/), this has greatly influenced prioritizing such a decision.

Even if it does not have some of the practicalities of the functional paradigm like guards and pattern matching, it may be that if it is possible to emulate such designs they are added to the language.

## TODO

- Right now? Everything, nothing is current working
- Create a linter for the language so that more flexible patterns are placed and directed to the community to configure them in the way they think best
- Create code analyzer to perform duplicity analysis and other things just as [codeclimate](https://codeclimate.com/) already performs
- Write another packet to transpile the TypeR code into [Julia](https://julialang.org/) code
- Help out [Romain](https://community.rstudio.com/t/running-go-code-from-r/2340/3) write a Go to R integration package or even allow such integration into TypeR itself
- Much more

## Author

As the idea is to actually leave this repository just for discussions related to code already present: issues and pull requests, any other questions like project schedule or ask to add new features, you can talk to me about it at:

- [Twitter](https://twitter.com/the_fznd)
- [Telegram](https://t.me/farmy)

## Reference

### Books

- [Writing An Interpreter In Go](https://interpreterbook.com/)
- [Writing A Compiler In Go](https://compilerbook.com/)

### Podcasts

#### Hipsters
- [Um Pouco de Compiladores](https://hipsters.tech/um-pouco-de-compiladores-hipsters-ponto-tech-105/)
- [Linguagens Funcionais](https://hipsters.tech/linguagens-funcionais-hipsters-91/)
- [Grandes Livros de Tecnologia](https://hipsters.tech/grandes-livros-de-tecnologia-hipsters-113/)

### Videos

- [Inkscape Tutorial: Abstract Galaxy Logo](https://youtu.be/AgbsozDUyTs)
- [Lexical Scanning in Go - Rob Pike](https://youtu.be/HxaD_trXwRE)
- [Sydney Golang Meetup - Rob Pike - Go 2 Draft Specifications](https://youtu.be/RIvL2ONhFBI)
- [Why Go Is Successful (aka Simplicity is Complicated)](https://youtu.be/k9Zbuuo51go)
