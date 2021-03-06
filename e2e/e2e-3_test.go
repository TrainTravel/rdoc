// e2e tests map to the examples of the original paper
package rdoc

import (
	"fmt"
	"github.com/gpestana/rdoc"
	n "github.com/gpestana/rdoc/node"
	op "github.com/gpestana/rdoc/operation"
	"testing"
)

//Case C:  Two replicas concurrently create ordered lists under the same map keu
func TestCaseC(t *testing.T) {
	id1, id2 := "1", "2"
	doc1 := rdoc.Init(id1)
	doc2 := rdoc.Init(id2)

	// doc1 adds key "grosseries" to map at root level
	curDoc1 := op.NewCursor("groceries", op.MapKey{"groceries"})
	mutDoc1, _ := op.NewMutation(op.Noop, nil, nil)
	opList1, _ := op.New(id1+".1", []string{}, curDoc1, mutDoc1)
	doc1.ApplyOperation(*opList1)

	// doc1 adds "eggs" and "ham" entries to list
	curDoc1 = op.NewCursor("groceries", op.MapKey{"groceries"})
	mutDoc1, _ = op.NewMutation(op.Insert, 0, "eggs")
	opEggs, _ := op.New(id1+".2", []string{id1 + ".1"}, curDoc1, mutDoc1)
	_, err := doc1.ApplyOperation(*opEggs)
	if err != nil {
		t.Error(err)
	}

	curDoc1 = op.NewCursor("groceries", op.MapKey{"groceries"})
	mutDoc1, _ = op.NewMutation(op.Insert, 1, "ham")
	opHam, _ := op.New(id1+".3", []string{id1 + ".1", id1 + ".2"}, curDoc1, mutDoc1)
	_, err = doc1.ApplyOperation(*opHam)
	if err != nil {
		t.Error(err)
	}

	// doc2 adds key "grosseries" to map at root level
	curDoc2 := op.NewCursor("groceries", op.MapKey{"groceries"})
	mutDoc2, _ := op.NewMutation(op.Noop, nil, nil)
	opList2, _ := op.New(id2+".1", []string{}, curDoc2, mutDoc2)
	doc2.ApplyOperation(*opList2)

	// doc2 adds "milk" and "flour" entries to list
	curDoc2 = op.NewCursor("groceries", op.MapKey{"groceries"})
	mutDoc2, _ = op.NewMutation(op.Insert, 0, "milk")
	opMilk, _ := op.New(id2+".2", []string{id2 + ".1"}, curDoc2, mutDoc2)
	_, err = doc2.ApplyOperation(*opMilk)
	if err != nil {
		t.Error(err)
	}

	curDoc2 = op.NewCursor("groceries", op.MapKey{"groceries"})
	mutDoc2, _ = op.NewMutation(op.Insert, 1, "flour")
	opFlour, _ := op.New(id2+".3", []string{id2 + ".1", id2 + ".2"}, curDoc2, mutDoc2)
	_, err = doc2.ApplyOperation(*opFlour)
	if err != nil {
		t.Error(err)
	}

	// applies remote operations in both replicas
	doc1.ApplyRemoteOperation(*opList2)
	_, err = doc1.ApplyRemoteOperation(*opMilk)
	if err != nil {
		t.Fatal(err)
	}
	_, err = doc1.ApplyRemoteOperation(*opFlour)
	if err != nil {
		t.Fatal(err)
	}

	doc2.ApplyRemoteOperation(*opList1)
	_, err = doc2.ApplyRemoteOperation(*opEggs)
	if err != nil {
		t.Fatal(err)
	}
	_, err = doc2.ApplyRemoteOperation(*opHam)
	if err != nil {
		t.Fatal(err)
	}

	// verifications
	doc1If, _ := doc1.Head.Map().Get("groceries")
	doc1Groceries := doc1If.(*n.Node).List()
	if doc1Groceries.Size() != 4 {
		t.Error(fmt.Sprintf("Doc1 grosseries list should have 4 items after applying remote operations, got %v", doc1Groceries.Size()))
	}

	doc2If, _ := doc2.Head.Map().Get("groceries")
	doc2Groceries := doc2If.(*n.Node).List()
	if doc2Groceries.Size() != 4 {
		t.Error(fmt.Sprintf("Doc2 grosseries list should have 4 items after applying remote operations, got %v", doc2Groceries.Size()))
	}

	if doc1Groceries.Size() != doc2Groceries.Size() {
		t.Error(fmt.Sprintf("List must have same number of elements in doc1 (got: %v) and doc2 (got: %v)", doc1Groceries.Size(), doc2Groceries.Size()))
	}

	// compares list elements order
	var list1Keys string
	var list2Keys string
	for i := 0; i < doc1Groceries.Size(); i++ {
		el1, _ := doc1Groceries.Get(i)
		el2, _ := doc2Groceries.Get(i)
		el1Keys := el1.(*n.Node).Reg().Keys()
		el2Keys := el2.(*n.Node).Reg().Keys()
		list1Keys = fmt.Sprintf("%v %v", list1Keys, el1Keys[0].(string))
		list2Keys = fmt.Sprintf("%v %v", list2Keys, el2Keys[0].(string))
	}
	if list1Keys != list2Keys {
		t.Error(fmt.Sprintf("Final list state did not converge in both replicas: %v != %v", list1Keys, list2Keys))
	}
}
