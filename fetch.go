package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"

    "github.com/PuerkitoBio/goquery"
)

func main() {
    // Define flags for the input and output file names
    asnFile := flag.String("f", "", "The file containing the list of ASNs in the format one per line AS123")
    outputFile := flag.String("o", "", "The output file to write the routes to")
    flag.Parse()

    if *asnFile == "" || *outputFile == "" {
        log.Fatal("You must specify an AS input file and an output file")
    }

    file, err := os.Open(*asnFile)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    out, err := os.Create(*outputFile)
    if err != nil {
        log.Fatal(err)
    }
    defer out.Close()

    writer := bufio.NewWriter(out)

    client := &http.Client{}

    userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        asn := scanner.Text()
        url := fmt.Sprintf("https://bgp.he.net/%s#_prefixes", asn)

        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            log.Printf("Error creating request for %s: %v", asn, err)
            continue
        }
        req.Header.Set("User-Agent", userAgent)

        res, err := client.Do(req)
        if err != nil {
            log.Printf("Error fetching URL for %s: %v", asn, err)
            continue
        }

        doc, err := goquery.NewDocumentFromReader(res.Body)
        if err != nil {
            log.Printf("Error parsing HTML for %s: %v", asn, err)
            res.Body.Close()
            continue
        }
        res.Body.Close()

        doc.Find("a[href*='/net/']").Each(func(i int, s *goquery.Selection) {
            href, exists := s.Attr("href")
            if exists {
                prefix := strings.TrimPrefix(href, "/net/")
                if _, err := writer.WriteString(prefix + "\n"); err != nil {
                    log.Printf("Error writing to output file: %v", err)
                }
            }
        })

        writer.Flush()
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Completed writing IP routes to %s\n", *outputFile)
}
