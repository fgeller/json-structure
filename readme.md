# json-structure

Reads a JSON value from stdin and converts it to a JSON value that describes the structure of the original value.

## Example

```
% json-structure -help
Usage of json-structure:
  -flatten
        flatten schemata in arrays and combine objects
  -schema
        output a json schema
% cat sample.json 
{
  "productId": 1,
  "productName": "A green door",
  "price": 12.50,
  "tags": [ "home", "green" ]
}
% cat sample.json | json-structure | jq
{
  "price": "number",
  "productId": "number",
  "productName": "string",
  "tags": [
    "string",
    "string"
  ]
}
% cat sample.json | json-structure -flatten | jq
{
  "price": "number",
  "productId": "number",
  "productName": "string",
  "tags": [
    "string"
  ]
}
% cat sample.json | json-structure -flatten -schema | jq
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "properties": {
    "price": {
      "type": "number"
    },
    "productId": {
      "type": "number"
    },
    "productName": {
      "type": "string"
    },
    "tags": {
      "type": "array",
      "contains": {
        "type": "string"
      }
    }
  }
}
```

## Installation

```
% go install github.com/fgeller/json-structure@latest
```
