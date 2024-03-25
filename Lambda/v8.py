import boto3
import gzip
import json
import base64
import os

bedrock = boto3.client("bedrock-runtime")

def handler(event, context):

    # Check if the event is from CloudWatch Logs
    if 'awslogs' in event:
        
        #读取日志信息到log_events
        # Decode the base64-encoded data
        compressed_payload = base64.b64decode(event['awslogs']['data'])
        # Decompress the gzip-compressed data
        uncompressed_payload = gzip.decompress(compressed_payload)
        # Parse the uncompressed payload as JSON
        log_events = json.loads(uncompressed_payload)
        
        #Step1:Prompt Engineering预处理Massage：
        #------将log加载到userPrompt、定义system_prompt、assistant_massage、output_Format等-----
        userPrompt_prefix = "Here is some AWS IoTCore log data, please help me to analyze them"
        userPrompt = "\n".join([userPrompt_prefix, log_events['logEvents'][0]['message']]) 
        
        assistant_message =  {"role": "assistant", "content": "<emoji>"}
        
        system_prompt = """
            Human: You are an AI-trained MQTT Operations Expert. Please follow the instructions below and help me analyze my MQTT log:
                1. Please help analyze the reason for any errors.
                2. Locate the error log lines.
                3. According to the reason, list several solutions.
                
            Please output in JSON format as follows:
            {
                "Reason": "",
                "Solution": ""
            }
        """
        
        Prompt_Input = json.dumps({
            "anthropic_version": "bedrock-2023-05-31",
            "max_tokens": 1000,
            "system":system_prompt,
            "messages": [
                {
                    "role": "user", 
                    "content": userPrompt
                },
                #{
                 #   "role": "assistant", 
                  #  "content": assistant_message
                #}
            ],
            "temperature": 0.8,
            "top_p": 0.999,
            #"top_k": int,
            #"stop_sequences": [string]
        })
        print(f"Step-1 Input: {Prompt_Input}")
        
        #Step2:Call the claude3 on the Amazon Bedrock
        response = bedrock.invoke_model(
            body=Prompt_Input, 
            modelId="anthropic.claude-3-sonnet-20240229-v1:0"
            )
        response_body = json.loads(response.get('body').read())
        response_body2 = json.dumps(response_body)
        print(f"Step-2 Claude_response: {response_body2}")
        
        #Step3: 从response中读取content信息
        text_value = response_body["content"][0]["text"]
        json_text = json.loads(text_value)
        results = json.dumps(json_text)
        print(results)

            # Add additional processing or transformation logic here
    else:
        # Handle other types of events
        print(f"Received event: {event}")

        # Add additional processing logic for other event types
    #返回结果,发布到SNS
    return results