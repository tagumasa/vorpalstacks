package iot

import "testing"

// Deleting a job must also remove its jobExecution/* records; executions of
// other jobs must survive.
func TestDeleteJobCascadesExecutions(t *testing.T) {
	st := newIotStore(t)

	if _, err := st.CreateJob(&Job{JobID: "job-1"}); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateJob(&Job{JobID: "job-2"}); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{
		"jobExecution/job-1/thing-a",
		"jobExecution/job-1/thing-b",
		"jobExecution/job-2/thing-a",
	} {
		if err := st.PutGeneric(key, map[string]interface{}{"status": "QUEUED"}); err != nil {
			t.Fatal(err)
		}
	}

	if err := st.DeleteJob("job-1"); err != nil {
		t.Fatal(err)
	}

	remaining, err := st.ListGeneric("jobExecution/job-1/")
	if err != nil {
		t.Fatal(err)
	}
	if len(remaining) != 0 {
		t.Fatalf("executions of the deleted job survived: %v", remaining)
	}

	other, err := st.ListGeneric("jobExecution/job-2/")
	if err != nil {
		t.Fatal(err)
	}
	if len(other) != 1 {
		t.Fatalf("executions of job-2 disturbed: %v", other)
	}

	if _, err := st.GetJob("job-1"); err == nil {
		t.Fatal("job itself was not deleted")
	}
}
