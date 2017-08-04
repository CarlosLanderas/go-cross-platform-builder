**Simple GO wrapper above GO CLI for Cross Platform compilation**

Install GO enviroment and build the go file:

**go build crossbuild.go**

configure your target platforms and architectures in targets.json. For example:

```json

{
  "targets" : [
      {"id": "windows" , "arch": ["386", "amd64"]},
      {"id": "linux" , "arch": ["386", "amd64"]},
      {"id": "darwin" , "arch": ["386", "amd64"]}
  ]
}
```

And then ...

Execute the CLI with -compile {file} to compile the target file.

Example:

**crossbuild.exe -compile test\test.go**




