package main

import (
    "bufio"
    "flag"
    "fmt"
    "io"
    "net/url"
    "os"
    "sort"
    "strings"
)

func main() {
    var appendMode bool
    var inputFile string
    flag.BoolVar(&appendMode, "a", false, "Append the value instead of replacing it")
    flag.StringVar(&inputFile, "i", "", "Input file")
    flag.Parse()

    seen := make(map[string]bool)

    // read URLs on stdin, then replace the values in the query string
    // with some user-provided value
    f, err := os.Open(inputFile)
    defer f.Close()
    rd := bufio.NewReader(f)
//  sc := bufio.NewScanner(os.Stdin)
    for {
        line, err := rd.ReadString('\n')
        if err != nil || io.EOF == err {
            break
        }
        u, err := url.Parse(line)
        if err != nil {
            fmt.Fprintf(os.Stderr, "failed to parse url %s [%s]\n", line, err)
            continue
        }

        // Go's maps aren't ordered, but we want to use all the param names
        // as part of the key to output only unique requests. To do that, put
        // them into a slice and then sort it.
        pp := make([]string, 0)
        for p, _ := range u.Query() {
            pp = append(pp, p)
        }
        sort.Strings(pp)

        key := fmt.Sprintf("%s%s?%s", u.Hostname(), u.EscapedPath(), strings.Join(pp, "&"))

        // Only output each host + path + params combination once
        if _, exists := seen[key]; exists {
            continue
        }
        seen[key] = true

        qs := url.Values{}
        for param, vv := range u.Query() {
            if appendMode {
                qs.Set(param, vv[0]+flag.Arg(0))
            } else {
                qs.Set(param, flag.Arg(0))
            }
        }

        u.RawQuery = qs.Encode()

        fmt.Printf("%s\n", u)

    }
}
