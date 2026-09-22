package cloudwatchlogs

import (
	"testing"
)

// Log stream names legally contain '/' and '\' (every character but ':'
// and '*'), and the stream component of the stored key is raw — a caller
// prefix carrying those separators must still match the streams it
// prefixes, including under a group name whose own component is escaped.
func TestListLogStreamsPrefixMatchesNamesWithSeparators(t *testing.T) {
	s := newLogsTestStore(t)

	for _, group := range []string{"plain", "team/app"} {
		if err := s.CreateLogGroup(&LogGroup{Name: group}); err != nil {
			t.Fatal(err)
		}
		for _, stream := range []string{"svc/api", `svc\web`, "other"} {
			if err := s.CreateLogStream(&LogStream{LogGroupName: group, Name: stream}); err != nil {
				t.Fatal(err)
			}
		}
	}

	cases := []struct {
		group  string
		prefix string
		want   []string
	}{
		{"plain", "svc/", []string{"svc/api"}},
		{"plain", "svc", []string{`svc\web`, "svc/api"}},
		{"plain", "svc\\", []string{`svc\web`}},
		{"plain", "no/x", nil},
		{"team/app", "svc/", []string{"svc/api"}},
		{"team/app", "", []string{`svc\web`, "other", "svc/api"}},
	}
	for _, tc := range cases {
		streams, _, err := s.ListLogStreams(tc.group, tc.prefix, "", 50)
		if err != nil {
			t.Fatalf("group %q prefix %q: %v", tc.group, tc.prefix, err)
		}
		var got []string
		for _, ls := range streams {
			got = append(got, ls.Name)
		}
		if len(got) != len(tc.want) {
			t.Fatalf("group %q prefix %q: got %v, want %v", tc.group, tc.prefix, got, tc.want)
		}
		for _, want := range tc.want {
			found := false
			for _, name := range got {
				if name == want {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("group %q prefix %q: got %v, want %v", tc.group, tc.prefix, got, tc.want)
			}
		}
	}
}

// Log stream names may contain '%' (every character but ':' and '*'), and
// the chunk index key escapes the stream name — a name whose literal bytes
// already look like an escape sequence ("a%2fb") would collide with the
// escaped form of a separator-bearing name ("a/b") unless '%' itself is
// escaped: the collision lets one stream's deletion remove the other's
// chunks and one stream's reads serve the other's events.
func TestStreamNamesLookingLikeEscapesStayDistinct(t *testing.T) {
	s := newLogsTestStore(t)

	const group = "g"
	if err := s.CreateLogGroup(&LogGroup{Name: group}); err != nil {
		t.Fatal(err)
	}
	for _, stream := range []string{"a/b", "a%2fb"} {
		if err := s.CreateLogStream(&LogStream{LogGroupName: group, Name: stream}); err != nil {
			t.Fatal(err)
		}
		if _, err := s.PutLogEvents(group, stream, []LogEntry{{
			Timestamp: 1700000000000,
			Message:   "event-of:" + stream,
		}}); err != nil {
			t.Fatal(err)
		}
	}

	for _, stream := range []string{"a/b", "a%2fb"} {
		events, _, _, err := s.GetLogEvents(group, stream, 0, 1700000100000, 50, true, "")
		if err != nil {
			t.Fatalf("stream %q: %v", stream, err)
		}
		if len(events) != 1 || events[0].Message != "event-of:"+stream {
			t.Fatalf("stream %q: got %v, want exactly its own event", stream, events)
		}
	}

	if err := s.DeleteLogStream(group, "a/b"); err != nil {
		t.Fatal(err)
	}
	events, _, _, err := s.GetLogEvents(group, "a%2fb", 0, 1700000100000, 50, true, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 || events[0].Message != "event-of:a%2fb" {
		t.Fatalf("after deleting a/b: got %v, want a%%2fb's own event intact", events)
	}
}

// Log-group names may contain '#', and the metric-filter key joins the
// escaped group name and the filter name with ':': the escaped group
// component contains no separator byte, so one group's scan prefix can
// never reach another group's records and no (group, filter) pair can
// alias another's exact key — including pairs whose names differ only
// through '#' ("g" + "#h" vs "g#" + "h").
func TestListMetricFiltersScopedToRequestedGroup(t *testing.T) {
	s := newLogsTestStore(t)

	for _, group := range []string{"plain", "plain#nested", "g", "g#"} {
		if err := s.CreateLogGroup(&LogGroup{Name: group}); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.PutMetricFilter(NewMetricFilter("plain", "root", "ERROR", nil)); err != nil {
		t.Fatal(err)
	}
	if err := s.PutMetricFilter(NewMetricFilter("plain#nested", "leaf", "WARN", nil)); err != nil {
		t.Fatal(err)
	}
	// The exact-key aliasing pair: under a raw '#' join these two filters
	// share one key and the second put overwrites the first.
	if err := s.PutMetricFilter(NewMetricFilter("g", "#h", "ALPHA", nil)); err != nil {
		t.Fatal(err)
	}
	if err := s.PutMetricFilter(NewMetricFilter("g#", "h", "BETA", nil)); err != nil {
		t.Fatal(err)
	}

	plain, _, err := s.ListMetricFilters("plain", "", "", 50)
	if err != nil {
		t.Fatal(err)
	}
	if len(plain) != 1 || plain[0].Name != "root" {
		t.Fatalf("group %q: got %v, want one filter %q", "plain", plain, "root")
	}

	nested, _, err := s.ListMetricFilters("plain#nested", "", "", 50)
	if err != nil {
		t.Fatal(err)
	}
	if len(nested) != 1 || nested[0].Name != "leaf" {
		t.Fatalf("group %q: got %v, want one filter %q", "plain#nested", nested, "leaf")
	}

	hashGroup, err := s.GetMetricFilter("g", "#h")
	if err != nil {
		t.Fatal(err)
	}
	if hashGroup.FilterPattern != "ALPHA" {
		t.Fatalf("filter (g, #h): got pattern %q, want its own ALPHA", hashGroup.FilterPattern)
	}
	hashTail, err := s.GetMetricFilter("g#", "h")
	if err != nil {
		t.Fatal(err)
	}
	if hashTail.FilterPattern != "BETA" {
		t.Fatalf("filter (g#, h): got pattern %q, want its own BETA", hashTail.FilterPattern)
	}
}
