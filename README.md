# Elastistack

A simple tool for taking Golang stack trace data and pushing it as structured text into Elasticsearch.

### Usage
```
Given a textual Golang stack trace, the import
command will parse the input file and insert the stack
trace data into Elasticsearch for further analysis.

Usage:
  elastistack import [flags]

  Flags:
    -e, --host string    Hostname for Elasticsearch endpoint (default "localhost")
    -i, --input string   Input filename containing Golang stack trace data
    -p, --port int       Port for Elasticsearch endpoint (default 9200)

  Global Flags:
    --log-level string   set the logging level (info,warn,err,debug) (default "warn")
```

### A Full Example

See more details on how I use `elastistack` in this [blog post](https://integratedcode.us/2016/05/25/taming-the-golang-stack-trace/)
