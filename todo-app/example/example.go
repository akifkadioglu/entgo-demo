package example

import (
	"context"
	"fmt"
	"log"
	"todo/ent"
	"todo/ent/todo"

	"entgo.io/ent/dialect"
	_ "github.com/mattn/go-sqlite3"
)

func Example_Todo() {
	// Create an ent.Client with in-memory SQLite database.
	client, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	ctx := context.Background()
	// Run the automatic migration tool to create all schema resources.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	task1, err := client.Todo.Create().SetText("Nabersin").Save(ctx)
	if err != nil {
		log.Fatalf("failed creating a todo: %v", err)
	}
	task2, err := client.Todo.Create().SetText("Iyiyim").Save(ctx)
	if err != nil {
		log.Fatalf("failed creating a todo: %v", err)
	}
	fmt.Println(task1)
	fmt.Println(task2)
	if err := task2.Update().SetParent(task1).Exec(ctx); err != nil {
		log.Fatalf("failed creating a todo: %v", err)
	}
	fmt.Println(task2)
	items, err := client.Todo.Query().All(ctx)
	if err != nil {
		log.Fatalf("failed querying todos: %v", err)
	}
	for _, t := range items {
		fmt.Printf("%d: %q\n", t.ID, t.Text)
	}

	// Query all todo items that don't depend on other items and have items that depend them.
	items, err = client.Todo.Query().
		Where(
			todo.Not(
				todo.HasParent(),
			),
			todo.HasChildren(),
		).
		All(ctx)
	if err != nil {
		log.Fatalf("failed querying todos: %v", err)
	}
	for _, t := range items {
		fmt.Println(t)
	}

	// Get a parent item through its children and expect the
	// query to return exactly one item.
	parent, err := client.Todo.Query(). // Query all todos.
						Where(todo.HasParent()). // Filter only those with parents.
						QueryParent().           // Continue traversals to the parents.
						Only(ctx)                // Expect exactly one item.
	if err != nil {
		log.Fatalf("failed querying todos: %v", err)
	}
	fmt.Printf("%d: %q\n", parent.ID, parent.Text)
}
