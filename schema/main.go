package main

import (
	"context"
	"fmt"
	"log"
	"schemaadvance/ent"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	head, _ := client.Node.Create().SetValue(0).Save(context.Background())
	curr := head
	// Generate the following linked-list: 1<->2<->3<->4<->5.
	for i := 0; i < 4; i++ {
		curr, _ = client.Node.
			Create().
			SetValue(curr.Value + 1).
			
			SetPrev(curr).
			Save(context.Background())
	}
	curr.Update().SetNextID(1).Save(context.Background())
	circular, _ := client.Node.Query().All(context.Background())
	for _, v := range circular {
		va, _ := v.QueryNext().First(context.Background())
		fmt.Println(v.ID, va.ID)
	}
}
