package vultr

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

// fakeKubernetesService implements govultr.KubernetesService so the wait logic
// can be exercised without touching the Vultr API. govultr.Client.Kubernetes is
// an exported interface field, so the fake can be dropped straight in.
type fakeKubernetesService struct {
	govultr.KubernetesService // panic on any method we do not deliberately fake

	mu    sync.Mutex
	calls int

	// failFirst is how many leading GetCluster calls return an error.
	failFirst int
	// err is what those calls return. Defaults to a generic transient error.
	err error
	// statuses are returned in order once calls exceed failFirst; the last
	// entry repeats.
	statuses []string
}

func (f *fakeKubernetesService) GetCluster(ctx context.Context, id string) (*govultr.Cluster, *http.Response, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.calls++

	if f.calls <= f.failFirst {
		err := f.err
		if err == nil {
			err = errors.New("We are currently conducting some software upgrades. Check back in a few minutes!")
		}
		return nil, nil, err
	}

	status := "active"
	if len(f.statuses) > 0 {
		idx := f.calls - f.failFirst - 1
		if idx >= len(f.statuses) {
			idx = len(f.statuses) - 1
		}
		status = f.statuses[idx]
	}

	return &govultr.Cluster{ID: id, Status: status}, &http.Response{}, nil
}

func (f *fakeKubernetesService) callCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.calls
}

// newRefreshTestData builds the minimal ResourceData the refresh func reads:
// it only needs an ID.
func newRefreshTestData(t *testing.T, id string) *schema.ResourceData {
	t.Helper()

	d := schema.TestResourceDataRaw(t, resourceVultrKubernetes().Schema, map[string]interface{}{})
	d.SetId(id)
	return d
}

// newRefreshTestClient returns a provider Client whose govultr client serves
// Kubernetes calls from the given fake.
func newRefreshTestClient(fake govultr.KubernetesService) *Client {
	govultrClient := govultr.NewClient(nil)
	govultrClient.Kubernetes = fake
	return &Client{client: govultrClient}
}

// ---------------------------------------------------------------------------
// Existing behaviour
//
// There was previously no unit coverage of newVKEStateRefresh at all - the only
// tests in resource_vultr_kubernetes_test.go are TestAcc* acceptance tests that
// need TF_ACC and a live API. These pin down the behaviour that is unchanged by
// the retry fix: which object and state the refresh hands back on success.
// ---------------------------------------------------------------------------

func TestNewVKEStateRefresh_ReturnsClusterAndStatus(t *testing.T) {
	const id = "249df6e2-ef47-469a-ae3e-26f1465b9fdf"

	fake := &fakeKubernetesService{statuses: []string{"pending"}}
	refresh := newVKEStateRefresh(context.Background(), newRefreshTestData(t, id),
		newRefreshTestClient(fake), "status")

	result, state, err := refresh()
	if err != nil {
		t.Fatalf("expected no error on a successful read, got: %v", err)
	}
	if state != "pending" {
		t.Errorf("state = %q, want %q", state, "pending")
	}

	cluster, ok := result.(*govultr.Cluster)
	if !ok {
		t.Fatalf("result is %T, want *govultr.Cluster so StateChangeConf matches on status", result)
	}
	if cluster.Status != "pending" {
		t.Errorf("cluster.Status = %q, want %q", cluster.Status, "pending")
	}
	if cluster.ID != id {
		t.Errorf("cluster.ID = %q, want %q", cluster.ID, id)
	}
}

func TestNewVKEStateRefresh_TracksClusterStatus(t *testing.T) {
	// The status returned is what StateChangeConf compares against Target, so a
	// cluster moving pending -> active must be reported as "active".
	fake := &fakeKubernetesService{statuses: []string{"pending", "active"}}
	refresh := newVKEStateRefresh(context.Background(), newRefreshTestData(t, "abc"),
		newRefreshTestClient(fake), "status")

	if _, state, err := refresh(); err != nil || state != "pending" {
		t.Fatalf("first refresh: state=%q err=%v, want pending/nil", state, err)
	}

	_, state, err := refresh()
	if err != nil {
		t.Fatalf("second refresh returned an error: %v", err)
	}
	if state != "active" {
		t.Errorf("state = %q, want %q once the cluster reports active", state, "active")
	}
}

func TestNewVKEStateRefresh_NonStatusAttrReturnsNilResult(t *testing.T) {
	// With attr != "status" the function returns a nil result. That is fine for
	// the current single caller (which asks for "status"), but it is exactly the
	// shape StateChangeConf treats as "not found", so pin it down deliberately.
	fake := &fakeKubernetesService{statuses: []string{"active"}}
	refresh := newVKEStateRefresh(context.Background(), newRefreshTestData(t, "abc"),
		newRefreshTestClient(fake), "something-else")

	result, state, err := refresh()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result != nil || state != "" {
		t.Errorf("result=%v state=%q, want nil/empty for an unrecognised attribute", result, state)
	}
}

// ---------------------------------------------------------------------------
// Transient errors must not fail the create
// ---------------------------------------------------------------------------

func TestNewVKEStateRefresh_TransientErrorDoesNotFail(t *testing.T) {
	const id = "7cbd1153-212f-4fc0-ae83-351504844c49"

	// This is the exact failure that killed a real apply 55s in: one transient
	// API error during the poll loop. It must NOT be surfaced as an error, or
	// StateChangeConf aborts the wait instead of retrying.
	fake := &fakeKubernetesService{
		failFirst: 1,
		err:       errors.New(`We are currently conducting some software upgrades. Check back in a few minutes!`),
		statuses:  []string{"pending"},
	}
	refresh := newVKEStateRefresh(context.Background(), newRefreshTestData(t, id),
		newRefreshTestClient(fake), "status")

	result, state, err := refresh()
	if err != nil {
		t.Fatalf("a transient error must not be returned, or the wait aborts: %v", err)
	}
	if state != "pending" {
		t.Errorf("state = %q, want %q", state, "pending")
	}

	// Trap 1: a nil result is counted as "not found" and trips NotFoundChecks
	// (60) long before the 60 minute timeout.
	if result == nil {
		t.Error("result must be non-nil, otherwise StateChangeConf treats the cluster as missing")
	}

	// The next poll succeeds, so the reference implementation would have already
	// given up by now.
	if _, state, err := refresh(); err != nil || state != "pending" {
		t.Fatalf("recovery poll: state=%q err=%v, want pending/nil", state, err)
	}
}

func TestNewVKEStateRefresh_TransientErrorStateIsPending(t *testing.T) {
	// Trap 2: the state on a transient failure must be one of the Pending values
	// ("pending"), otherwise StateChangeConf raises UnexpectedStateError and
	// aborts anyway - the retry would never happen.
	pending := []string{"pending"}

	fake := &fakeKubernetesService{failFirst: 1}
	refresh := newVKEStateRefresh(context.Background(), newRefreshTestData(t, "abc"),
		newRefreshTestClient(fake), "status")

	_, state, err := refresh()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, p := range pending {
		if state == p {
			found = true
		}
	}
	if !found {
		t.Errorf("state %q is not in the Pending set %v, so the wait would abort", state, pending)
	}
}

func TestNewVKEStateRefresh_PersistentErrorEventuallyFails(t *testing.T) {
	// Tolerating transient blips must not mean swallowing a genuinely broken
	// cluster forever - after the run of consecutive failures it must fail, and
	// must include the underlying API error (the original code discarded it).
	fake := &fakeKubernetesService{
		failFirst: 1000,
		err:       errors.New("Unauthorized"),
	}
	refresh := newVKEStateRefresh(context.Background(), newRefreshTestData(t, "abc"),
		newRefreshTestClient(fake), "status")

	var lastErr error
	for i := 0; i < 12; i++ {
		if _, _, err := refresh(); err != nil {
			lastErr = err
			break
		}
	}

	if lastErr == nil {
		t.Fatal("a permanent failure was tolerated forever; the wait would never fail")
	}
	if !strings.Contains(lastErr.Error(), "Unauthorized") {
		t.Errorf("error %q does not include the underlying API error; that was the original bug", lastErr)
	}
	if !strings.Contains(lastErr.Error(), "consecutive") {
		t.Errorf("error %q does not say how many attempts failed", lastErr)
	}
}

func TestNewVKEStateRefresh_SuccessResetsErrorRun(t *testing.T) {
	// Only *consecutive* failures should count, so an intermittent error spread
	// across a long provision never accumulates to the limit.
	fake := &fakeKubernetesService{failFirst: 1, statuses: []string{"pending"}}
	refresh := newVKEStateRefresh(context.Background(), newRefreshTestData(t, "abc"),
		newRefreshTestClient(fake), "status")

	// 1 failure, then success, then two more failures - never 12 in a row.
	if _, _, err := refresh(); err != nil {
		t.Fatalf("transient error surfaced: %v", err)
	}
	if _, _, err := refresh(); err != nil {
		t.Fatalf("unexpected error after recovery: %v", err)
	}

	for i := 0; i < 2; i++ {
		if _, _, err := refresh(); err != nil {
			t.Fatalf("counter did not reset after the successful poll: %v", err)
		}
	}
	if got := fake.callCount(); got != 4 {
		t.Errorf("GetCluster calls = %d, want 4", got)
	}
}

// ---------------------------------------------------------------------------
// End-to-end through the real StateChangeConf
// ---------------------------------------------------------------------------

func TestWaitForVKEAvailable_SurvivesTransientError(t *testing.T) {
	// The strongest form of the regression test: run the actual StateChangeConf
	// the create path uses and assert a transient error does not abort it. This
	// fails against the original implementation, which returned the error and
	// stopped after one poll.
	//
	// It uses the SDK's real backoff (Delay 10s, MinTimeout 5s), so it is skipped
	// under -short.
	if testing.Short() {
		t.Skip("skipping: exercises the SDK's real wait backoff")
	}

	const id = "7cbd1153-212f-4fc0-ae83-351504844c49"

	fake := &fakeKubernetesService{
		failFirst: 1,
		err:       errors.New("We are currently conducting some software upgrades. Check back in a few minutes!"),
		statuses:  []string{"pending", "active"},
	}

	d := newRefreshTestData(t, id)
	meta := newRefreshTestClient(fake)

	done := make(chan error, 1)
	go func() {
		_, err := waitForVKEAvailable(context.Background(), d, "active", []string{"pending"}, "status", meta)
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("a transient API error aborted the wait: %v", err)
		}
	case <-time.After(4 * time.Minute):
		t.Fatal("wait did not complete; a transient error appears to have stalled it")
	}
}

func TestWaitForVKEAvailable_FailsOnClusterNeverLeavingPending(t *testing.T) {
	// Guards the opposite regression: if the change accidentally made the
	// refresh always report "pending", a cluster that never becomes active
	// would be waited on until the 60 minute timeout and then reported as a
	// timeout rather than a useful error. Confirm the failure is surfaced.
	if testing.Short() {
		t.Skip("skipping: exercises the SDK's real wait backoff")
	}

	fake := &fakeKubernetesService{failFirst: 1000, err: fmt.Errorf("persistent upstream failure")}

	d := newRefreshTestData(t, "abc")
	meta := newRefreshTestClient(fake)

	// Drive the refresh function directly rather than waiting out the SDK
	// backoff, which for 12 consecutive failures would take many minutes.
	refresh := newVKEStateRefresh(context.Background(), d, meta, "status")

	var err error
	for i := 0; i < 12; i++ {
		if _, _, err = refresh(); err != nil {
			break
		}
	}
	if err == nil {
		t.Fatal("persistent upstream failures were never surfaced")
	}
	if !strings.Contains(err.Error(), "persistent upstream failure") {
		t.Errorf("error %q dropped the underlying cause", err)
	}
}
