package crobat

import (
    "bufio"
    "context"
    "crypto/tls"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strings"
    "sync"

    "github.com/pkg/errors"
    "golang.org/x/sync/errgroup"

    crobat "github.com/cgboal/sonarsearch/proto"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
)

type CrobatClient struct {
    conn   *grpc.ClientConn
    client crobat.CrobatClient
}

type UniqueStringSlice map[string]struct{}

func (slice *UniqueStringSlice) Append(value string) bool {
    if _, ok := (*slice)[value]; !ok {
        (*slice)[value] = struct{}{}
        return true
    }
    return false
}

func ProcessArg(arg string) (args []string) {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        log.Fatal(err)
    }

    fileName := fmt.Sprintf("%s/%s", dir, arg)
    if _, err := os.Stat(fileName); err == nil {
        file, _ := os.Open(fileName)
        defer file.Close()
        scanner := bufio.NewScanner(file)

        for scanner.Scan() {
            args = append(args, scanner.Text())
        }
    } else {
        args = strings.Split(arg, " ")
    }

    return args
}

func NewCrobatClient() CrobatClient {
    config := &tls.Config{}
    conn, err := grpc.Dial("crobat-rpc.omnisint.io:443", grpc.WithTransportCredentials(credentials.NewTLS(config)))
    if err != nil {
        log.Fatal(err)
    }

    client := crobat.NewCrobatClient(conn)
    return CrobatClient{
        conn:   conn,
        client: client,
    }
}

func (c *CrobatClient) GetSubdomains(arg string, resultsChan chan string) {
    args := ProcessArg(arg)
    defer close(resultsChan)
    for _, domain := range args {
        query := &crobat.QueryRequest{
            Query: domain,
        }

        stream, err := c.client.GetSubdomains(context.Background(), query)
        if err != nil {
            log.Fatal(err)
        }

        for {
            domain, err := stream.Recv()
            if err == io.EOF {
                break
            }

            if err != nil {
                log.Fatal(err)
            }
            resultsChan <- domain.Domain
        }
    }

}

func (c *CrobatClient) GetTlds(arg string, resultsChan chan string) {
    args := ProcessArg(arg)
    defer close(resultsChan)
    for _, domain := range args {

        query := &crobat.QueryRequest{
            Query: domain,
        }

        stream, err := c.client.GetTLDs(context.Background(), query)
        if err != nil {
            log.Fatal(err)
        }

        for {
            domain, err := stream.Recv()
            if err == io.EOF {
                break
            }

            if err != nil {
                log.Fatal(err)
            }
            resultsChan <- domain.Domain
        }
    }

}

func (c *CrobatClient) ReverseDNS(arg string, resultsChan chan string) {
    args := ProcessArg(arg)
    defer close(resultsChan)
    for _, ipv4 := range args {

        query := &crobat.QueryRequest{
            Query: ipv4,
        }

        stream, err := c.client.ReverseDNS(context.Background(), query)
        if err != nil {
            log.Fatal(err)
        }

        for {
            domain, err := stream.Recv()
            if err == io.EOF {
                break
            }

            if err != nil {
                log.Fatal(err)
            }
            resultsChan <- domain.Domain
        }

    }
}

func (c *CrobatClient) ReverseDNSRange(arg string, resultsChan chan string) {
    args := ProcessArg(arg)
    defer close(resultsChan)
    for _, ipv4 := range args {
        query := &crobat.QueryRequest{
            Query: ipv4,
        }

        stream, err := c.client.ReverseDNSRange(context.Background(), query)
        if err != nil {
            log.Fatal(err)
        }

        for {
            domain, err := stream.Recv()
            if err == io.EOF {
                break
            }

            if err != nil {
                log.Fatal(err)
            }

            if domain == nil {
                continue
            }
            resultsChan <- domain.Domain
        }

    }
}

type InputType struct {
    // From https://github.com/Cgboal/SonarSearch/cmd/crobat/main.go
    DomainSub string  // "Get subdomains for this value. Supports files and quoted lists"
    DomainTLD string  // "Get tlds for this value. Supports files and quoted lists"
    ReverseDNS string // "Perform reverse lookup on IP address or CIDR range. Supports files and quoted lists"
    UniqueSort bool   // "Ensures results are unique, may cause instability on large queries due to RAM requirements"
}

func Run(ctx context.Context, in InputType, results chan string) error {
    defer close(results)
    resultsChan := make(chan string)
    var wg sync.WaitGroup

    wg.Add(1)
    go func() {
        defer wg.Done()
        uniqueSlice := UniqueStringSlice{}
        for result := range resultsChan {
            if in.UniqueSort {
                if uniqueSlice.Append(result) {
                    results <- result
                }
            } else {
                results <- result
            }
        }
    }()

    client := NewCrobatClient()
    if in.DomainSub != "" {
        client.GetSubdomains(in.DomainSub, resultsChan)
    } else if in.DomainTLD != "" {
        client.GetTlds(in.DomainTLD, resultsChan)
    } else if in.ReverseDNS != "" {
        if !strings.Contains(in.ReverseDNS, "/") {
            client.ReverseDNS(in.ReverseDNS, resultsChan)
        } else {
            client.ReverseDNSRange(in.ReverseDNS, resultsChan)
        }
    } else {
        return errors.New("must provide domain sub, TLD or reverse DNS inputs")
    }

    wg.Wait()

    return nil
}

func FetchSubDomains(ctx context.Context, in InputType) ([]string, error) {
    resultsChan := make(chan string)
    errs, ctx := errgroup.WithContext(ctx)
    var empty []string
    errs.Go(func() error {
        if err := Run(ctx, in, resultsChan); err != nil {
            return errors.WithStack(err)
        }
        return nil
    })

    var sdms []string
    for sd := range resultsChan {
        if sd != in.DomainSub {
            sdms = append(sdms, sd)
        }
    }

    if err := errs.Wait(); err != nil {
        return empty, errors.WithStack(err)
    }

    return sdms, nil
}
