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

When compilations ends we will get a folder structure based on our settings.
For example:

windows/386/test.exe
windows/amd64/test.exe
linux/386/test
linux/amd64/test
darwin/386/test
darwin/amd64/test




