package main

import (
	"context"
	"fmt"
	"net/http"

	"achapromo.com/sirv/sirv"
)

func main() {

	c := &http.Client{}
	s := sirv.NewClient(c, sirv.FreePlan)

	p := sirv.AuthPayload{
		ClientId:     "",
		ClientSecret: "",
	}

	tr, err := s.GetToken(context.Background(), p)
	if err != nil {
		panic(fmt.Sprintf("\n%s\n%+v\n", "GetToken", err))

	}
	fmt.Printf("\nGetToken:\n%+v\n\n", tr)

	ai, err := s.GetAccountInfo(context.Background())
	if err != nil {
		panic(fmt.Sprintf("\n%s\n%+v\n", "GetAccountInfo", err))
	}
	fmt.Printf("\nGetAccountInfo:\n%+v\n\n", ai)

	al, err := s.GetAPILimits(context.Background())
	if err != nil {
		panic(fmt.Sprintf("\n%s\n%+v\n", "GetAPILimits", err))
	}
	fmt.Printf("\nGetAPILimits:\n%+v\n\n", al)

	si, err := s.GetStorageInfo(context.Background())
	if err != nil {
		panic(fmt.Sprintf("\n%s\n%+v\n", "GetStorageInfo", err))
	}
	fmt.Printf("\nGetStorageInfo:\n%+v\n\n", si)

	fsp := sirv.FileSearchPayload{
		Query:  "extension:.jpg",
		Scroll: true,
	}
	printAllFiles(s, fsp)

	x, err := s.ReadFolderContents(context.Background(), "/", "")
	if err != nil {
		panic(fmt.Sprintf("\n%s\n%+v\n", "ReadFolderContents", err))
	}

	for _, f := range x.Contents {
		fmt.Printf("%+v\n", f.Filename)
	}
}

func printAllFiles(s *sirv.Client, fsp sirv.FileSearchPayload) error {
	// initial search
	fsr, err := s.SearchFiles(context.Background(), fsp)
	if err != nil {
		return err
	}

	// print files
	for _, f := range fsr.Hits {
		fmt.Printf("%+v --- %+v\n\n", f.Source.Dirname, f.Source.Filename)
	}

	// break if no more files
	if len(fsr.Hits) == 0 {
		return nil
	}

	for {

		// scroll for more files
		fss := sirv.FileSearchScrollPayload{ScrollId: fsr.ScrollId}
		fsr, err := s.ScrollFilesSearch(context.Background(), fss)
		if err != nil {
			return err
		}

		// print files
		for _, f := range fsr.Hits {
			fmt.Printf("%+v --- %+v\n\n", f.Source.Dirname, f.Source.Filename)
		}

		// break if no more files
		if len(fsr.Hits) == 0 {
			break
		}

	}

	return nil
}
