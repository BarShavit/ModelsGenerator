package main

import (
	"reflect"
	"testing"
)

func Test_isList(t *testing.T) {
	type args struct {
		memberType string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 string
	}{
		{
			name:  "Bool list",
			args:  args{memberType: "list<bool>"},
			want:  true,
			want1: "bool",
		},
		{
			name:  "string list",
			args:  args{memberType: "list<string>"},
			want:  true,
			want1: "string",
		},
		{
			name:  "not list",
			args:  args{memberType: "bool"},
			want:  false,
			want1: "",
		},
		{
			name:  "list without type parameter",
			args:  args{memberType: "list"},
			want:  false,
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := isList(tt.args.memberType)
			if got != tt.want {
				t.Errorf("isList() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("isList() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_isMap(t *testing.T) {
	type args struct {
		memberType string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 string
		want2 string
	}{
		{
			name:  "Valid string key bool value map",
			args:  args{memberType: "map<string,bool>"},
			want:  true,
			want1: "string",
			want2: "bool",
		},
		{
			name:  "Valid int key string value map",
			args:  args{memberType: "map<int,string>"},
			want:  true,
			want1: "int",
			want2: "string",
		},
		{
			name:  "Not map",
			args:  args{memberType: "bool"},
			want:  false,
			want1: "",
			want2: "",
		},
		{
			name:  "Map without types",
			args:  args{memberType: "map"},
			want:  false,
			want1: "",
			want2: "",
		},
		{
			name:  "Map with one type",
			args:  args{memberType: "map<bool>"},
			want:  false,
			want1: "",
			want2: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := isMap(tt.args.memberType)
			if got != tt.want {
				t.Errorf("isMap() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("isMap() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("isMap() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_toCamelCase(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Valid Camel Test",
			args: args{value: "TestString"},
			want: "testString",
		},
		{
			name: "Valid string without need of a change",
			args: args{value: "testString"},
			want: "testString",
		},
		{
			name: "Empty string",
			args: args{value: ""},
			want: "",
		},
		{
			name: "One capital char",
			args: args{value: "T"},
			want: "t",
		},
		{
			name: "One lowercase char",
			args: args{value: "t"},
			want: "t",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toCamelCase(tt.args.value); got != tt.want {
				t.Errorf("toCamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toFirstCharUpper(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Valid first to upper test",
			args: args{value: "testString"},
			want: "TestString",
		},
		{
			name: "Valid string without need of a change",
			args: args{value: "TestString"},
			want: "TestString",
		},
		{
			name: "Empty string",
			args: args{value: ""},
			want: "",
		},
		{
			name: "One capital char",
			args: args{value: "T"},
			want: "T",
		},
		{
			name: "One lowercase char",
			args: args{value: "t"},
			want: "T",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toFirstCharUpper(tt.args.value); got != tt.want {
				t.Errorf("toFirstCharUpper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_appendUnique(t *testing.T) {
	type args struct {
		strings []string
		str     string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Unique",
			args: args{
				strings: []string{},
				str:     "new",
			},
			want: []string{"new"},
		},
		{
			name: "Not memberUnique",
			args: args{
				strings: []string{"a"},
				str:     "a",
			},
			want: []string{"a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := appendUnique(tt.args.strings, tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendUnique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memberUnique(t *testing.T) {
	type args struct {
		members []*dataMember
		member  *dataMember
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Unique member",
			args: args{
				members: []*dataMember{},
				member: &dataMember{
					memberType: "int",
					name:       "test",
				},
			},
			want: true,
		},
		{
			name: "Not memberUnique member",
			args: args{
				members: []*dataMember{
					{
						memberType: "string",
						name:       "test",
					},
				},
				member: &dataMember{
					memberType: "int",
					name:       "test",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := memberUnique(tt.args.members, tt.args.member); got != tt.want {
				t.Errorf("memberUnique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_enumValueUnique(t *testing.T) {
	type args struct {
		values []*enumValue
		value  *enumValue
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Unique enum value",
			args: args{
				values: []*enumValue{},
				value: &enumValue{
					name:  "test",
					value: 1,
				},
			},
			want: true,
		},
		{
			name: "Not unique enum value",
			args: args{
				values: []*enumValue{
					{
						name:  "test",
						value: 5,
					},
				},
				value: &enumValue{
					name:  "test",
					value: 1,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := enumValueUnique(tt.args.values, tt.args.value); got != tt.want {
				t.Errorf("enumValueUnique() = %v, want %v", got, tt.want)
			}
		})
	}
}
