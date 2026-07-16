# VTL (Velocity Template Language) Engine

## Overview

This package provides a Velocity Template Language (VTL) engine implementation for processing API Gateway request and response templates. VTL is a template language used by Amazon API Gateway to transform and manipulate data during API requests and responses.

## Features

### Context Variables
- `$context.*` - Access to request context information such as authoriser, identity, and request/response overrides
- `$stageVariables.*` - Access to deployment stage variables
- `$input.body`, `$input.json()`, `$input.path()`, `$input.params()` - Access to request input data
- Simple variable interpolation using `$variable` syntax

### Control Flow
- `#if`, `#elseif`, `#else`, `#end` - Conditional branching
- `#foreach`, `#end` - Loop iteration over collections
- `#set($var = value)` - Variable assignment

### Utility Functions
- `$util.urlEncode()` - URL encoding
- `$util.urlDecode()` - URL decoding
- `$util.base64Encode()` - Base64 encoding
- `$util.base64Decode()` - Base64 decoding
- `$util.parseJson()` - JSON parsing
- `$util.escapeJavaScript()` - JavaScript string escaping

### Comparison Operators
- Equality: `==`, `!=`
- Relational: `>`, `<`, `>=`, `<=`
- Logical: `&&`, `||`, `!`

## Usage

```go
engine := vtl.NewEngine()

ctx := &vtl.Context{
    Body:      `{"message": "Hello, World!"}`,
    Params:    map[string]string{"name": "John"},
    StageVars: map[string]string{"env": "prod"},
    Context:   map[string]interface{}{"requestId": "12345"},
}
engine.SetContext(ctx)

template := `#if($input.params('name'))
Hello, $input.params('name')!
#else
Hello, World!
#end`

result, err := engine.Transform(template)
if err != nil {
    log.Fatal(err)
}
```

## Installation

```bash
go get github.com/vorpalstacks/vorpalstacks/pkg/vtl
```

## License

Copyright 2026 Vorpalstacks. Licensed under the Apache License, Version 2.0. See the LICENSE file for details.
