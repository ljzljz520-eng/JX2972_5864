package main

import (
	"emergency-claim-code/internal/httpapi"
	"emergency-claim-code/internal/model"
	"emergency-claim-code/internal/repository"
	"emergency-claim-code/internal/service"
	"emergency-claim-code/internal/store"
	"emergency-claim-code/internal/workflow"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	path := flag.String("db", "claims.db", "database path")
	listen := flag.String("listen", "", "http listen address")
	flag.Parse()
	if err := store.EnsureParent(*path); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	db, err := store.Open(*path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()
	repo := repository.New(db)
	svc := service.New(repo)
	engine := workflow.New(svc, repo)
	if *listen != "" {
		fmt.Printf("claim-code server listening on %s\n", *listen)
		_ = http.ListenAndServe(*listen, httpapi.New(engine, svc).Handler())
		return
	}
	if err := runCommand(flag.Args(), engine, svc); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runCommand(args []string, engine *workflow.Engine, svc *service.Service) error {
	if len(args) == 0 {
		fmt.Println("claimctl ready: use create, submit, approve, archive, get, list, recall-retry, serve")
		return nil
	}
	switch args[0] {
	case "create":
		return commandCreate(args[1:], engine)
	case "submit":
		return commandStatus(args[1:], svc, "submitted")
	case "approve":
		return commandStatus(args[1:], svc, "approved")
	case "archive":
		return commandStatus(args[1:], svc, "archived")
	case "get":
		return commandGet(args[1:], svc)
	case "list":
		return commandList(args[1:], svc)
	case "recall-retry":
		return commandRecall(args[1:], engine)
	case "serve":
		return fmt.Errorf("use -listen address")
	default:
		return fmt.Errorf("unknown command %s", args[0])
	}
}

func commandCreate(args []string, engine *workflow.Engine) error {
	if len(args) < 3 {
		return fmt.Errorf("create requires batch applicant quantity")
	}
	var quantity int
	if _, err := fmt.Sscanf(args[2], "%d", &quantity); err != nil {
		return err
	}
	item, err := engine.Register(structRecord(args[0], args[1], quantity), "cli", "cli")
	if err != nil {
		return err
	}
	return printJSON(item)
}

func structRecord(batch, applicant string, quantity int) model.Record {
	return model.Record{BatchID: batch, Applicant: applicant, Quantity: quantity}
}

func commandStatus(args []string, svc *service.Service, target string) error {
	if len(args) < 1 {
		return fmt.Errorf("id required")
	}
	item, err := svc.ChangeStatus(args[0], target, "cli", "cli")
	if err != nil {
		return err
	}
	return printJSON(item)
}

func commandGet(args []string, svc *service.Service) error {
	if len(args) < 1 {
		return fmt.Errorf("id required")
	}
	item, err := svc.GetRecord(args[0])
	if err != nil {
		return err
	}
	return printJSON(item)
}

func commandList(args []string, svc *service.Service) error {
	term := ""
	if len(args) > 0 {
		term = args[0]
	}
	items, err := svc.ListRecords(term, "")
	if err != nil {
		return err
	}
	return printJSON(items)
}

func commandRecall(args []string, engine *workflow.Engine) error {
	if len(args) < 1 {
		return fmt.Errorf("id required")
	}
	item, err := engine.RecallRetry(args[0], "cli", "cli")
	if err != nil {
		return err
	}
	return printJSON(item)
}

func printJSON(value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err == nil {
		fmt.Println(string(data))
	}
	return err
}
