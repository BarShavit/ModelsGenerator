package main

import (
	"reflect"
	"testing"
)

func getValidParse() (string, string, string, []middleware) {
	validContent := "class someclass\n{\nsomemember string\nanother int\n}\n" +
		"enum someEnum\n{\nfirst 5\nsecond 8\n}"

	spacedValidContent := "class    someclass\n{\nsomemember    string\n     another int\n}\n" +
		"enum someEnum\n{\nfirst 5\nsecond 8\n}"

	validWithEmptyLinesContent := "class someclass\n\n{\nsomemember string\n\n\nanother int\n}\n" +
		"enum someEnum\n{\nfirst 5\nsecond 8\n}"

	expectedValidContent := []middleware{
		&class{
			name: "someclass",
			dataMembers: []*dataMember{
				{
					memberType: "string",
					name:       "somemember",
				},
				{
					memberType: "int",
					name:       "another",
				},
			},
		},
		&enum{
			name: "someEnum",
			enumValues: []*enumValue{
				{
					name:  "first",
					value: 5,
				},
				{
					name:  "second",
					value: 8,
				},
			},
		},
	}

	return validContent, spacedValidContent, validWithEmptyLinesContent, expectedValidContent
}

func Test_parse(t *testing.T) {
	validContent, spacedContent, withEmptyLinesContent, expectedMeddlers := getValidParse()

	type args struct {
		fileContent string
	}
	tests := []struct {
		name    string
		args    args
		want    []middleware
		wantErr bool
	}{
		{
			name:    "Valid content",
			args:    struct{ fileContent string }{fileContent: validContent},
			want:    expectedMeddlers,
			wantErr: false,
		},
		{
			name: "Without declare class row",
			args: args{
				fileContent: "{\nsomemember int\n}",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Without { row",
			args: args{
				fileContent: "class test\nsomemember int\n}",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Without } closer",
			args: args{
				fileContent: "class test\n{\nsomemember int",
			},
			want:    make([]middleware, 0),
			wantErr: false,
		},
		{
			name:    "Valid spaced content",
			args:    args{fileContent: spacedContent},
			want:    expectedMeddlers,
			wantErr: false,
		},
		{
			name:    "With empty lines",
			args:    struct{ fileContent string }{fileContent: withEmptyLinesContent},
			want:    expectedMeddlers,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse(tt.args.fileContent)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readMiddlewareDeclare(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    middleware
		wantErr bool
	}{
		{name: "Creating valid class",
			args:    args{line: "class bla"},
			want:    newClass("bla"),
			wantErr: false},
		{name: "Creating valid enum",
			args:    args{line: "enum bla"},
			want:    newEnum("bla"),
			wantErr: false},
		{name: "Empty line",
			args:    args{line: ""},
			want:    nil,
			wantErr: true},
		{name: "Class without name",
			args:    args{line: "class"},
			want:    nil,
			wantErr: true},
		{name: "Enum without name",
			args:    args{line: "enum"},
			want:    nil,
			wantErr: true},
		{name: "Ignore data member",
			args:    args{line: "int bla"},
			want:    nil,
			wantErr: true},
		{name: "Ignore enum value",
			args:    args{line: "SOME_VALUE 5"},
			want:    nil,
			wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readMiddlewareDeclare(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("readMiddlewareDeclare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readMiddlewareDeclare() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readMiddlewareValue(t *testing.T) {
	classArg := newClass("testClass")
	enumArg := newEnum("testEnum")

	type args struct {
		middleware middleware
		line       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Valid data member in class",
			args:    args{middleware: classArg, line: "test int"},
			wantErr: false},

		{name: "Valid enum value",
			args:    args{middleware: enumArg, line: "test 5"},
			wantErr: false},

		{name: "Invalid enum value",
			args:    args{middleware: enumArg, line: "test fdgdf"},
			wantErr: true},

		{name: "Invalid line",
			args:    args{middleware: enumArg, line: "test"},
			wantErr: true},

		{name: "Empty line",
			args:    args{middleware: enumArg, line: ""},
			wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := readMiddlewareValue(tt.args.middleware, tt.args.line); (err != nil) != tt.wantErr {
				t.Errorf("readMiddlewareValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	expectedClass := &class{
		name: "testClass",
		dataMembers: []*dataMember{
			{
				memberType: "int",
				name:       "test",
			},
		},
	}

	expectedEnum := &enum{
		name: "testEnum",
		enumValues: []*enumValue{
			{
				name:  "test",
				value: 5,
			},
		},
	}

	if !reflect.DeepEqual(classArg, expectedClass) {
		t.Errorf("readMiddlewareDeclare() got unexpected class data members."+
			" result got %v members and we expected %v",
			classArg.dataMembers[0].name, expectedClass.dataMembers[0].name)
	}

	if !reflect.DeepEqual(enumArg, expectedEnum) {
		t.Errorf("readMiddlewareDeclare() got unexpected class data members")
	}
}

func Test_trimContent(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Spaces from all sides", args: args{content: "   word  word  "}, want: "word word"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trimContent(tt.args.content); got != tt.want {
				t.Errorf("trimContent() = %v, want %v", got, tt.want)
			}
		})
	}
}
