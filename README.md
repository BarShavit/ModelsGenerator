# ModelsGenerator

## Overview

 In microservices architecture, each service may be develop in a different programming language.<br/>
 Those services communicate with data structures, which serialized to a given format (ex. JSON).<br/>
 As a result, you need to build the same model in few different languages and make it exactly the same.<br/>
 
 This tool gets as a parameter a data structure written in a "gen" file, and knows to convert it to different languages.<br/>
 The supported languages at the moment are: Go, Typescript, Kotlin and C#.
 
 ## Install
 Copy the release file to a directory and put the directory in the path enivorment variable.
 
 ## Usage
 The tool is using "gen" file type. The basic sysntax will be specified later.<br/>
 To use it, open a terminal in the gen file directory.<br/>
 Generate the structures by the command ```fileName.gen language:packageName language:packagename ...```<br/>
 While "language" can be go, c#, typescript or kotlin. Package name is an extra data that can generate the files within the given package.<br/>
 It won't effect typescript. You can avoid writing ```:packageName``` and the files will be created without package / namespace - and you will need to add it yourself.<br/>
 
 ## Output
 After running a generation on a gen file, the output will be in the given language folder within a folder with the generate time.<br/>
 Go output will be in "go" folder, Kotlin in "kotlin", typescript in "typescript" and C# in "c#".<br/>
 The classes will have JSON annotations if needed so all the JSON serialize will be into camelcase names. The managed JSON library in C# is the standard Newtownsoft.JSON.
 
 # Gen File
 
 Gen file is a simple text file represent a data structures. A data structure can be a class (struct) or an enum.<br/>
 
 ### Supported Types
 Gen files support primitives types, maps, lists and another given types (which suppose to be another classes / enums).<br/>
 The supported primitive types are ```bool, int, string, double, float, char, byte and date.
 
 ### File Structure
 Classes will be represented like:
 ```
 class className
 {
    memberName type
    memberName type
 }
 ```
 
 Enums will be represented like:
 ```
 enum enumName
 {
    valueName value
    valueName value
 }
 ```
 
 ## Examples
 
 We will use the following gen file, which contains 2 classes with a connection between them and an enum.<br/>
 For each data structure in the file, the generator will create us a file with the output.
 
 ```
class test
{
	first int
	second string
	third bool
	listType list<int>
	mapType map<int,string>
}

class anotherClass
{
    inner test
    someList list<int>
    randomMap map<int,test>
}

enum randomEnum
{
	value 5
	another 8
}
 ```
 * Notice: The generator will handle the uppercase / camelcase names according to the language conventions.
 Now, we will execute the generate by the command ```mg file.gen go:packageName c#:packageName kotlin:packageName typescript```
 
 ### Go
 In Go, we handle data members from another struct types as a pointer.<br/>
 The output:
 ##### test.go
 ```
package packageName

// **********************************
//	Generated by ModelsGenerator
//	2020-04-25T13:51:15+03:00
// **********************************

type Test struct {
	First int `json:"first"`
	Second string `json:"second"`
	Third bool `json:"third"`
	ListType []int `json:"listType"`
	MapType map[int]string `json:"mapType"`
}
 ```
 ##### anotherClass.go
```
package packageName

// **********************************
//	Generated by ModelsGenerator
//	2020-04-25T13:59:40+03:00
// **********************************

type AnotherClass struct {
	Inner *test `json:"inner"`
	SomeList []int `json:"someList"`
	RandomMap map[int]*test `json:"randomMap"`
}
```
 ##### randomEnum.go
 ```
package packageName

// **********************************
//	Generated by ModelsGenerator
//	2020-04-25T13:59:40+03:00
// **********************************

type randomEnum int

const (
	randomEnumValue = randomEnum(5)
	randomEnumAnother = randomEnum(8)
)
 ```
 
 ### Typescript
 If you create a class that contains another one as a data member, it will be imported automaticly.<br/>
 The output:
 ##### test.ts
 ```
// **********************************
//	Generated by ModelsGenerator
//	2020-04-25T13:59:40+03:00
// **********************************

export class Test {
	first: number;
	second: string;
	third: boolean;
	listType: number[];
	mapType: Map<number, string>;
}
 ```
 ##### anotherClass.ts
```
import { Test } from "./test";

// **********************************
//	Generated by ModelsGenerator
//	2020-04-25T13:59:40+03:00
// **********************************

export class AnotherClass {
	inner: Test;
	someList: number[];
	randomMap: Map<number, Test>;
}
```
 ##### randomEnum.ts
 ```
// **********************************
//	Generated by ModelsGenerator
//	2020-04-25T13:59:40+03:00
// **********************************

export enum RandomEnum {
	Value = 5,
	Another = 8
}
 ```
 
 ### Kotlin
 The data structures will be converted to data classes.<br/>
 The output:
 ##### test.kt
 ```
package packageName

// **********************************
//	Generated by ModelsGenerator
//	2020-04-25T13:59:40+03:00
// **********************************

data class Test(val first: Int, val second: String, val third: Boolean, val listType: List<Int>, val mapType: HashMap<Int, String>)
 ```
 ##### anotherClass.kt
```
package packageName

// **********************************
//	Generated by ModelsGenerator
//	2020-04-25T13:59:40+03:00
// **********************************

data class AnotherClass(val inner: Test, val someList: List<Int>, val randomMap: HashMap<Int, Test>)
```
 ##### randomEnum.kt
 ```
package packageName

// **********************************
//	Generated by ModelsGenerator
//	2020-04-25T13:59:40+03:00
// **********************************

enum class RandomEnum(val value: Int) {
	VALUE(5),
	ANOTHER(8)
}
 ```
 
 ### C#
 The JSON package by default is Newtonsoft.JSON.<br/>
 The output:
 ##### test.cs
 ```
using Newtonsoft.Json;
using System.Collections.Generic;

// **********************************
//	Generated by ModelsGenerator
//	2020-04-25T13:59:40+03:00
// **********************************

namespace packageName
{
	public class Test
	{
		[JsonProperty(PropertyName = "first")]
		public int First { get; set; }
		[JsonProperty(PropertyName = "second")]
		public string Second { get; set; }
		[JsonProperty(PropertyName = "third")]
		public bool Third { get; set; }
		[JsonProperty(PropertyName = "listType")]
		public List<int> ListType { get; set; }
		[JsonProperty(PropertyName = "mapType")]
		public Dictionary<int, string> MapType { get; set; }
	}
}
 ```
 ##### anotherClass.cs
```
using Newtonsoft.Json;
using System.Collections.Generic;

// **********************************
//	Generated by ModelsGenerator
//	2020-04-25T13:59:40+03:00
// **********************************

namespace packageName
{
	public class AnotherClass
	{
		[JsonProperty(PropertyName = "inner")]
		public Test Inner { get; set; }
		[JsonProperty(PropertyName = "someList")]
		public List<int> SomeList { get; set; }
		[JsonProperty(PropertyName = "randomMap")]
		public Dictionary<int, Test> RandomMap { get; set; }
	}
}
```
 ##### randomEnum.cs
 ```
// **********************************
//	Generated by ModelsGenerator
//	2020-04-25T13:59:40+03:00
// **********************************

namespace packageName
{
	public enum RandomEnum
	{
		Value = 5,
		Another = 8
	}
}
 ```
