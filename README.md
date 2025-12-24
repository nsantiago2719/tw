# tw

A CLI tool for managing Terraform resources with simplified workflows.

## Features

- Register and manage multiple Terraform resources
- Execute Terraform plan and apply operations with variable files
- Interactive command-line interface with support for user input
- Centralized configuration management

## Installation

```bash
make build
```

## Usage

Initialize the configuration:
```bash
./bin/tw init
```

Register a resource:
```bash
./bin/tw register --name my-resource --path ./path/to/terraform --var-files values.tfvars
```

List registered resources:
```bash
./bin/tw list-resources
```

Plan changes:
```bash
./bin/tw plan my-resource
```

Apply changes:
```bash
./bin/tw run my-resource
```

## License

MIT License

Copyright (c) 2025

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## AI DISCLAIMER
This README file was generated with the assistance of AI technology. While efforts have been made to ensure accuracy and clarity, please review the content for any potential errors or omissions. The author assumes no responsibility for any issues arising from the use of this documentation.

Some of the code snippets and explanations may have been created or enhanced using AI tools. Users are advised to verify the information and adapt it as necessary for their specific use cases.
