// e2e tests map to the examples of the original paper
package rdoc

import (
	op "github.com/gpestana/rdoc/operation"
	"testing"
)

// Case A: different value assignment of a register in different replicas
func TestCaseA(t *testing.T) {
	id1 := "1"
	doc1 := Init(id1)

	id2 := "2"
	doc2 := Init(id2)
	emptyC := op.NewEmptyCursor()

	// contructs operation to initially populate the docs
	mut0, _ := op.NewMutation(op.Assign, "key", "A")
	op0, _ := op.New(id1+".0", []string{}, emptyC, mut0) // using id1 means that the operation was generated by id1

	doc1.ApplyOperation(*op0)
	doc2.ApplyRemoteOperation(*op0)

	// constructs and applies locally operation from replica 1
	mut1, _ := op.NewMutation(op.Assign, "key", "B")
	op1, _ := op.New(id1+".1", []string{}, emptyC, mut1)

	doc1.ApplyOperation(*op1)

	// constructs and applies locally operation for replica 2
	mut2, _ := op.NewMutation(op.Assign, "key", "C")
	op2, _ := op.New(id2+".1", []string{}, emptyC, mut2)

	doc2.ApplyOperation(*op2)

	// network communication: cross-apply operations in replica 1 and 2
	doc1.ApplyRemoteOperation(*op2)
	doc2.ApplyRemoteOperation(*op1)

	// verify result of the merging
	val1, exist1 := doc1.Head.hmap.Get("key")
	val2, exist2 := doc2.Head.hmap.Get("key")

	if exist1 == false || exist2 == false {
		t.Error("expected keys do not exist: ", val1, val2)
	}

}
