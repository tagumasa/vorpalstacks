// Package storage provides storage functionality for vorpalstacks.
package storage

import (
	"context"
	"fmt"
)

type twoPhaseTxn struct {
	storage    Storage
	validators []Validator
	executors  []Executor
}

// NewTwoPhaseTransaction creates a new two-phase transaction for performing
// validate-then-execute operations atomically.
func NewTwoPhaseTransaction(storage Storage) TwoPhaseTransaction {
	return &twoPhaseTxn{
		storage:    storage,
		validators: make([]Validator, 0),
		executors:  make([]Executor, 0),
	}
}

// AddValidator adds a validator to the transaction.
func (t *twoPhaseTxn) AddValidator(v Validator) {
	t.validators = append(t.validators, v)
}

// AddExecutor adds an executor to the transaction.
func (t *twoPhaseTxn) AddExecutor(e Executor) {
	t.executors = append(t.executors, e)
}

// Commit executes all validators first, then all executors in a single transaction.
func (t *twoPhaseTxn) Commit(ctx context.Context) error {
	return t.storage.Update(ctx, func(txn Transaction) error {
		for i, validator := range t.validators {
			if err := validator.Validate(ctx, txn); err != nil {
				return fmt.Errorf("validation failed at step %d: %w", i+1, err)
			}
		}

		for i, executor := range t.executors {
			if err := executor.Execute(ctx, txn); err != nil {
				return fmt.Errorf("execution failed at step %d: %w", i+1, err)
			}
		}

		return nil
	})
}

// Clear removes all validators and executors from the transaction.
func (t *twoPhaseTxn) Clear() {
	t.validators = make([]Validator, 0)
	t.executors = make([]Executor, 0)
}

// ValidatorCount returns the number of validators in the transaction.
func (t *twoPhaseTxn) ValidatorCount() int {
	return len(t.validators)
}

// ExecutorCount returns the number of executors in the transaction.
func (t *twoPhaseTxn) ExecutorCount() int {
	return len(t.executors)
}

type multiItemTransaction struct {
	storage    Storage
	operations []transactionalOperation
	conditions []transactionalCondition
}

type transactionalOperation struct {
	bucket string
	key    []byte
	opType operationType
	value  []byte
}

type transactionalCondition struct {
	bucket    string
	key       []byte
	condition Condition
}

type operationType int

const (
	opPut operationType = iota
	opDelete
)

// NewMultiItemTransaction creates a new multi-item transaction for performing
// batch operations with conditional checks across multiple items.
func NewMultiItemTransaction(storage Storage) MultiItemTransaction {
	return &multiItemTransaction{
		storage:    storage,
		operations: make([]transactionalOperation, 0),
		conditions: make([]transactionalCondition, 0),
	}
}

// Put adds a put operation to the transaction.
func (t *multiItemTransaction) Put(bucket string, key, value []byte) {
	t.operations = append(t.operations, transactionalOperation{
		bucket: bucket,
		key:    key,
		opType: opPut,
		value:  value,
	})
}

// Delete adds a delete operation to the transaction.
func (t *multiItemTransaction) Delete(bucket string, key []byte) {
	t.operations = append(t.operations, transactionalOperation{
		bucket: bucket,
		key:    key,
		opType: opDelete,
	})
}

// ConditionCheck adds a conditional check to the transaction.
// The check is performed before the operations are executed.
func (t *multiItemTransaction) ConditionCheck(bucket string, key []byte, condition Condition) {
	t.conditions = append(t.conditions, transactionalCondition{
		bucket:    bucket,
		key:       key,
		condition: condition,
	})
}

// Commit executes all condition checks and operations in a single transaction.
func (t *multiItemTransaction) Commit(ctx context.Context) error {
	twoPhase := NewTwoPhaseTransaction(t.storage).(*twoPhaseTxn)

	for _, cond := range t.conditions {
		cond := cond
		twoPhase.AddValidator(ValidatorFunc(func(ctx context.Context, txn Transaction) error {
			bucket := txn.Bucket(cond.bucket)
			current, err := bucket.Get(cond.key)
			if err != nil {
				return fmt.Errorf("condition check failed for key %s: %w", string(cond.key), err)
			}
			if !cond.condition(current) {
				return NewConditionalCheckFailedError(cond.key, current)
			}
			return nil
		}))
	}

	for _, op := range t.operations {
		op := op
		twoPhase.AddExecutor(ExecutorFunc(func(ctx context.Context, txn Transaction) error {
			bucket := txn.Bucket(op.bucket)
			switch op.opType {
			case opPut:
				return bucket.Put(op.key, op.value)
			case opDelete:
				return bucket.Delete(op.key)
			}
			return nil
		}))
	}

	return twoPhase.Commit(ctx)
}

// Clear removes all operations and conditions from the transaction.
func (t *multiItemTransaction) Clear() {
	t.operations = make([]transactionalOperation, 0)
	t.conditions = make([]transactionalCondition, 0)
}

// OperationCount returns the number of operations in the transaction.
func (t *multiItemTransaction) OperationCount() int {
	return len(t.operations)
}

// ConditionCount returns the number of conditions in the transaction.
func (t *multiItemTransaction) ConditionCount() int {
	return len(t.conditions)
}
