package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
)

func (tc *athenaTestContext) testCapacityReservations() []TestResult {
	var results []TestResult
	client := tc.client
	ctx := tc.ctx

	// Per the Smithy model, CreateCapacityReservation declares only
	// InternalServerException and InvalidRequestException — a duplicate
	// name must be rejected with InvalidRequestException.
	results = append(results, tc.runner.RunTest("athena", "CreateCapacityReservation_DuplicateRejected", func() error {
		crName := tc.uniqueName("dup-cr")
		_, err := client.CreateCapacityReservation(ctx, &athena.CreateCapacityReservationInput{
			Name:       aws.String(crName),
			TargetDpus: aws.Int32(24),
		})
		if err != nil {
			return fmt.Errorf("first create: %w", err)
		}
		defer func() {
			// Delete requires the CANCELLED state, so cancel first.
			_, _ = client.CancelCapacityReservation(ctx, &athena.CancelCapacityReservationInput{Name: aws.String(crName)})
			_, _ = client.DeleteCapacityReservation(ctx, &athena.DeleteCapacityReservationInput{Name: aws.String(crName)})
		}()

		_, err = client.CreateCapacityReservation(ctx, &athena.CreateCapacityReservationInput{
			Name:       aws.String(crName),
			TargetDpus: aws.Int32(48),
		})
		if err := AssertErrorContains(err, "InvalidRequestException"); err != nil {
			return fmt.Errorf("expected InvalidRequestException for duplicate capacity reservation, got: %v", err)
		}
		return nil
	}))

	return results
}
