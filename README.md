# cyto

**Cyto** is a Go library for cytometry data analysis. It provides tools to parse, process, and analyze Flow Cytometry Standard (FCS) files.

## 📦 Features

- Read and parse `.fcs` files
- Access FCS metadata and event data
- Command-line interface for basic operations

## 🛠️ Installation

### Using Go Get

```bash
go get github.com/immunoconductor/cyto@latest
```

The binary will be placed in $GOBIN (default: $HOME/go/bin).
Make sure that directory is in your PATH:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Clone Repository

```bash
git clone https://github.com/immunoconductor/cyto.git
cd cyto
go build
```

### Import as Module

Add it to your module by importing:

```go
import "github.com/immunoconductor/cyto/fcs"
```

Then run:

```bash
go mod tidy
```

## 📋 Usage

### Command Line Interface

Convert FCS files to CSV format:

```bash
cyto fcs -i example.fcs -o example.csv --transform --shortnames
```

Options:
- `-i`: Input FCS file path
- `-o`: Output CSV file path
- `--transform`: Apply arcsinh transformation to the data
- `--shortnames`: Use short names (concise identifiers for each parameter in the data file) commonly used in flow cytometry experiments

### Go Library Example

```go
import "github.com/immunoconductor/cyto/fcs"

func main() {
    fcs, err := fcs.Read("fcs3.1.fcs", false)
    if err != nil {
        t.Errorf(err.Error())
    }

    fmt.Println(fcs.Names())
}
```

## 📊 Documentation

Documentation is currently under development. For now, please refer to the code comments and examples in this README.

## 🔄 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the GPL-3 License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Contributors to the project
- The Flow Cytometry community
