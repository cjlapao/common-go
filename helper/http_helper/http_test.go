package http_helper

import "testing"

func Test_joinUrl(t *testing.T) {
	type args struct {
		element []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no elements",
			args: args{},
			want: "/",
		},
		{
			name: "mixed elements",
			args: args{
				element: []string{
					"/foo",
					"bar",
				},
			},
			want: "/foo/bar",
		},
		{
			name: "mixed elements suffix",
			args: args{
				element: []string{
					"/foo",
					"bar/",
				},
			},
			want: "/foo/bar",
		},
		{
			name: "mixed elements prefix",
			args: args{
				element: []string{
					"/foo/",
					"bar",
				},
			},
			want: "/foo/bar",
		},
		{
			name: "mixed elements prefix with all empty",
			args: args{
				element: []string{
					"",
					"",
					"",
					"",
				},
			},
			want: "/",
		},
		{
			name: "mixed elements prefix with empty middle and end",
			args: args{
				element: []string{
					"/foo/",
					"",
					"bar",
					"",
				},
			},
			want: "/foo/bar",
		},
		{
			name: "mixed elements prefix with empty start",
			args: args{
				element: []string{
					"",
					"/foo/",
					"bar",
				},
			},
			want: "/foo/bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JoinUrl(tt.args.element...); got != tt.want {
				t.Errorf("joinUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
