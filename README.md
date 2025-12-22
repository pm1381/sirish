# Sirish

**Sirish** is a Go code-generation tool that automatically creates **Elastic APM–instrumented wrappers** for Go interfaces.

It allows you to add deep observability to your services without touching your core business logic or polluting your code with boilerplate tracing calls.

<p align="center">
  <img src="assets/logo.png" alt="sirish logo" width="260"/>
</p>

---

## 🚀 Key Features

* 🔍 **AST-based Discovery**: Parses Go source files using Go’s AST (no reflection, no runtime hacks).
* 🧩 **Non-Intrusive**: Generates wrappers that implement your interfaces, keeping your original implementations clean.
* 📡 **Native Elastic APM**: Automatically adds spans (and transactions) to every interface method.
* 🧠 **Context-Aware**:
    * Uses existing `context.Context` for distributed tracing.
    * Safely creates a transaction if no context exists (ideal for background jobs).
* 🧬 **Generics Support**: Full support for Go 1.18+ generic interfaces (e.g., `Repo[T any]`).
* 📦 **Smart Imports**: Uses `golang.org/x/tools/imports` to handle and format imports automatically.
* 🔁 **`go:generate` Ready**: Designed to fit perfectly into your existing Go build workflow.

---

## 🔍 How it Works

Sirish takes your interface and generates a struct that "wraps" your real implementation. Every time a method is called, the wrapper starts an Elastic APM span, records errors if they occur, and then calls your actual logic.

### Before (Your Code)
```go
type TestModule interface {
    DoTest1(ctx context.Context, req DoTest1Request, span int) (string, error)
}
```
### After (Generated Wrapper)
```go
// Automatically generated
func (w *TestModuleSirishWrapperImpl) DoTest1(ctx_0_0 context.Context, req DoTest1Request, span int) (string, error) {
    var DoTest1SpncSoZ_2_0 *apm.Span

    DoTest1SpncSoZ_2_0, ctx_0_0 = apm.StartSpan(ctx_0_0, "TestModule.DoTest1", w.tagType)
    defer DoTest1SpncSoZ_2_0.End()

    DoTest1ResUnwkcU_0_0, DoTest1ResUnmXMw_1_0 := w.wrapped.DoTest1(ctx_0_0, req, span)
    return DoTest1ResUnwkcU_0_0, DoTest1ResUnmXMw_1_0
}
```
---
## Installation

### Install the binary (recommended)
Run the following command to install the `sirish` binary into your `$GOBIN` (or `$GOPATH/bin`):
```bash
  go install github.com/pm1381/sirish@latest
```

### Verify the installation
```bash
  sirish --help
```
---
## Usage

Sirish is designed to be seamless and is typically invoked via the standard `go:generate` tool.

### 1️⃣ Mark your interface
Add a comment with the `sirish:` prefix followed by the interface name to identify the target for instrumentation.
another way for this is using directive -t flag where you can specify what interfaces you need sirish for

```go
// sirish:TestModule
type TestModule interface {
    DoSomething(ctx context.Context) error
    ProcessData(data string)
}
```
### 2️⃣ Add a go:generate directive
Insert the sirish command into your source file.
You can place this at the top of the file or in
a dedicated generate.go file within the same package.
```go
//go:generate sirish -f module.go -t TestModule
```
you can check the full flags using ```sirish --help```

### 3️⃣ Run generation
From your project root, run the standard Go generate command
```bash
  go generate ./...
```
Sirish will create an APM-instrumented wrapper file in
the same directory (e.g., module.sirish.go). If multiple
interfaces are defined in the same file, Sirish handles
the naming and imports deterministically.
---
## 📖 Examples

Check the ```examples/``` directory for a full implementation featuring:

* Echo Framework integration.
* Elastic APM middleware setup.
* Context propagation across layers.
---
## 🤔 When should I use Sirish?

Sirish is the right choice for your project if you:

* **Want APM visibility without pollution**: You want to track performance and errors in Elastic APM but don't want to clutter your core business logic with `apm.StartSpan` or `apm.CaptureError` calls.
* **Rely on Clean Architecture**: Your project uses interfaces to decouple layers, and you want a clean way to "plug in" observability as a decorator.
* **Need consistency**: You want to ensure that every method call in a specific service layer is traced identically across your entire team or organization.
* **Prefer code generation over "magic"**: You prefer explicit, type-safe Go code that you can read and debug over runtime reflection or complex proxy patterns.
