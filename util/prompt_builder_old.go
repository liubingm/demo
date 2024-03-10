package util

/*
func BuildClaudePromt(codes []string) string {

	const systemPrompt = `You should anlysis the problem follow below steps:
	1. please explain each public interface and parameters.
	2. please make sure output only jsonformat.`

	const promptStart = `

Human: Below is a java language project floder, plase read the code and help me to analyze it.
Help me to find class definition and fill in the "className" of the jsonformat.
In the class, there might be public interface and class methods, please fill in the "interfaceName" and "methodName" of the jsonformat.
In the class, there might be class fields, please fill in the "fieldName" of the jsonformat.
Please fill in the "methodAnalysis" with your analysis of the code.
Each code block contains a java source file, please anlysis it one by one and out put your anlysis result together.

`
	const promptEnd = `Output Json Format:
<jsonformat>
{
	"files": [
		{
		  "fileName": "ExampleClass.java",
		  "className": "ExampleClass",
		  "packageName": "com.example",
		  "interfaces": [
			{
			  "interfaceName": "ExampleInterface",
			  "methods": [
				{
				  "methodName": "exampleMethod",
				  "returnType": "void",
				  "parameters": [
					{
					  "paramName": "param1",
					  "paramType": "String"
					},
					{
					  "paramName": "param2",
					  "paramType": "int"
					}
				  ],
				  "accessModifier": "public"
				  "methodAnalysis": "xxx"
				}
			  ]
			}
		  ],
		  "classMethods": [
			{
			  "methodName": "classMethod",
			  "returnType": "int",
			  "parameters": [
				{
				  "paramName": "input",
				  "paramType": "String"
				}
			  ],
			  "accessModifier": "private"
			  "methodAnalysis": "xxx"
			}
		  ],
		  "classFields": [
			{
			  "fieldName": "exampleField",
			  "fieldType": "String",
			  "accessModifier": "private"
			}
		  ]
		},
		{
		  "fileName": "AnotherClass.java",
		  "className": "AnotherClass",
		  "packageName": "com.example.other",
		  "interfaces": [],
		  "classMethods": [
			{
			  "methodName": "anotherMethod",
			  "returnType": "void",
			  "parameters": [
				{
				  "paramName": "map",
				  "paramType": "Map<String, Integer>"
				}
			  ],
			  "accessModifier": "public"
			  "methondAnalysis": "xxx"
			}
		  ],
		  "classFields": [
			{
			  "fieldName": "anotherField",
			  "fieldType": "int",
			  "accessModifier": "protected"
			}
		  ]
		}
		// you may add more java code
	  ]
}
</jsonformat>


Assistant:{
	"files": [`

	var codePrompt string

	for _, code := range codes {
		codePrompt += "<code>" + code + "</code>\n"
	}

	prompt := systemPrompt + promptStart + codePrompt + promptEnd

	fmt.Println(prompt)

	return prompt
}*/

func BuildGPTPrompt(codes []string) (string, string) {

	// System prompt preparation
	systemPrompt := "You are a expert of software developer, you can help to review the source code document.\n\n"

	userPrompt := "Please analyis below source code, put your analysis result in json in design_doc block:"

	for _, code := range codes {

		userPrompt += "<code>/n " + code + "<code/n>"
	}

	userPrompt += `

	You should anlysis the problem follow below steps:
	1. Please list all the public function in the provided source code files.
	2. Please list how many source code row in the source code files.
	
	
Output Json Format:
<design_doc>
{
	"filename1": {
	  "functions": [
		"Functions1",
		"Functions2",
		...
	  ],
	  "row_No": xxx
	},
	"filename2": {
	  "functions": [
		"Functions1",
		"Functions3",
		...
	  ],
	  "row_No": xxx
	},
	...
  }
</design_doc>
	`

	return systemPrompt, userPrompt
}
