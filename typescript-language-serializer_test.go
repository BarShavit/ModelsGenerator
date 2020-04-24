package main

import (
	"reflect"
	"strings"
	"testing"
)

func Test_typescriptLanguageSerializer_generateCode(t *testing.T) {
	g := newTypescriptLanguageSerializer()

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

func Test_typescriptLanguageSerializer_getType(t *testing.T) {
	type fields struct {
		typesMap map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   languageType
	}{
		{name: "Valid type", fields: fields{typesMap: map[string]string{}}, want: LanguageTypeTypescript},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &typescriptLanguageSerializer{
				typesMap: tt.fields.typesMap,
			}
			if got := g.getType(); got != tt.want {
				t.Errorf("getType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_typescriptLanguageSerializer_getTypeName(t *testing.T) {
	type fields struct {
		typesMap map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "Valid type name", fields: fields{typesMap: map[string]string{}}, want: "typescript"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &typescriptLanguageSerializer{
				typesMap: tt.fields.typesMap,
			}
			if got := g.getTypeName(); got != tt.want {
				t.Errorf("getTypeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_typescriptLanguageSerializer_serializeClass(t *testing.T) {
	type args struct {
		class   *class
		imports []string
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
				imports: []string{},
			},
			want: &generatedCode{
				fileName: "test.ts",
				code:     "export class Test {\n\ta: string;\n\tb: number;\n}",
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
				imports: []string{},
			},
			want: &generatedCode{
				fileName: "test.ts",
				code:     "export class Test {\n\ta: string;\n\tb: number[];\n}",
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
				imports: []string{"bla"},
			},
			want: &generatedCode{
				fileName: "test.ts",
				code:     "export class Test {\n\ta: string;\n\tb: Bla[];\n}",
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
				imports: []string{},
			},
			want: &generatedCode{
				fileName: "test.ts",
				code:     "export class Test {\n\ta: string;\n\tb: Map<number, string>;\n}",
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
				imports: []string{"Bla"},
			},
			want: &generatedCode{
				fileName: "test.ts",
				code:     "export class Test {\n\ta: string;\n\tb: Map<number, Bla>;\n}",
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
				imports: []string{"Bla"},
			},
			want: &generatedCode{
				fileName: "test.ts",
				code:     "export class Test {\n\ta: string;\n\tb: Bla;\n}",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newTypescriptLanguageSerializer()
			got, err := g.serializeClass(tt.args.class)

			if got != nil {
				// delete the header
				got.code = strings.Replace(got.code, g.serializeDeclaration(tt.args.imports), "", -1)
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

func Test_typescriptLanguageSerializer_serializeDeclaration(t *testing.T) {
	type fields struct {
		typesMap map[string]string
	}
	type args struct {
		imports []string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          []string
		shouldContain bool
	}{
		{
			name: "With imports",
			fields: fields{
				typesMap: map[string]string{},
			},
			args: args{imports: []string{
				"bla",
				"other",
			}},
			want: []string{
				"import { Bla } from \"./bla\";",
				"import { Other } from \"./other\";",
			},
			shouldContain: true,
		},
		{
			name: "Without imports",
			fields: fields{
				typesMap: map[string]string{},
			},
			args:          args{imports: []string{}},
			want:          []string{"import"},
			shouldContain: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &typescriptLanguageSerializer{
				typesMap: tt.fields.typesMap,
			}
			got := g.serializeDeclaration(tt.args.imports)

			for _, str := range tt.want {
				if (tt.shouldContain && !strings.Contains(got, str)) ||
					(!tt.shouldContain && strings.Contains(got, str)) {
					t.Errorf("serializeDeclaration() = %v, want %v", got, tt.want)
				}
			}

		})
	}
}

func Test_typescriptLanguageSerializer_serializeEnum(t *testing.T) {
	type fields struct {
		typesMap map[string]string
	}
	type args struct {
		enum *enum
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
			},
			want: &generatedCode{
				fileName: "test.ts",
				code:     "export enum Test {\n\tFirst = 5,\n\tSecond = 8\n}",
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
			},
			want: &generatedCode{
				fileName: "test.ts",
				code:     "export enum Test {\n}",
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
			},
			want: &generatedCode{
				fileName: "test.ts",
				code:     "export enum Test {\n\tFirst = 5,\n\tSecond = 8\n}",
			},
			wantErr: false,
		},
	}

	var imports []string

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &typescriptLanguageSerializer{
				typesMap: tt.fields.typesMap,
			}
			got, err := g.serializeEnum(tt.args.enum)

			if got != nil {
				// delete the header
				got.code = strings.Replace(got.code, g.serializeDeclaration(imports), "", -1)
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

func Test_typescriptLanguageSerializer_serializeMiddleware(t *testing.T) {
	g := newTypescriptLanguageSerializer()

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
	classCode, _ := g.serializeClass(testClass)

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
	enumCode, _ := g.serializeEnum(testEnum)

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
			got, err := g.serializeMiddleware(tt.args.middleware)
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
