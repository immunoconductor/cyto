# cyto

Cyto is a CLI library for the analysis of cytometry data.

# Cyto

**Cyto** is a Go library for cytometry data analysis. It provides tools to parse, process, and analyze Flow Cytometry Standard (FCS) files.

> GitHub: [immunoconductor/cyto](https://github.com/immunoconductor/cyto)

---

## 📦 Features

- Read and parse `.fcs` files
- Access FCS metadata and event data
- Command-line interface for basic operations

---

## 🛠️ Installation

```bash
go get github.com/immunoconductor/cyto@latest


The binary will be placed in $GOBIN (default: $HOME/go/bin).
Make sure that directory is in your PATH:

`export PATH="$PATH:$(go env GOPATH)/bin"`

Or clone directly:

```
git clone https://github.com/immunoconductor/cyto.git
cd cyto
go build
```

Or add it to your module by importing and running:

`go mod tidy`


```
import "github.com/immunoconductor/cyto/fcs"

func main() {
    fcs, err := fcs.Read("fcs3.1.fcs", false)
        if err != nil {
            t.Errorf(err.Error())
        }

        fmt.Println(fcs.Names())
}

```





