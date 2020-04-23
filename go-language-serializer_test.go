package main

import (
	"reflect"
	"strings"
	"testing"
)

func Test_goLanguageSerializer_generateCode(t *testing.T) {
	g := newGoLanguageSerializer()

	testClass := &class{
		name: "test",
		dataMembers: []*dataMember{
			{
				memberType: "a",
				name:       "int",
			},
		},
	}

	testEnum := &enum{
		name: "test",
		enumValues: []*enumValue{
			{
				name:  "a",
				value: 5,
			},
			{
				name:  "b",
				value: 8,
			},
		},
	}

	meddlers := []middleware{testClass, testEnum}

	generatedCode, err := g.generateCode(meddlers, &serializerInfo{packageName: "bla"})
	if err != nil {
		t.Errorf("generateCode() error = %v, wantErr %v", err, false)
		return
	}

	if len(generatedCode) != 2 {
		t.Errorf("generateCode() generated %v files. expected 2", len(generatedCode))
		return
	}
}

func Test_goLanguageSerializer_getType(t *testing.T) {
	type fields struct {
		typesMap map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   languageType
	}{
		{name: "Valid type", fields: fields{typesMap: map[string]string{}}, want: LanguageTypeGo},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &goLanguageSerializer{
				typesMap: tt.fields.typesMap,
			}
			if got := g.getType(); got != tt.want {
				t.Errorf("getType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_goLanguageSerializer_getTypeName(t *testing.T) {
	type fields struct {
		typesMap map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "Valid type name", fields: fields{typesMap: map[string]string{}}, want: "go"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &goLanguageSerializer{
				typesMap: tt.fields.typesMap,
			}
			if got := g.getTypeName(); got != tt.want {
				t.Errorf("getTypeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_goLanguageSerializer_serializeClass(t *testing.T) {
	type args struct {
		class          *class
		serializerInfo *serializerInfo
	}
	tests := []struct {
		name    string
		args    args
		want    *generatedCode
		wantErr bool
	}{
		{
			name: "Primitive class serialize",
			args: args{
				class: &class{
					name: "test",
					dataMembers: []*dataMember{
						{
							memberType: "string",
							name:       "a",
						},
						{
							memberType: "double",
							name:       "b",
						},
					},
				},
				serializerInfo: &serializerInfo{packageName: "bla"},
			},
			want: &generatedCode{
				fileName: "test.go",
				code:     "type Test struct {\n\tA string `json:\"a\"`\n\tB float64 `json:\"b\"`\n}",
			},
			wantErr: false,
		},
		{
			name: "Class with list serialize",
			args: args{
				class: &class{
					name: "test",
					dataMembers: []*dataMember{
						{
							memberType: "string",
							name:       "a",
						},
						{
							memberType: "list<int>",
							name:       "b",
						},
					},
				},
				serializerInfo: &serializerInfo{packageName: "bla"},
			},
			want: &generatedCode{
				fileName: "test.go",
				code:     "type Test struct {\n\tA string `json:\"a\"`\n\tB []int `json:\"b\"`\n}",
			},
			wantErr: false,
		},
		{
			name: "Class with list of struct serialize",
			args: args{
				class: &class{
					name: "test",
					dataMembers: []*dataMember{
						{
							memberType: "string",
							name:       "a",
						},
						{
							memberType: "list<Bla>",
							name:       "b",
						},
					},
				},
				serializerInfo: &serializerInfo{packageName: "bla"},
			},
			want: &generatedCode{
				fileName: "test.go",
				code:     "type Test struct {\n\tA string `json:\"a\"`\n\tB []*Bla `json:\"b\"`\n}",
			},
			wantErr: false,
		},
		{
			name: "Class with basic map",
			args: args{
				class: &class{
					name: "test",
					dataMembers: []*dataMember{
						{
							memberType: "string",
							name:       "a",
						},
						{
							memberType: "map<int,string>",
							name:       "b",
						},
					},
				},
				serializerInfo: &serializerInfo{packageName: "bla"},
			},
			want: &generatedCode{
				fileName: "test.go",
				code:     "type Test struct {\n\tA string `json:\"a\"`\n\tB map[int]string `json:\"b\"`\n}",
			},
			wantErr: false,
		},
		{
			name: "Class with struct map",
			args: args{
				class: &class{
					name: "test",
					dataMembers: []*dataMember{
						{
							memberType: "string",
							name:       "a",
						},
						{
							memberType: "map<int,Bla>",
							name:       "b",
						},
					},
				},
				serializerInfo: &serializerInfo{packageName: "bla"},
			},
			want: &generatedCode{
				fileName: "test.go",
				code:     "type Test struct {\n\tA string `json:\"a\"`\n\tB map[int]*Bla `json:\"b\"`\n}",
			},
			wantErr: false,
		},
		{
			name: "Class with struct",
			args: args{
				class: &class{
					name: "test",
					dataMembers: []*dataMember{
						{
							memberType: "string",
							name:       "a",
						},
						{
							memberType: "Bla",
							name:       "b",
						},
					},
				},
				serializerInfo: &serializerInfo{packageName: "bla"},
			},
			want: &generatedCode{
				fileName: "test.go",
				code:     "type Test struct {\n\tA string `json:\"a\"`\n\tB *Bla `json:\"b\"`\n}",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newGoLanguageSerializer()
			got, err := g.serializeClass(tt.args.class, tt.args.serializerInfo)

			if got != nil {
				// delete the header
				got.code = strings.Replace(got.code, g.serializeDeclaration(tt.args.serializerInfo), "", -1)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("serializeClass() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("serializeClass() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_goLanguageSerializer_serializeDeclaration(t *testing.T) {
	type fields struct {
		typesMap map[string]string
	}
	type args struct {
		serializerInfo *serializerInfo
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          string
		shouldContain bool
	}{
		{
			name: "With package",
			fields: fields{
				typesMap: map[string]string{},
			},
			args: args{serializerInfo: &serializerInfo{
				packageName: "bla",
			}},
			want:          "package bla",
			shouldContain: true,
		},
		{
			name: "Without package",
			fields: fields{
				typesMap: map[string]string{},
			},
			args: args{serializerInfo: &serializerInfo{
				packageName: "",
			}},
			want:          "package",
			shouldContain: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &goLanguageSerializer{
				typesMap: tt.fields.typesMap,
			}
			if got := g.serializeDeclaration(tt.args.serializerInfo); (tt.shouldContain && !strings.Contains(got, tt.want)) ||
				(!tt.shouldContain && strings.Contains(got, tt.want)) {
				t.Errorf("serializeDeclaration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_goLanguageSerializer_serializeEnum(t *testing.T) {
	type fields struct {
		typesMap map[string]string
	}
	type args struct {
		enum           *enum
		serializerInfo *serializerInfo
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *generatedCode
		wantErr bool
	}{
		{
			name:   "Valid enum generator",
			fields: fields{typesMap: map[string]string{}},
			args: args{
				enum: &enum{
					name: "test",
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
				serializerInfo: &serializerInfo{
					packageName: "test",
				},
			},
			want: &generatedCode{
				fileName: "test.go",
				code:     "type test int\n\nconst (\n\ttestFirst = test(5)\n\ttestSecond = test(8)\n)",
			},
			wantErr: false,
		},
		{
			name:   "Valid enum generator without values",
			fields: fields{typesMap: map[string]string{}},
			args: args{
				enum: &enum{
					name:       "test",
					enumValues: []*enumValue{},
				},
				serializerInfo: &serializerInfo{
					packageName: "test",
				},
			},
			want: &generatedCode{
				fileName: "test.go",
				code:     "type test int\n\nconst (\n)",
			},
			wantErr: false,
		},
		{
			name:   "Valid enum generator with uppercase values",
			fields: fields{typesMap: map[string]string{}},
			args: args{
				enum: &enum{
					name: "test",
					enumValues: []*enumValue{
						{
							name:  "First",
							value: 5,
						},
						{
							name:  "second",
							value: 8,
						},
					},
				},
				serializerInfo: &serializerInfo{
					packageName: "test",
				},
			},
			want: &generatedCode{
				fileName: "test.go",
				code:     "type test int\n\nconst (\n\ttestFirst = test(5)\n\ttestSecond = test(8)\n)",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &goLanguageSerializer{
				typesMap: tt.fields.typesMap,
			}
			got, err := g.serializeEnum(tt.args.enum, tt.args.serializerInfo)

			if got != nil {
				// delete the header
				got.code = strings.Replace(got.code, g.serializeDeclaration(tt.args.serializerInfo), "", -1)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("serializeEnum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("serializeEnum() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_goLanguageSerializer_serializeMiddleware(t *testing.T) {
	g := newGoLanguageSerializer()

	info := &serializerInfo{packageName: "bla"}
	testClass := &class{
		name: "test",
		dataMembers: []*dataMember{
			{
				memberType: "a",
				name:       "int",
			},
		},
	}
	classCode, _ := g.serializeClass(testClass, info)

	testEnum := &enum{
		name: "test",
		enumValues: []*enumValue{
			{
				name:  "a",
				value: 5,
			},
			{
				name:  "b",
				value: 8,
			},
		},
	}
	enumCode, _ := g.serializeEnum(testEnum, info)

	type args struct {
		middleware     middleware
		serializerInfo *serializerInfo
	}
	tests := []struct {
		name    string
		args    args
		want    *generatedCode
		wantErr bool
	}{
		{
			name: "Serialize class",
			args: args{
				middleware:     testClass,
				serializerInfo: info,
			},
			want:    classCode,
			wantErr: false,
		},
		{
			name: "Serialize enum",
			args: args{
				middleware:     testEnum,
				serializerInfo: info,
			},
			want:    enumCode,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := g.serializeMiddleware(tt.args.middleware, tt.args.serializerInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("serializeMiddleware() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("serializeMiddleware() got = %v, want %v", got, tt.want)
			}
		})
	}
}
