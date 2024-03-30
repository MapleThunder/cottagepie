# CottagePie Language

This project is a language I am creating while reading the book [Writing an Interpreter in GO](https://interpreterbook.com).

List of original features from the book:

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

List of features I would like to add if I get the time:

- [ ] Function declaration support
- [ ] Comment support ("//")
- [ ] Add line numbers (and maybe column numbers) to error messages

---

## Usage

Here is how to bind values to names in CottagePie:

```js
bake age to 22;
bake name to "CottagePie";
bake result to 10 * (20 / 2);
```

Besides integers, booleans and strings, the CottagePie interpreter also support arrays and hashes. Here’s what binding an array of integers to a name looks like:

```js
bake my_array to [1, 2, 3, 4, 5];
```

And here is a hash, where values are associated with keys:

```js
bake niko to {"name": "Niko", "age": 22};
```

Accessing the elements in arrays and hashes is done with index expressions:

```js
myArray[0]       // => 1
niko["name"]     // => "Niko"
```

The bake statements can also be used to bind recipes (functions) to names. Here’s a small recipe that adds two numbers:

```js
bake add to recipe(a, b) { serves a + b; };
```

CottagePie not only supports serves (return) statements, implicit serves values are also possible ! Which means we can leave out the serves if we want to:

```js
bake add to rc(a, b) { a + b; };
```

And calling a recipe is as easy as you’d expect:

```js
add(1, 2);
```
