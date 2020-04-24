package main

import (
	"reflect"
	"strings"
	"testing"
)

func Test_csharpLanguageSerializer_generateCode(t *testing.T) {
	g := newCsharpLanguageSerializer()

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

func Test_csharpLanguageSerializer_getType(t *testing.T) {
	type fields struct {
		typesMap map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   languageType
	}{
		{name: "Valid type", fields: fields{typesMap: map[string]string{}}, want: LanguageTypeCSharp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &csharpLanguageSerializer{
				typesMap: tt.fields.typesMap,
			}
			if got := g.getType(); got != tt.want {
				t.Errorf("getType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_csharpLanguageSerializer_getTypeName(t *testing.T) {
	type fields struct {
		typesMap map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "Valid type name", fields: fields{typesMap: map[string]string{}}, want: "c#"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &csharpLanguageSerializer{
				typesMap: tt.fields.typesMap,
			}
			if got := g.getTypeName(); got != tt.want {
				t.Errorf("getTypeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_csharpLanguageSerializer_serializeClass(t *testing.T) {
	info := &serializerInfo{packageName: "main"}

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
				imports: []string{"Newtonsoft.Json"},
			},
			want: &generatedCode{
				fileName: "test.cs",
				code: "\tpublic class Test\n\t{\n\t\t[JsonProperty(PropertyName = \"a\")]\n\t\tpublic string A { get; set; }\n" +
					"\t\t[JsonProperty(PropertyName = \"b\")]\n\t\tpublic double B { get; set; }\n\t}\n}",
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
				imports: []string{"Newtonsoft.Json", "System.Collections.Generic"},
			},
			want: &generatedCode{
				fileName: "test.cs",
				code: "\tpublic class Test\n\t{\n\t\t[JsonProperty(PropertyName = \"a\")]\n\t\tpublic string A { get; set; }\n" +
					"\t\t[JsonProperty(PropertyName = \"b\")]\n" +
					"\t\tpublic List<int> B { get; set; }\n\t}\n}",
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
				imports: []string{"Newtonsoft.Json", "System.Collections.Generic"},
			},
			want: &generatedCode{
				fileName: "test.cs",
				code: "\tpublic class Test\n\t{\n\t\t[JsonProperty(PropertyName = \"a\")]\n\t\tpublic string A { get; set; }\n" +
					"\t\t[JsonProperty(PropertyName = \"b\")]\n" +
					"\t\tpublic List<Bla> B { get; set; }\n\t}\n}",
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
				imports: []string{"Newtonsoft.Json", "System.Collections.Generic"},
			},
			want: &generatedCode{
				fileName: "test.cs",
				code: "\tpublic class Test\n\t{\n\t\t[JsonProperty(PropertyName = \"a\")]\n\t\tpublic string A { get; set; }\n" +
					"\t\t[JsonProperty(PropertyName = \"b\")]\n" +
					"\t\tpublic Dictionary<int, string> B { get; set; }\n\t}\n}",
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
				imports: []string{"Newtonsoft.Json", "System.Collections.Generic"},
			},
			want: &generatedCode{
				fileName: "test.cs",
				code: "\tpublic class Test\n\t{\n\t\t[JsonProperty(PropertyName = \"a\")]\n\t\tpublic string A { get; set; }\n" +
					"\t\t[JsonProperty(PropertyName = \"b\")]\n" +
					"\t\tpublic Dictionary<int, Bla> B { get; set; }\n\t}\n}",
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
				imports: []string{"Newtonsoft.Json"},
			},
			want: &generatedCode{
				fileName: "test.cs",
				code: "\tpublic class Test\n\t{\n\t\t[JsonProperty(PropertyName = \"a\")]\n\t\tpublic string A { get; set; }\n" +
					"\t\t[JsonProperty(PropertyName = \"b\")]\n" +
					"\t\tpublic Bla B { get; set; }\n\t}\n}",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newCsharpLanguageSerializer()
			got, err := g.serializeClass(tt.args.class, info)

			if got != nil {
				// delete the header
				got.code = strings.Replace(got.code, g.serializeDeclaration(tt.args.imports, info), "", -1)
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

func Test_csharpLanguageSerializer_serializeDeclaration(t *testing.T) {
	info := &serializerInfo{packageName: "main"}

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
				"using bla;",
				"using other;",
			},
			shouldContain: true,
		},
		{
			name: "Without imports",
			fields: fields{
				typesMap: map[string]string{},
			},
			args:          args{imports: []string{}},
			want:          []string{"using"},
			shouldContain: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &csharpLanguageSerializer{
				typesMap: tt.fields.typesMap,
			}
			got := g.serializeDeclaration(tt.args.imports, info)

			for _, str := range tt.want {
				if (tt.shouldContain && !strings.Contains(got, str)) ||
					(!tt.shouldContain && strings.Contains(got, str)) {
					t.Errorf("serializeDeclaration() = %v, want %v", got, tt.want)
				}
			}

		})
	}
}

func Test_csharpLanguageSerializer_serializeEnum(t *testing.T) {
	info := &serializerInfo{packageName: "main"}
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
				fileName: "test.cs",
				code:     "\tpublic enum Test\n\t{\n\t\tFirst = 5,\n\t\tSecond = 8\n\t}\n}",
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
				fileName: "test.cs",
				code:     "\tpublic enum Test\n\t{\n\t}\n}",
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
				fileName: "test.cs",
				code:     "\tpublic enum Test\n\t{\n\t\tFirst = 5,\n\t\tSecond = 8\n\t}\n}",
			},
			wantErr: false,
		},
	}

	var imports []string

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &csharpLanguageSerializer{
				typesMap: tt.fields.typesMap,
			}
			got, err := g.serializeEnum(tt.args.enum, info)

			if got != nil {
				// delete the header
				got.code = strings.Replace(got.code, g.serializeDeclaration(imports, info), "", -1)
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

func Test_csharpLanguageSerializer_serializeMiddleware(t *testing.T) {
	g := newCsharpLanguageSerializer()

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
			got, err := g.serializeMiddleware(tt.args.middleware, info)
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
