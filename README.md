# CottagePie Language

This project is a language I am creating while reading the book [Writing an Interpreter in GO](https://interpreterbook.com).

List of features:

- C-like syntax
- variable bindings
- integers & booleans
- arithmetic expressions
- built-in functions
- first-class & higher-order functions
- closures
- multiple data structures
  - string
  - array
  - hash

Need to Build:

- [ ] Lexer
- [ ] Parser
- [ ] Abstract Syntax Tree (AST)
- [ ] Internal object system
- [ ] Evaluator

---

Here is how to bind values to names in CottagePie:

```js
let age = 22;
let name = "CottagePie";
let result = 10 * (20 / 2);
```

Besides integers, booleans and strings, the CottagePie interpreter also support arrays and hashes. Here’s what binding an array of integers to a name looks like:

```js
let my_array = [1, 2, 3, 4, 5];
```

And here is a hash, where values are associated with keys:

```js
let niko = {"name": "Niko", "age": 22};
```

Accessing the elements in arrays and hashes is done with index expressions:

```js
myArray[0]       // => 1
niko["name"]     // => "Niko"
```

The let statements can also be used to bind functions to names. Here’s a small function that adds two numbers:

```js
let add = fn(a, b) { return a + b; };
```

But CottagePie not only supports return statements. Implicit return values are also possible, which means we can leave out the return if we want to:

```js
let add = fn(a, b) { a + b; };
```

And calling a function is as easy as you’d expect:

```js
add(1, 2);
```
